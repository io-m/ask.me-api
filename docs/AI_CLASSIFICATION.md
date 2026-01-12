# AI Classification System

## Overview

When a user creates a post, the **backend** analyzes the content using an AI classifier (LLM). The classification enriches the post with structured metadata enabling:

- **Smart feed curation** - Show relevant content based on preferences
- **Content moderation** - Flag potentially harmful content
- **Tag normalization** - Convert free-text to canonical tags
- **Depth matching** - Connect users seeking similar conversation depth

---

## AI Prompt Structure

The backend sends the post content to an LLM with this prompt:

```
Analyze this post and classify:

1. category: Choose the best fit from [career, relationships, tech, health, 
   finance, fun, opinion, lifestyle, education, other]

2. intent: Describe what the user wants in 2-4 words using verb-noun pattern 
   (e.g., "seeking-advice", "sharing-story", "asking-question", "seeking-opinion", 
   "sharing-tip", "starting-discussion")

3. depth: Choose from [casual, neutral, serious]

4. tags: Suggest 2-5 specific tags

5. confidence: Your confidence in this classification (0.0-1.0)

6. risk: Content risk assessment [low, medium, high]

7. flags: Any content flags (empty array if none)

Post text: "{user_post_text}"
Post type: "{postType}" (text or poll)

Respond in JSON format.
```

**Important:** The AI receives `postType` for context but does **not** determine it — the user explicitly chooses text or poll.

---

## Classification Pipeline

```
┌─────────────────────────────────────────────────────────────────────┐
│                        POST CREATION FLOW                           │
└─────────────────────────────────────────────────────────────────────┘

     ① User writes post text
        (chooses text/poll)
               │
               ▼
       ┌───────────────┐
       │  Client App   │
       │(React Native) │
       └───────┬───────┘
               │
               │ ② POST /posts { postType, text }
               ▼
       ┌───────────────┐
       │   API Server  │
       │     (Go)      │
       └───────┬───────┘
               │
               │ ③ Call LLM with post text + prompt
               ▼
       ┌───────────────┐
       │   LLM API     │
       │(OpenAI/Claude)│
       └───────┬───────┘
               │
               │ ④ Returns aiRaw JSON classification
               ▼
       ┌───────────────┐
       │   API Server  │
       │     (Go)      │
       └───────┬───────┘
               │
               │ ⑤ Validate & Normalize
               │
               ├──────────────────┬──────────────────┐
               │                  │                  │
               ▼                  ▼                  ▼
       ┌───────────────┐  ┌───────────────┐  ┌───────────────┐
       │   Validate    │  │   Normalize   │  │   Normalize   │
       │   Category    │  │     Depth     │  │     Tags      │
       │  (enum check) │  │  (enum check) │  │  (slugify)    │
       └───────┬───────┘  └───────┬───────┘  └───────┬───────┘
               │                  │                  │
               └──────────────────┼──────────────────┘
                                  │
                                  ▼
                          ┌───────────────┐
                          │   ArangoDB    │
                          │  ⑥ Store Post │
                          │  Create Edges │
                          └───────────────┘
```

### Why Backend Calls the LLM

| Benefit | Description |
|---------|-------------|
| **Secure API keys** | LLM credentials never exposed to client |
| **Consistent classification** | Same model/prompt version for all posts |
| **Caching & batching** | Can optimize LLM calls, cache similar content |
| **Moderation control** | Can reject high-risk content before storage |
| **Audit trail** | Full logging of classification decisions |

---

## Storage Flow

### Step 1: Frontend Sends Request

Frontend sends only the post content — **no AI classification**:

```json
POST /posts
{
  "authorId": "u1",
  "postType": "text",
  "text": "How do I switch from frontend to backend?"
}
```

For polls:
```json
POST /posts/poll
{
  "authorId": "u5",
  "postType": "poll",
  "text": "Which frontend framework do you prefer?",
  "pollOptions": ["React", "Vue", "Angular", "Svelte"]
}
```

### Step 2: Backend Calls LLM

Backend sends post to LLM and receives `aiRaw`:

```json
{
  "category": "career",
  "intent": "seeking-advice",
  "depth": "serious",
  "tags": ["career switch", "backend dev"],
  "confidence": 0.92,
  "risk": "low",
  "flags": []
}
```

### Step 3: Backend Validates & Normalizes

| Field | Validation | Fallback |
|-------|------------|----------|
| `postType` | Must be `text` or `poll` | Error if invalid |
| `category` | Must match enum | Map to `other` |
| `intent` | Store as-is | Allows natural evolution of intent types |
| `depth` | Must match enum | Map to `neutral` |
| `tags` | Slugify & canonicalize | Create new tags if needed |
| `risk` | Check threshold | Reject if `high` with certain flags |

### Step 4: Create Document & Edges

**Post Document Stored:**
```json
{
  "_key": "p1",
  "authorId": "users/u1",
  "postType": "text",
  "text": "How do I switch from frontend to backend?",
  "category": "career",
  "intent": "seeking-advice",
  "depth": "serious",
  "tags": ["career-change", "backend-dev"],
  "aiRaw": {
    "category": "career",
    "intent": "seeking-advice",
    "depth": "serious",
    "tags": ["career switch", "backend dev"],
    "confidence": 0.92,
    "risk": "low",
    "flags": []
  },
  "createdAt": 1736000000000
}
```

**Edges Created:**
```
users/u1 ──[created]──▶ posts/p1
posts/p1 ──[post_has_tag]──▶ tags/career-change
posts/p1 ──[post_has_tag]──▶ tags/backend-dev
```

### Poll Storage Example

For polls, the flow is identical but includes poll-specific fields:

```json
{
  "_key": "p3",
  "authorId": "users/u5",
  "postType": "poll",
  "text": "Which frontend framework do you prefer?",
  "pollOptions": ["React", "Vue", "Angular", "Svelte"],
  "pollVotes": {},
  "category": "tech",
  "intent": "seeking-opinion",
  "depth": "casual",
  "tags": ["frontend", "frameworks"],
  "aiRaw": {
    "category": "tech",
    "intent": "seeking-opinion",
    "depth": "casual",
    "tags": ["frontend frameworks", "developer tools"],
    "confidence": 0.95,
    "risk": "low",
    "flags": []
  },
  "createdAt": 1736000000000
}
```

---

## Classification Fields

### `aiRaw` - Raw AI Response

The AI classifier returns this structure (sent from client):

```json
{
  "category": "career",
  "intent": "seeking-advice",
  "depth": "serious",
  "tags": ["career change", "backend development", "learning path"],
  "confidence": 0.92,
  "risk": "low",
  "flags": []
}
```

| Field | Type | Description |
|-------|------|-------------|
| `category` | string | Primary topic category |
| `intent` | string | What the user is seeking (verb-noun pattern) |
| `depth` | string | Conversation seriousness level |
| `tags` | string[] | Free-form topic tags (2-5 suggestions) |
| `confidence` | float | AI confidence score (0-1) |
| `risk` | string | Content risk assessment |
| `flags` | string[] | Content warning flags |

---

## Categories

Categories group posts by topic area. Users can set interests and blocked topics.

| Category | Description | Example |
|----------|-------------|---------|
| `career` | Work, jobs, professional growth | "How do I negotiate salary?" |
| `relationships` | Dating, family, friendships | "How to handle difficult roommate?" |
| `tech` | Technology, programming, gadgets | "Best IDE for Python?" |
| `health` | Physical/mental wellness | "Tips for better sleep?" |
| `finance` | Money, investing, budgeting | "Should I max out 401k first?" |
| `fun` | Entertainment, games, hobbies | "Best co-op games in 2024?" |
| `opinion` | Hot takes, debates, preferences | "Is remote work better?" |
| `lifestyle` | Daily life, routines, habits | "Morning routine ideas?" |
| `education` | Learning, courses, self-improvement | "Best way to learn Go?" |
| `other` | Doesn't fit elsewhere | Fallback category |

### Category Normalization

The backend normalizes various AI outputs to canonical categories:

```go
func NormalizeCategory(raw string) PostCategory {
    mapping := map[string]PostCategory{
        "career":        CategoryCareer,
        "job":           CategoryCareer,
        "work":          CategoryCareer,
        "professional":  CategoryCareer,
        "tech":          CategoryTech,
        "technology":    CategoryTech,
        "programming":   CategoryTech,
        // ... more mappings
    }
    if cat, ok := mapping[strings.ToLower(raw)]; ok {
        return cat
    }
    return CategoryOther
}
```

---

## Intent

Intent captures what the user is looking for from responses.

| Intent | Description | Example |
|--------|-------------|---------|
| `asking-question` | Seeking factual information | "What's the capital of France?" |
| `seeking-advice` | Wants guidance/recommendations | "Should I take this job offer?" |
| `seeking-opinion` | Wants perspectives/viewpoints | "What do you think about X?" |
| `sharing` | Sharing experience/information | "Here's what worked for me..." |
| `venting` | Expressing frustration | "I can't believe this happened..." |
| `debating` | Wants discussion/argument | "Change my mind: X is better than Y" |

### Intent in Feed Algorithm

Intent influences how posts are ranked:

- `seeking-advice` → Prioritize users with relevant experience
- `seeking-opinion` → Show to diverse audience
- `venting` → Show to empathetic responders
- `debating` → Show to users who enjoy discussion

---

## Depth

Depth indicates how serious/deep the conversation should be.

| Depth | Description | Tone | Example |
|-------|-------------|------|---------|
| `casual` | Light, fun, quick responses OK | 😄 Playful | "Favorite pizza topping?" |
| `neutral` | Balanced, helpful responses | 🙂 Friendly | "Best laptop for coding?" |
| `serious` | Thoughtful, considered responses | 🤔 Empathetic | "How do I cope with burnout?" |

### Depth Normalization

```go
func NormalizeDepth(raw string) PostDepth {
    switch strings.ToLower(raw) {
    case "casual", "light", "fun":
        return DepthCasual
    case "serious", "deep", "heavy":
        return DepthSerious
    default:
        return DepthNeutral
    }
}
```

### Depth in Feed Algorithm

Users can filter feed by depth:
- Some users want only casual content
- Others prefer serious discussions
- Depth mismatch leads to poor conversations

---

## Tags

Tags provide fine-grained topic classification.

### Raw vs Normalized Tags

**AI returns free-form tags:**
```json
["career change", "backend development", "learning path"]
```

**Backend normalizes to slugs:**
```json
["career-change", "backend-dev", "learning"]
```

### Tag Normalization Process

```
┌─────────────────┐
│  Raw AI Tags    │
│ "career change" │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Slug Conversion │
│ "career-change" │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Alias Lookup   │
│  (tags collection)│
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
 Found    Not Found
    │         │
    ▼         ▼
┌────────┐ ┌────────────┐
│Use _key│ │Create new  │
│        │ │tag record  │
└────────┘ └────────────┘
```

### Tag Entity

```json
{
  "_key": "career-change",
  "label": "Career Change",
  "aliases": [
    "changing careers",
    "career switch",
    "job change",
    "career transition"
  ],
  "usageCount": 47
}
```

When AI returns "changing careers", it maps to `career-change`.

---

## Risk Assessment

The AI evaluates content risk for moderation.

| Risk Level | Action |
|------------|--------|
| `low` | Publish immediately |
| `medium` | Publish with warning |
| `high` | Queue for review |

### Flags

Specific content warnings:

| Flag | Description |
|------|-------------|
| `self-harm` | Mentions of self-harm |
| `hate-speech` | Discriminatory content |
| `explicit` | Adult content |
| `spam` | Promotional/spam content |
| `misinformation` | Potentially false claims |

---

## Confidence Score

AI confidence (0.0 - 1.0) indicates classification certainty.

| Score | Interpretation |
|-------|----------------|
| 0.9+ | High confidence, use classification |
| 0.7-0.9 | Moderate confidence, use with caution |
| < 0.7 | Low confidence, may need manual review |

---

## Data Flow Example

### 1. User Creates Post

Client sends only the content:
```json
POST /posts
{
  "authorId": "u1",
  "postType": "text",
  "text": "How do I transition from frontend to backend engineering?"
}
```

### 2. Backend Calls LLM

```go
// service.go - CreatePost

// 1. Call LLM API with post text
aiRaw := llmClient.Classify(req.Text, req.PostType)
// Returns: { category: "career", intent: "seeking-advice", depth: "serious", ... }

// 2. Normalize category
category := domain.NormalizeCategory(aiRaw.Category)  // → "career"

// 3. Normalize depth
depth := domain.NormalizeDepth(aiRaw.Depth)  // → "serious"

// 4. Normalize tags
tags := normalizeTags(aiRaw.Tags)  // → ["career-change", "backend-dev", "frontend"]

// 5. Check risk level
if aiRaw.Risk == "high" && containsFlag(aiRaw.Flags, "self-harm", "hate-speech") {
    return ErrContentRejected
}

// 6. Create post document
post := &Post{
    AuthorID: "users/" + req.AuthorID,
    PostType: req.PostType,
    Text:     req.Text,
    Category: category,
    Intent:   aiRaw.Intent,
    Depth:    depth,
    Tags:     tags,
    AIRaw:    aiRaw,  // Store original for audit
}
```

### 3. Database Storage

```
┌─────────────────────────────────────────────────────┐
│                    posts/p1                         │
├─────────────────────────────────────────────────────┤
│ authorId: "users/u1"                                │
│ category: "career"                                  │
│ intent: "seeking-advice"                            │
│ depth: "serious"                                    │
│ tags: ["career-change", "backend-dev", "frontend"]  │
│ aiRaw: { ... original AI response ... }             │
└─────────────────────────────────────────────────────┘
          │
          │ [post_has_tag]
          │
          ├──────────────▶ tags/career-change
          ├──────────────▶ tags/backend-dev
          └──────────────▶ tags/frontend
```

### 4. Feed Retrieval

When user requests feed:

```go
// feed/service.go

// 1. Get user interests
user := getUserInterests(userId)  // ["tech", "career"]

// 2. Query posts matching interests
posts := queryPostsByCategories(user.Interests)

// 3. Filter by depth preference (if set)
posts = filterByDepth(posts, user.PreferredDepth)

// 4. Rank by relevance
posts = rankByRelevance(posts, user)

// 5. Return feed
return posts
```

---

## Future Enhancements

### Planned Features

1. **Real-time Classification Preview**
   - WebSocket endpoint for draft classification
   - Show suggested tags as user types
   - Preview category/depth before posting

2. **Learning System**
   - Track which classifications lead to good conversations
   - Adjust confidence thresholds
   - Improve tag suggestions based on engagement

3. **Semantic Search**
   - Enable searching by meaning, not just keywords
   - "Posts about career advice" finds related content

4. **Advanced Moderation**
   - Queue system for high-risk content review
   - User reputation system affecting classification trust
   - Appeal process for rejected content
   - Moderator dashboard for flagged posts

5. **Classification Caching**
   - Cache similar post classifications
   - Reduce LLM API calls for common patterns
   - Batch processing for off-peak optimization
