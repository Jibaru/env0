package envdiff

// ChangeType represents the type of change in an environment variable
type ChangeType string

const (
	Added    ChangeType = "ADDED"
	Modified ChangeType = "MODIFIED"
	Deleted  ChangeType = "DELETED"
)

// Change represents a single variable change
type Change struct {
	Name     string
	Type     ChangeType
	OldValue interface{}
	NewValue interface{}
}

// DiffResult contains all changes between two environment states
type DiffResult struct {
	SafeToMerge bool
	Changes     []Change
}

// ConflictReport represents conflicts found during merge
type ConflictReport struct {
	Environment string
	Conflicts   []Change
	Timestamp   string
}
