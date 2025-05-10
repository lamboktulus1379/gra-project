#!/bin/bash
# full_update.sh
# This script updates all core package references to use the external gra framework

echo "====== Full GRA Framework Migration ======"

# Create backup directory with timestamp
BACKUP_DIR="backup_full_$(date +%Y%m%d%H%M%S)"
mkdir -p "$BACKUP_DIR"
echo "Creating backups in $BACKUP_DIR"

# Backup all Go files
find . -name "*.go" -type f -exec cp {} "$BACKUP_DIR"/{} \;
echo "Backup complete"

# 1. Update imports across all files
echo "1. Updating imports..."

# Find all Go files that import the internal core package
CORE_FILES=$(grep -l "github.com/lamboktulussimamora/gra-project/internal/core" $(find . -name "*.go"))

# For each file with core imports
for file in $CORE_FILES; do
    echo "Processing $file..."
    
    # Check if the file contains core.Context
    if grep -q "core\.Context" "$file"; then
        echo "  - Replacing core.Context with context.Context"
        sed -i '' 's/core\.Context/context\.Context/g' "$file"
        
        # Add context import if not present
        if ! grep -q '"github.com/lamboktulussimamora/gra/context"' "$file"; then
            sed -i '' '/^import (/a\'$'\n\t"github.com/lamboktulussimamora/gra/context"' "$file"
        fi
    fi
    
    # Check if the file contains core.Router
    if grep -q "core\.Router\|core\.New(" "$file"; then
        echo "  - Replacing core.Router with router.Router"
        sed -i '' 's/core\.Router/router\.Router/g' "$file"
        sed -i '' 's/core\.New(/router\.New(/g' "$file"
        
        # Add router import if not present
        if ! grep -q '"github.com/lamboktulussimamora/gra/router"' "$file"; then
            sed -i '' '/^import (/a\'$'\n\t"github.com/lamboktulussimamora/gra/router"' "$file"
        fi
    fi
    
    # Check if the file contains middleware functions
    if grep -q "core\.LoggingMiddleware\|core\.RecoveryMiddleware\|core\.CORSMiddleware\|core\.AuthMiddleware" "$file"; then
        echo "  - Replacing core middleware with middleware package"
        sed -i '' 's/core\.LoggingMiddleware/middleware\.LoggingMiddleware/g' "$file"
        sed -i '' 's/core\.RecoveryMiddleware/middleware\.RecoveryMiddleware/g' "$file"
        sed -i '' 's/core\.CORSMiddleware/middleware\.CORSMiddleware/g' "$file"
        sed -i '' 's/core\.AuthMiddleware/middleware\.AuthMiddleware/g' "$file"
        
        # Add middleware import if not present
        if ! grep -q '"github.com/lamboktulussimamora/gra/middleware"' "$file"; then
            sed -i '' '/^import (/a\'$'\n\t"github.com/lamboktulussimamora/gra/middleware"' "$file"
        fi
    fi
    
    # Check if the file contains validator functions
    if grep -q "core\.Validator\|core\.NewValidator" "$file"; then
        echo "  - Replacing core validator with validator package"
        sed -i '' 's/core\.Validator/validator\.Validator/g' "$file"
        sed -i '' 's/core\.NewValidator/validator\.New/g' "$file"
        
        # Add validator import if not present
        if ! grep -q '"github.com/lamboktulussimamora/gra/validator"' "$file"; then
            sed -i '' '/^import (/a\'$'\n\t"github.com/lamboktulussimamora/gra/validator"' "$file"
        fi
    fi
    
    # Check if the file contains logger functions
    if grep -q "core\.Logger\|core\.NewLogger" "$file"; then
        echo "  - Replacing core logger with logger package"
        sed -i '' 's/core\.Logger/logger\.Logger/g' "$file"
        sed -i '' 's/core\.NewLogger/logger\.New/g' "$file"
        
        # Add logger import if not present
        if ! grep -q '"github.com/lamboktulussimamora/gra/logger"' "$file"; then
            sed -i '' '/^import (/a\'$'\n\t"github.com/lamboktulussimamora/gra/logger"' "$file"
        fi
    fi
    
    # Check if the file contains adapter functions
    if grep -q "core\.Adapter\|core\.NewAdapter" "$file"; then
        echo "  - Replacing core adapter with adapter package"
        sed -i '' 's/core\.Adapter/adapter\.Adapter/g' "$file"
        sed -i '' 's/core\.NewAdapter/adapter\.New/g' "$file"
        
        # Add adapter import if not present
        if ! grep -q '"github.com/lamboktulussimamora/gra/adapter"' "$file"; then
            sed -i '' '/^import (/a\'$'\n\t"github.com/lamboktulussimamora/gra/adapter"' "$file"
        fi
    fi
    
    # Remove the internal core import
    echo "  - Removing internal core import"
    sed -i '' '/github.com\/lamboktulussimamora\/gra-project\/internal\/core/d' "$file"
done

# 2. Update the go.mod file to ensure the external gra dependency is properly included
echo "2. Updating go.mod dependencies..."
go mod tidy

# 3. Verify there are no remaining references to the internal core package
echo "3. Checking for remaining references to internal/core package..."
REMAINING=$(grep -r "github.com/lamboktulussimamora/gra-project/internal/core" --include="*.go" .)
if [ -n "$REMAINING" ]; then
    echo "WARNING: Found remaining references to the internal core package:"
    echo "$REMAINING"
else
    echo "No remaining references found to internal core package."
fi

# 4. Final steps
echo "4. Final verification..."
if go build ./cmd/core-api; then
    echo "Build successful! The migration is complete."
    echo ""
    echo "If all tests pass, you can safely remove the internal/core directory with:"
    echo "rm -rf ./internal/core"
else
    echo "Build failed. There may be some references that need manual fixing."
fi

echo "====== Migration Complete ======"
