package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	*mongo.Collection
}

type IndexFilter interface {
	Build() bson.D
}

func (c *Collection) Find(ctx context.Context, filter IndexFilter, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return c.Collection.Find(ctx, filter.Build(), opts...)
}
