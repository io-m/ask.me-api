package tag

// Tag represents a canonical tag in the system
type Tag struct {
	Key        string   `json:"_key,omitempty"`
	Label      string   `json:"label"`
	Aliases    []string `json:"aliases,omitempty"`
	UsageCount int      `json:"usageCount,omitempty"`
	CreatedAt  int64    `json:"createdAt"`
}
