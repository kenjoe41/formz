# Formz

Formz is a command-line tool written in Go that checks URLs for the presence of HTML forms. It supports reading URLs from a file or standard input, and optionally follows HTTP redirects. It detects if a website is protected by Cloudflare and skips such URLs.

## Features

- **Input Modes:**
  - Read URLs from a specified file.
  - Accept URLs via standard input.

- **HTTP Handling:**
  - Supports HTTP redirects (optionally configurable).
  - Uses retryable HTTP client for robustness against transient errors.

- **Cloudflare Detection:**
  - Skips URLs that are protected by Cloudflare, indicated by "Attention Required" and "Cloudflare" in the HTML content.

- **Form Detection:**
  - Identifies URLs that contain HTML forms and outputs them to standard output.

## Usage

### Command-line Flags

- `-file`: Path to a file containing URLs to check.
- `-follow-redirects`: (Optional) Whether to follow HTTP redirects (default: true).

### Example Usage

#### Read URLs from a file:
```sh
formz -file urls.txt
```

#### Read URLs from standard input:
```sh
cat urls.txt | formz
```

#### Disable following redirects:
```sh
formz -file urls.txt -follow-redirects=false
```

## Dependencies

- [github.com/hashicorp/go-retryablehttp](https://pkg.go.dev/github.com/hashicorp/go-retryablehttp): Retryable HTTP client for resilient HTTP requests.
- [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html): HTML parsing library for detecting forms in HTML content.

## Installation

To install and use Form, make sure you have Go installed. Then, run:

```sh
go install github.com/kenjoe41/formz@latest
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.