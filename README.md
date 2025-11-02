# User Manager (Go + AWS Lambda)

A **serverless user management API** written in Go and deployed on **AWS Lambda** using the **AWS CDK** (TypeScript).  
It supports user registration, login, password reset (SES), JWT authentication, and role-based access control with DynamoDB persistence.

---

## Features

- **Go AWS Lambda**
- **JWT Authentication**
- **User Registration / Login / Deactivation / External user check**
- **Role-based Access Control**
- **Email via AWS SES for password reset**
- **DynamoDB**
- **API Gateway**
- **CDK Infrastructure-as-Code**
- **In-memory + DynamoDB repositories**

---

## Project Structure
(Main folders)

```
user-manager/
├── cmd/
│   └── lambda/               # Lambda entrypoint (main.go)
├── infra/                    # AWS CDK Stack (TypeScript)
│   ├── bin/
│   │   └── user-manager.ts   # CDK app entrypoint
│   └── lib/
│       └── user-manager-stack.ts
├── internal/
│   ├── config/               # App config
│   ├── ddb/                  # DynamoDB client
│   ├── errs/                 # Custom error definitions
│   ├── httpx/                # Server start, Middleware (logger, auth, recover)
│   ├── mailer/               # SES + Local (SMTP) mailer services
│   ├── mocks/                # Mock mailer for tests
│   ├── routes/               # Route setup with auth
│   ├── ses/                  # SES initialization and client
│   └── user/
│       ├── dtos/             # Request/Response DTOs
│       ├── handler/          # HTTP handlers
│       ├── jwt/              # JWT management
│       ├── model/            # User + Role models
│       ├── repository/       # DynamoDB + in-memory repositories
│       ├── service/          # Business logic layer
│       └── password_hasher/  # Password hashing utility
└── internal/test/            # Integration tests (httptest)
```

---

## Local Development

### 1. Build and Run Locally
```bash
make run_dev
```
This runs the Lambda locally using Go’s native HTTP server on  `http://localhost:8080`, using Air for hot reload.

### 2. Run Tests
```bash
make test
```

---

## Cloud Deployment (AWS Lambda + CDK)

### Prerequisites
- Node.js 18+
- AWS CLI configured (`aws configure`)
- CDK installed globally:
  ```bash
  npm install -g aws-cdk
  ```
- AWS SES email verified (`MAIL_FROM_EMAIL`)
- set Parameters in AWS Systems Manager (SSM):
  
  /user-manager/app/jwt-secret (base 64 64-bit string - can be generated in a terminal `openssl rand -base64 64`)
  /user-manager/app/api-key

### Build Lambda binary
```bash
make bootstrap
```

### Deploy via CDK
```bash
make deploy
```

---

## Example API Requests

### Health check
```bash
curl https://<api-url>/health
```

### Register
```bash
curl -X POST https://<api-url>/register   -H "Content-Type: application/json"   -d '{"email":"user@example.com","password":"StrongP@ssw0rd12345","roles":["user"]}'
```

### Login
```bash
curl -X POST https://<api-url>/login   -H "Content-Type: application/json"   -d '{"email":"user@example.com","password":"StrongP@ssw0rd12345"}'
```

### Get User Data
```bash
curl -H "Authorization: Bearer <JWT_TOKEN>" https://<api-url>/users/data
```

### Admin List All Users
```bash
curl -H "Authorization: Bearer <ADMIN_JWT>" https://<api-url>/users
```

---

##  Makefile Commands

| Command | Description |
|----------|--------------|
| `make run_dev` | Runs the project locally with **Air** (live reload for Go). |
| `make bootstrap` | Builds the Go Lambda binary (`bootstrap`) for **Linux ARM64**, stripping debug info for AWS deployment. |
| `make zip` | Packages the built binary into a **ZIP** file for manual Lambda upload. |
| `make package` | Alias for `make zip`. |
| `make clean` | Removes the compiled binary and ZIP artifacts. |
| `make deploy` | Builds and deploys the Lambda + infrastructure via **AWS CDK**. |
| `make test` | Runs all Go unit and integration tests. |

---

## 📄 License

MIT © 2025 Danilo Barolo Martins de Lima  
All rights reserved.
