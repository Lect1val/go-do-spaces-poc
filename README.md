# go-do-spaces-poc

A proof-of-concept Go application demonstrating file upload, delete, and list operations with DigitalOcean Spaces using the AWS S3-compatible API.

## Features

- ✅ Upload files to DigitalOcean Spaces
- ✅ Delete files from DigitalOcean Spaces
- ✅ List all objects in a bucket
- ✅ Set lifecycle policies to auto-delete inactive files
- ✅ Manage lifecycle policies (list, delete)
- ✅ RESTful API with Gin framework
- ✅ AWS SDK v2 for S3-compatible operations
- ✅ Environment-based configuration

## Prerequisites

- Go 1.21 or higher
- DigitalOcean Spaces account with:
  - Access Key
  - Secret Key
  - Bucket name
  - Region and endpoint information

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-do-spaces-poc
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the project root:
```bash
cp .env.example .env
```

4. Configure your `.env` file with your DigitalOcean Spaces credentials:
```env
DO_SPACES_KEY=your_access_key_here
DO_SPACES_SECRET=your_secret_key_here
DO_SPACES_ENDPOINT=https://sgp1.digitaloceanspaces.com
DO_SPACES_REGION=sgp1
DO_SPACES_BUCKET=your_bucket_name
```

## Configuration

| Environment Variable | Description | Example |
|---------------------|-------------|---------|
| `DO_SPACES_KEY` | DigitalOcean Spaces access key | `YOUR_ACCESS_KEY` |
| `DO_SPACES_SECRET` | DigitalOcean Spaces secret key | `YOUR_SECRET_KEY` |
| `DO_SPACES_ENDPOINT` | Spaces endpoint URL | `https://nyc3.digitaloceanspaces.com` |
| `DO_SPACES_REGION` | AWS region (required by SDK) | `us-east-1` |
| `DO_SPACES_BUCKET` | Your Spaces bucket name | `my-bucket` |

### Available Regions

- `nyc3.digitaloceanspaces.com` - New York 3
- `sfo3.digitaloceanspaces.com` - San Francisco 3
- `ams3.digitaloceanspaces.com` - Amsterdam 3
- `sgp1.digitaloceanspaces.com` - Singapore 1
- `fra1.digitaloceanspaces.com` - Frankfurt 1

## Usage

### Run with Make

```bash
# Run the application
make run

# Build the binary
make build

# Clean build artifacts
make clean

# Run tests
make test

# Download and tidy dependencies
make deps
```

### Run directly

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Upload File

**Endpoint:** `POST /upload`

**Content-Type:** `multipart/form-data`

**Parameters:**
- `file` (required): The file to upload
- `path` (optional): Custom folder path (default: `uploads/`)

**Example using cURL:**
```bash
curl -X POST http://localhost:8080/upload \
  -F "file=@/path/to/your/file.jpg" \
  -F "path=images/"
```

**Success Response:**
```json
{
  "url": "https://nyc3.digitaloceanspaces.com/your-bucket/images/file.jpg"
}
```

### 2. Delete File

**Endpoint:** `DELETE /delete/:key`

**Parameters:**
- `key` (required): The object key/path to delete (URL parameter)

**Example using cURL:**
```bash
curl -X DELETE http://localhost:8080/delete/uploads/file.jpg
```

**Success Response:**
```json
{
  "message": "deleted successfully"
}
```

### 3. List Files

**Endpoint:** `GET /list`

**Example using cURL:**
```bash
curl http://localhost:8080/list
```

**Success Response:**
```json
{
  "files": [
    "uploads/file1.jpg",
    "uploads/file2.png",
    "images/photo.jpg"
  ]
}
```

### 4. Set Lifecycle Policy (Auto-Delete Inactive Files)

**Endpoint:** `POST /lifecycle/set`

**Content-Type:** `application/json`

**Parameters:**
- `prefix` (required): The folder/path prefix for files (e.g., `"temp/"`, `"uploads/"`)
- `expiration_days` (required): Number of days after which files will be automatically deleted (minimum: 1)
- `rule_id` (optional): Custom identifier for the rule (auto-generated if not provided)

**Example using cURL:**
```bash
curl -X POST http://localhost:8080/lifecycle/set \
  -H "Content-Type: application/json" \
  -d '{
    "prefix": "temp/",
    "expiration_days": 7,
    "rule_id": "delete-temp-files"
  }'
```

**Success Response:**
```json
{
  "message": "lifecycle policy set successfully",
  "rule_id": "delete-temp-files",
  "prefix": "temp/",
  "expiration_days": 7
}
```

**Note:** This creates a lifecycle policy that automatically deletes all files in the specified folder after the specified number of days. For example, setting `expiration_days: 7` for the `"temp/"` folder will delete all files in that folder that are older than 7 days.

### 5. List Lifecycle Policies

**Endpoint:** `GET /lifecycle/list`

**Example using cURL:**
```bash
curl http://localhost:8080/lifecycle/list
```

**Success Response:**
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
      "id": "delete-old-uploads",
      "status": "Enabled",
      "prefix": "uploads/old/",
      "expiration_days": 30
    }
  ]
}
```

### 6. Delete Lifecycle Policy

**Endpoint:** `DELETE /lifecycle/delete/:ruleId`

**Parameters:**
- `ruleId` (required): The ID of the lifecycle rule to delete (URL parameter)

**Example using cURL:**
```bash
curl -X DELETE http://localhost:8080/lifecycle/delete/delete-temp-files
```

**Success Response:**
```json
{
  "message": "lifecycle policy deleted successfully",
  "rule_id": "delete-temp-files"
}
```

## Project Structure

```
go-do-spaces-poc/
├── cmd/
│   └── main.go                 # Application entry point
│
├── config/
│   └── config.go               # Environment configuration loader
│
├── storage/
│   ├── client.go               # S3-compatible Spaces client initialization
│   ├── upload.go               # File upload logic
│   ├── delete.go               # File deletion logic
│   ├── list.go                 # List objects functionality
│   └── lifecycle.go            # Lifecycle policy management
│
├── handler/
│   └── storage_handler.go      # HTTP request handlers
│
├── router/
│   └── router.go               # Route registration
│
├── .env                        # Environment variables (create from .env.example)
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── Makefile                    # Build and run commands
└── README.md                   # This file
```
