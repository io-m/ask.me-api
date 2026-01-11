package chat

// AQL queries for chat operations
const (
	// GetChatByID retrieves a chat by its key
	GetChatByID = `
		FOR c IN chats
		FILTER c._key == @key
		RETURN c
	`

	// GetChatMessages retrieves messages for a chat
	GetChatMessages = `
		FOR m IN messages
		FILTER m.chatId == @chatId
		SORT m.createdAt ASC
		LIMIT @offset, @limit
		RETURN m
	`

	// UpdateMessageStatus updates a message's delivery status
	UpdateMessageStatus = `
		FOR m IN messages
		FILTER m._key == @key
		UPDATE m WITH { status: @status } IN messages
	`

	// GetChatUnreadCount counts unread messages for a user in a chat
	GetChatUnreadCount = `
		RETURN LENGTH(
			FOR m IN messages
			FILTER m.chatId == @chatId
			   AND m.senderId != @userId
			   AND m.status != 'seen'
			RETURN 1
		)
	`

	// UpdateParticipation updates a user's participation status
	UpdateParticipation = `
		FOR e IN participates_in
		FILTER e._from == @from AND e._to == @to
		UPDATE e WITH { 
			status: @status, 
			notificationsEnabled: @notificationsEnabled,
			joinedAt: @joinedAt
		} IN participates_in
	`

	// GetParticipation retrieves a user's participation in a chat
	GetParticipation = `
		FOR e IN participates_in
		FILTER e._from == @from AND e._to == @to
		RETURN e
	`

	// GetChatParticipants retrieves all participants of a chat
	GetChatParticipants = `
		FOR edge IN participates_in
		FILTER edge._to == @chatId
		FOR user IN users
		FILTER user._id == edge._from
		RETURN {
			id: user._key,
			username: user.username,
			avatarUrl: user.avatarUrl,
			role: edge.role,
			status: edge.status
		}
	`

	// GetUserChatThreads retrieves all chat threads for a user
	GetUserChatThreads = `
		FOR edge IN participates_in
		FILTER edge._from == @userId
		
		LET chat = DOCUMENT(edge._to)
		LET post = DOCUMENT(chat.postId)
		
		// Get partner (other participant)
		LET partner = FIRST(
			FOR otherEdge IN participates_in
			FILTER otherEdge._to == chat._id
			   AND otherEdge._from != @userId
			FOR user IN users
			FILTER user._id == otherEdge._from
			RETURN user
		)
		
		// Get last message
		LET lastMsg = FIRST(
			FOR m IN messages
			FILTER m.chatId == chat._id
			SORT m.createdAt DESC
			LIMIT 1
			RETURN m
		)
		
		// Count unread
		LET unreadCount = LENGTH(
			FOR m IN messages
			FILTER m.chatId == chat._id
			   AND m.senderId != @userId
			   AND m.status != 'seen'
			RETURN 1
		)
		
		FILTER lastMsg != null
		SORT lastMsg.createdAt DESC
		LIMIT @limit
		
		RETURN {
			id: chat._key,
			question: {
				id: post._key,
				text: post.text,
				authorId: post.authorId,
				createdAt: post.createdAt
			},
			partner: {
				id: partner._key,
				username: partner.username,
				avatarUrl: partner.avatarUrl
			},
			lastMessage: {
				id: lastMsg._key,
				text: lastMsg.text,
				senderId: lastMsg.senderId,
				createdAt: lastMsg.createdAt
			},
			unreadCount: unreadCount,
			hasUnread: unreadCount > 0
		}
	`

	// GetChatForPostAndUser finds an existing chat for a post and user
	GetChatForPostAndUser = `
		FOR c IN chats
		FILTER c.postId == @postId
		FOR edge IN participates_in
		FILTER edge._to == c._id
		   AND edge._from == @userId
		RETURN c
	`
)
