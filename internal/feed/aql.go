package feed

// AQL queries for feed operations
const (
	// GetRecommendedPosts retrieves personalized posts for a user's feed
	GetRecommendedPosts = `
		// Get user's interaction history for personalization
		LET userTags = (
			FOR edge IN responded
			FILTER edge._from == @userId
			FOR postEdge IN post_has_tag
			FILTER postEdge._from == edge._to
			RETURN DISTINCT postEdge._to
		)
		
		LET userCategories = (
			FOR edge IN responded
			FILTER edge._from == @userId
			FOR post IN posts
			FILTER post._id == edge._to
			COLLECT category = post.category WITH COUNT INTO cnt
			SORT cnt DESC
			LIMIT 5
			RETURN category
		)
		
		// Get recommended posts
		FOR p IN posts
			// Filter by category if specified
			FILTER @category == '' OR p.category == @category
			FILTER @depth == '' OR p.depth == @depth
			
			// Get author
			LET author = FIRST(
				FOR edge IN created
				FILTER edge._to == p._id
				FOR user IN users
				FILTER user._id == edge._from
				RETURN user
			)
			
			// Check if user has a chat for this post
			LET userChat = FIRST(
				FOR c IN chats
				FILTER c.postId == p._id
				FOR edge IN participates_in
				FILTER edge._from == @userId AND edge._to == c._id
				RETURN c
			)
			
			// Get last message if chat exists
			LET lastMsg = userChat ? FIRST(
				FOR m IN messages
				FILTER m.chatId == userChat._id
				SORT m.createdAt DESC
				LIMIT 1
				RETURN m
			) : null
			
			// Get user's reaction to last message (if any)
			LET myReaction = lastMsg ? FIRST(
				FOR e IN reacted
				FILTER e._from == @userId AND e._to == lastMsg._id
				RETURN e.emoji
			) : null
			
			// Count unread messages (messages not from user and not seen)
			LET unreadCount = userChat ? LENGTH(
				FOR m IN messages
				FILTER m.chatId == userChat._id
				   AND m.senderId != @userId
				   AND m.status != 'seen'
				RETURN 1
			) : 0
			
			// Get tags
			LET postTags = (
				FOR edge IN post_has_tag
				FILTER edge._from == p._id
				FOR tag IN tags
				FILTER tag._id == edge._to
				RETURN tag._key
			)
			
			// Calculate relevance score
			LET tagMatch = LENGTH(INTERSECTION(postTags, userTags))
			LET categoryMatch = p.category IN userCategories ? 1 : 0
			LET recency = (DATE_NOW() - p.createdAt) / (1000 * 60 * 60 * 24)
			
			LET score = (categoryMatch * 40) + (tagMatch * 20) + (100 - MIN([recency, 100]) * 0.1)
			
			SORT score DESC, p.createdAt DESC
			LIMIT @limit
			
			RETURN {
				id: p._key,
				postType: p.postType,
				text: p.text,
				pollOptions: p.pollOptions,
				category: p.category,
				intent: p.intent,
				depth: p.depth,
				tags: postTags,
				author: {
					id: author._key,
					username: author.username,
					avatarUrl: author.avatarUrl
				},
				chatId: userChat ? userChat._key : null,
				lastMessage: lastMsg ? {
					id: lastMsg._key,
					text: lastMsg.text,
					senderId: LAST(SPLIT(lastMsg.senderId, "/")),
					status: lastMsg.status,
					createdAt: lastMsg.createdAt,
					myReaction: myReaction
				} : null,
				unreadCount: unreadCount,
				createdAt: p.createdAt
			}
	`

	// GetUserInteractionTags retrieves tags from posts user has interacted with
	GetUserInteractionTags = `
		FOR edge IN responded
		FILTER edge._from == @userId
		FOR postEdge IN post_has_tag
		FILTER postEdge._from == edge._to
		FOR tag IN tags
		FILTER tag._id == postEdge._to
		COLLECT tagKey = tag._key WITH COUNT INTO cnt
		SORT cnt DESC
		LIMIT 20
		RETURN tagKey
	`

	// GetUserCategories retrieves categories user frequently engages with
	GetUserCategories = `
		FOR edge IN responded
		FILTER edge._from == @userId
		FOR post IN posts
		FILTER post._id == edge._to
		COLLECT category = post.category WITH COUNT INTO cnt
		SORT cnt DESC
		LIMIT 5
		RETURN category
	`

	// GetUserIntents retrieves intents user frequently responds to
	GetUserIntents = `
		FOR edge IN responded
		FILTER edge._from == @userId
		FOR post IN posts
		FILTER post._id == edge._to
		COLLECT intent = post.intent WITH COUNT INTO cnt
		SORT cnt DESC
		LIMIT 10
		RETURN intent
	`
)
