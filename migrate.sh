#!/bin/bash

# Migration script to update from internal core to external gra framework
# Run this script after publishing the gra framework to GitHub

echo "GRA Framework Migration Script"
echo "============================="

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print colored text
print_color() {
    local color="$1"
    local message="$2"
    
    if command_exists tput; then
        case "$color" in
            "red")    tput setaf 1 ;;
            "green")  tput setaf 2 ;;
            "yellow") tput setaf 3 ;;
            "blue")   tput setaf 4 ;;
        esac
        echo -e "$message"
        tput sgr0
    else
        echo -e "$message"
    fi
}

# Function to print step header
print_step() {
    print_color "blue" "\n=== $1 ==="
}

# Detect if we're running in CI/CD or terminal
CI=${CI:-false}
if [ "$CI" = "true" ]; then
    INTERACTIVE=false
else
    INTERACTIVE=true
fi

# Get timestamp
TIMESTAMP=$(date +%Y%m%d%H%M%S)

# Create log file
LOG_FILE="/tmp/gra_migration_$TIMESTAMP.log"
touch "$LOG_FILE"

# Function to log messages
log() {
    local message="$1"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $message" >> "$LOG_FILE"
    echo "$message"
}

# Function to ask for confirmation
confirm() {
    if [ "$INTERACTIVE" != "true" ]; then
        return 0
    fi
    
    local message="$1"
    read -p "$message [y/N] " response
    case "$response" in
        [yY][eE][sS]|[yY]) 
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Check if git is installed
if ! command_exists git; then
    print_color "red" "Error: git is not installed. Please install git and try again."
    exit 1
fi

# Check if git is clean
if [ -n "$(git status --porcelain)" ]; then
    print_color "yellow" "Warning: You have uncommitted changes."
    confirm "Do you want to continue anyway?" || exit 1
fi

print_step "Checking if gra framework is available"

# Check if gra framework is available
if ! go get github.com/lamboktulussimamora/gra@v1.0.0 2>/dev/null; then
    print_color "red" "Error: Could not find github.com/lamboktulussimamora/gra@v1.0.0"
    print_color "yellow" "You need to publish the framework to GitHub first with a v1.0.0 tag."
    print_color "yellow" "Alternatively, you can create a fork and update the package name accordingly."
    
    read -p "Do you want to continue with a local version for testing? [y/N] " local_response
    case "$local_response" in
        [yY][eE][sS]|[yY]) 
            print_color "yellow" "Continuing with a local version."
            ;;
        *)
            exit 1
            ;;
    esac
else
    print_color "green" "Found gra framework v1.0.0"
fi

print_step "Creating backup"

# Backup current project
BACKUP_DIR="$(pwd)/backup_$TIMESTAMP"
mkdir -p "$BACKUP_DIR"
cp -R ./* "$BACKUP_DIR/"

log "Backup created at $BACKUP_DIR"
print_color "green" "Backup created at $BACKUP_DIR"

print_step "Updating go.mod"

# Update the go.mod file
if grep -q "github.com/lamboktulussimamora/gra" go.mod; then
    log "Dependency already exists in go.mod"
    print_color "green" "Dependency already exists in go.mod"
else
    log "Adding dependency to go.mod"
    print_color "yellow" "Adding dependency to go.mod"
    go get github.com/lamboktulussimamora/gra@v1.0.0
    go mod tidy
    print_color "green" "Dependency added to go.mod"
fi

print_step "Updating imports"

# Update imports in main.go
log "Updating cmd/core-api/main.go"
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS requires an empty string after -i
    sed -i '' 's|"github.com/lamboktulussimamora/gra-project/internal/core"|"github.com/lamboktulussimamora/gra/core"|g' cmd/core-api/main.go
else
    # Linux/Unix doesn't need the empty string
    sed -i 's|"github.com/lamboktulussimamora/gra-project/internal/core"|"github.com/lamboktulussimamora/gra/core"|g' cmd/core-api/main.go
fi

if [ $? -eq 0 ]; then
    print_color "green" "Updated cmd/core-api/main.go"
else
    print_color "red" "Failed to update cmd/core-api/main.go"
fi

# Check for other files that need updating
print_step "Checking for other files using internal core"

OTHER_FILES=$(grep -r "\"github.com/lamboktulussimamora/gra-project/internal/core\"" --include="*.go" . --exclude="./cmd/core-api/main.go" --exclude="./internal/core/*" | cut -d ":" -f 1)

if [ -n "$OTHER_FILES" ]; then
    log "Found other files using internal core:"
    for file in $OTHER_FILES; do
        log "  - $file"
        print_color "yellow" "Updating $file"
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS requires an empty string after -i
            sed -i '' 's|"github.com/lamboktulussimamora/gra-project/internal/core"|"github.com/lamboktulussimamora/gra/core"|g' "$file"
        else
            # Linux/Unix doesn't need the empty string
            sed -i 's|"github.com/lamboktulussimamora/gra-project/internal/core"|"github.com/lamboktulussimamora/gra/core"|g' "$file"
        fi
        
        if [ $? -eq 0 ]; then
            print_color "green" "Updated $file"
        else
            print_color "red" "Failed to update $file"
        fi
    done
else
    log "No other files using internal core directly"
    print_color "green" "No other files using internal core directly"
fi

# Update API function names if needed
print_step "Updating API function names"

log "Updating NewRouter() to New()"
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS requires an empty string after -i
    find . -type f -name "*.go" -not -path "./internal/core/*" | xargs sed -i '' 's|core.NewRouter()|core.New()|g'
else
    # Linux/Unix doesn't need the empty string
    find . -type f -name "*.go" -not -path "./internal/core/*" | xargs sed -i 's|core.NewRouter()|core.New()|g'
fi
print_color "green" "Updated function names"

print_step "Running tests"

# Run tests
go test ./... 2>&1 | tee -a "$LOG_FILE"
if [ ${PIPESTATUS[0]} -eq 0 ]; then
    print_color "green" "All tests passed!"
else
    print_color "red" "Some tests failed. Please review the changed files."
    print_color "yellow" "See $LOG_FILE for details."
fi

print_step "Cleaning up"

# Clean up any existing .bak files
log "Cleaning up backup files"
print_color "yellow" "Cleaning up .bak files..."
BAK_FILES_COUNT=$(find . -name "*.bak" | wc -l | tr -d '[:space:]')

if [ "$BAK_FILES_COUNT" -gt 0 ]; then
    find . -name "*.bak" -delete
    print_color "green" "Removed $BAK_FILES_COUNT backup files"
    log "Removed $BAK_FILES_COUNT backup files"
else
    print_color "green" "No backup files found"
    log "No backup files found"
fi

print_step "Migration Summary"

log "Migration completed!"
print_color "green" "Migration completed!"
print_color "blue" "Next steps:"
print_color "yellow" "1. Review the changes and ensure everything works correctly"
print_color "yellow" "2. If you're satisfied with the changes, you can remove the internal core package:"
print_color "blue" "   rm -rf internal/core"
print_color "yellow" "3. Commit the changes:"
print_color "blue" "   git add ."
print_color "blue" "   git commit -m \"Migrated from internal core to external GRA framework\""

echo ""
print_color "blue" "Log file: $LOG_FILE"
print_color "blue" "Backup directory: $BACKUP_DIR"
echo "============================"
