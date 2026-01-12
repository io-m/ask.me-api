# ask.me API Documentation

## Overview

ask.me is a question-driven social API where users post questions, receive responses, and engage in private conversations. The backend uses:

- **Go 1.25** with stdlib routing
- **ArangoDB** for graph-based data storage
- **Feature-first architecture** with repository pattern and DI

## Base URL

```
http://localhost:8080
```

## Response Format

All responses follow this structure:

```json
{
  "success": true,
  "data": { ... }
}
```

Error responses:

```json
{
  "success": false,
  "error": "error message"
}
```

---

## Users

### GET /users/{userId}

Get user profile by ID.

**Response:**

```json
{
  "success": true,
  "data": {
    "_key": "u1",
    "username": "alex_dev",
    "createdAt": 1736000000000,
    "interests": ["tech", "career"],
    "blockedTopics": [],
    "settings": {
      "allowDMs": true,
      "allowTagging": true
    },
    "stats": {
      "postsCreated": 5,
      "responsesGiven": 12
    }
  }
}
```

### POST /users

Create a new user.

**Request:**

```json
{
  "username": "new_user",
  "interests": ["tech", "finance"],
  "settings": {
    "allowDMs": true,
    "allowTagging": true
  }
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "_key": "u123",
    "createdAt": 1736000000000
  }
}
```

### POST /users/{userId}/follow

Follow another user.

**Request:**

```json
{
  "followUserId": "u2"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "success": true,
    "followId": "edge_key"
  }
}
```

---

## Posts

### GET /posts/{postId}

Get post details with AI classification data.

**Response (Text Post):**

```json
{
  "success": true,
  "data": {
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
}
```

**Response (Poll Post):**

```json
{
  "success": true,
  "data": {
    "_key": "p3",
    "authorId": "users/u5",
    "postType": "poll",
    "text": "Which frontend framework do you prefer?",
    "pollOptions": ["React", "Vue", "Angular", "Svelte"],
    "category": "tech",
    "intent": "seeking-opinion",
    "depth": "casual",
    "tags": ["frontend", "frameworks"],
    "createdAt": 1736000000000
  }
}
```

### POST /posts

Create a text post. The backend automatically classifies the content using AI.

**Request:**

```json
{
  "authorId": "u1",
  "postType": "text",
  "text": "What's the best way to learn system design?"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "_key": "p123",
    "category": "education",
    "tags": ["system-design", "learning"],
    "createdAt": 1736000000000
  }
}
```

### POST /posts/poll

Create a poll post.

**Request:**

```json
{
  "authorId": "u2",
  "postType": "poll",
  "text": "What's your preferred code editor?",
  "pollOptions": ["VS Code", "Neovim", "JetBrains", "Sublime"]
}
```

**Note:** The backend automatically classifies the poll using AI (category, intent, depth, tags).

### POST /posts/{postId}/respond

Respond to a post (starts a chat).

**Request:**

```json
{
  "userId": "u2",
  "text": "I would recommend starting with building side projects!",
  "chatType": "direct"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "chatId": "c123",
    "messageId": "m001",
    "createdAt": 1736000000000
  }
}
```

### POST /posts/{postId}/vote

Vote on a poll.

**Request:**

```json
{
  "userId": "u1",
  "option": "React"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "postId": "p3",
    "option": "React",
    "votes": {
      "React": 3,
      "Vue": 1,
      "Angular": 0,
      "Svelte": 1
    }
  }
}
```

---

## Chats

### GET /users/{userId}/chats

Get all chat threads for a user.

**Query Parameters:**

- `limit` (optional): Max threads to return (default: 50)
- `cursor` (optional): Pagination cursor

**Response (Direct Chat):**

```json
{
  "success": true,
  "data": {
    "threads": [
      {
        "id": "c1",
        "type": "direct",
        "question": {
          "id": "p1",
          "text": "How do I switch from frontend to backend?",
          "authorId": "users/u1",
          "createdAt": 1736000000000,
          "formattedTime": "7 days ago"
        },
        "partner": {
          "id": "u3",
          "username": "john_doe",
          "avatarUrl": null
        },
        "lastMessage": {
          "id": "m3",
          "text": "Mostly building projects. I'd recommend starting with a REST API.",
          "senderId": "users/u3",
          "createdAt": 1736000000000,
          "formattedTime": "5 days ago"
        },
        "unreadCount": 1,
        "hasUnread": true
      }
    ],
    "nextCursor": null
  }
}
```

**Response (Group Chat):**

For group chats, the response includes a `participants` array with all chat members:

```json
{
  "success": true,
  "data": {
    "threads": [
      {
        "id": "c5",
        "type": "group",
        "question": {
          "id": "p10",
          "text": "Tech folks! @alex @nina - What stack are you using?",
          "authorId": "users/u7",
          "createdAt": 1736000000000,
          "formattedTime": "2 days ago"
        },
        "partner": {
          "id": "u7",
          "username": "david_tech",
          "avatarUrl": null
        },
        "participants": [
          {
            "id": "u7",
            "username": "david_tech",
            "avatarUrl": null,
            "role": "author",
            "status": "active"
          },
          {
            "id": "u1",
            "username": "alex_dev",
            "avatarUrl": null,
            "role": "invited",
            "status": "active"
          },
          {
            "id": "u8",
            "username": "nina_design",
            "avatarUrl": null,
            "role": "invited",
            "status": "pending"
          }
        ],
        "lastMessage": {
          "id": "m20",
          "text": "I switched to Bun + Hono for backend!",
          "senderId": "users/u1",
          "createdAt": 1736000000000,
          "formattedTime": "45 min ago"
        },
        "unreadCount": 8,
        "hasUnread": true
      }
    ],
    "nextCursor": null
  }
}
```

**Notes:**

- The `type` field indicates whether the chat is `direct` (1:1) or `group` (3+ participants)
- The `partner` field always contains the primary partner (question author for answered questions, or first responder for your questions)
- The `participants` array is only included for group chats and contains all members
- Each participant has `role` (`author`, `responder`, `invited`) and `status` (`active`, `pending`, `muted`)

### GET /chats/{chatId}

Get chat with all messages.

**Response:**

```json
{
  "success": true,
  "data": {
    "_key": "c1",
    "postId": "posts/p1",
    "type": "direct",
    "messages": [
      {
        "_key": "m1",
        "senderId": "users/u3",
        "text": "I made the switch last year! Start with Go or Node.js.",
        "status": "seen",
        "createdAt": 1736000000000
      },
      {
        "_key": "m2",
        "senderId": "users/u1",
        "text": "Thanks! Did you take any courses or just learn by building?",
        "status": "seen",
        "createdAt": 1736000000000
      }
    ],
    "createdAt": 1736000000000
  }
}
```

### POST /chats/{chatId}/message

Send a message in a chat.

**Request:**

```json
{
  "senderId": "u1",
  "text": "That's helpful! I'll try that approach."
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "messageId": "m123",
    "createdAt": 1736000000000
  }
}
```

### GET /chats/{chatId}/participants

Get chat participants.

**Response:**

```json
{
  "success": true,
  "data": {
    "chatId": "c1",
    "type": "direct",
    "participants": [
      {
        "id": "u1",
        "username": "alex_dev",
        "avatarUrl": null,
        "role": "author",
        "status": "active"
      },
      {
        "id": "u3",
        "username": "john_doe",
        "avatarUrl": null,
        "role": "responder",
        "status": "active"
      }
    ]
  }
}
```

### POST /chats/{chatId}/accept

Accept a group chat invite.

**Request:**

```json
{
  "userId": "u3"
}
```

### POST /chats/{chatId}/mute

Mute chat notifications.

**Request:**

```json
{
  "userId": "u1"
}
```

---

## Feed

### GET /feed

Get personalized feed for a user.

**Query Parameters:**

- `userId` (required): User ID
- `limit` (optional): Max items (default: 20, max: 50)
- `cursor` (optional): Pagination cursor
- `category` (optional): Filter by category
- `depth` (optional): Filter by depth (casual, neutral, serious)

**Response (Text Post):**

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "p1",
        "postType": "text",
        "text": "How do I switch from frontend to backend engineering?",
        "category": "career",
        "intent": "seeking-advice",
        "depth": "serious",
        "tags": ["career-change", "backend-dev"],
        "author": {
          "id": "u1",
          "username": "alex_dev",
          "avatarUrl": null
        },
        "chatId": "c1",
        "lastMessage": {
          "id": "m3",
          "text": "Mostly building projects...",
          "senderId": "users/u3",
          "status": "delivered",
          "createdAt": 1736000000000,
          "formattedTime": "5 days ago"
        },
        "createdAt": 1736000000000
      }
    ],
    "nextCursor": null
  }
}
```

**Response (Poll Post):**

For poll posts, the response includes a `pollOptions` array:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "p11",
        "postType": "poll",
        "text": "Which frontend framework do you prefer for new projects?",
        "pollOptions": ["React", "Vue", "Angular", "Svelte"],
        "category": "tech",
        "intent": "seeking-opinion",
        "depth": "casual",
        "tags": ["frontend", "frameworks"],
        "author": {
          "id": "u5",
          "username": "devmaster",
          "avatarUrl": "https://example.com/avatar.jpg"
        },
        "createdAt": 1736000000000
      }
    ],
    "nextCursor": null
  }
}
```

**Notes:**

- The `pollOptions` field is only present for posts with `postType: "poll"`
- The `chatId` and `lastMessage` fields are only present if the user has an existing conversation on the post
- The `formattedTime` field is computed server-side (e.g., "5 days ago", "1 week ago")

---

## Tags

### GET /tags/{tagId}

Get tag by ID.

**Response:**

```json
{
  "success": true,
  "data": {
    "_key": "career-change",
    "label": "Career Change",
    "aliases": ["changing careers", "career switch"],
    "usageCount": 15,
    "createdAt": 1736000000000
  }
}
```

### GET /tags

List or search tags.

**Query Parameters:**

- `limit` (optional): Max tags (default: 50)
- `offset` (optional): Pagination offset
- `q` (optional): Search query

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "_key": "remote-work",
      "label": "Remote Work",
      "aliases": ["wfh", "work from home"],
      "usageCount": 25,
      "createdAt": 1736000000000
    },
    {
      "_key": "learning",
      "label": "Learning",
      "aliases": ["education", "self-improvement"],
      "usageCount": 22,
      "createdAt": 1736000000000
    }
  ]
}
```

---

## Enums

### Post Categories

- `career`, `relationships`, `tech`, `health`, `finance`, `fun`, `opinion`, `lifestyle`, `education`, `other`

### Post Depth

- `casual`, `neutral`, `serious`

### Post Types

- `text`, `poll`

### Chat Types

- `direct`, `group`

### Message Status

- `sending`, `sent`, `delivered`, `seen`, `failed`

### Participant Role

- `author`, `responder`, `invited`

### Participant Status

- `active`, `pending`, `muted`

---

## Error Codes

| Status | Error |
|--------|-------|
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists |
| 403 | Forbidden - Not allowed |
| 500 | Internal Server Error |
