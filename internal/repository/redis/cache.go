package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/derticom/doc-store/internal/domain/document"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const timeout = 10 * time.Second

type Cache struct {
	client *redis.Client
}

func New(ctx context.Context, addr string) (*Cache, error) {
	opts, err := redis.ParseURL(addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to redis.ParseURL")
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp := client.Ping(ctx)
	err = resp.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to client.Ping")
	}

	return &Cache{client: client}, nil
}

func (c *Cache) Get(id string) (*document.Document, bool) {
	data, err := c.client.Get(context.Background(), "doc:"+id).Bytes()
	if err != nil {
		return nil, false
	}
	var wrapper struct {
		Doc *document.Document
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, false
	}
	return wrapper.Doc, true
}

func (c *Cache) Set(id string, doc *document.Document) {
	data, _ := json.Marshal(struct {
		Doc *document.Document
	}{doc})
	c.client.Set(context.Background(), "doc:"+id, data, time.Minute*10)
}

func (c *Cache) Invalidate(id string) {
	c.client.Del(context.Background(), "doc:"+id)
}
