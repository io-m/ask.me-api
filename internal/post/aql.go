package post

// AQL queries for post operations
const (
	// GetPostByID retrieves a post by its key
	GetPostByID = `
		FOR p IN posts
		FILTER p._key == @key
		RETURN p
	`

	// GetPostTags retrieves all tags for a post
	GetPostTags = `
		FOR edge IN post_has_tag
		FILTER edge._from == @postId
		FOR tag IN tags
		FILTER tag._id == edge._to
		RETURN tag._key
	`

	// GetPollVotes aggregates votes for a poll
	GetPollVotes = `
		FOR edge IN voted
		FILTER edge._to == @postId
		COLLECT option = edge.option WITH COUNT INTO count
		RETURN { option, count }
	`

	// CheckUserVoted checks if a user has voted on a poll
	CheckUserVoted = `
		FOR edge IN voted
		FILTER edge._from == @userId AND edge._to == @postId
		RETURN true
	`

	// CheckUserResponded checks if a user has responded to a post
	CheckUserResponded = `
		FOR edge IN responded
		FILTER edge._from == @userId AND edge._to == @postId
		RETURN true
	`

	// GetPostAuthor retrieves the author of a post
	GetPostAuthor = `
		FOR edge IN created
		FILTER edge._to == @postId
		FOR user IN users
		FILTER user._id == edge._from
		RETURN {
			id: user._key,
			username: user.username,
			avatarUrl: user.avatarUrl
		}
	`
)
