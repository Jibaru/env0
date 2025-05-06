# env0

`env0` is a command-line interface for interacting with the **env0 API**, a zero-trusted environment variables storage service.

## Features

* **Zero-Trust Encryption:** All environment variable values are encrypted at rest.
* **Multi-Environment Support:** Store variables under arbitrary environment names (dev, prod, staging, etc.).
* **Access Control:** Grant or revoke per-user access to each app.
* **RESTful API:** Simple endpoints for CRUD operations on users, apps, and envs.
* **JSON-Based Storage:** Flexibly store any number of key/value pairs per environment.

## Installation

Ensure you have Go 1.24+ installed. Then:

```bash
git clone https://github.com/Jibaru/env0.git
cd env0-cli
go build -o env0 ./cmd/env0
# Optional: Move to your PATH
env0 ~/.env0 /usr/local/bin/
```

## Examples

```bash
# Sign up for a new account
env0 signup alice alice@example.com supersecret

# Log in to your account
env0 login alice supersecret

# Initialize a new app named 'myapp'
env0 init myapp

# Clone an existing app 'alice/myapp' to local .env files
env0 clone alice/myapp

# Pull latest environment variables
env0 pull

# Push all env files or just 'dev'
env0 push
env0 push dev

# Add a new user 'bob' to your app
env0 adduser bob

# Remove user 'bob' from your app
env0 deluser bob
```
