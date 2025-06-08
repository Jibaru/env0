package envdiff

import (
	"fmt"
	"strings"
)

// DeletedValue is a special type to represent deleted values
type DeletedValue struct{}

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

// FormatGitStyleConflict formats a variable conflict with git-style markers
func FormatGitStyleConflict(key string, localValue, remoteValue interface{}) string {
	var sb strings.Builder
	sb.WriteString("<<<<<<< LOCAL\n")
	if _, isDeleted := localValue.(DeletedValue); !isDeleted {
		sb.WriteString(fmt.Sprintf("%s=%v\n", key, localValue))
	}
	sb.WriteString("=======\n")
	if _, isDeleted := remoteValue.(DeletedValue); !isDeleted {
		sb.WriteString(fmt.Sprintf("%s=%v\n", key, remoteValue))
	}
	sb.WriteString(">>>>>>> REMOTE\n")
	return sb.String()
}
