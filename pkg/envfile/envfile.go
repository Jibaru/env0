package envfile

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Parser represents the environment file parser
type Parser struct {
	filename string
}

// NewParser creates a new environment file parser
func NewParser(filename string) *Parser {
	return &Parser{
		filename: filename,
	}
}

// Parse reads and parses an environment file, returning a map of key-value pairs
func (p *Parser) Parse() (map[string]interface{}, error) {
	file, err := os.Open(p.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open environment file %s: %v", p.filename, err)
	}
	defer file.Close()

	vars := make(map[string]interface{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle export prefix if present
		if strings.HasPrefix(line, "export") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export"))
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		value = strings.Trim(value, `"'`)

		// Skip if key is empty
		if key == "" {
			continue
		}

		vars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading environment file: %v", err)
	}

	return vars, nil
}

// Writer represents the environment file writer
type Writer struct {
	filename string
}

// NewWriter creates a new environment file writer
func NewWriter(filename string) *Writer {
	return &Writer{
		filename: filename,
	}
}

// Write writes the environment variables to a file
func (w *Writer) Write(vars map[string]interface{}) error {
	file, err := os.Create(w.filename)
	if err != nil {
		return fmt.Errorf("failed to create environment file %s: %v", w.filename, err)
	}
	defer file.Close()

	// Sort keys for consistent output
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Write sorted key-value pairs
	for _, k := range keys {
		if _, err := fmt.Fprintf(file, "%s=%v\n", k, vars[k]); err != nil {
			return fmt.Errorf("failed to write to environment file: %v", err)
		}
	}

	return nil
}

func WriteEnvFile(fileName string, vars map[string]interface{}) error {
	writer := NewWriter(fileName)
	return writer.Write(vars)
}

func ParseEnvFile(filename string) (map[string]interface{}, error) {
	parser := NewParser(filename)
	return parser.Parse()
}
