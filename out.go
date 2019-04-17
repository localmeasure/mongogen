package mongogen

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Users struct {
	db mongo.Database
}

func UserWithAccessToken(accessToken string) bson.D {
	return bson.D{{Key: "access_token", Value: accessToken}}
}

func UserWithEmail(email string) bson.D {
	return bson.D{{Key: "email", Value: email}}
}

func UserWithAvailableMerchantIDs(merchantIDs []string) bson.D {
	return bson.D{{Key: "available_merchant_ids", Value: bson.D{{Key: "$in", Value: merchantIDs}}}}
}

// Find returns multiple documents
// Usage:
//		ctx := context.Background()
//		s.Find(ctx, UserWithAccessToken("token"))
//		s.Find(ctx, UserWithAccessToken("token"), options.Find().SetLimit(10).SetSkip(0))
func (s *Users) Find(ctx context.Context, filter bson.D, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return s.db.Collection("users").Find(ctx, filter, opts...)
}

// FindOne returns up to one document
// Usage:
//		ctx := context.Background()
//      s.FindOne(ctx, UserWithEmail("email"))
func (s *Users) FindOne(ctx context.Context, filter bson.D, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return s.db.Collection("users").FindOne(ctx, filter, opts...)
}
