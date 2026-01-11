package user

// User represents a user document in ArangoDB
type User struct {
	Key           string       `json:"_key,omitempty"`
	Username      string       `json:"username"`
	CreatedAt     int64        `json:"createdAt"`
	Interests     []string     `json:"interests,omitempty"`
	BlockedTopics []string     `json:"blockedTopics,omitempty"`
	Settings      UserSettings `json:"settings,omitempty"`
	Stats         UserStats    `json:"stats,omitempty"`
}

type UserSettings struct {
	AllowDMs     bool `json:"allowDMs"`
	AllowTagging bool `json:"allowTagging"`
}

type UserStats struct {
	PostsCreated   int `json:"postsCreated"`
	ResponsesGiven int `json:"responsesGiven"`
}

// FollowsEdge represents a follows relationship between users
type FollowsEdge struct {
	From      string `json:"_from"`
	To        string `json:"_to"`
	CreatedAt int64  `json:"createdAt"`
}

// CreateUserRequest is the request payload for creating a user
type CreateUserRequest struct {
	Username  string       `json:"username"`
	Interests []string     `json:"interests,omitempty"`
	Settings  UserSettings `json:"settings,omitempty"`
}

// CreateUserResponse is the response payload for creating a user
type CreateUserResponse struct {
	Key       string `json:"_key"`
	CreatedAt int64  `json:"createdAt"`
}

// FollowUserRequest is the request payload for following a user
type FollowUserRequest struct {
	FollowUserID string `json:"followUserId"`
}

// FollowUserResponse is the response payload for following a user
type FollowUserResponse struct {
	Success  bool   `json:"success"`
	FollowID string `json:"followId"`
}
