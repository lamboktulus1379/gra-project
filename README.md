# Go Clean Architecture Project with JWT Authentication

A RESTful API built with Go following Clean Architecture principles and implementing JWT-based authentication.

## Project Structure

The project follows Clean Architecture principles with the following layers:

```
.
├── cmd/                     # Application entry points
│   └── api/                # Main API server
│       └── main.go
├── internal/                # Private application code
│   ├── domain/             # Enterprise business rules
│   │   ├── auth/          # Authentication domain
│   │   │   ├── jwt.go
│   │   │   └── password.go
│   │   └── user/          # User domain
│   │       ├── repository.go
│   │       └── user.go
│   ├── interface/          # Interface adapters
│   │   ├── common/        # Shared utilities
│   │   │   └── response.go
│   │   ├── handler/       # HTTP handlers
│   │   │   ├── common.go
│   │   │   ├── hello_handler.go
│   │   │   ├── protected_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/    # HTTP middleware
│   │   │   └── auth_middleware.go
│   │   └── repository/    # Data storage implementations
│   │       └── memory_user_repository.go
│   └── usecase/           # Application business rules
│       └── user_usecase.go
```

## Features

- **Clean Architecture**: Separation of concerns with domain, use cases, and interface layers
- **JWT Authentication**: Secure authentication using JSON Web Tokens
- **Password Hashing**: Argon2id algorithm for secure password storage
- **Protected Routes**: Middleware-based route protection
- **RESTful API**: Standard HTTP methods and status codes

## Authentication Flow

### Registration

1. Client sends a POST request to `/register` with user details
2. Password is securely hashed using Argon2id
3. User is persisted in the repository
4. Client receives the created user information

### Login

1. Client sends a POST request to `/login` with email and password
2. Server verifies the credentials
3. On success, server generates and returns a JWT token
4. Client stores the token for future authenticated requests

### Accessing Protected Resources

1. Client sends a request to a protected endpoint with the JWT token in the Authorization header
2. Auth middleware validates the token
3. If valid, the request proceeds to the handler with user claims in the context
4. If invalid, an appropriate error response is returned

## API Endpoints

| Method | Endpoint   | Description                  | Authentication |
|--------|------------|------------------------------|----------------|
| GET    | /hello     | Simple hello world endpoint  | Public         |
| POST   | /register  | User registration            | Public         |
| POST   | /login     | User authentication          | Public         |
| GET    | /profile   | User profile information     | Protected      |

## JWT Implementation

- **Token Generation**: Creates tokens with user data embedded as claims
- **Token Validation**: Verifies token integrity and expiration
- **Secret Key**: Configurable signing key (stored securely in production)
- **Expiration**: Configurable token lifetime

## Password Security

- **Argon2id**: Modern, secure password hashing algorithm
- **Salt Generation**: Unique random salt for each password
- **Configurable Parameters**: Memory, iterations, parallelism, salt and key length

## Development

### Prerequisites

- Go 1.16+

### Running the Application

```bash
# From the root directory
go run cmd/api/main.go
```

The server will start on port 8080.

### Example Requests

#### Register a User

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "secure_password_123"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "secure_password_123"
  }'
```

#### Access Protected Profile

```bash
curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Future Enhancements

- Database integration
- Refresh token mechanism
- Role-based authorization
- API rate limiting
- Request validation middleware
