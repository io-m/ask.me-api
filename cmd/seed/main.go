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

	log.Println("Seeding database with frontend mock-aligned data...")

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

func (s *Seeder) truncate(collection string) error {
	url := fmt.Sprintf("%s/collection/%s/truncate", s.baseURL, collection)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.password)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("truncate failed for %s: %s", collection, resp.Status)
	}

	log.Printf("Truncated collection: %s", collection)
	return nil
}

func (s *Seeder) insert(collection string, doc any) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}

	// Use overwriteMode=replace to upsert (insert or replace if exists)
	url := fmt.Sprintf("%s/document/%s?overwriteMode=replace", s.baseURL, collection)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("insert failed for %s: %s", collection, resp.Status)
	}

	return nil
}

func (s *Seeder) SeedAll() error {
	// Truncate all collections first to ensure clean state
	collections := []string{
		// Document collections
		"users",
		"posts",
		"chats",
		"messages",
		"tags",
		// Edge collections
		"created",
		"responded",
		"participates_in",
		"post_has_tag",
		"follows",
		"tagged",
		"voted",
		"reacted",
	}

	log.Println("Truncating all collections...")
	for _, collection := range collections {
		if err := s.truncate(collection); err != nil {
			// Log but don't fail - collection might not exist yet
			log.Printf("Warning: could not truncate %s: %v", collection, err)
		}
	}

	now := time.Now().UnixMilli()
	day := int64(24 * 60 * 60 * 1000)
	hour := int64(60 * 60 * 1000)
	minute := int64(60 * 1000)

	// ============================================
	// USERS - Match frontend mock data exactly
	// ============================================
	users := []map[string]any{
		// Main logged-in user
		{
			"_key":          "u-johndoe",
			"username":      "johndoe",
			"avatarUrl":     "https://images.unsplash.com/photo-1633332755192-727a05c4013d?w=400&h=400&fit=crop&crop=face",
			"createdAt":     now - 60*day,
			"interests":     []string{"tech", "career", "lifestyle"},
			"blockedTopics": []string{},
			"settings":      map[string]any{"allowDMs": true, "allowTagging": true},
			"stats":         map[string]any{"postsCreated": 5, "responsesGiven": 15},
		},
		// Post authors from posts.ts
		{
			"_key":      "u-sandro",
			"username":  "sandrogasparella",
			"avatarUrl": "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 50*day,
			"interests": []string{"relationships"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-maria",
			"username":  "mariachen",
			"avatarUrl": "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 45*day,
			"interests": []string{"career", "tech"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-alex",
			"username":  "alexthompson",
			"avatarUrl": nil,
			"createdAt": now - 40*day,
			"interests": []string{"health", "work"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-jordan",
			"username":  "jordanlee",
			"avatarUrl": "https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 35*day,
			"interests": []string{"lifestyle", "wisdom"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-sam",
			"username":  "samrivera",
			"avatarUrl": nil,
			"createdAt": now - 30*day,
			"interests": []string{"career", "finance"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-emma",
			"username":  "emmawatson",
			"avatarUrl": "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 28*day,
			"interests": []string{"lifestyle", "work"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-lucas",
			"username":  "lucasmartinez",
			"avatarUrl": "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 25*day,
			"interests": []string{"lifestyle", "lessons"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-sophia",
			"username":  "sophiakim",
			"avatarUrl": nil,
			"createdAt": now - 22*day,
			"interests": []string{"relationships"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-oliver",
			"username":  "oliverbrown",
			"avatarUrl": "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 20*day,
			"interests": []string{"lifestyle", "habits"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-ava",
			"username":  "avajohnson",
			"avatarUrl": "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 18*day,
			"interests": []string{"health", "motivation"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-dev",
			"username":  "devmaster",
			"avatarUrl": "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 15*day,
			"interests": []string{"tech", "frameworks"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-nina",
			"username":  "ninarodriguez",
			"avatarUrl": "https://images.unsplash.com/photo-1517841905240-472988babdf9?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 12*day,
			"interests": []string{"career", "remote"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		// Additional users from chats.ts
		{
			"_key":      "u-alice",
			"username":  "alicewonders",
			"avatarUrl": "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 40*day,
			"interests": []string{"productivity"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-bob",
			"username":  "bobthebuilder",
			"avatarUrl": "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 38*day,
			"interests": []string{"remote work"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-sarah",
			"username":  "sarahtravel",
			"avatarUrl": "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 36*day,
			"interests": []string{"travel"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-mike",
			"username":  "mikeadventure",
			"avatarUrl": "https://images.unsplash.com/photo-1570295999919-56ceb5ecca61?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 34*day,
			"interests": []string{"travel", "adventure"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-david",
			"username":  "davidtech",
			"avatarUrl": "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 32*day,
			"interests": []string{"tech"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-lisa",
			"username":  "lisareads",
			"avatarUrl": "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 30*day,
			"interests": []string{"books", "reading"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-tom",
			"username":  "tombooks",
			"avatarUrl": "https://images.unsplash.com/photo-1552058544-f2b08422138a?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 28*day,
			"interests": []string{"books", "sci-fi"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-rachel",
			"username":  "rachelwrites",
			"avatarUrl": "https://images.unsplash.com/photo-1517841905240-472988babdf9?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 26*day,
			"interests": []string{"books", "writing"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-maya",
			"username":  "mayacafe",
			"avatarUrl": "https://images.unsplash.com/photo-1487412720507-e7ab37603c6f?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 24*day,
			"interests": []string{"coffee", "tea"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-chris",
			"username":  "chrisbarista",
			"avatarUrl": "https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 22*day,
			"interests": []string{"coffee"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
		{
			"_key":      "u-anna",
			"username":  "annatea",
			"avatarUrl": "https://images.unsplash.com/photo-1508214751196-bcfd4ca60f91?w=400&h=400&fit=crop&crop=face",
			"createdAt": now - 20*day,
			"interests": []string{"tea"},
			"settings":  map[string]any{"allowDMs": true, "allowTagging": true},
		},
	}

	log.Println("Seeding users...")
	for _, u := range users {
		if err := s.insert("users", u); err != nil {
			log.Printf("Warning inserting user %s: %v", u["_key"], err)
		}
	}

	// ============================================
	// TAGS
	// ============================================
	tags := []map[string]any{
		{"_key": "relationships", "label": "Relationships", "aliases": []string{"dating", "love", "trust"}, "usageCount": 25, "createdAt": now - 60*day},
		{"_key": "career", "label": "Career", "aliases": []string{"work", "job", "profession"}, "usageCount": 30, "createdAt": now - 60*day},
		{"_key": "tech", "label": "Tech", "aliases": []string{"technology", "programming"}, "usageCount": 20, "createdAt": now - 60*day},
		{"_key": "health", "label": "Health", "aliases": []string{"wellness", "mental health"}, "usageCount": 15, "createdAt": now - 60*day},
		{"_key": "lifestyle", "label": "Lifestyle", "aliases": []string{"life", "living"}, "usageCount": 18, "createdAt": now - 60*day},
		{"_key": "advice", "label": "Advice", "aliases": []string{"tips", "guidance"}, "usageCount": 22, "createdAt": now - 60*day},
		{"_key": "growth", "label": "Growth", "aliases": []string{"self-improvement", "development"}, "usageCount": 16, "createdAt": now - 60*day},
		{"_key": "anxiety", "label": "Anxiety", "aliases": []string{"stress", "worry"}, "usageCount": 10, "createdAt": now - 55*day},
		{"_key": "decisions", "label": "Decisions", "aliases": []string{"choices", "decision-making"}, "usageCount": 12, "createdAt": now - 55*day},
		{"_key": "work-life-balance", "label": "Work-Life Balance", "aliases": []string{"balance", "boundaries"}, "usageCount": 14, "createdAt": now - 50*day},
		{"_key": "habits", "label": "Habits", "aliases": []string{"routine", "self-improvement"}, "usageCount": 11, "createdAt": now - 50*day},
		{"_key": "motivation", "label": "Motivation", "aliases": []string{"inspiration", "drive"}, "usageCount": 13, "createdAt": now - 45*day},
		{"_key": "frontend", "label": "Frontend", "aliases": []string{"ui", "web development"}, "usageCount": 8, "createdAt": now - 40*day},
		{"_key": "frameworks", "label": "Frameworks", "aliases": []string{"libraries", "tools"}, "usageCount": 9, "createdAt": now - 40*day},
		{"_key": "remote-work", "label": "Remote Work", "aliases": []string{"wfh", "work from home"}, "usageCount": 7, "createdAt": now - 35*day},
	}

	log.Println("Seeding tags...")
	for _, t := range tags {
		if err := s.insert("tags", t); err != nil {
			log.Printf("Warning inserting tag %s: %v", t["_key"], err)
		}
	}

	// ============================================
	// POSTS - Match posts.ts exactly
	// ============================================
	posts := []map[string]any{
		// p1 - Has chat c1, has lastMessage from logged-in user
		{
			"_key":        "p1",
			"authorId":    "users/u-sandro",
			"postType":    "text",
			"text":        "I have an issue with my boyfriend. He is cheating on me. Did you ever end up in that situation, and what did you do?",
			"category":    "relationships",
			"intent":      "seeking-advice",
			"depth":       "serious",
			"pollOptions": nil,
			"createdAt":   now - 5*day, // matches mock's 1704240000000 relative position
		},
		// p2 - Has chat c2, no lastMessage in mock
		{
			"_key":        "p2",
			"authorId":    "users/u-maria",
			"postType":    "text",
			"text":        "What was the most difficult decision you ever had to make in your career?",
			"category":    "career",
			"intent":      "asking-question",
			"depth":       "neutral",
			"pollOptions": nil,
			"createdAt":   now - 10*day,
		},
		// p3 - No chat, no lastMessage
		{
			"_key":        "p3",
			"authorId":    "users/u-alex",
			"postType":    "text",
			"text":        "How do you deal with anxiety before important meetings or presentations? I often feel my heart racing and my palms getting sweaty, and sometimes I even struggle to find the right words when speaking in front of a large group of people. This has been affecting my professional growth significantly, and I find myself avoiding opportunities that would require public speaking or leading team meetings. I have tried breathing exercises and meditation, but the physical symptoms still overwhelm me. Have you found any techniques that genuinely help calm your nerves and allow you to perform confidently under pressure? I would really appreciate any advice or personal experiences you could share, as this is becoming a major obstacle in my career development and I really want to overcome this fear once and for all.",
			"category":    "health",
			"intent":      "seeking-advice",
			"depth":       "serious",
			"pollOptions": nil,
			"createdAt":   now - 7*day,
		},
		// p4 - Has chat c4, has lastMessage from post author (jordanlee)
		{
			"_key":        "p4",
			"authorId":    "users/u-jordan",
			"postType":    "text",
			"text":        "What advice would you give to your 20-year-old self?",
			"category":    "lifestyle",
			"intent":      "asking-question",
			"depth":       "neutral",
			"pollOptions": nil,
			"createdAt":   now - 8*day,
		},
		// p5 - No chat
		{
			"_key":        "p5",
			"authorId":    "users/u-sam",
			"postType":    "text",
			"text":        "Have you ever completely changed your career path after spending many years in one field? How did you know it was the right time to make such a drastic change, and what gave you the courage to take that leap of faith into the unknown? I have been working in finance for over a decade now, and while it pays well, I feel completely unfulfilled and disconnected from my work. The thought of starting over in a completely different industry both excites and terrifies me. I keep wondering if I am being naive or if I should just accept that work is meant to be tolerable, not necessarily meaningful. What signs did you recognize that told you it was time to pivot, and how did you manage the financial uncertainty that comes with such a transition? Did you have a safety net, or did you take the risk without one? I am curious about the practical aspects of making such a change, especially when you have responsibilities and commitments like a mortgage and family to support. Any insights from people who have been through this would be incredibly valuable to me right now.",
			"category":    "career",
			"intent":      "seeking-advice",
			"depth":       "serious",
			"pollOptions": nil,
			"createdAt":   now - 6*day,
		},
		// p6 - Has chat c6, has lastMessage from logged-in user
		{
			"_key":        "p6",
			"authorId":    "users/u-emma",
			"postType":    "text",
			"text":        "How do you maintain work-life balance when you love what you do?",
			"category":    "lifestyle",
			"intent":      "asking-question",
			"depth":       "neutral",
			"pollOptions": nil,
			"createdAt":   now - 5*day,
		},
		// p7 - Has chat c7, has lastMessage from post author (lucas)
		{
			"_key":        "p7",
			"authorId":    "users/u-lucas",
			"postType":    "text",
			"text":        "What is the biggest mistake you made in your 20s that taught you a valuable lesson?",
			"category":    "lifestyle",
			"intent":      "asking-question",
			"depth":       "neutral",
			"pollOptions": nil,
			"createdAt":   now - 14*day,
		},
		// p8 - No chat
		{
			"_key":        "p8",
			"authorId":    "users/u-sophia",
			"postType":    "text",
			"text":        "How do you handle criticism from people you care about?",
			"category":    "relationships",
			"intent":      "seeking-advice",
			"depth":       "serious",
			"pollOptions": nil,
			"createdAt":   now - 3*day,
		},
		// p9 - No chat
		{
			"_key":        "p9",
			"authorId":    "users/u-oliver",
			"postType":    "text",
			"text":        "What habit changed your life the most and how long did it take to develop?",
			"category":    "lifestyle",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"createdAt":   now - 2*day,
		},
		// p10 - No chat
		{
			"_key":        "p10",
			"authorId":    "users/u-ava",
			"postType":    "text",
			"text":        "How do you stay motivated when everything seems to be going wrong?",
			"category":    "health",
			"intent":      "seeking-advice",
			"depth":       "serious",
			"pollOptions": nil,
			"createdAt":   now - 1*day,
		},
		// p11 - Poll
		{
			"_key":        "p11",
			"authorId":    "users/u-dev",
			"postType":    "poll",
			"text":        "Which frontend framework do you prefer for new projects?",
			"pollOptions": []string{"React", "Vue", "Angular", "Svelte"},
			"category":    "tech",
			"intent":      "seeking-opinion",
			"depth":       "casual",
			"createdAt":   now - 12*hour,
		},
		// p12 - Sharing story
		{
			"_key":        "p12",
			"authorId":    "users/u-nina",
			"postType":    "text",
			"text":        "I just landed my first remote job after 6 months of searching! 🎉",
			"category":    "career",
			"intent":      "sharing-story",
			"depth":       "casual",
			"pollOptions": nil,
			"createdAt":   now - 6*hour,
		},
		// Additional posts from chats.ts (user's own questions)
		{
			"_key":        "q-mine-1",
			"authorId":    "users/u-johndoe",
			"postType":    "text",
			"text":        "How do you stay focused when working from home?",
			"category":    "lifestyle",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"createdAt":   now - 1*day,
		},
		{
			"_key":        "q-mine-2",
			"authorId":    "users/u-johndoe",
			"postType":    "text",
			"text":        "Anyone else struggle with imposter syndrome in tech?",
			"category":    "career",
			"intent":      "asking-question",
			"depth":       "neutral",
			"pollOptions": nil,
			"createdAt":   now - 12*hour,
		},
		{
			"_key":        "q-mine-3",
			"authorId":    "users/u-johndoe",
			"postType":    "text",
			"text":        "Planning a trip to Japan next month! @sarah @mike @emma - want to join? Any recommendations?",
			"category":    "lifestyle",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"createdAt":   now - 6*hour,
		},
		{
			"_key":        "q-mine-4",
			"authorId":    "users/u-johndoe",
			"postType":    "text",
			"text":        "Book club time! @lisa @tom @rachel - What should we read next? I am thinking sci-fi 📚",
			"category":    "fun",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"createdAt":   now - 4*day,
		},
		// Group chat posts from others (chats.ts)
		{
			"_key":        "p-group-1",
			"authorId":    "users/u-david",
			"postType":    "text",
			"text":        "Tech folks! @johndoe @alex @nina - What stack are you using for your side projects in 2024?",
			"category":    "tech",
			"intent":      "asking-question",
			"depth":       "casual",
			"pollOptions": nil,
			"createdAt":   now - 2*day,
		},
		{
			"_key":        "p-group-2",
			"authorId":    "users/u-maya",
			"postType":    "text",
			"text":        "Coffee vs Tea debate! @johndoe @chris @anna - Which team are you on? ☕🍵",
			"category":    "fun",
			"intent":      "seeking-opinion",
			"depth":       "casual",
			"pollOptions": nil,
			"createdAt":   now - 10*day,
		},
	}

	log.Println("Seeding posts...")
	for _, p := range posts {
		if err := s.insert("posts", p); err != nil {
			log.Printf("Warning inserting post %s: %v", p["_key"], err)
		}
	}

	// ============================================
	// CHATS - Consolidated (no duplicates per post+user)
	// Posts.ts chatIds: c1 (p1), c2 (p2), c4 (p4), c6 (p6), c7 (p7)
	// Chats.ts: chat-1,2 (q-mine-1), chat-7,8 (q-mine-2), groups
	// ============================================
	chats := []map[string]any{
		// Direct chats for posts johndoe answered (from posts.ts)
		{"_key": "c1", "postId": "posts/p1", "type": "direct", "createdAt": now - 4*day, "participantCount": 2},   // sandro's relationship post
		{"_key": "c4", "postId": "posts/p4", "type": "direct", "createdAt": now - 7*day, "participantCount": 2},   // jordan's 20-year-old advice
		{"_key": "c6", "postId": "posts/p6", "type": "direct", "createdAt": now - 4*day, "participantCount": 2},   // emma's work-life balance
		{"_key": "c7", "postId": "posts/p7", "type": "direct", "createdAt": now - 13*day, "participantCount": 2},  // lucas's 20s mistake

		// Direct chats for johndoe's questions (multiple responders to same question)
		{"_key": "chat-1", "postId": "posts/q-mine-1", "type": "direct", "createdAt": now - 1*day, "participantCount": 2},  // alice answered
		{"_key": "chat-2", "postId": "posts/q-mine-1", "type": "direct", "createdAt": now - 1*day, "participantCount": 2},  // bob answered
		{"_key": "chat-7", "postId": "posts/q-mine-2", "type": "direct", "createdAt": now - 12*hour, "participantCount": 2}, // maria answered imposter syndrome
		{"_key": "chat-8", "postId": "posts/q-mine-2", "type": "direct", "createdAt": now - 12*hour, "participantCount": 2}, // oliver answered imposter syndrome

		// Group chats
		{"_key": "chat-group-1", "postId": "posts/q-mine-3", "type": "group", "createdAt": now - 6*hour, "participantCount": 4},  // Japan trip
		{"_key": "chat-group-2", "postId": "posts/p-group-1", "type": "group", "createdAt": now - 2*day, "participantCount": 4},  // Tech stack (david's question)
		{"_key": "chat-group-3", "postId": "posts/q-mine-4", "type": "group", "createdAt": now - 4*day, "participantCount": 4},  // Book club
		{"_key": "chat-group-4", "postId": "posts/p-group-2", "type": "group", "createdAt": now - 10*day, "participantCount": 4}, // Coffee vs Tea
	}

	log.Println("Seeding chats...")
	for _, c := range chats {
		if err := s.insert("chats", c); err != nil {
			log.Printf("Warning inserting chat %s: %v", c["_key"], err)
		}
	}

	// ============================================
	// MESSAGES - Match lastMessage data from mocks
	// ============================================
	messages := []map[string]any{
		// Messages for c1 (p1 - sandro's post, johndoe responded)
		{"_key": "m1-1", "chatId": "chats/c1", "senderId": "users/u-johndoe", "text": "That sounds really tough. Have you tried talking to him about it?", "status": "seen", "createdAt": now - 4*day},
		{"_key": "m1-2", "chatId": "chats/c1", "senderId": "users/u-sandro", "text": "Thank you so much for your advice, it really helped me see things clearly 💙", "status": "seen", "createdAt": now - 2*hour},

		// Messages for c4 (p4 - jordan's post about 20-year-old self)
		{"_key": "m4-1", "chatId": "chats/c4", "senderId": "users/u-johndoe", "text": "Invest early, even small amounts. Compound interest is magical!", "status": "seen", "createdAt": now - 5*day},
		{"_key": "m4-2", "chatId": "chats/c4", "senderId": "users/u-jordan", "text": "That's great advice! I wish I had started earlier.", "status": "seen", "createdAt": now - 4*day},
		{"_key": "m4-3", "chatId": "chats/c4", "senderId": "users/u-jordan", "text": "Wow, that's actually really insightful advice. Thank you for sharing! 🙏", "status": "seen", "createdAt": now - 3*day},

		// Messages for c6 (p6 - emma's work-life balance)
		{"_key": "m6-1", "chatId": "chats/c6", "senderId": "users/u-johndoe", "text": "I struggle with this too! Setting boundaries has helped me a lot.", "status": "delivered", "createdAt": now - 2*day},

		// Messages for c7 (p7 - lucas's 20s mistake)
		{"_key": "m7-1", "chatId": "chats/c7", "senderId": "users/u-johndoe", "text": "Not saving enough money early on. Learned the hard way!", "status": "seen", "createdAt": now - 10*day},
		{"_key": "m7-2", "chatId": "chats/c7", "senderId": "users/u-lucas", "text": "That really resonates with me. Thanks for sharing! 🙌", "status": "seen", "createdAt": now - 7*day},

		// Messages for chat-1 (johndoe's q-mine-1, alice responded)
		{"_key": "m-chat1-1", "chatId": "chats/chat-1", "senderId": "users/u-alice", "text": "I use the Pomodoro technique! 25 min work, 5 min break. Game changer 🍅", "status": "seen", "createdAt": now - 5*minute},

		// Messages for chat-2 (johndoe's q-mine-1, bob responded)
		{"_key": "m-chat2-1", "chatId": "chats/chat-2", "senderId": "users/u-bob", "text": "Having a dedicated workspace is key. Even if it is just a corner of your room!", "status": "delivered", "createdAt": now - 30*minute},

		// Messages for chat-7 (johndoe's imposter syndrome, maria responded)
		{"_key": "m-chat7-1", "chatId": "chats/chat-7", "senderId": "users/u-maria", "text": "Every. Single. Day. But remember, you earned your spot! 💪", "status": "delivered", "createdAt": now - 5*hour},

		// Messages for chat-8 (johndoe's imposter syndrome, oliver responded)
		{"_key": "m-chat8-1", "chatId": "chats/chat-8", "senderId": "users/u-oliver", "text": "The seniors feel it too. It never fully goes away, but you learn to manage it.", "status": "seen", "createdAt": now - 8*hour},

		// Group chat messages - chat-group-1 (Japan trip)
		{"_key": "m-g1-1", "chatId": "chats/chat-group-1", "senderId": "users/u-sarah", "text": "Yes! I have so many recommendations. Tokyo is amazing in spring 🌸", "status": "seen", "createdAt": now - 15*minute},

		// Group chat messages - chat-group-2 (Tech stack)
		{"_key": "m-g2-1", "chatId": "chats/chat-group-2", "senderId": "users/u-maria", "text": "Next.js + Prisma + PostgreSQL for me!", "status": "seen", "createdAt": now - 2*day + 1*hour},
		{"_key": "m-g2-2", "chatId": "chats/chat-group-2", "senderId": "users/u-alex", "text": "I switched to Bun + Hono for backend. So fast! 🚀", "status": "delivered", "createdAt": now - 45*minute},

		// Group chat messages - chat-group-3 (Book club)
		{"_key": "m-g3-1", "chatId": "chats/chat-group-3", "senderId": "users/u-tom", "text": "Project Hail Mary by Andy Weir! Its amazing 🚀", "status": "seen", "createdAt": now - 3*hour},

		// Group chat messages - chat-group-4 (Coffee vs Tea)
		{"_key": "m-g4-1", "chatId": "chats/chat-group-4", "senderId": "users/u-chris", "text": "Team coffee all the way! ☕", "status": "seen", "createdAt": now - 8*day},
		{"_key": "m-g4-2", "chatId": "chats/chat-group-4", "senderId": "users/u-anna", "text": "Team matcha all the way! 🍵✨", "status": "seen", "createdAt": now - 5*day},
	}

	log.Println("Seeding messages...")
	for _, m := range messages {
		if err := s.insert("messages", m); err != nil {
			log.Printf("Warning inserting message %s: %v", m["_key"], err)
		}
	}

	// ============================================
	// EDGES: created (users -> posts)
	// ============================================
	created := []map[string]any{
		{"_from": "users/u-sandro", "_to": "posts/p1", "createdAt": now - 5*day},
		{"_from": "users/u-maria", "_to": "posts/p2", "createdAt": now - 10*day},
		{"_from": "users/u-alex", "_to": "posts/p3", "createdAt": now - 7*day},
		{"_from": "users/u-jordan", "_to": "posts/p4", "createdAt": now - 8*day},
		{"_from": "users/u-sam", "_to": "posts/p5", "createdAt": now - 6*day},
		{"_from": "users/u-emma", "_to": "posts/p6", "createdAt": now - 5*day},
		{"_from": "users/u-lucas", "_to": "posts/p7", "createdAt": now - 14*day},
		{"_from": "users/u-sophia", "_to": "posts/p8", "createdAt": now - 3*day},
		{"_from": "users/u-oliver", "_to": "posts/p9", "createdAt": now - 2*day},
		{"_from": "users/u-ava", "_to": "posts/p10", "createdAt": now - 1*day},
		{"_from": "users/u-dev", "_to": "posts/p11", "createdAt": now - 12*hour},
		{"_from": "users/u-nina", "_to": "posts/p12", "createdAt": now - 6*hour},
		{"_from": "users/u-johndoe", "_to": "posts/q-mine-1", "createdAt": now - 1*day},
		{"_from": "users/u-johndoe", "_to": "posts/q-mine-2", "createdAt": now - 12*hour},
		{"_from": "users/u-johndoe", "_to": "posts/q-mine-3", "createdAt": now - 6*hour},
		{"_from": "users/u-johndoe", "_to": "posts/q-mine-4", "createdAt": now - 4*day},
		{"_from": "users/u-david", "_to": "posts/p-group-1", "createdAt": now - 2*day},
		{"_from": "users/u-maya", "_to": "posts/p-group-2", "createdAt": now - 10*day},
	}

	log.Println("Seeding created edges...")
	for _, e := range created {
		if err := s.insert("created", e); err != nil {
			log.Printf("Warning inserting created edge: %v", err)
		}
	}

	// ============================================
	// EDGES: responded (users -> posts)
	// ============================================
	responded := []map[string]any{
		// Johndoe responded to other people's posts
		{"_from": "users/u-johndoe", "_to": "posts/p1", "chatId": "chats/c1", "createdAt": now - 4*day},   // sandro's post
		{"_from": "users/u-johndoe", "_to": "posts/p4", "chatId": "chats/c4", "createdAt": now - 7*day},   // jordan's post
		{"_from": "users/u-johndoe", "_to": "posts/p6", "chatId": "chats/c6", "createdAt": now - 4*day},   // emma's post
		{"_from": "users/u-johndoe", "_to": "posts/p7", "chatId": "chats/c7", "createdAt": now - 13*day},  // lucas's post
		// Others responded to johndoe's posts
		{"_from": "users/u-alice", "_to": "posts/q-mine-1", "chatId": "chats/chat-1", "createdAt": now - 1*day},
		{"_from": "users/u-bob", "_to": "posts/q-mine-1", "chatId": "chats/chat-2", "createdAt": now - 1*day},
		{"_from": "users/u-maria", "_to": "posts/q-mine-2", "chatId": "chats/chat-7", "createdAt": now - 12*hour},
		{"_from": "users/u-oliver", "_to": "posts/q-mine-2", "chatId": "chats/chat-8", "createdAt": now - 12*hour},
	}

	log.Println("Seeding responded edges...")
	for _, e := range responded {
		if err := s.insert("responded", e); err != nil {
			log.Printf("Warning inserting responded edge: %v", err)
		}
	}

	// ============================================
	// EDGES: participates_in (users -> chats)
	// ============================================
	participatesIn := []map[string]any{
		// c1: sandro (author) and johndoe (responder) - for p1
		{"_from": "users/u-sandro", "_to": "chats/c1", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 5*day},
		{"_from": "users/u-johndoe", "_to": "chats/c1", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 4*day},

		// c4: jordan (author) and johndoe (responder) - for p4
		{"_from": "users/u-jordan", "_to": "chats/c4", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 8*day},
		{"_from": "users/u-johndoe", "_to": "chats/c4", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 7*day},

		// c6: emma (author) and johndoe (responder) - for p6
		{"_from": "users/u-emma", "_to": "chats/c6", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 5*day},
		{"_from": "users/u-johndoe", "_to": "chats/c6", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 4*day},

		// c7: lucas (author) and johndoe (responder) - for p7
		{"_from": "users/u-lucas", "_to": "chats/c7", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 14*day},
		{"_from": "users/u-johndoe", "_to": "chats/c7", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 13*day},

		// chat-1: johndoe (author) and alice (responder) - for q-mine-1
		{"_from": "users/u-johndoe", "_to": "chats/chat-1", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 1*day},
		{"_from": "users/u-alice", "_to": "chats/chat-1", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 1*day},

		// chat-2: johndoe (author) and bob (responder) - for q-mine-1
		{"_from": "users/u-johndoe", "_to": "chats/chat-2", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 1*day},
		{"_from": "users/u-bob", "_to": "chats/chat-2", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 1*day},

		// chat-7: johndoe (author) and maria (responder) - for q-mine-2
		{"_from": "users/u-johndoe", "_to": "chats/chat-7", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 12*hour},
		{"_from": "users/u-maria", "_to": "chats/chat-7", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 12*hour},

		// chat-8: johndoe (author) and oliver (responder) - for q-mine-2
		{"_from": "users/u-johndoe", "_to": "chats/chat-8", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 12*hour},
		{"_from": "users/u-oliver", "_to": "chats/chat-8", "role": "responder", "status": "active", "notificationsEnabled": true, "joinedAt": now - 12*hour},

		// chat-group-1: Japan trip (johndoe author, sarah/mike invited, emma pending) - for q-mine-3
		{"_from": "users/u-johndoe", "_to": "chats/chat-group-1", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 6*hour},
		{"_from": "users/u-sarah", "_to": "chats/chat-group-1", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 6*hour},
		{"_from": "users/u-mike", "_to": "chats/chat-group-1", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 6*hour},
		{"_from": "users/u-emma", "_to": "chats/chat-group-1", "role": "invited", "status": "pending", "notificationsEnabled": false, "joinedAt": nil},

		// chat-group-2: Tech stack (david author, johndoe/alex/nina invited) - for p-group-1
		{"_from": "users/u-david", "_to": "chats/chat-group-2", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},
		{"_from": "users/u-johndoe", "_to": "chats/chat-group-2", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},
		{"_from": "users/u-alex", "_to": "chats/chat-group-2", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},
		{"_from": "users/u-nina", "_to": "chats/chat-group-2", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 2*day},

		// chat-group-3: Book club (johndoe author, lisa/tom/rachel invited, rachel muted) - for q-mine-4
		{"_from": "users/u-johndoe", "_to": "chats/chat-group-3", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 4*day},
		{"_from": "users/u-lisa", "_to": "chats/chat-group-3", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 4*day},
		{"_from": "users/u-tom", "_to": "chats/chat-group-3", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 4*day},
		{"_from": "users/u-rachel", "_to": "chats/chat-group-3", "role": "invited", "status": "muted", "notificationsEnabled": false, "joinedAt": now - 4*day},

		// chat-group-4: Coffee vs Tea (maya author, johndoe/chris/anna invited) - for p-group-2
		{"_from": "users/u-maya", "_to": "chats/chat-group-4", "role": "author", "status": "active", "notificationsEnabled": true, "joinedAt": now - 10*day},
		{"_from": "users/u-johndoe", "_to": "chats/chat-group-4", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 10*day},
		{"_from": "users/u-chris", "_to": "chats/chat-group-4", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 10*day},
		{"_from": "users/u-anna", "_to": "chats/chat-group-4", "role": "invited", "status": "active", "notificationsEnabled": true, "joinedAt": now - 10*day},
	}

	log.Println("Seeding participates_in edges...")
	for _, e := range participatesIn {
		if err := s.insert("participates_in", e); err != nil {
			log.Printf("Warning inserting participates_in edge: %v", err)
		}
	}

	// ============================================
	// EDGES: post_has_tag (posts -> tags)
	// ============================================
	postHasTag := []map[string]any{
		{"_from": "posts/p1", "_to": "tags/relationships", "confidence": 0.95, "source": "ai"},
		{"_from": "posts/p1", "_to": "tags/advice", "confidence": 0.88, "source": "ai"},
		{"_from": "posts/p2", "_to": "tags/career", "confidence": 0.92, "source": "ai"},
		{"_from": "posts/p2", "_to": "tags/decisions", "confidence": 0.85, "source": "ai"},
		{"_from": "posts/p3", "_to": "tags/anxiety", "confidence": 0.90, "source": "ai"},
		{"_from": "posts/p3", "_to": "tags/health", "confidence": 0.88, "source": "ai"},
		{"_from": "posts/p4", "_to": "tags/lifestyle", "confidence": 0.91, "source": "ai"},
		{"_from": "posts/p4", "_to": "tags/advice", "confidence": 0.86, "source": "ai"},
		{"_from": "posts/p5", "_to": "tags/career", "confidence": 0.93, "source": "ai"},
		{"_from": "posts/p5", "_to": "tags/decisions", "confidence": 0.89, "source": "ai"},
		{"_from": "posts/p6", "_to": "tags/work-life-balance", "confidence": 0.94, "source": "ai"},
		{"_from": "posts/p6", "_to": "tags/lifestyle", "confidence": 0.87, "source": "ai"},
		{"_from": "posts/p7", "_to": "tags/lifestyle", "confidence": 0.90, "source": "ai"},
		{"_from": "posts/p7", "_to": "tags/growth", "confidence": 0.85, "source": "ai"},
		{"_from": "posts/p8", "_to": "tags/relationships", "confidence": 0.91, "source": "ai"},
		{"_from": "posts/p9", "_to": "tags/habits", "confidence": 0.93, "source": "ai"},
		{"_from": "posts/p9", "_to": "tags/lifestyle", "confidence": 0.86, "source": "ai"},
		{"_from": "posts/p10", "_to": "tags/motivation", "confidence": 0.92, "source": "ai"},
		{"_from": "posts/p10", "_to": "tags/health", "confidence": 0.84, "source": "ai"},
		{"_from": "posts/p11", "_to": "tags/frontend", "confidence": 0.95, "source": "ai"},
		{"_from": "posts/p11", "_to": "tags/frameworks", "confidence": 0.90, "source": "ai"},
		{"_from": "posts/p12", "_to": "tags/career", "confidence": 0.88, "source": "ai"},
		{"_from": "posts/p12", "_to": "tags/remote-work", "confidence": 0.92, "source": "ai"},
	}

	log.Println("Seeding post_has_tag edges...")
	for _, e := range postHasTag {
		if err := s.insert("post_has_tag", e); err != nil {
			log.Printf("Warning inserting post_has_tag edge: %v", err)
		}
	}

	// ============================================
	// EDGES: follows (users -> users)
	// ============================================
	follows := []map[string]any{
		{"_from": "users/u-johndoe", "_to": "users/u-sandro", "createdAt": now - 30*day},
		{"_from": "users/u-johndoe", "_to": "users/u-maria", "createdAt": now - 28*day},
		{"_from": "users/u-johndoe", "_to": "users/u-emma", "createdAt": now - 25*day},
		{"_from": "users/u-maria", "_to": "users/u-johndoe", "createdAt": now - 27*day},
		{"_from": "users/u-sandro", "_to": "users/u-johndoe", "createdAt": now - 29*day},
		{"_from": "users/u-emma", "_to": "users/u-johndoe", "createdAt": now - 24*day},
		{"_from": "users/u-alice", "_to": "users/u-johndoe", "createdAt": now - 20*day},
		{"_from": "users/u-bob", "_to": "users/u-johndoe", "createdAt": now - 18*day},
	}

	log.Println("Seeding follows edges...")
	for _, e := range follows {
		if err := s.insert("follows", e); err != nil {
			log.Printf("Warning inserting follows edge: %v", err)
		}
	}

	// ============================================
	// EDGES: tagged (posts -> users) for @mentions
	// ============================================
	tagged := []map[string]any{
		{"_from": "posts/q-mine-3", "_to": "users/u-sarah", "createdAt": now - 6*hour},
		{"_from": "posts/q-mine-3", "_to": "users/u-mike", "createdAt": now - 6*hour},
		{"_from": "posts/q-mine-3", "_to": "users/u-emma", "createdAt": now - 6*hour},
		{"_from": "posts/q-mine-4", "_to": "users/u-lisa", "createdAt": now - 4*day},
		{"_from": "posts/q-mine-4", "_to": "users/u-tom", "createdAt": now - 4*day},
		{"_from": "posts/q-mine-4", "_to": "users/u-rachel", "createdAt": now - 4*day},
		{"_from": "posts/p-group-1", "_to": "users/u-johndoe", "createdAt": now - 2*day},
		{"_from": "posts/p-group-1", "_to": "users/u-alex", "createdAt": now - 2*day},
		{"_from": "posts/p-group-1", "_to": "users/u-nina", "createdAt": now - 2*day},
		{"_from": "posts/p-group-2", "_to": "users/u-johndoe", "createdAt": now - 10*day},
		{"_from": "posts/p-group-2", "_to": "users/u-chris", "createdAt": now - 10*day},
		{"_from": "posts/p-group-2", "_to": "users/u-anna", "createdAt": now - 10*day},
	}

	log.Println("Seeding tagged edges...")
	for _, e := range tagged {
		if err := s.insert("tagged", e); err != nil {
			log.Printf("Warning inserting tagged edge: %v", err)
		}
	}

	// ============================================
	// EDGES: voted (users -> posts) for polls
	// ============================================
	voted := []map[string]any{
		{"_from": "users/u-johndoe", "_to": "posts/p11", "option": "React", "createdAt": now - 10*hour},
		{"_from": "users/u-maria", "_to": "posts/p11", "option": "Vue", "createdAt": now - 9*hour},
		{"_from": "users/u-alex", "_to": "posts/p11", "option": "React", "createdAt": now - 8*hour},
		{"_from": "users/u-emma", "_to": "posts/p11", "option": "Svelte", "createdAt": now - 7*hour},
	}

	log.Println("Seeding voted edges...")
	for _, e := range voted {
		if err := s.insert("voted", e); err != nil {
			log.Printf("Warning inserting voted edge: %v", err)
		}
	}

	return nil
}
