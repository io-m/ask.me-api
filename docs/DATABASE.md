# ask.me Database Architecture

## Overview

ask.me uses **ArangoDB**, a multi-model database that combines document and graph capabilities. This enables:

- **Document storage** for entities (users, posts, messages)
- **Graph traversal** for relationships (follows, responses, tags)
- **Efficient querying** with AQL (ArangoDB Query Language)

---

## Document Collections

### `users`

User profiles and preferences.

```json
{
  "_key": "u1",
  "username": "alex_dev",
  "createdAt": 1736000000000,
  "interests": ["tech", "career"],
  "blockedTopics": ["politics"],
  "settings": {
    "allowDMs": true,
    "allowTagging": true
  },
  "stats": {
    "postsCreated": 5,
    "responsesGiven": 12
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `_key` | string | Unique user ID |
| `username` | string | Display name |
| `createdAt` | int64 | Unix timestamp (ms) |
| `interests` | string[] | Categories user follows |
| `blockedTopics` | string[] | Categories to hide |
| `settings.allowDMs` | bool | Accept direct messages |
| `settings.allowTagging` | bool | Can be tagged in posts |
| `stats.postsCreated` | int | Total posts authored |
| `stats.responsesGiven` | int | Total responses made |

---

### `posts`

Questions and polls created by users.

**Text Post:**
```json
{
  "_key": "p1",
  "authorId": "users/u1",
  "postType": "text",
  "text": "How do I switch from frontend to backend engineering?",
  "category": "career",
  "intent": "seeking-advice",
  "depth": "serious",
  "tags": ["career-change", "backend-dev"],
  "aiRaw": {
    "category": "career",
    "intent": "seeking-advice",
    "depth": "serious",
    "tags": ["career change", "backend dev"],
    "confidence": 0.92,
    "risk": "low",
    "flags": []
  },
  "createdAt": 1736000000000
}
```

**Poll Post:**
```json
{
  "_key": "p3",
  "authorId": "users/u5",
  "postType": "poll",
  "text": "Which frontend framework do you prefer?",
  "pollOptions": ["React", "Vue", "Angular", "Svelte"],
  "pollVotes": {
    "React": 15,
    "Vue": 8,
    "Angular": 3,
    "Svelte": 6
  },
  "category": "tech",
  "intent": "seeking-opinion",
  "depth": "casual",
  "tags": ["frontend", "frameworks"],
  "createdAt": 1736000000000
}
```

| Field | Type | Description |
|-------|------|-------------|
| `_key` | string | Unique post ID |
| `authorId` | string | Reference to users collection |
| `postType` | enum | `text` or `poll` |
| `text` | string | Question content |
| `pollOptions` | string[] | Poll choices (polls only) |
| `pollVotes` | map | Vote counts per option |
| `category` | enum | AI-classified category |
| `intent` | enum | AI-classified intent |
| `depth` | enum | AI-classified depth |
| `tags` | string[] | Normalized tag keys |
| `aiRaw` | object | Raw AI classification data |
| `createdAt` | int64 | Unix timestamp (ms) |

---

### `tags`

Canonical tags with aliases for normalization.

```json
{
  "_key": "career-change",
  "label": "Career Change",
  "aliases": ["changing careers", "career switch", "job change"],
  "usageCount": 47,
  "createdAt": 1736000000000
}
```

| Field | Type | Description |
|-------|------|-------------|
| `_key` | string | Normalized tag ID (slug) |
| `label` | string | Display name |
| `aliases` | string[] | Alternative names |
| `usageCount` | int | Times used in posts |
| `createdAt` | int64 | First usage timestamp |

---

### `chats`

Container for conversation threads.

```json
{
  "_key": "c1",
  "postId": "posts/p1",
  "type": "direct",
  "createdAt": 1736000000000
}
```

| Field | Type | Description |
|-------|------|-------------|
| `_key` | string | Unique chat ID |
| `postId` | string | Reference to originating post |
| `type` | enum | `direct` (1:1) or `group` |
| `createdAt` | int64 | Unix timestamp (ms) |

---

### `messages`

Individual messages within chats.

```json
{
  "_key": "m1",
  "chatId": "chats/c1",
  "senderId": "users/u3",
  "text": "I made the switch last year! Start with Go or Node.js.",
  "status": "seen",
  "createdAt": 1736000000000
}
```

| Field | Type | Description |
|-------|------|-------------|
| `_key` | string | Unique message ID |
| `chatId` | string | Reference to chat |
| `senderId` | string | Reference to sender |
| `text` | string | Message content |
| `status` | enum | Message delivery status |
| `createdAt` | int64 | Unix timestamp (ms) |

**Message Status Flow:**
```
sending → sent → delivered → seen
              ↘ failed
```

---

## Edge Collections

Edges connect documents and enable graph traversals.

### `created`

Links users to posts they authored.

```
users/u1 ──[created]──▶ posts/p1
```

```json
{
  "_from": "users/u1",
  "_to": "posts/p1",
  "createdAt": 1736000000000
}
```

**Use case:** Get all posts by a user.

---

### `responded`

Links users to posts they responded to.

```
users/u3 ──[responded]──▶ posts/p1
```

```json
{
  "_from": "users/u3",
  "_to": "posts/p1",
  "chatId": "chats/c1",
  "createdAt": 1736000000000
}
```

**Use case:** Track who has responded, prevent duplicate responses.

---

### `post_has_tag`

Links posts to their tags.

```
posts/p1 ──[post_has_tag]──▶ tags/career-change
posts/p1 ──[post_has_tag]──▶ tags/backend-dev
```

```json
{
  "_from": "posts/p1",
  "_to": "tags/career-change"
}
```

**Use case:** Find posts by tag, get related tags.

---

### `follows`

Links users who follow each other.

```
users/u1 ──[follows]──▶ users/u3
```

```json
{
  "_from": "users/u1",
  "_to": "users/u3",
  "createdAt": 1736000000000
}
```

**Use case:** Build personalized feed, show content from followed users.

---

### `participates_in`

Links users to chats they're part of.

```
users/u1 ──[participates_in]──▶ chats/c1
users/u3 ──[participates_in]──▶ chats/c1
```

```json
{
  "_from": "users/u1",
  "_to": "chats/c1",
  "role": "author",
  "status": "active",
  "joinedAt": 1736000000000
}
```

| Field | Type | Values |
|-------|------|--------|
| `role` | enum | `author`, `responder`, `invited` |
| `status` | enum | `active`, `pending`, `muted` |

**Use case:** Get user's chats, check permissions.

---

### `tagged`

Links posts to users mentioned/tagged.

```
posts/p4 ──[tagged]──▶ users/u1
```

```json
{
  "_from": "posts/p4",
  "_to": "users/u1",
  "taggedAt": 1736000000000
}
```

**Use case:** Notify users when mentioned, show tagged posts.

---

### `voted`

Links users to polls they voted on.

```
users/u1 ──[voted]──▶ posts/p3
```

```json
{
  "_from": "users/u1",
  "_to": "posts/p3",
  "option": "React",
  "votedAt": 1736000000000
}
```

**Use case:** Prevent duplicate votes, show user's vote.

---

## Graph Visualization

```
                    ┌─────────┐
                    │  users  │
                    └────┬────┘
                         │
          ┌──────────────┼──────────────┐
          │              │              │
     [created]      [follows]    [participates_in]
          │              │              │
          ▼              ▼              ▼
     ┌─────────┐    ┌─────────┐    ┌─────────┐
     │  posts  │    │  users  │    │  chats  │
     └────┬────┘    └─────────┘    └────┬────┘
          │                             │
    ┌─────┼─────┐                       │
    │     │     │                       │
[post_has_tag] [responded]         [messages]
    │     │     │                       │
    ▼     │     ▼                       ▼
┌──────┐  │  ┌─────────┐          ┌──────────┐
│ tags │  │  │  users  │          │ messages │
└──────┘  │  └─────────┘          └──────────┘
          │
       [voted]
          │
          ▼
     ┌─────────┐
     │  users  │
     └─────────┘
```

---

## Common Graph Queries

### Get user's posts with tags
```aql
FOR post IN 1..1 OUTBOUND @userId created
  LET tags = (
    FOR tag IN 1..1 OUTBOUND post post_has_tag
    RETURN tag.label
  )
  RETURN MERGE(post, { tagLabels: tags })
```

### Get feed from followed users
```aql
FOR followed IN 1..1 OUTBOUND @userId follows
  FOR post IN 1..1 OUTBOUND followed created
    SORT post.createdAt DESC
    LIMIT @limit
    RETURN post
```

### Get chat participants
```aql
FOR user IN 1..1 INBOUND @chatId participates_in
  RETURN user
```

### Find posts by tag
```aql
FOR post IN 1..1 INBOUND @tagId post_has_tag
  SORT post.createdAt DESC
  RETURN post
```
