package user

// AQL queries for user operations
const (
	// GetUserByID retrieves a user by their key
	GetUserByID = `
		FOR u IN users
		FILTER u._key == @key
		RETURN u
	`

	// DeleteFollowEdge removes a follow relationship
	DeleteFollowEdge = `
		FOR e IN follows
		FILTER e._from == @from AND e._to == @to
		REMOVE e IN follows
	`

	// CheckIsFollowing checks if one user follows another
	CheckIsFollowing = `
		FOR e IN follows
		FILTER e._from == @from AND e._to == @to
		RETURN true
	`

	// CheckMutualFollowers checks if two users follow each other
	CheckMutualFollowers = `
		LET follows1 = (
			FOR e IN follows
			FILTER e._from == @user1 AND e._to == @user2
			RETURN true
		)
		LET follows2 = (
			FOR e IN follows
			FILTER e._from == @user2 AND e._to == @user1
			RETURN true
		)
		RETURN LENGTH(follows1) > 0 AND LENGTH(follows2) > 0
	`

	// GetFollowerCount counts how many followers a user has
	GetFollowerCount = `
		RETURN LENGTH(
			FOR e IN follows
			FILTER e._to == @user
			RETURN 1
		)
	`

	// GetFollowingCount counts how many users someone follows
	GetFollowingCount = `
		RETURN LENGTH(
			FOR e IN follows
			FILTER e._from == @user
			RETURN 1
		)
	`
)
