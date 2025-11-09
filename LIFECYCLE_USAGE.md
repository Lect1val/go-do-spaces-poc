# Lifecycle Policy Management - Usage Guide

This document provides examples and usage instructions for the lifecycle policy management endpoints.

## What are Lifecycle Policies?

Lifecycle policies automatically delete files from your DigitalOcean Space after a specified number of days. This is useful for:
- Temporary files that should be auto-cleaned
- Old uploads that are no longer needed
- Expired content
- Managing storage costs

## API Endpoints

### 1. Set Lifecycle Policy

Creates or updates a lifecycle policy for a specific folder.

```bash
curl -X POST http://localhost:8080/lifecycle/set \
  -H "Content-Type: application/json" \
  -d '{
    "prefix": "temp/",
    "expiration_days": 7,
    "rule_id": "delete-temp-files"
  }'
```

**Parameters:**
- `prefix`: The folder path (e.g., "temp/", "uploads/old/", "archives/")
- `expiration_days`: Number of days before deletion (minimum: 1)
- `rule_id`: Optional custom ID (auto-generated if not provided)

**Example Use Cases:**

Delete temporary files after 1 day:
```bash
curl -X POST http://localhost:8080/lifecycle/set \
  -H "Content-Type: application/json" \
  -d '{
    "prefix": "temp/",
    "expiration_days": 1
  }'
```

Delete old backups after 90 days:
```bash
curl -X POST http://localhost:8080/lifecycle/set \
  -H "Content-Type: application/json" \
  -d '{
    "prefix": "backups/",
    "expiration_days": 90,
    "rule_id": "cleanup-old-backups"
  }'
```

### 2. List All Lifecycle Policies

View all active lifecycle policies in your Space.

```bash
curl http://localhost:8080/lifecycle/list
```

**Response:**
```json
{
  "rules": [
    {
      "id": "delete-temp-files",
      "status": "Enabled",
      "prefix": "temp/",
      "expiration_days": 7
    },
    {
      "id": "cleanup-old-backups",
      "status": "Enabled",
      "prefix": "backups/",
      "expiration_days": 90
    }
  ]
}
```

### 3. Delete Lifecycle Policy

Remove a specific lifecycle policy by its rule ID.

```bash
curl -X DELETE http://localhost:8080/lifecycle/delete/delete-temp-files
```

**Response:**
```json
{
  "message": "lifecycle policy deleted successfully",
  "rule_id": "delete-temp-files"
}
```

## Important Notes

1. **Prefix Matching**: The policy applies to all files that start with the specified prefix.
   - `"temp/"` matches all files in the temp folder
   - `"uploads/"` matches all files in the uploads folder
   - `""` (empty string) matches ALL files in the bucket

2. **Expiration Time**: Files are deleted after they are older than the specified number of days based on their creation/modification time.

3. **Automatic Cleanup**: Once a policy is set, DigitalOcean Spaces automatically handles the deletion. You don't need to do anything else.

4. **Rule Updates**: Setting a policy with the same prefix or rule_id will update the existing rule.

5. **Multiple Policies**: You can have multiple policies for different folders.

## Testing the Implementation

### Step 1: Start the server
```bash
make run
# or
go run cmd/main.go
```

### Step 2: Set a test policy
```bash
curl -X POST http://localhost:8080/lifecycle/set \
  -H "Content-Type: application/json" \
  -d '{
    "prefix": "test/",
    "expiration_days": 1,
    "rule_id": "test-policy"
  }'
```

### Step 3: Verify the policy was created
```bash
curl http://localhost:8080/lifecycle/list
```

### Step 4: Delete the test policy
```bash
curl -X DELETE http://localhost:8080/lifecycle/delete/test-policy
```

### Step 5: Confirm deletion
```bash
curl http://localhost:8080/lifecycle/list
```

## Best Practices

1. **Use descriptive rule IDs**: Makes it easier to identify and manage policies
2. **Test with short expiration times first**: Start with 1-2 days for testing
3. **Be careful with empty prefixes**: An empty prefix applies to ALL files
4. **Document your policies**: Keep track of what policies are active and why
5. **Regular audits**: Periodically review active policies with `/lifecycle/list`

## Troubleshooting

### Policy not working?
- Check that your DigitalOcean Spaces account supports lifecycle policies
- Verify the prefix is correct (includes trailing slash if needed)
- Ensure expiration_days is at least 1

### Can't delete a policy?
- Make sure you're using the correct rule_id
- List all policies first to see the exact rule_id

### Files not being deleted?
- Remember that deletion happens automatically after the specified days
- The policy applies to files based on their creation/modification date
- It may take up to 24 hours for DigitalOcean to process lifecycle rules

