package autofill

// FieldRef allows referencing a field by column header or position
// in a CSV row.
type FieldRef struct {
	Header   string // e.g.. "account number"
	Position int    // 1-based position in the row
}

// ConditionType enumerates the supported conditional functions
type ConditionType string

const (
	CondEquals ConditionType = "equals"
	// TODO: Support more condition types
	// // Strings:
	// CondContains    ConditionType = "contains"
	// CondRegex       ConditionType = "regex"
	// CondContainsAny ConditionType = "containsAny"
	// CondConainsAll  ConditionType = "containsAll"
	// CondStartsWith  ConditionType = "startsWith"
	// CondEndsWith    ConditionType = "endsWith"
	// // Dates:
	// CondBefore  ConditionType = "before"
	// CondAfter   ConditionType = "after"
	// CondBetween ConditionType = "between"
	// // Numbers:
	// CondGreaterThan ConditionType = "greaterThan"
	// CondLessThan    ConditionType = "lessThan"
)

// Condition represents a single condition or a logical combination
// that must be met.
type Condition struct {
	Field FieldRef      // Field to check condition against
	Type  ConditionType // Type of condition to apply
	Value string        // Value to compare against
	And   []*Condition  // Additional conditions to combine with AND
	Or    []*Condition  // Additional conditions to combine with OR
	Not   *Condition    // Negate this condition
}

type OverwriteBehavior string

const (
	OverwriteAlways OverwriteBehavior = "always" // always overwrite existing data
	OverwriteNever  OverwriteBehavior = "never"  // never overwrite existing data
	OverwritePrompt OverwriteBehavior = "prompt" // prompt user before overwriting existing data
)

// Rule defines a single autofill rule.
type Rule struct {
	Name      string            // Name of the rule
	Format    string            // A format name identifying a bankcsv.Format
	Target    FieldRef          // The field to fill in the CSV row
	Value     any               // The value to fill in the target field
	Overwrite OverwriteBehavior // Behavior when target field already has a value
	Condition Condition         // Condition that must be met for this rule to apply
}
