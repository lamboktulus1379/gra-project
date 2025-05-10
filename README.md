# Go Clean Architecture Project with JWT Authentication

A RESTful API built with Go following Clean Architecture principles and implementing JWT-based authentication.

## Framework Migration Complete

This project has been successfully migrated from using an internal core framework to using the external [GRA Framework](https://github.com/lamboktulussimamora/gra). This migration provides better maintainability, version control, and easier updates.

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

## Framework Migration

This project initially used an internal core framework package located at `internal/core`. We've now migrated to using the external [GRA Framework](https://github.com/lamboktulussimamora/gra) which provides the same functionality.

### Migration Benefits

- Shared framework development and maintenance
- Better separation of concerns
- Easier updates and bug fixes
- Standardized API across projects

### Migration Steps

1. We added the GRA framework as a dependency:
   ```bash
   go get github.com/lamboktulussimamora/gra@v1.0.0
   ```

2. Updated imports from internal to external package:
   ```go
   // Before
   import "github.com/lamboktulussimamora/gra-project/internal/core"
   
   // After
   import (
       "github.com/lamboktulussimamora/gra/context"
       "github.com/lamboktulussimamora/gra/middleware"
       "github.com/lamboktulussimamora/gra/router"
       "github.com/lamboktulussimamora/gra/validator"
   )
   ```

3. Added a compatibility layer to smooth migration:
   ```go
   // Compatibility helper package for migration
   import "github.com/lamboktulussimamora/gra-project/internal/compatibility"
   ```

4. Updated function and type references:
   ```go
   // Before
   router := core.New()
   router.Use(core.LoggingMiddleware())
   
   // After
   r := router.New()
   r.Use(middleware.Logger())
   ```

5. Removed the internal core package after verification

### Running the Migration

If you're still using the internal core package, you can run the migration script:

```bash
./migrate.sh
```

For a more comprehensive migration, use the full update script which handles all import changes:

```bash
./full_update.sh
```

### Testing the Migration

We provide several testing methods to verify the migration:

1. Unit tests in `tests/unit/gra_test.go`
2. Integration tests in `tests/integration/`
3. Manual testing checklist in `TEST_PLAN.md`

To run the integration tests:

```bash
cd tests/integration
./run_api_tests.sh
```

## Development Tools

### Cleanup Scripts

During the migration process, various `.bak` files are created as backups. To clean these up, you can use:

```bash
# Clean up .bak files with interactive prompt
./cleanup.sh

# Or run the command directly:
find . -name "*.bak" -delete
```

The `.gitignore` file has been updated to exclude these backup files from version control.

### Backup Management

The migration process creates a full backup of the project before making changes. These backups are stored in:

```
backup_YYYYMMDDHHMMSS/
```

To restore from a backup:

```bash
cp -R ./backup_YYYYMMDDHHMMSS/* ./
```

Backup directories are excluded from version control via `.gitignore`.
