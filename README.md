# go-do-spaces-poc

A proof-of-concept Go application demonstrating file upload, delete, and list operations with DigitalOcean Spaces using the AWS S3-compatible API.

## Features

- ✅ Upload files to DigitalOcean Spaces
- ✅ Delete files from DigitalOcean Spaces
- ✅ List all objects in a bucket
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
│   └── list.go                 # List objects functionality
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
