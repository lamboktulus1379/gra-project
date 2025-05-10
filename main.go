package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Find the path to the cmd/api directory
	dir, _ := os.Getwd()
	cmdPath := filepath.Join(dir, "cmd", "api")

	// Run the main.go in cmd/api
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = cmdPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
