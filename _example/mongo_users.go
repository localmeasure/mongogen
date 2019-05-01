// Code generated by MongoGen. DO NOT EDIT.
// Collection: users

package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserFilter struct {
	Filter bson.D
}

type Users struct {
	db *mongo.Database
}

func NewUsers(db *mongo.Database) *Users {
	return &Users{db}
}

func UserWithID(id primitive.ObjectID) UserFilter {
	return UserFilter{bson.D{{Key: "_id", Value: id}}}
}

func UserWithIDs(ids []primitive.ObjectID) UserFilter {
	return UserFilter{bson.D{{Key: "_id", Value: bson.M{"$in": ids}}}}
}

func UserWithGroupId(groupId primitive.ObjectID) UserFilter {
	return UserFilter{bson.D{{Key: "group_id", Value: groupId}}}
}

func UserWithGroupIdName(groupId primitive.ObjectID, name string) UserFilter {
	return UserFilter{bson.D{{Key: "group_id", Value: groupId}, {Key: "name", Value: name}}}
}

func UserWithTeamId(teamId primitive.ObjectID) UserFilter {
	return UserFilter{bson.D{{Key: "team_id", Value: teamId}}}
}

func UserWithTeamIdLastSeen(teamId primitive.ObjectID, lastSeen time.Time) UserFilter {
	return UserFilter{bson.D{{Key: "team_id", Value: teamId}, {Key: "last_seen", Value: lastSeen}}}
}

func (s *Users) Find(ctx context.Context, filter UserFilter, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return s.db.Collection("users").Find(ctx, filter.Filter, opts...)
}

func (s *Users) FindWithIDs(ctx context.Context, ids []primitive.ObjectID, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return s.db.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": ids}}, opts...)
}

func (s *Users) FindOne(ctx context.Context, filter UserFilter, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return s.db.Collection("users").FindOne(ctx, filter.Filter, opts...)
}

func (s *Users) FindOneWithID(ctx context.Context, id primitive.ObjectID, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return s.db.Collection("users").FindOne(ctx, bson.M{"_id": id}, opts...)
}

func (s *Users) Count(ctx context.Context, filter UserFilter, opts ...*options.CountOptions) (int64, error) {
	return s.db.Collection("users").CountDocuments(ctx, filter.Filter, opts...)
}

