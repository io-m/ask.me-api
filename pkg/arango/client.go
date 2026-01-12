package arango

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"github.com/arangodb/go-driver/v2/connection"

	"github.com/askme/api/internal/config"
)

// Client wraps the ArangoDB connection
type Client struct {
	db arangodb.Database
}

// NewClient creates a new ArangoDB client
func NewClient(cfg config.ArangoDBConfig) (*Client, error) {
	endpoint := connection.NewRoundRobinEndpoints([]string{cfg.Endpoint})
	auth := connection.NewBasicAuth(cfg.Username, cfg.Password)

	conn := connection.NewHttpConnection(connection.HttpConfiguration{
		Endpoint:       endpoint,
		Authentication: auth,
	})

	client := arangodb.NewClient(conn)

	db, err := client.Database(context.Background(), cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("get database: %w", err)
	}

	return &Client{db: db}, nil
}

// Database returns the underlying database
func (c *Client) Database() arangodb.Database {
	return c.db
}

// Collection represents an ArangoDB collection
type Collection string

const (
	CollectionUsers    Collection = "users"
	CollectionPosts    Collection = "posts"
	CollectionTags     Collection = "tags"
	CollectionChats    Collection = "chats"
	CollectionMessages Collection = "messages"
)

// Edge collection names
const (
	EdgeCreated        Collection = "created"
	EdgeResponded      Collection = "responded"
	EdgePostHasTag     Collection = "post_has_tag"
	EdgeFollows        Collection = "follows"
	EdgeParticipatesIn Collection = "participates_in"
	EdgeTagged         Collection = "tagged"
	EdgeVoted          Collection = "voted"
)

// Query executes an AQL query and returns results
func Query[T any](ctx context.Context, client *Client, query string, bindVars map[string]any) ([]T, error) {
	cursor, err := client.db.Query(ctx, query, &arangodb.QueryOptions{
		BindVars: bindVars,
	})
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer cursor.Close()

	var results []T
	for cursor.HasMore() {
		var raw json.RawMessage
		_, err := cursor.ReadDocument(ctx, &raw)
		if err != nil {
			return nil, fmt.Errorf("read document failed: %w", err)
		}

		var doc T
		if err := json.Unmarshal(raw, &doc); err != nil {
			return nil, fmt.Errorf("unmarshal document failed: %w", err)
		}
		results = append(results, doc)
	}

	return results, nil
}

// QueryOne executes a query expecting a single result
func QueryOne[T any](ctx context.Context, client *Client, query string, bindVars map[string]any) (*T, error) {
	results, err := Query[T](ctx, client, query, bindVars)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return &results[0], nil
}

// InsertDocument inserts a document into a collection
func InsertDocument[T any](ctx context.Context, client *Client, collection Collection, doc T) (string, error) {
	col, err := client.db.Collection(ctx, string(collection))
	if err != nil {
		return "", fmt.Errorf("get collection: %w", err)
	}

	meta, err := col.CreateDocument(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("create document: %w", err)
	}

	return meta.Key, nil
}

// UpdateDocument updates a document in a collection
func UpdateDocument[T any](ctx context.Context, client *Client, collection Collection, key string, doc T) error {
	col, err := client.db.Collection(ctx, string(collection))
	if err != nil {
		return fmt.Errorf("get collection: %w", err)
	}

	_, err = col.UpdateDocument(ctx, key, doc)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}

	return nil
}

// GetDocument retrieves a document by key
func GetDocument[T any](ctx context.Context, client *Client, collection Collection, key string) (*T, error) {
	col, err := client.db.Collection(ctx, string(collection))
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}

	var doc T
	_, err = col.ReadDocument(ctx, key, &doc)
	if err != nil {
		if shared.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read document: %w", err)
	}

	return &doc, nil
}

// DeleteDocument removes a document by key
func DeleteDocument(ctx context.Context, client *Client, collection Collection, key string) error {
	col, err := client.db.Collection(ctx, string(collection))
	if err != nil {
		return fmt.Errorf("get collection: %w", err)
	}

	_, err = col.DeleteDocument(ctx, key)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	return nil
}
