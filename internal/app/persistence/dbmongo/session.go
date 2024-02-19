package dbmongo

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"api-hotel-booking/internal/app/persistence"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *SessionRepo) Insert(ctx context.Context, userId, key string, timeout time.Time) (string, error) {
	sessionId := strings.Join([]string{
		"session",
		userId,
		fmt.Sprintf("%020d", time.Now().UnixNano()),
		fmt.Sprintf("%04d", rand.Intn(9999)),
	}, "_")

	if err := s.client.Insert(ctx, dbName, collSession, persistence.Session{
		SessionId: sessionId,
		UserId:    userId,
		Key:       key,
		Timeout:   timeout,
	}); err != nil {
		return "", castKnownError(err)
	}

	return sessionId, nil
}

func (s *SessionRepo) Get(ctx context.Context, sessionId string) (persistence.Session, error) {
	session := persistence.Session{}
	query := bson.D{{"sessionId", sessionId}}

	err := s.client.SearchOne(ctx, dbName, collSession, query, &session)
	if err != nil {
		return session, castKnownError(err)
	}
	return session, nil
}

func (s *SessionRepo) UpdateTimeout(ctx context.Context, sessionId string, timeout time.Time) error {
	query := bson.D{{"sessionId", sessionId}}
	update := bson.D{
		{"$set", bson.D{
			{"timeout", timeout},
		}},
	}
	err := s.client.Update(ctx, dbName, collSession, query, update)
	if err != nil {
		return castKnownError(err)
	}
	return nil
}

func (s *SessionRepo) Delete(sessionId string) error {
	query := bson.D{{"sessionId", sessionId}}

	if result, err := s.client.GetCollection(dbName, collSession).DeleteOne(context.TODO(), query); err != nil {
		return castKnownError(err)
	} else if result.DeletedCount == 0 {
		return persistence.NotFoundError
	}
	return nil
}
