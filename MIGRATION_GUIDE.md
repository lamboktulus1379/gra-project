# Migration Guide: Internal Core to GRA Framework

This guide details the steps to migrate from using the internal core package to the external GRA framework.

## Prerequisites

Before starting the migration, ensure that:
1. The GRA framework has been published to GitHub
2. You have created a release with the tag v1.0.0

## Migration Steps

### 1. Add the GRA Framework Dependency

```bash
# Ensure you have the latest code
git pull

# Add the dependency
go get github.com/lamboktulussimamora/gra@v1.0.0
go mod tidy
```

### 2. Update Import Statements

The following files need to be updated:

#### cmd/core-api/main.go

Change:
```go
import (
    "github.com/lamboktulussimamora/gra-project/internal/core"
)
```

To:
```go
import (
    "github.com/lamboktulussimamora/gra/core"
)
```

#### internal/interface/handler/example_handler.go

This file has already been updated to use the external framework.

### 3. Update API Function References

The `cmd/core-api/main.go` file uses several functions from the core package that might have different names in the external framework:

| Internal Function | External Function |
| ---------------- | ----------------- |
| `core.NewRouter()` | `core.New()` |
| `core.Chain()` | `core.Chain()` (same) |

Update these function calls as necessary:

```go
// Change this:
router := core.NewRouter()

// To this:
router := core.New()
```

### 4. Test the Changes

After making the updates, run all tests to ensure everything still works:

```bash
go test ./...
```

### 5. Update Documentation

Update any internal documentation that references the internal core package to reflect the use of the external GRA framework.

### 6. Remove the Internal Core Package (Optional)

Once you've verified that everything works with the external framework, you can remove the internal core package:

```bash
# Remove the directory
rm -rf internal/core

# Commit the changes
git add .
git commit -m "Migrated from internal core to external GRA framework"
```

## Manual Migration (Alternative)

If you encounter issues with the automated script, you can follow these manual steps:

1. Update go.mod to include the GRA dependency:
   ```bash
   go get github.com/lamboktulussimamora/gra@v1.0.0
   go mod tidy
   ```

2. Update import statements in:
   - cmd/core-api/main.go
   - Any other files using the internal core package

3. Update function calls if necessary (as mentioned in step 3 above)

4. Test your changes

## Supporting Files

The `migration-samples/` directory contains examples of how the files should look after migration:

- `main.go`: Updated version of the cmd/core-api/main.go file
- `example_handler.go`: Updated version of the internal/interface/handler/example_handler.go file
- `go.mod`: Updated version of the go.mod file with the new dependency

## Troubleshooting

If you encounter issues during migration:

1. Ensure that the GRA framework is properly published to GitHub with the v1.0.0 tag
2. Check that all function calls match the external API
3. Verify that all imports have been updated
4. Look for any type mismatches between the internal and external packages
