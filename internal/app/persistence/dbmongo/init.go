package dbmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/thirdparty/mongo"
	mongolib "go.mongodb.org/mongo-driver/mongo"
)

type (
	CarRepo struct {
		client mongoClient
	}

	CompanyRepo struct {
		client mongoClient
	}

	SessionRepo struct {
		client mongoClient
	}

	UserRepo struct {
		client mongoClient
	}

	mongoClient interface {
		Close() error
		Search(ctx context.Context, db, collection string, query, result interface{}) error
		SearchOne(ctx context.Context, db, collection string, query, result interface{}) error
		SearchWithPageable(ctx context.Context, db, collection string, page mongo.Pageable, query, result interface{}) error
		SearchWithProjection(ctx context.Context, db, collection string, query, projection, result interface{}) error
		Insert(ctx context.Context, db, collection string, data interface{}) error
		Update(ctx context.Context, db, collection string, query, data interface{}) error
		UpsertWithOperationResult(ctx context.Context, db, collection string, query, into interface{}) (int64, int64, error)
		FindOneAndUpdate(ctx context.Context, db, collection string, query, update, projection, result interface{}) error
		Count(ctx context.Context, db, collection string, query interface{}) (int64, error)
		CreateUniqueIndex(ctx context.Context, db, collection, field string) error
		GetCollection(db, collection string) *mongolib.Collection
		UpdateAll(ctx context.Context, db, collection string, query, into interface{}) error
	}
)

var dbName = "car_ai_core"

const (
	collUserProfile    = "user_profile"
	collSession        = "session"
	collCompanyProfile = "company_profile"
	collCarProfile     = "car_profile"
	collCarApproveLog  = "car_approve_log"
)

const (
	fieldEmail     = "email"
	fieldUniqueTag = "uniqueTag"
	fieldDeletedDt = "deletedDt"
)

func NewSessionRepo(databaseName string, client mongoClient) persistence.SessionDB {
	if databaseName != "" {
		dbName = databaseName
	}
	return &SessionRepo{client: client}
}

func NewUserRepo(databaseName string, client mongoClient) persistence.UserDB {
	if databaseName != "" {
		dbName = databaseName
	}

	collection := client.GetCollection(dbName, collUserProfile)
	indexModel := mongolib.IndexModel{
		Keys:    bson.D{{Key: fieldEmail, Value: 1}, {Key: fieldDeletedDt, Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Panicf("failed to create unique index: %v", err.Error())
	}
	return &UserRepo{client: client}
}
