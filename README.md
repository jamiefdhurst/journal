# Journal

![License](https://img.shields.io/github/license/jamiefdhurst/journal.svg)
[![Build](https://github.com/jamiefdhurst/journal/actions/workflows/build.yml/badge.svg)](https://github.com/jamiefdhurst/journal/actions/workflows/build.yml)
[![Latest Version](https://img.shields.io/github/release/jamiefdhurst/journal.svg)](https://github.com/jamiefdhurst/journal/releases)

A simple web-based journal written in Go. You can post, edit and view entries,
with the addition of an API.

It makes use of a SQLite database to store the journal entries.

[API Documentation](api/README.md) - also available via `openapi.yml` as a URL
when deployed.

[Installation Guide](docs/installation.md) - full installation guide covering
all methods, configuration reference, and reverse proxy setup.

[User Guide](docs/user-guide.md) - creating and editing entries, and
navigating the journal.

![Screenshot of Journal](/docs/screenshot.png)

## Purpose

Journal serves as an easy-to-read and simple Golang program for new developers 
to try their hand at modifying, extending and playing with. It deliberately has 
only one dependency to ensure that the full end-to-end flow of the system can 
be understood through standard Golang libraries.

It's also a nice little Journal that you can use to keep your thoughts in, or 
as a basic blog platform.

## Installation

See the [Installation Guide](docs/installation.md) for full details, including
configuration, running as a service, and setting up a reverse proxy.

### Homebrew (macOS)

```bash
brew tap jamiefdhurst/journal
brew install journal
```

### Docker / Container Runtime

```bash
docker run -d \
  --name journal \
  -p 3000:3000 \
  -v /var/lib/journal:/go/data \
  jamiefdhurst/journal:latest
```

Images are also published to the GitHub Container Registry as
`ghcr.io/jamiefdhurst/journal:latest`.

### Debian / Ubuntu (apt)

```bash
curl -fsSL https://jamiefdhurst.github.io/packages/journal.asc \
  | sudo tee /usr/share/keyrings/journal.asc > /dev/null

echo "deb [signed-by=/usr/share/keyrings/journal.asc] \
  https://jamiefdhurst.github.io/packages stable main" \
  | sudo tee /etc/apt/sources.list.d/journal.list

sudo apt update && sudo apt install journal
```

### CentOS / RHEL / Fedora (yum/dnf)

```bash
sudo tee /etc/yum.repos.d/journal.repo > /dev/null <<'EOF'
[journal]
name=Journal
baseurl=https://jamiefdhurst.github.io/packages/yum
enabled=1
gpgcheck=1
gpgkey=https://jamiefdhurst.github.io/packages/journal.asc
EOF

sudo yum install journal
```

### ZIP archive (all platforms)

Pre-built archives for Linux and macOS (amd64 and arm64) are attached to every
[GitHub release](https://github.com/jamiefdhurst/journal/releases). Download
the archive for your platform, extract it, and run the `journal` binary inside.

### Build from source

```bash
git clone https://github.com/jamiefdhurst/journal.git
cd journal
go mod download
go build -o journal ./cmd/journal
./journal
```

## Configuration through Environment Variables

The application uses environment variables to configure all aspects.

You can optionally supply these through a `.env` file that will be parsed before
any additional environment variables.

### General Configuration

* `J_CREATE` - Set to `0` to disable post creation
* `J_WEB_PATH` - Override the directory used to locate web assets (templates, static files, themes). Defaults to the directory containing the binary, or the current working directory.
* `J_DB_PATH` - Path to SQLite DB - default is `./data/journal.db`
* `J_DESCRIPTION` - Set the HTML description of the Journal
* `J_EDIT` - Set to `0` to disable post modification
* `J_EXCERPT_WORDS` - The length of the post shown as a preview/excerpt in the index, default `50`
* `J_GA_CODE` - Google Analytics tag value, starts with `UA-`, or ignore to disable Google Analytics
* `J_PORT` - Port to expose over HTTP, default is `3000`
* `J_POSTS_PER_PAGE` - Posts to display per page, default `20`
* `J_THEME` - Theme to use from within the _web/themes_ folder, defaults to `default`
* `J_TITLE` - Set the title of the Journal

### SSL/TLS Configuration

* `J_SSL_CERT` - Path to SSL certificate file for HTTPS (enables SSL when set)
* `J_SSL_KEY` - Path to SSL private key file for HTTPS

### Session and Cookie Security

* `J_SESSION_KEY` - 32-byte encryption key for session data (AES-256). Must be exactly 32 printable ASCII characters. If not set, a random key is generated on startup (sessions won't persist across restarts).
* `J_SESSION_NAME` - Cookie name for sessions, default `journal-session`
* `J_COOKIE_DOMAIN` - Domain restriction for cookies, default is current domain only
* `J_COOKIE_MAX_AGE` - Cookie expiry time in seconds, default `2592000` (30 days)
* `J_COOKIE_HTTPONLY` - Set to `0` or `false` to allow JavaScript access to cookies (not recommended). Default is `true` for XSS protection.

**Note:** When `J_SSL_CERT` is configured, session cookies automatically use the `Secure` flag to prevent transmission over unencrypted connections.

## Layout

The project layout follows the standard set out in the following document:
[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

* `/api` - API documentation
* `/internal/app/controller` - Controllers for the main application
* `/internal/app/model` - Models for the main application
* `/internal/app/router` - Implementation of router for given app
* `/pkg/adapter` - Adapters for connecting to external services
* `/pkg/controller` - Controller logic
* `/pkg/database` - Database connection logic
* `/pkg/router` - Router for handling services
* `/test` - API tests
* `/test/data` - Test data
* `/test/mocks` - Mock files for testing
* `/web/static` - Compiled static public assets
* `/web/templates` - View templates
* `/web/themes` - Front-end themes, a default theme is included

## Development

### Back-end

The back-end can be extended and modified following the folder structure above. 
Tests for each file live alongside and are designed to be easy to read and as 
functionally complete as possible.

The easiest way to develop incrementally is to use a local go installation and
run your Journal as follows:

```bash
go run ./cmd/journal
```

Naturally, any changes to the logic or functionality will require a restart of 
the binary itself.

#### Dependencies

The application has the following dependencies (using go.mod and go.sum):

- [github.com/ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3)
- [github.com/gomarkdown/markdown](https://github.com/gomarkdown/markdown)

This can be installed using the following commands from the journal folder:

```bash
go get -v ./...
```

#### Templates

The templates are in `html/template` format in _web/templates_ and are used 
within each of the controllers. These can be modified while the binary stays 
loaded, as they are loaded on the fly by the application as it runs and serves 
content.

### Front-end

The front-end source files are intended to be divided into themes within the
_web/themes_ folder. Each theme can include icons and a CSS stylesheet.

A simple, basic and minimalist "default" theme is included, but any other 
themes can be built and modified.

### Building/Testing

All pushed code is currently built using GitHub Actions to test PRs, build 
packages and create releases.

To test locally, simply use:

```bash
go test -v ./...
```

