#!/bin/zsh

# Script to clean up .bak files in the project

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}====== Backup File Cleanup ======${NC}"
echo -e "${YELLOW}Searching for .bak files...${NC}"

BAK_FILES=$(find . -name "*.bak")
BAK_COUNT=$(echo "$BAK_FILES" | grep -c "")

if [ "$BAK_COUNT" -eq 0 ]; then
  echo -e "${GREEN}No .bak files found.${NC}"
  exit 0
fi

echo -e "${YELLOW}Found $BAK_COUNT .bak files:${NC}"
echo "$BAK_FILES" | sed 's/^/  /'

read -p "Do you want to remove these files? (y/N) " response
case "$response" in
  [yY][eE][sS]|[yY])
    find . -name "*.bak" -delete
    echo -e "${GREEN}Removed $BAK_COUNT .bak files.${NC}"
    ;;
  *)
    echo -e "${YELLOW}No files were removed.${NC}"
    ;;
esac

echo -e "${BLUE}=============================${NC}"
