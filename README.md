
# env0

`env0` is a command-line interface for interacting with the **env0 API**, a zero-trusted environment variables storage service.


[![GitHub release](https://img.shields.io/github/v/release/jibaru/env0.svg?style=flat-square)](https://github.com/jibaru/env0/releases/latest) [![Build Status](https://img.shields.io/github/actions/workflow/status/docker/compose/ci.yml?label=ci&logo=github&style=flat-square)](https://github.com/docker/compose/actions?query=workflow%3Aci)

---

## Table of Contents

- [env0](#env0)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
    - [Download the latest release](#download-the-latest-release)
    - [Using go](#using-go)
  - [Getting Started](#getting-started)
  - [Commands](#commands)
    - [Authentication](#authentication)
    - [App Management](#app-management)
    - [Environment Operations](#environment-operations)
    - [User Management](#user-management)
  - [Configuration](#configuration)
  - [Examples](#examples)
  - [Contributing](#contributing)

---

## Features

* **Zero-Trust Encryption**: Client-side encryption ensures all values are encrypted before transmission and never stored in plaintext.
* **Multi-Environment Support**: Organize variables under named environments (e.g., `dev`, `staging`, `prod`).
* **Access Control**: Grant or revoke per-user permissions at the app or environment level.
* **RESTful API**: Simple CRUD endpoints for users, apps, and environments.
* **JSON-Based Storage**: Flexible schema allowing arbitrary key/value pairs per environment.

---

## Installation

### Download the latest release

[See releases](https://github.com/Jibaru/env0/releases)

### Using go

Ensure you have Go 1.24 or newer installed.

```bash
go install github.com/Jibaru/env0/cmd/env0@v
```

Or you can build it by yourself.

```bash
git clone https://github.com/Jibaru/env0.git
cd env0/cli
go build -o env0 ./cmd/env0
# (Optional) Move binary into your PATH:
mv env0 /usr/local/bin/
```

---

## Getting Started

Before interacting with your first app, you need to create an account and authenticate.

```bash
# Sign up for a new account
env0 signup <username> <email> <password>

# Log in to env0 (stores credentials locally)
env0 login <username> <password>
```

Once authenticated, you can initialize, clone, and manage apps.

---

## Commands

### Authentication

| Command  | Description                      |
| -------- | -------------------------------- |
| `signup` | Create a new user account        |
| `login`  | Authenticate an existing account |
| `logout` | Remove local credentials         |

### App Management

| Command       | Description                          |
| ------------- | ------------------------------------ |
| `init <app>`  | Create a new app repository          |
| `clone <app>` | Download an existing app's env files |

### Environment Operations

> Important: environment called "default" is reserved for `.env` file.

| Command               | Description                                    |
| --------------------- | ---------------------------------------------- |
| `pull [<env>]`        | Fetch latest variables to local `.env` files       |
| `push [<env>]`        | Upload local `.env` files to remote service.    |

### User Management

| Command              | Description                                |
| -------------------- | ------------------------------------------ |
| `adduser <username>` | Grant a user access to the current app |
| `deluser <username>` | Revoke a user's access to the current app                     |

---

## Configuration

By default, `env0` stores credentials data in `$HOME/.env0/`.

---

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

# Pull ".env" file only
env0 pull default

# Pull ".env.prod" file only
env0 pull prod

# Push all env files
env0 push

# Push ".env.dev" file only
env0 push dev

# Push only contents from ".env" file
env0 push default

# Add a new user 'bob' to your app
env0 adduser bob

# Remove user 'bob' from your app
env0 deluser bob
```

---

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/foo`)
3. Commit your changes (`git commit -m "feat: add foo"`)
4. Push to the branch (`git push origin feature/foo`)
5. Open a Pull Request

