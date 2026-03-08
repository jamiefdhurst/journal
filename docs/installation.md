# Installation Guide

This guide covers how to install and configure Journal for self-hosting.

## Prerequisites

Journal stores its data in a SQLite database file. Before running the
application, decide where you want that file to live on the host — you will
need to make that path available to the process.

---

## Option 1: Binary (Linux x86-64)

Pre-built binaries for Linux x86-64 are attached to every
[GitHub release](https://github.com/jamiefdhurst/journal/releases).

### Download and run

Replace `<version>` with the release you want (e.g. `1.0.0`):

```bash
# Download the binary
curl -L -o journal \
  https://github.com/jamiefdhurst/journal/releases/download/v<version>/journal-bin_linux_x64-v<version>

# Make it executable
chmod +x journal

# Run it
./journal
```

The application listens on port `3000` by default. Open
`http://localhost:3000` in your browser.

### Persistent data

By default Journal writes its SQLite database to `$GOPATH/data/journal.db`.
Set `J_DB_PATH` to an absolute path to store it wherever you like:

```bash
J_DB_PATH=/var/lib/journal/journal.db ./journal
```

### Running as a service (systemd)

```ini
# /etc/systemd/system/journal.service
[Unit]
Description=Journal
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/journal
Restart=on-failure
Environment=J_DB_PATH=/var/lib/journal/journal.db
Environment=J_TITLE=My Journal

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now journal
```

---

## Option 2: Docker / Container Runtime

A Docker image is published to the GitHub Container Registry on every release:

```
ghcr.io/jamiefdhurst/journal:latest
ghcr.io/jamiefdhurst/journal:v<version>
```

### Docker

```bash
docker run -d \
  --name journal \
  -p 3000:3000 \
  -v /var/lib/journal:/go/data \
  ghcr.io/jamiefdhurst/journal:latest
```

Pass configuration via `-e` flags:

```bash
docker run -d \
  --name journal \
  -p 3000:3000 \
  -v /var/lib/journal:/go/data \
  -e J_TITLE="My Journal" \
  -e J_DESCRIPTION="A place for my thoughts" \
  -e J_CREATE=1 \
  -e J_EDIT=1 \
  ghcr.io/jamiefdhurst/journal:latest
```

### Docker Compose

```yaml
services:
  journal:
    image: ghcr.io/jamiefdhurst/journal:latest
    restart: unless-stopped
    ports:
      - "3000:3000"
    volumes:
      - journal_data:/go/data
    environment:
      J_TITLE: My Journal
      J_DESCRIPTION: A place for my thoughts
      J_CREATE: "1"
      J_EDIT: "1"
      J_SESSION_KEY: "a-random-32-character-string-here"

volumes:
  journal_data:
```

```bash
docker compose up -d
```

### Podman

```bash
podman run -d \
  --name journal \
  -p 3000:3000 \
  -v /var/lib/journal:/go/data:Z \
  ghcr.io/jamiefdhurst/journal:latest
```

---

## Configuration

All configuration is done through environment variables. You can also place
them in a `.env` file in the working directory — the application reads this
file on startup before any environment variables are applied.

### Example `.env` file

```env
J_TITLE=My Journal
J_DESCRIPTION=A place for my thoughts
J_PORT=3000
J_DB_PATH=/var/lib/journal/journal.db
J_CREATE=1
J_EDIT=1
J_SESSION_KEY=a-random-32-character-string-here
```

### General

| Variable | Description | Default |
|---|---|---|
| `J_TITLE` | Title displayed in the journal UI | _(empty)_ |
| `J_DESCRIPTION` | HTML description shown in the journal UI | _(empty)_ |
| `J_PORT` | HTTP port to listen on | `3000` |
| `J_DB_PATH` | Path to the SQLite database file | `$GOPATH/data/journal.db` |
| `J_CREATE` | Set to `0` to disable creating new entries | _(enabled)_ |
| `J_EDIT` | Set to `0` to disable editing entries | _(enabled)_ |
| `J_POSTS_PER_PAGE` | Number of entries shown per page | `20` |
| `J_EXCERPT_WORDS` | Word count for entry previews on the index | `50` |
| `J_THEME` | Theme name from the `web/themes/` folder | `default` |
| `J_GA_CODE` | Google Analytics tag (e.g. `UA-XXXXX-X`) — omit to disable | _(disabled)_ |

### SSL/TLS

| Variable | Description |
|---|---|
| `J_SSL_CERT` | Path to PEM certificate file. Setting this enables HTTPS. |
| `J_SSL_KEY` | Path to PEM private key file. |

When `J_SSL_CERT` is set, session cookies automatically gain the `Secure`
flag so they are never sent over plain HTTP.

### Session and Cookie Security

| Variable | Description | Default |
|---|---|---|
| `J_SESSION_KEY` | 32-character ASCII encryption key for session data (AES-256). If unset, a random key is generated each startup — sessions will not survive restarts. | _(random)_ |
| `J_SESSION_NAME` | Name of the session cookie | `journal-session` |
| `J_COOKIE_DOMAIN` | Restricts cookies to a specific domain | _(current domain)_ |
| `J_COOKIE_MAX_AGE` | Cookie lifetime in seconds | `2592000` (30 days) |
| `J_COOKIE_HTTPONLY` | Set to `0` to allow JavaScript access to cookies (not recommended) | `true` |

> **Security tip:** Always set `J_SESSION_KEY` to a stable, randomly-generated
> 32-character string in production so that user sessions survive application
> restarts. You can generate one with:
>
> ```bash
> LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 32; echo
> ```

---

## Putting it behind a reverse proxy

It is recommended to run Journal behind a reverse proxy such as Nginx or
Caddy so that you can terminate TLS centrally and use a standard port.

### Caddy example

```
journal.example.com {
    reverse_proxy localhost:3000
}
```

Caddy will obtain and renew a TLS certificate automatically.

### Nginx example

```nginx
server {
    listen 443 ssl;
    server_name journal.example.com;

    ssl_certificate     /etc/ssl/certs/journal.crt;
    ssl_certificate_key /etc/ssl/private/journal.key;

    location / {
        proxy_pass         http://localhost:3000;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
    }
}
```

#### Making your journal private with HTTP basic authentication

If you want to restrict access to your journal, Nginx can prompt visitors for
a username and password before forwarding requests to the application.

1. Install the `apache2-utils` package (Debian/Ubuntu) or `httpd-tools`
   (RHEL/Fedora) to get the `htpasswd` command:

   ```bash
   sudo apt install apache2-utils   # Debian/Ubuntu
   sudo dnf install httpd-tools     # RHEL/Fedora
   ```

2. Create a password file and add a user:

   ```bash
   sudo htpasswd -c /etc/nginx/.htpasswd yourname
   ```

   To add more users later, omit the `-c` flag (it would overwrite the file):

   ```bash
   sudo htpasswd /etc/nginx/.htpasswd anotheruser
   ```

3. Add `auth_basic` directives to your server block:

   ```nginx
   server {
       listen 443 ssl;
       server_name journal.example.com;

       ssl_certificate     /etc/ssl/certs/journal.crt;
       ssl_certificate_key /etc/ssl/private/journal.key;

       location / {
           auth_basic           "Journal";
           auth_basic_user_file /etc/nginx/.htpasswd;

           proxy_pass         http://localhost:3000;
           proxy_set_header   Host $host;
           proxy_set_header   X-Real-IP $remote_addr;
       }
   }
   ```

4. Reload Nginx to apply the change:

   ```bash
   sudo nginx -t && sudo systemctl reload nginx
   ```

Visitors will now see a browser login prompt before they can access any page.
If you want to keep the API public while protecting the web UI, use separate
`location` blocks:

```nginx
location /api/ {
    proxy_pass         http://localhost:3000;
    proxy_set_header   Host $host;
    proxy_set_header   X-Real-IP $remote_addr;
}

location / {
    auth_basic           "Journal";
    auth_basic_user_file /etc/nginx/.htpasswd;

    proxy_pass         http://localhost:3000;
    proxy_set_header   Host $host;
    proxy_set_header   X-Real-IP $remote_addr;
}
```
