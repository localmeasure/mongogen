// Code generated by MongoGen. DO NOT EDIT.
// Collection: users

package users

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type useGroupId struct {
	isSet                   bool
	groupId                 bson.M
	name                    bson.M
}

func UseGroupId() *useGroupId {
	return &useGroupId{}
}

func (use *useGroupId) Build() bson.D {
	filter := bson.D{primitive.E{Key: "group_id", Value: use.groupId}}
	if use.name != nil {
		filter = append(filter, primitive.E{Key: "name", Value: use.name})
	}
	return filter
}

func (use *useGroupId) GroupId(value primitive.ObjectID) *useGroupId {
	use.groupId = bson.M{"$eq": value}
	return use
}

func (use *useGroupId) GroupIdNe(value primitive.ObjectID) *useGroupId {
	use.groupId = bson.M{"$ne": value}
	return use
}

func (use *useGroupId) GroupIdIn(values []primitive.ObjectID) *useGroupId {
	use.groupId = bson.M{"$in": values}
	return use
}

func (use *useGroupId) GroupIdNin(values []primitive.ObjectID) *useGroupId {
	use.groupId = bson.M{"$nin": values}
	return use
}

func (use *useGroupId) GroupIdGt(value primitive.ObjectID) *useGroupId {
	if use.isSet {
		use.groupId["$gt"] = value
	} else {
		use.groupId = bson.M{"$gt": value}
		use.isSet = true
	}
	return use
}

func (use *useGroupId) GroupIdGte(value primitive.ObjectID) *useGroupId {
	if use.isSet {
		use.groupId["$gte"] = value
	} else {
		use.groupId = bson.M{"$gte": value}
		use.isSet = true
	}
	return use
}

func (use *useGroupId) GroupIdLt(value primitive.ObjectID) *useGroupId {
	if use.isSet {
		use.groupId["$lt"] = value
	} else {
		use.groupId = bson.M{"$lt": value}
		use.isSet = true
	}
	return use
}

func (use *useGroupId) GroupIdLte(value primitive.ObjectID) *useGroupId {
	if use.isSet {
		use.groupId["$lte"] = value
	} else {
		use.groupId = bson.M{"$lte": value}
		use.isSet = true
	}
	return use
}

func (use *useGroupId) Name(value string) *useGroupId {
	use.name = bson.M{"$eq": value}
	return use
}

func (use *useGroupId) NameNe(value string) *useGroupId {
	use.name = bson.M{"$ne": value}
	return use
}

func (use *useGroupId) NameIn(values []string) *useGroupId {
	use.name = bson.M{"$in": values}
	return use
}

func (use *useGroupId) NameNin(values []string) *useGroupId {
	use.name = bson.M{"$nin": values}
	return use
}

type useTeamId struct {
	isSet                   bool
	teamId                  bson.M
	lastSeen                bson.M
}

func UseTeamId() *useTeamId {
	return &useTeamId{}
}

func (use *useTeamId) Build() bson.D {
	filter := bson.D{primitive.E{Key: "team_id", Value: use.teamId}}
	if use.lastSeen != nil {
		filter = append(filter, primitive.E{Key: "last_seen", Value: use.lastSeen})
	}
	return filter
}

func (use *useTeamId) TeamId(value primitive.ObjectID) *useTeamId {
	use.teamId = bson.M{"$eq": value}
	return use
}

func (use *useTeamId) TeamIdNe(value primitive.ObjectID) *useTeamId {
	use.teamId = bson.M{"$ne": value}
	return use
}

func (use *useTeamId) TeamIdIn(values []primitive.ObjectID) *useTeamId {
	use.teamId = bson.M{"$in": values}
	return use
}

func (use *useTeamId) TeamIdNin(values []primitive.ObjectID) *useTeamId {
	use.teamId = bson.M{"$nin": values}
	return use
}

func (use *useTeamId) TeamIdGt(value primitive.ObjectID) *useTeamId {
	if use.isSet {
		use.teamId["$gt"] = value
	} else {
		use.teamId = bson.M{"$gt": value}
		use.isSet = true
	}
	return use
}

func (use *useTeamId) TeamIdGte(value primitive.ObjectID) *useTeamId {
	if use.isSet {
		use.teamId["$gte"] = value
	} else {
		use.teamId = bson.M{"$gte": value}
		use.isSet = true
	}
	return use
}

func (use *useTeamId) TeamIdLt(value primitive.ObjectID) *useTeamId {
	if use.isSet {
		use.teamId["$lt"] = value
	} else {
		use.teamId = bson.M{"$lt": value}
		use.isSet = true
	}
	return use
}

func (use *useTeamId) TeamIdLte(value primitive.ObjectID) *useTeamId {
	if use.isSet {
		use.teamId["$lte"] = value
	} else {
		use.teamId = bson.M{"$lte": value}
		use.isSet = true
	}
	return use
}

func (use *useTeamId) LastSeenGt(value time.Time) *useTeamId {
	if use.isSet {
		use.lastSeen["$gt"] = value
	} else {
		use.lastSeen = bson.M{"$gt": value}
		use.isSet = true
	}
	return use
}

func (use *useTeamId) LastSeenGte(value time.Time) *useTeamId {
	if use.isSet {
		use.lastSeen["$gte"] = value
	} else {
		use.lastSeen = bson.M{"$gte": value}
		use.isSet = true
	}
	return use
}

func (use *useTeamId) LastSeenLt(value time.Time) *useTeamId {
	if use.isSet {
		use.lastSeen["$lt"] = value
	} else {
		use.lastSeen = bson.M{"$lt": value}
		use.isSet = true
	}
	return use
}

func (use *useTeamId) LastSeenLte(value time.Time) *useTeamId {
	if use.isSet {
		use.lastSeen["$lte"] = value
	} else {
		use.lastSeen = bson.M{"$lte": value}
		use.isSet = true
	}
	return use
}

