package domain

// PostType defines the structural format of the post (frontend-determined)
type PostType string

const (
	PostTypeText PostType = "text"
	PostTypePoll PostType = "poll"
)

// PostCategory defines what the post is about (backend enum)
type PostCategory string

const (
	CategoryCareer        PostCategory = "career"
	CategoryRelationships PostCategory = "relationships"
	CategoryTech          PostCategory = "tech"
	CategoryHealth        PostCategory = "health"
	CategoryFinance       PostCategory = "finance"
	CategoryFun           PostCategory = "fun"
	CategoryOpinion       PostCategory = "opinion"
	CategoryLifestyle     PostCategory = "lifestyle"
	CategoryEducation     PostCategory = "education"
	CategoryOther         PostCategory = "other"
)

// ValidCategory checks if a category is valid
func ValidCategory(c string) bool {
	switch PostCategory(c) {
	case CategoryCareer, CategoryRelationships, CategoryTech, CategoryHealth,
		CategoryFinance, CategoryFun, CategoryOpinion, CategoryLifestyle,
		CategoryEducation, CategoryOther:
		return true
	}
	return false
}

// NormalizeCategory returns a valid category or "other" as fallback
func NormalizeCategory(c string) PostCategory {
	if ValidCategory(c) {
		return PostCategory(c)
	}
	return CategoryOther
}

// PostDepth defines content seriousness level (backend enum)
type PostDepth string

const (
	DepthCasual  PostDepth = "casual"
	DepthNeutral PostDepth = "neutral"
	DepthSerious PostDepth = "serious"
)

// NormalizeDepth returns a valid depth or "neutral" as fallback
func NormalizeDepth(d string) PostDepth {
	switch PostDepth(d) {
	case DepthCasual, DepthNeutral, DepthSerious:
		return PostDepth(d)
	}
	return DepthNeutral
}

// MessageStatus defines the delivery status of a message
type MessageStatus string

const (
	MessageStatusSending   MessageStatus = "sending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusSeen      MessageStatus = "seen"
	MessageStatusFailed    MessageStatus = "failed"
)

// ChatType defines the type of chat
type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup  ChatType = "group"
)

// ParticipantRole defines the role of a user in a chat
type ParticipantRole string

const (
	RoleAuthor    ParticipantRole = "author"
	RoleResponder ParticipantRole = "responder"
	RoleInvited   ParticipantRole = "invited"
)

// ParticipantStatus defines the participation status
type ParticipantStatus string

const (
	StatusActive  ParticipantStatus = "active"
	StatusPending ParticipantStatus = "pending"
	StatusMuted   ParticipantStatus = "muted"
)

// AIRawData represents raw AI classification output
type AIRawData struct {
	Category   string   `json:"category"`
	Intent     string   `json:"intent"`
	Depth      string   `json:"depth,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Confidence float64  `json:"confidence,omitempty"`
	Risk       string   `json:"risk,omitempty"`
	Flags      []string `json:"flags,omitempty"`
}
