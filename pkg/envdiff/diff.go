package envdiff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CompareMaps compares two environment maps and returns a DiffResult
func CompareMaps(original, new map[string]interface{}) DiffResult {
	var changes []Change
	safeToMerge := true

	// Check for modifications and additions
	for key, newValue := range new {
		if oldValue, exists := original[key]; exists {
			if oldValue != newValue {
				changes = append(changes, Change{
					Name:     key,
					Type:     Modified,
					OldValue: oldValue,
					NewValue: newValue,
				})
				safeToMerge = false
			}
		} else {
			changes = append(changes, Change{
				Name:     key,
				Type:     Added,
				NewValue: newValue,
			})
		}
	}

	// Check for deletions
	for key, oldValue := range original {
		if _, exists := new[key]; !exists {
			changes = append(changes, Change{
				Name:     key,
				Type:     Deleted,
				OldValue: oldValue,
			})
			safeToMerge = false
		}
	}

	return DiffResult{
		SafeToMerge: safeToMerge,
		Changes:     changes,
	}
}

// CreateBackup creates a backup of the environment file
func CreateBackup(envFile string) (string, error) {
	content, err := os.ReadFile(envFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No backup needed if file doesn't exist
		}
		return "", fmt.Errorf("failed to read env file: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")
	backupDir := filepath.Join(".env0", "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %v", err)
	}

	backupFile := filepath.Join(backupDir, fmt.Sprintf("%s.%s.bak", filepath.Base(envFile), timestamp))
	if err := os.WriteFile(backupFile, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file: %v", err)
	}

	return backupFile, nil
}

// SaveConflictReport saves a conflict report to a file
func SaveConflictReport(report ConflictReport) error {
	reportsDir := filepath.Join(".env0", "conflicts")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create conflicts directory: %v", err)
	}

	filename := filepath.Join(reportsDir, fmt.Sprintf("conflicts_%s_%s.json",
		report.Environment,
		report.Timestamp.Format("20060102150405")))

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conflict report: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write conflict report: %v", err)
	}

	return nil
}

// MergeMaps merges two environment maps based on the diff result
func MergeMaps(original, new map[string]interface{}, diff DiffResult) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy original map
	for k, v := range original {
		result[k] = v
	}

	// Apply only safe changes
	if diff.SafeToMerge {
		for _, change := range diff.Changes {
			if change.Type == Added {
				result[change.Name] = change.NewValue
			}
		}
	}

	return result
}
