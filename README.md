# **triple-s**

`triple-s` is a minimal local object storage service inspired by Amazon S3. It offers basic RESTful operations over a filesystem using Go.

### Project Structure

```plaintext
triple-s/
├── base/
│   └── data/              # Working storage directory
│       └── metadata.csv   # Metadata for buckets and objects
├── source/                # Source code
│   ├── flagcheck.go
│   ├── help.go
│   ├── methodhandler/
│   ├── delete.go
│   ├── get.go
│   ├── put.go
│   ├── router/
│   ├── structure/
│   └── tools/
├── main.go                # Entry point
└── README.md
```
## Source
Main directory containing all the logic and structure of the app.

source/structure/
Stores global variables and flag definitions.

Common struct definitions used across packages.

source/flagcheck.go
Validates CLI flags (like --dir, --port, --help).

Ensures no invalid directory names or empty values.

source/help.go
Provides the help text for CLI usage.

source/tools/
Contains InitSSS and helpers for file operations and metadata management (e.g., working with metadata.csv).

Can be extended to support additional tools like logging or analytics.

source/router/
Defines HTTP routes and connects them to handlers (GET, PUT, DELETE).

Uses a multiplexer to map endpoints like /bucket/object.

source/methodhandler/ (optional or combined in future)
If you choose to split logic per HTTP method, this folder can host files like get.go, put.go, delete.go.


### 🚀 HOW TO RUN
Start the server:
go run main.go --port 8080 --dir data

Server will run at: http://localhost:8080/
All buckets and objects will be saved inside base/data/.


### POSTMAN COMMANDS

1. Create a bucket
Method: PUT
URL: localhost:8080/bucketname

Replace bucketname with your bucket name, e.g. localhost:8080/photos

2. Upload an object
Method: PUT
URL: localhost:8080/bucketname/objectname
Body: Choose raw or binary to upload a file

Example: localhost:8080/photos/image.png

3. Download an object
Method: GET
URL: localhost:8080/bucketname/objectname

Example: localhost:8080/photos/image.png

4. Delete an object
Method: DELETE
URL: localhost:8080/bucketname/objectname

Example: localhost:8080/photos/image.png


### CLI Flags

| Flag     | Description                          | Example      |
| -------- | ------------------------------------ | ------------ |
| `--dir`  | Directory for storage inside `base/` | `--dir data` |
| `--help` | Show help screen                     | `--help`     |


### 📌 Notes
All paths must be relative. Absolute paths or parent directory traversal (e.g. ..) are not allowed.

All bucket and object metadata is saved to metadata.csv.

You can use tools like curl, Postman, or any HTTP client to test the API.

# Made by [aradilkha](https://platform.alem.school/git/aradilkha)