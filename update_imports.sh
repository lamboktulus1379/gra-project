#!/bin/bash
# update_imports.sh
# This script updates all imports from the internal core package to the external gra packages

# Set the backup directory with timestamp
BACKUP_DIR="backup_imports_$(date +%Y%m%d%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "====== Import Path Migration ======"
echo "Creating backups in $BACKUP_DIR"

# First, back up all Go files
find . -name "*.go" -type f -exec cp {} "$BACKUP_DIR"/{} \;

# Update imports in all Go files
echo "Updating import paths..."

# Detect OS for sed compatibility
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS requires an empty string after -i
  SED_CMD="sed -i ''"
else
  # Linux/Unix doesn't need the empty string
  SED_CMD="sed -i"
fi

# 1. Replace core.Context with context.Context
find . -name "*.go" -type f -exec $SED_CMD 's|"github.com/lamboktulussimamora/gra/core"|"github.com/lamboktulussimamora/gra/context"|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|\*core\.Context|\*context\.Context|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|core\.NewValidator|context\.NewValidator|g' {} \;

# 2. Update router related imports
find . -name "*.go" -type f -exec $SED_CMD 's|core\.New()|router\.New()|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|\*core\.Router|\*router\.Router|g' {} \;

# 3. Add router import where needed
grep -l "router\.Router\|router\.New" $(find . -name "*.go" -type f) | xargs -I{} $SED_CMD '/import (/a\
\t"github.com/lamboktulussimamora/gra/router"' {}

# 4. Update middleware related imports
find . -name "*.go" -type f -exec $SED_CMD 's|core\.LoggingMiddleware|middleware\.LoggingMiddleware|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|core\.RecoveryMiddleware|middleware\.RecoveryMiddleware|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|core\.CORSMiddleware|middleware\.CORSMiddleware|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|core\.AuthMiddleware|middleware\.AuthMiddleware|g' {} \;

# 5. Add middleware import where needed
grep -l "middleware\.LoggingMiddleware\|middleware\.RecoveryMiddleware\|middleware\.CORSMiddleware\|middleware\.AuthMiddleware" $(find . -name "*.go" -type f) | xargs -I{} $SED_CMD '/import (/a\
\t"github.com/lamboktulussimamora/gra/middleware"' {}

# 6. Update validator related imports
find . -name "*.go" -type f -exec $SED_CMD 's|core\.Validator|validator\.Validator|g' {} \;

# 7. Add validator import where needed
grep -l "validator\.Validator" $(find . -name "*.go" -type f) | xargs -I{} $SED_CMD '/import (/a\
\t"github.com/lamboktulussimamora/gra/validator"' {}

# 8. Update logger related imports
find . -name "*.go" -type f -exec $SED_CMD 's|core\.Logger|logger\.Logger|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|core\.NewLogger|logger\.NewLogger|g' {} \;

# 9. Add logger import where needed
grep -l "logger\.Logger\|logger\.NewLogger" $(find . -name "*.go" -type f) | xargs -I{} $SED_CMD '/import (/a\
\t"github.com/lamboktulussimamora/gra/logger"' {}

# 10. Update adapter related imports
find . -name "*.go" -type f -exec $SED_CMD 's|core\.Adapter|adapter\.Adapter|g' {} \;
find . -name "*.go" -type f -exec $SED_CMD 's|core\.NewAdapter|adapter\.NewAdapter|g' {} \;

# 11. Add adapter import where needed
grep -l "adapter\.Adapter\|adapter\.NewAdapter" $(find . -name "*.go" -type f) | xargs -I{} $SED_CMD '/import (/a\
\t"github.com/lamboktulussimamora/gra/adapter"' {}

# Ensure proper Go module deps
echo "Updating Go module dependencies..."
go mod tidy

echo "Import path migration complete!"
echo "You can verify the changes and run 'go build ./...' to check for compilation errors."
echo "If everything works correctly, you can safely remove the internal/core package."
echo "========================================="
