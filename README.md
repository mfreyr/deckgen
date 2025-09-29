# Resume Adapter Backend

This project is a Go-based backend service that takes a resume and a job advertisement as input, uses an LLM to adapt the resume to the job's requirements, and outputs a structured JSON presentation.

## Features

- **Extensible Parser System**: Uses a `Parser` interface to allow for easy addition of new document formats (e.g., `.docx`, raw text). The initial implementation supports PDF files.
- **Pluggable LLM Clients**: Features an `LLM` interface, making it simple to switch between different language model providers (e.g., OpenAI, Anthropic).
- **Simple API**: Exposes a single `POST /generate` endpoint for processing resumes.

## Prerequisites

This application requires the following to be installed on your system:

1.  **Go**: Version 1.22 or later.
2.  **`poppler-utils`**: The PDF parser relies on the `pdftotext` command-line tool, which is part of the `poppler-utils` package.

    You can install it on Debian-based systems (like Ubuntu) using:
    ```bash
    sudo apt-get update && sudo apt-get install -y poppler-utils
    ```

## Installation & Running

1.  **Clone the repository**:
    ```bash
    git clone <repository-url>
    cd resume-adapter-backend
    ```

2.  **Install Go dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Set up your environment**:
    You need to provide an OpenAI API key. The application reads it from the `OPENAI_API_KEY` environment variable.
    ```bash
    export OPENAI_API_KEY="your-secret-api-key"
    ```

4.  **Run the server**:
    ```bash
    go run cmd/server/main.go
    ```
    The server will start on port `8080`.

## API Usage

Send a `POST` request to the `/generate` endpoint with a `multipart/form-data` payload.

-   **Endpoint**: `http://localhost:8080/generate`
-   **Method**: `POST`
-   **Body**: `multipart/form-data` with two fields:
    -   `resume`: The resume file (e.g., a PDF).
    -   `jobAd`: The plain text of the job advertisement.

### Example `curl` command:

```bash
curl -X POST \
  -F "resume=@/path/to/your/resume.pdf" \
  -F "jobAd=We are hiring a Senior Go Developer..." \
  http://localhost:8080/generate
```