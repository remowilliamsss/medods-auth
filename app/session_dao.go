package app

import (
	"context"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionDAO struct {
	collection *mongo.Collection
}

func NewSessionDAO(client *mongo.Client) *SessionDAO {
	return &SessionDAO{
		collection: client.Database("medods-auth-db").Collection("session"),
	}
}

func (dao *SessionDAO) FindById(ctx context.Context, uuid *uuid.UUID) (*Session, error) {
	filter := bson.D{{"_id", uuid}}
	var session Session
	err := dao.collection.FindOne(ctx, filter).Decode(&session)
	switch err {
	case nil:
		return &session, nil
	case mongo.ErrNoDocuments:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (dao *SessionDAO) Update(ctx context.Context, session *Session) error {
	filter := bson.D{{"_id", session.UUID}}
	updateResult, err := dao.collection.ReplaceOne(ctx, filter, session)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
