package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	endpoint := os.Getenv("ARANGO_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8529"
	}
	database := os.Getenv("ARANGO_DATABASE")
	if database == "" {
		database = "askme"
	}
	username := os.Getenv("ARANGO_USERNAME")
	if username == "" {
		username = "root"
	}
	password := os.Getenv("ARANGO_PASSWORD")
	if password == "" {
		password = "rootpassword"
	}

	baseURL := fmt.Sprintf("%s/_db/%s/_api", endpoint, database)

	seeder := &Seeder{
		baseURL:  baseURL,
		username: username,
		password: password,
		client:   &http.Client{Timeout: 10 * time.Second},
	}

	log.Println("Seeding database...")

	if err := seeder.SeedAll(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding complete!")
}

type Seeder struct {
	baseURL  string
	username string
	password string
	client   *http.Client
}

func (s *Seeder) insert(collection string, doc any) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Use overwriteMode=replace to upsert (insert or replace if exists)
	url := fmt.Sprintf("%s/document/%s?overwriteMode=replace", s.baseURL, collection)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.SetBasicAuth(s.username, s.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("insert failed: %s", resp.Status)
	}

	return nil
}

func (s *Seeder) SeedAll() error {
	now := time.Now().UnixMilli()
	day := int64(24 * 60 * 60 * 1000)

	// Seed Users
	users := []map[string]any{
		{
			"_key":          "u1",
			"username":      "alex_dev",
			"avatarUrl":     "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop&crop=face",
			"createdAt":     now - 30*day,
			"interests":     []string{"tech", "career"},
			"blockedTopics": []string{},
			"settings":      map[string]any{"allowDMs": true, "allowTagging": true},
			"stats":         map[string]any{"postsCreated": 5, "responsesGiven": 12},
		},
		{
			"_key":          "u2",
			"username":      "maria_chen",
			"avatarUrl":     "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400&h=400&fit=crop&crop=face",
			"createdAt":     now - 25*day,
			"interests":     []string{"health", "lifestyle"},
			"blockedTopics": []string{"politics"},
			"settings":      map[string]any{"allowDMs": true, "allowTagging": true},
			"stats":         map[string]any{"postsCreated": 3, "responsesGiven": 8},
		},
		{
			"_key":          "u3",
			"username":      "john_doe",
			"avatarUrl":     "https://images.unsplash.com/photo-1633332755192-727a05c4013d?w=400&h=400&fit=crop&crop=face",
			"createdAt":     now - 20*day,
			"interests":     []string{"finance", "career"},
			"blockedTopics": []string{},
			"settings":      map[string]any{"allowDMs": true, "allowTagging": false},
			"stats":         map[string]any{"postsCreated": 2, "responsesGiven": 15},
		},
		{
			"_key":          "u4",
			"username":      "sarah_k",
			"avatarUrl":     "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=400&h=400&fit=crop&crop=face",
			"createdAt":     now - 15*day,
			"interests":     []string{"relationships", "fun"},
			"blockedTopics": []string{},
			"settings":      map[string]any{"allowDMs": true, "allowTagging": true},
			"stats":         map[string]any{"postsCreated": 4, "responsesGiven": 6},
		},
		{
			"_key":          "u5",
			"username":      "dev_master",
			"avatarUrl":     "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=400&h=400&fit=crop&crop=face",
			"createdAt":     now - 10*day,
			"interests":     []string{"tech", "education"},
			"blockedTopics": []string{},
			"settings":      map[string]any{"allowDMs": false, "allowTagging": true},
			"stats":         map[string]any{"postsCreated": 7, "responsesGiven": 20},
		},
	}

	log.Println("Seeding users...")
	for _, u := range users {
		if err := s.insert("users", u); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Tags
	tags := []map[string]any{
		{"_key": "career-change", "label": "Career Change", "aliases": []string{"changing careers", "career switch"}, "usageCount": 15, "createdAt": now - 30*day},
		{"_key": "backend-dev", "label": "Backend Development", "aliases": []string{"backend", "server-side"}, "usageCount": 12, "createdAt": now - 30*day},
		{"_key": "frontend", "label": "Frontend", "aliases": []string{"frontend dev", "ui development"}, "usageCount": 18, "createdAt": now - 30*day},
		{"_key": "remote-work", "label": "Remote Work", "aliases": []string{"wfh", "work from home"}, "usageCount": 25, "createdAt": now - 28*day},
		{"_key": "productivity", "label": "Productivity", "aliases": []string{"being productive", "efficiency"}, "usageCount": 20, "createdAt": now - 28*day},
		{"_key": "relationships", "label": "Relationships", "aliases": []string{"dating", "love"}, "usageCount": 10, "createdAt": now - 25*day},
		{"_key": "health-tips", "label": "Health Tips", "aliases": []string{"wellness", "healthy living"}, "usageCount": 8, "createdAt": now - 25*day},
		{"_key": "investing", "label": "Investing", "aliases": []string{"investment", "stocks"}, "usageCount": 14, "createdAt": now - 20*day},
		{"_key": "frameworks", "label": "Frameworks", "aliases": []string{"web frameworks", "libraries"}, "usageCount": 16, "createdAt": now - 20*day},
		{"_key": "learning", "label": "Learning", "aliases": []string{"education", "self-improvement"}, "usageCount": 22, "createdAt": now - 15*day},
	}

	log.Println("Seeding tags...")
	for _, t := range tags {
		if err := s.insert("tags", t); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Posts
	posts := []map[string]any{
		{
			"_key":        "p1",
			"authorId":    "users/u1",
			"postType":    "text",
			"text":        "How do I switch from frontend to backend engineering?",
			"category":    "career",
			"intent":      "seeking-advice",
			"depth":       "serious",
			"pollOptions": nil,
			"aiRaw":       map[string]any{"category": "career", "intent": "seeking-advice", "depth": "serious", "tags": []string{"career change", "backend dev"}, "confidence": 0.92, "risk": "low", "flags": []string{}},
			"createdAt":   now - 7*day,
		},
		{
			"_key":        "p2",
			"authorId":    "users/u2",
			"postType":    "text",
			"text":        "What's your morning routine for staying productive while working from home?",
			"category":    "lifestyle",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"aiRaw":       map[string]any{"category": "lifestyle", "intent": "asking-question", "depth": "casual", "tags": []string{"remote work", "productivity"}, "confidence": 0.88, "risk": "low", "flags": []string{}},
			"createdAt":   now - 5*day,
		},
		{
			"_key":        "p3",
			"authorId":    "users/u5",
			"postType":    "poll",
			"text":        "Which frontend framework do you prefer?",
			"category":    "tech",
			"intent":      "seeking-opinion",
			"depth":       "casual",
			"pollOptions": []string{"React", "Vue", "Angular", "Svelte"},
			"aiRaw":       map[string]any{"category": "tech", "intent": "seeking-opinion", "depth": "casual", "tags": []string{"frontend", "frameworks"}, "confidence": 0.95, "risk": "low", "flags": []string{}},
			"createdAt":   now - 4*day,
		},
		{
			"_key":        "p4",
			"authorId":    "users/u3",
			"postType":    "text",
			"text":        "I just landed my first remote job after 6 months of searching! Here's what worked for me.",
			"category":    "career",
			"intent":      "sharing-story",
			"depth":       "neutral",
			"pollOptions": nil,
			"aiRaw":       map[string]any{"category": "career", "intent": "sharing-story", "depth": "neutral", "tags": []string{"remote work", "career change"}, "confidence": 0.85, "risk": "low", "flags": []string{}},
			"createdAt":   now - 3*day,
		},
		{
			"_key":        "p5",
			"authorId":    "users/u4",
			"postType":    "text",
			"text":        "How do you handle disagreements with your partner about finances?",
			"category":    "relationships",
			"intent":      "seeking-advice",
			"depth":       "serious",
			"pollOptions": nil,
			"aiRaw":       map[string]any{"category": "relationships", "intent": "seeking-advice", "depth": "serious", "tags": []string{"relationships", "investing"}, "confidence": 0.82, "risk": "low", "flags": []string{}},
			"createdAt":   now - 2*day,
		},
		{
			"_key":        "p6",
			"authorId":    "users/u1",
			"postType":    "text",
			"text":        "Best resources for learning Go in 2026?",
			"category":    "education",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"aiRaw":       map[string]any{"category": "education", "intent": "asking-question", "depth": "casual", "tags": []string{"backend dev", "learning"}, "confidence": 0.90, "risk": "low", "flags": []string{}},
			"createdAt":   now - 1*day,
		},
		{
			"_key":        "p7",
			"authorId":    "users/u2",
			"postType":    "poll",
			"text":        "How many hours of sleep do you get on average?",
			"category":    "health",
			"intent":      "seeking-opinion",
			"depth":       "casual",
			"pollOptions": []string{"Less than 6", "6-7 hours", "7-8 hours", "More than 8"},
			"aiRaw":       map[string]any{"category": "health", "intent": "seeking-opinion", "depth": "casual", "tags": []string{"health tips", "productivity"}, "confidence": 0.91, "risk": "low", "flags": []string{}},
			"createdAt":   now - 12*60*60*1000,
		},
		// Group chat post (tagged multiple users)
		{
			"_key":        "p8",
			"authorId":    "users/u1",
			"postType":    "text",
			"text":        "Tech folks! @maria_chen @john_doe @sarah_k - What stack are you using for your side projects in 2026?",
			"category":    "tech",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"aiRaw":       map[string]any{"category": "tech", "intent": "asking-question", "depth": "casual", "tags": []string{"tech", "side projects"}, "confidence": 0.89, "risk": "low", "flags": []string{}},
			"createdAt":   now - 2*day,
		},
	}

	log.Println("Seeding posts...")
	for _, p := range posts {
		if err := s.insert("posts", p); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Chats
	chats := []map[string]any{
		{"_key": "c1", "postId": "posts/p1", "type": "direct", "createdAt": now - 6*day, "participantCount": 2},
		{"_key": "c2", "postId": "posts/p2", "type": "direct", "createdAt": now - 4*day, "participantCount": 2},
		{"_key": "c3", "postId": "posts/p4", "type": "direct", "createdAt": now - 2*day, "participantCount": 2},
		// Group chat
		{"_key": "c4", "postId": "posts/p8", "type": "group", "createdAt": now - 2*day, "participantCount": 4},
	}

	log.Println("Seeding chats...")
	for _, c := range chats {
		if err := s.insert("chats", c); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Messages
	messages := []map[string]any{
		{"_key": "m1", "chatId": "chats/c1", "senderId": "users/u3", "text": "I made the switch last year! Start with Go or Node.js for backend.", "status": "seen", "createdAt": now - 6*day + 60*60*1000},
		{"_key": "m2", "chatId": "chats/c1", "senderId": "users/u1", "text": "Thanks! Did you take any courses or just learn by building?", "status": "seen", "createdAt": now - 6*day + 2*60*60*1000},
		{"_key": "m3", "chatId": "chats/c1", "senderId": "users/u3", "text": "Mostly building projects. I'd recommend starting with a REST API.", "status": "delivered", "createdAt": now - 5*day},
		{"_key": "m4", "chatId": "chats/c2", "senderId": "users/u5", "text": "I wake up at 6am, exercise, then start work at 8. No meetings before 10!", "status": "seen", "createdAt": now - 4*day + 30*60*1000},
		{"_key": "m5", "chatId": "chats/c2", "senderId": "users/u2", "text": "That's impressive! I struggle with the exercise part.", "status": "delivered", "createdAt": now - 3*day},
		{"_key": "m6", "chatId": "chats/c3", "senderId": "users/u1", "text": "Congrats! What was the key factor in landing the job?", "status": "sent", "createdAt": now - 1*day},
		// Group chat messages
		{"_key": "m7", "chatId": "chats/c4", "senderId": "users/u2", "text": "I've been using Next.js + Prisma + PostgreSQL. Really enjoying it!", "status": "seen", "createdAt": now - 2*day + 60*60*1000},
		{"_key": "m8", "chatId": "chats/c4", "senderId": "users/u3", "text": "Bun + Hono for backend, React for frontend. So fast! 🚀", "status": "seen", "createdAt": now - 2*day + 2*60*60*1000},
		{"_key": "m9", "chatId": "chats/c4", "senderId": "users/u4", "text": "I'm still on the MERN stack but considering switching to T3.", "status": "delivered", "createdAt": now - 1*day},
		{"_key": "m10", "chatId": "chats/c4", "senderId": "users/u1", "text": "T3 is great! TypeScript everywhere makes life easier.", "status": "sent", "createdAt": now - 6*60*60*1000},
	}

	log.Println("Seeding messages...")
	for _, m := range messages {
		if err := s.insert("messages", m); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Edge: created (users -> posts)
	created := []map[string]any{
		{"_from": "users/u1", "_to": "posts/p1", "createdAt": now - 7*day},
		{"_from": "users/u2", "_to": "posts/p2", "createdAt": now - 5*day},
		{"_from": "users/u5", "_to": "posts/p3", "createdAt": now - 4*day},
		{"_from": "users/u3", "_to": "posts/p4", "createdAt": now - 3*day},
		{"_from": "users/u4", "_to": "posts/p5", "createdAt": now - 2*day},
		{"_from": "users/u1", "_to": "posts/p6", "createdAt": now - 1*day},
		{"_from": "users/u2", "_to": "posts/p7", "createdAt": now - 12*60*60*1000},
		{"_from": "users/u1", "_to": "posts/p8", "createdAt": now - 2*day}, // Group chat post
	}

	log.Println("Seeding created edges...")
	for _, e := range created {
		if err := s.insert("created", e); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Edge: follows (users -> users)
	follows := []map[string]any{
		{"_from": "users/u1", "_to": "users/u2", "createdAt": now - 20*day},
		{"_from": "users/u1", "_to": "users/u3", "createdAt": now - 18*day},
		{"_from": "users/u2", "_to": "users/u1", "createdAt": now - 19*day},
		{"_from": "users/u2", "_to": "users/u5", "createdAt": now - 15*day},
		{"_from": "users/u3", "_to": "users/u1", "createdAt": now - 17*day},
		{"_from": "users/u3", "_to": "users/u4", "createdAt": now - 14*day},
		{"_from": "users/u4", "_to": "users/u3", "createdAt": now - 13*day},
		{"_from": "users/u5", "_to": "users/u1", "createdAt": now - 10*day},
		{"_from": "users/u5", "_to": "users/u2", "createdAt": now - 10*day},
	}

	log.Println("Seeding follows edges...")
	for _, e := range follows {
		if err := s.insert("follows", e); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Edge: responded (users -> posts)
	responded := []map[string]any{
		{"_from": "users/u3", "_to": "posts/p1", "chatId": "chats/c1", "createdAt": now - 6*day},
		{"_from": "users/u5", "_to": "posts/p2", "chatId": "chats/c2", "createdAt": now - 4*day},
		{"_from": "users/u1", "_to": "posts/p4", "chatId": "chats/c3", "createdAt": now - 2*day},
	}

	log.Println("Seeding responded edges...")
	for _, e := range responded {
		if err := s.insert("responded", e); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Edge: participates_in (users -> chats)
	participatesIn := []map[string]any{
		{"_from": "users/u1", "_to": "chats/c1", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 7*day},
		{"_from": "users/u3", "_to": "chats/c1", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 6*day},
		{"_from": "users/u2", "_to": "chats/c2", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 5*day},
		{"_from": "users/u5", "_to": "chats/c2", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 4*day},
		{"_from": "users/u3", "_to": "chats/c3", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 3*day},
		{"_from": "users/u1", "_to": "chats/c3", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},
		// Group chat participants (c4)
		{"_from": "users/u1", "_to": "chats/c4", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},
		{"_from": "users/u2", "_to": "chats/c4", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},
		{"_from": "users/u3", "_to": "chats/c4", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},
		{"_from": "users/u4", "_to": "chats/c4", "role": "invited", "status": "pending", "notificationsEnabled": false, "joinedAt": nil}, // Not yet accepted
	}

	log.Println("Seeding participates_in edges...")
	for _, e := range participatesIn {
		if err := s.insert("participates_in", e); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Edge: post_has_tag (posts -> tags)
	postHasTag := []map[string]any{
		{"_from": "posts/p1", "_to": "tags/career-change", "confidence": 0.92, "source": "ai"},
		{"_from": "posts/p1", "_to": "tags/backend-dev", "confidence": 0.88, "source": "ai"},
		{"_from": "posts/p2", "_to": "tags/remote-work", "confidence": 0.90, "source": "ai"},
		{"_from": "posts/p2", "_to": "tags/productivity", "confidence": 0.85, "source": "ai"},
		{"_from": "posts/p3", "_to": "tags/frontend", "confidence": 0.95, "source": "ai"},
		{"_from": "posts/p3", "_to": "tags/frameworks", "confidence": 0.93, "source": "ai"},
		{"_from": "posts/p4", "_to": "tags/remote-work", "confidence": 0.87, "source": "ai"},
		{"_from": "posts/p4", "_to": "tags/career-change", "confidence": 0.82, "source": "ai"},
		{"_from": "posts/p5", "_to": "tags/relationships", "confidence": 0.90, "source": "ai"},
		{"_from": "posts/p5", "_to": "tags/investing", "confidence": 0.75, "source": "ai"},
		{"_from": "posts/p6", "_to": "tags/backend-dev", "confidence": 0.91, "source": "ai"},
		{"_from": "posts/p6", "_to": "tags/learning", "confidence": 0.88, "source": "ai"},
		{"_from": "posts/p7", "_to": "tags/health-tips", "confidence": 0.89, "source": "ai"},
		{"_from": "posts/p7", "_to": "tags/productivity", "confidence": 0.80, "source": "ai"},
		{"_from": "posts/p8", "_to": "tags/frontend", "confidence": 0.85, "source": "ai"},
		{"_from": "posts/p8", "_to": "tags/frameworks", "confidence": 0.82, "source": "ai"},
	}

	log.Println("Seeding post_has_tag edges...")
	for _, e := range postHasTag {
		if err := s.insert("post_has_tag", e); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Edge: voted (users -> posts) for polls
	voted := []map[string]any{
		{"_from": "users/u1", "_to": "posts/p3", "option": "React", "createdAt": now - 3*day},
		{"_from": "users/u2", "_to": "posts/p3", "option": "Vue", "createdAt": now - 3*day},
		{"_from": "users/u3", "_to": "posts/p3", "option": "React", "createdAt": now - 2*day},
		{"_from": "users/u4", "_to": "posts/p3", "option": "Svelte", "createdAt": now - 2*day},
		{"_from": "users/u1", "_to": "posts/p7", "option": "7-8 hours", "createdAt": now - 6*60*60*1000},
		{"_from": "users/u3", "_to": "posts/p7", "option": "6-7 hours", "createdAt": now - 4*60*60*1000},
		{"_from": "users/u5", "_to": "posts/p7", "option": "Less than 6", "createdAt": now - 2*60*60*1000},
	}

	log.Println("Seeding voted edges...")
	for _, e := range voted {
		if err := s.insert("voted", e); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	// Seed Edge: tagged (posts -> users) for @mentions in group chat
	tagged := []map[string]any{
		{"_from": "posts/p8", "_to": "users/u2", "createdAt": now - 2*day},
		{"_from": "posts/p8", "_to": "users/u3", "createdAt": now - 2*day},
		{"_from": "posts/p8", "_to": "users/u4", "createdAt": now - 2*day},
	}

	log.Println("Seeding tagged edges...")
	for _, e := range tagged {
		if err := s.insert("tagged", e); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	return nil
}
