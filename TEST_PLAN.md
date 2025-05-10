# Migration Test Plan

This document outlines the testing strategy for verifying the successful migration from the internal core package to the external GRA framework.

## 1. Unit Tests

### 1.1 Framework Component Tests

| Test | Description | Status |
|------|-------------|--------|
| Context Tests | Verify that `Context` objects behave the same in the new framework | �� |
| Router Tests | Test routing and path parameter matching | 🔄 |
| Middleware Tests | Check that middleware chaining works correctly | 🔄 |
| Validation Tests | Ensure validation rules work as before | 🔄 |

### 1.2 API Tests

| Test | Description | Status |
|------|-------------|--------|
| Handler Tests | Test all handlers with the new framework | 🔄 |
| Request/Response Tests | Verify correct handling of requests and responses | 🔄 |
| Authentication Tests | Check JWT handling | 🔄 |

## 2. Integration Tests

- Test all API endpoints with HTTP requests
- Verify status codes, response bodies, and headers
- Test authentication flows
- Test validation errors

## 3. Manual Testing Checklist

- [ ] Start the application using `go run cmd/core-api/main.go`
- [ ] Test the hello endpoint: `curl http://localhost:8080/`
- [ ] Test the user registration: `curl -X POST -H "Content-Type: application/json" -d '{"first_name":"John","last_name":"Doe","email":"john@example.com","password":"password123"}' http://localhost:8080/users`
- [ ] Test login endpoint
- [ ] Test protected endpoints with a valid JWT token
- [ ] Test validation errors by submitting invalid data

## 4. Performance Comparison

Compare the performance of the application before and after migration:

| Metric | Before | After | Difference |
|--------|--------|-------|------------|
| Response Time (avg) | TBD | TBD | TBD |
| Memory Usage | TBD | TBD | TBD |
| Throughput (req/s) | TBD | TBD | TBD |

## 5. Rollback Plan

In case of issues, follow these steps to roll back:

1. Restore from backup: `cp -R ./backup_<timestamp>/* ./`
2. Run tests to verify restoration: `go test ./...`
3. Remove the external dependency: Edit `go.mod` to remove `github.com/lamboktulussimamora/gra`
4. Update go.mod: `go mod tidy`
