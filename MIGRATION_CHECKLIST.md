# Migration Checklist

Use this checklist to track your progress in migrating from the internal core package to the external GRA framework.

## Framework Publishing

- [x] Create GitHub repository for GRA framework
- [x] Push framework code to GitHub
- [x] Create v1.0.0 release tag
- [x] Verify module can be imported with `go get`

## Project Dependencies

- [x] Add GRA framework to go.mod
  ```bash
  go get github.com/lamboktulussimamora/gra@v1.0.0
  ```
- [x] Run `go mod tidy` to update dependencies

## Code Updates

- [x] Update imports in cmd/core-api/main.go
- [x] Update imports in internal/interface/handler/example_handler.go
- [x] Update any other files using the internal core package
- [x] Fix any function name changes or API differences
- [x] Update any custom middleware to use the new framework
- [x] Remove internal/core package

## Testing

- [x] Run unit tests
  ```bash
  go test ./...
  ```
- [x] Run framework integration tests
  ```bash
  go test ./tests/unit/gra_test.go
  ```
- [x] Run API integration tests
  ```bash
  ./tests/integration/run_api_tests.sh
  ```
- [x] Complete manual testing checklist from TEST_PLAN.md

## Cleanup

- [ ] Remove internal/core package (once everything works)
  ```bash
  rm -rf internal/core
  ```
- [ ] Run final tests to confirm everything still works
- [ ] Commit changes
  ```bash
  git add .
  git commit -m "Migrate from internal core to external GRA framework"
  ```

## Documentation

- [ ] Update README.md with migration notes
- [ ] Update any internal documentation referencing the core package
- [ ] Document any issues encountered during migration
