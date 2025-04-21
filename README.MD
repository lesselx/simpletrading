# 🛡️ OAuth2 Microservices Trading Platform

This project is a microservices-based application consisting of three main services:

- **Authentication Service** – OAuth2-compliant login, registration, and token issuance.
- **Data Service** – Periodically generates and stores data, exposes it via an API.
- **Trade Service** – Allows authenticated users to place trades under specific business rules.


## 🔧 Technologies Used

- **Go (Golang)**
- **Gorilla Mux / net/http**
- **GORM (SQLite)**
- **JWT**
- **OAuth2 (Google Login)**


## 🔧 Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/lesselx/simpletrading.git

```

### 2. Configure .env.local for Each Service

auth-service/.env
```bash
PORT=8080

# OAuth2 Google Login
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/oauth2/callback

# JWT
JWT_SECRET=super-secret-key

# Database
DB_PATH=auth.db
```


data-service/.env
```bash
PORT=8081

# JWT
JWT_SECRET=super-secret-key

# Database
DB_PATH=data.db
```

trade-service/.env
```bash
PORT=8082

# JWT
JWT_SECRET=super-secret-key

# Database
DB_PATH=trade.db
```
### 3. Run each services

```bash
cd auth-service
go mod tidy
go run ./cmd/auth-service/main.go

cd data-service
go mod tidy
go run ./cmd/data-service/main.go

cd trade-service
go mod tidy
go run ./cmd/trade-service/main.go
