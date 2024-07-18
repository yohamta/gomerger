package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileContent represents the content of a single Go file.
type FileContent struct {
	PackageName string   // The package name of the file
	Imports     []string // List of import statements
	Code        []string // The actual code content
}

// files stores the content of all processed Go files.
var files []FileContent

func main() {
	// Check if a project directory is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run merger.go <project_directory>")
		os.Exit(1)
	}

	projectDir := os.Args[1]
	outputFile := "merged_project.go"

	// Walk through the project directory and process each file
	err := filepath.Walk(projectDir, processFile)
	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
		os.Exit(1)
	}

	// Write the merged content to the output file
	err = writeOutputFile(outputFile)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Merged Go files written to %s\n", outputFile)
}

// processFile is called for each file and directory in the project.
// It reads and parses Go files, skipping test files and directories.
func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// Process only non-directory .go files that are not test files
	if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		fileContent := parseGoFile(string(content))
		if fileContent.PackageName != "" {
			files = append(files, fileContent)
		}
	}

	return nil
}

// parseGoFile extracts package name, imports, and code from a Go file's content.
func parseGoFile(content string) FileContent {
	var fileContent FileContent
	scanner := bufio.NewScanner(strings.NewReader(content))
	inImportBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		// Extract package name
		if strings.HasPrefix(line, "package ") {
			fileContent.PackageName = strings.TrimSpace(strings.TrimPrefix(line, "package "))
			continue
		}

		// Detect start of import block
		if strings.HasPrefix(line, "import (") {
			inImportBlock = true
			continue
		}

		// Process imports
		if inImportBlock {
			if line == ")" {
				inImportBlock = false
			} else {
				fileContent.Imports = append(fileContent.Imports, strings.TrimSpace(line))
			}
		} else {
			// Add non-import lines to code
			fileContent.Code = append(fileContent.Code, line)
		}
	}

	return fileContent
}

// writeOutputFile writes the merged content to the output file.
func writeOutputFile(outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write package declaration
	writer.WriteString("package main\n\n")

	// Write imports
	writer.WriteString("import (\n")
	imports := make(map[string]bool)
	for _, f := range files {
		for _, imp := range f.Imports {
			if !imports[imp] {
				imports[imp] = true
				writer.WriteString("\t" + imp + "\n")
			}
		}
	}
	writer.WriteString(")\n\n")

	// Write code from all files
	for i, f := range files {
		if i > 0 {
			writer.WriteString("\n")
		}
		writer.WriteString(fmt.Sprintf("// Original package: %s\n", f.PackageName))
		for _, line := range f.Code {
			writer.WriteString(line + "\n")
		}
	}

	// Add a main function if it doesn't exist
	if !mainFunctionExists() {
		writer.WriteString("\nfunc main() {\n\t// TODO: Add main logic here\n}\n")
	}

	return writer.Flush()
}

// mainFunctionExists checks if a main function is present in any of the processed files.
func mainFunctionExists() bool {
	for _, f := range files {
		for _, line := range f.Code {
			if strings.HasPrefix(strings.TrimSpace(line), "func main()") {
				return true
			}
		}
	}
	return false
}
