package tag

// AQL queries for tag operations
const (
	// GetTagByID retrieves a tag by its key
	GetTagByID = `
		FOR t IN tags
		FILTER t._key == @key
		RETURN t
	`

	// GetTagByAlias finds a tag by alias or label match
	GetTagByAlias = `
		FOR t IN tags
		FILTER @alias IN t.aliases OR LOWER(t.label) == LOWER(@alias)
		RETURN t
	`

	// IncrementTagUsageCount increments the usage count for a tag
	IncrementTagUsageCount = `
		FOR t IN tags
		FILTER t._key == @key
		UPDATE t WITH { usageCount: t.usageCount + 1 } IN tags
	`

	// ListTagsByUsage lists tags ordered by usage count
	ListTagsByUsage = `
		FOR t IN tags
		SORT t.usageCount DESC
		LIMIT @offset, @limit
		RETURN t
	`

	// SearchTags searches tags by label or alias
	SearchTags = `
		FOR t IN tags
		FILTER CONTAINS(LOWER(t.label), LOWER(@query))
		   OR LENGTH(
		       FOR alias IN t.aliases
		       FILTER CONTAINS(LOWER(alias), LOWER(@query))
		       RETURN alias
		   ) > 0
		SORT t.usageCount DESC
		LIMIT @limit
		RETURN t
	`
)
