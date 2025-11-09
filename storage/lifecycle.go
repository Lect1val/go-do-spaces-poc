package storage

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// SetLifecyclePolicy creates a lifecycle rule to delete objects in a specific folder after specified days
func SetLifecyclePolicy(ctx context.Context, client *s3.Client, bucket, prefix string, expirationDays int32, ruleID string) error {
	// Create the lifecycle rule with prefix filter
	rule := types.LifecycleRule{
		ID:     aws.String(ruleID),
		Status: types.ExpirationStatusEnabled,
		Prefix: aws.String(prefix),
		Expiration: &types.LifecycleExpiration{
			Days: aws.Int32(expirationDays),
		},
	}

	// Get existing lifecycle configuration
	getLifecycleOutput, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	})

	var rules []types.LifecycleRule
	if err == nil && getLifecycleOutput.Rules != nil {
		// Keep existing rules, but filter out any with the same ID or prefix
		for _, existingRule := range getLifecycleOutput.Rules {
			// Skip if same ID or same prefix to avoid duplicates
			shouldKeep := true
			if existingRule.ID != nil && *existingRule.ID == ruleID {
				shouldKeep = false
			} else if existingRule.Prefix != nil && *existingRule.Prefix == prefix {
				shouldKeep = false
			}
			if shouldKeep {
				rules = append(rules, existingRule)
			}
		}
	}

	// Add the new rule
	rules = append(rules, rule)

	// Apply the lifecycle configuration
	_, err = client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: rules,
		},
	})

	return err
}

// GetLifecyclePolicy retrieves all lifecycle rules for a bucket
func GetLifecyclePolicy(ctx context.Context, client *s3.Client, bucket string) ([]types.LifecycleRule, error) {
	output, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}
	return output.Rules, nil
}

// DeleteLifecyclePolicy removes a specific lifecycle rule by ID
func DeleteLifecyclePolicy(ctx context.Context, client *s3.Client, bucket, ruleID string) error {
	// Get existing lifecycle configuration
	getLifecycleOutput, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	// Filter out the rule to delete
	var rules []types.LifecycleRule
	for _, existingRule := range getLifecycleOutput.Rules {
		if existingRule.ID != nil && *existingRule.ID != ruleID {
			rules = append(rules, existingRule)
		}
	}

	// If no rules remain, delete the entire lifecycle configuration
	if len(rules) == 0 {
		_, err = client.DeleteBucketLifecycle(ctx, &s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucket),
		})
		return err
	}

	// Otherwise, update with remaining rules
	_, err = client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: rules,
		},
	})

	return err
}

// IsNoLifecycleConfigError checks if the error is due to no lifecycle configuration existing
func IsNoLifecycleConfigError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "NoSuchLifecycleConfiguration")
}
