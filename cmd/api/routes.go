package main

import "net/http"

// RegisterRoutes registers all HTTP routes using stdlib mux
// This can be easily swapped for Chi, Gin, or other routers
func (a *App) RegisterRoutes(mux *http.ServeMux) {
	// User routes
	mux.HandleFunc("GET /users/{userId}", a.userHandler.GetUser)
	mux.HandleFunc("POST /users", a.userHandler.CreateUser)
	mux.HandleFunc("POST /users/{userId}/follow", a.userHandler.FollowUser)

	// Post routes
	mux.HandleFunc("GET /posts/{postId}", a.postHandler.GetPost)
	mux.HandleFunc("POST /posts", a.postHandler.CreatePost)
	mux.HandleFunc("POST /posts/poll", a.postHandler.CreatePoll)
	mux.HandleFunc("POST /posts/{postId}/respond", a.postHandler.RespondToPost)
	mux.HandleFunc("POST /posts/{postId}/vote", a.postHandler.Vote)

	// Chat routes
	mux.HandleFunc("GET /users/{userId}/chats", a.chatHandler.GetUserChats)
	mux.HandleFunc("GET /chats/{chatId}", a.chatHandler.GetChat)
	mux.HandleFunc("POST /chats/{chatId}/message", a.chatHandler.SendMessage)
	mux.HandleFunc("POST /chats/{chatId}/accept", a.chatHandler.AcceptChat)
	mux.HandleFunc("POST /chats/{chatId}/mute", a.chatHandler.MuteChat)
	mux.HandleFunc("GET /chats/{chatId}/participants", a.chatHandler.GetParticipants)

	// Feed routes
	mux.HandleFunc("GET /feed", a.feedHandler.GetFeed)

	// Tag routes
	mux.HandleFunc("GET /tags/{tagId}", a.tagHandler.GetTag)
	mux.HandleFunc("GET /tags", a.tagHandler.ListTags)
}
