package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoDocuments = mongo.ErrNoDocuments
	ErrUninitiated = errors.New("MongoDB's connection has not been initialized")
)

type IClient interface {
	GetCollection(db, collection string) *mongo.Collection
	Search(ctx context.Context, db, collection string, query, result interface{}) error
	SearchOne(ctx context.Context, db, collection string, query, result interface{}) error
	SearchWithProjection(ctx context.Context, db, collection string, query, projection, result interface{}) error
	SearchWithPageable(ctx context.Context, db, collection string, page Pageable, query, result interface{}) error
	SearchPagination(ctx context.Context, db, collection string, pageSize, pageNumber int64, query, result interface{}) error
	Count(ctx context.Context, db, collection string, query interface{}) (int64, error)
	Insert(ctx context.Context, db, collection string, data interface{}) error
	InsertAll(ctx context.Context, db, collection string, list []interface{}) error
	Remove(ctx context.Context, db, collection string, query interface{}) error
	RemoveAll(ctx context.Context, db, collection string, query interface{}) error
	Update(ctx context.Context, db, collection string, query, into interface{}) error
	FindOneAndUpdate(ctx context.Context, db, collection string, query, update, projection, result interface{}) error
	Upsert(ctx context.Context, db, collection string, query, into interface{}) error
	UpsertWithOperationResult(ctx context.Context, db, collection string, query, into interface{}) (int64, int64, error)
	UpdateAll(ctx context.Context, db, collection string, query, into interface{}) error
	Distinct(ctx context.Context, db, collection, key string, query interface{}) ([]interface{}, error)
	Aggregate(ctx context.Context, db, collection string, pipeline []bson.D, result interface{}) error
	CreateUniqueIndex(ctx context.Context, db, collection, field string) error
	Close() error
}

type (
	Client struct {
		timeout time.Duration
		client  *mongo.Client
		cfg     Config
	}

	DialInfo struct {
		Addrs          []string
		URI            string
		Username       string
		Password       string
		Timeout        time.Duration
		Database       string
		ReplicaSetName string
		Source         string
	}

	Config struct {
		DialInfo DialInfo `mapstructure:"DialInfo"`
		TLS      struct {
			Enabled  bool   `mapstructure:"Enabled"`
			CertPath string `mapstructure:"CertPath"`
		} `mapstructure:"TLS"`
	}

	Pageable struct {
		Size int64       // page size: number of records return in one page
		No   int64       // page number, start from 1
		Sort interface{} // sort, nil if default, or bson.M{"_id": 1}
	}

	ObjectID           = primitive.ObjectID
	UnstructuredObject = bson.M
)

func NewMongoClient(cfg Config) (*Client, error) {
	clientOptions := options.Client().
		SetTimeout(cfg.DialInfo.Timeout)
	if cfg.DialInfo.URI == "" {
		clientOptions.SetHosts(cfg.DialInfo.Addrs)
	} else {
		clientOptions.ApplyURI(cfg.DialInfo.URI)
	}
	if cfg.DialInfo.Username != "" && cfg.DialInfo.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username:   cfg.DialInfo.Username,
			Password:   cfg.DialInfo.Password,
			AuthSource: cfg.DialInfo.Source,
		})
	}
	if replicaSet := cfg.DialInfo.ReplicaSetName; replicaSet != "" {
		clientOptions.SetReplicaSet(replicaSet)
	}
	if cfg.TLS.Enabled {
		var tlsCfg = tls.Config{}
		if path := cfg.TLS.CertPath; path == "" {
			tlsCfg.InsecureSkipVerify = true
		} else if ca, err := ioutil.ReadFile(path); err == nil {
			roots := x509.NewCertPool()
			roots.AppendCertsFromPEM(ca)
			tlsCfg.RootCAs = roots
		}
		clientOptions.TLSConfig = &tlsCfg
	}
	/*// Cannot apply for now.
	// Equivalent to mgo.Monotonic
	clientOptions.SetReadPreference(readpref.Secondary()).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))*/
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &Client{
		timeout: cfg.DialInfo.Timeout,
		client:  client,
		cfg:     cfg,
	}, nil
}

func (c *Client) GetCollection(db, collection string) *mongo.Collection {
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	return c.client.Database(db).Collection(collection)
}

func (c *Client) Search(ctx context.Context, db, collection string, query, result interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	curr, err := c.client.Database(db).Collection(collection).Find(ctx, query)
	if err != nil {
		return err
	}
	if curr.Err() != nil {
		return curr.Err()
	}
	if !curr.TryNext(ctx) {
		return ErrNoDocuments
	}
	return curr.All(ctx, result)
}

func (c *Client) SearchOne(ctx context.Context, db, collection string, query, result interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	return c.client.Database(db).Collection(collection).FindOne(ctx, query).Decode(result)
}

func (c *Client) SearchWithProjection(ctx context.Context, db, collection string, query, projection, result interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	opts := options.Find().SetProjection(projection)
	curr, err := c.client.Database(db).Collection(collection).Find(ctx, query, opts)
	if err != nil {
		return err
	}
	if curr.Err() != nil {
		return curr.Err()
	}
	defer func() {
		_ = curr.Close(ctx)
	}()
	if !curr.TryNext(ctx) {
		return ErrNoDocuments
	}
	return curr.All(ctx, result)
}

func (c *Client) SearchWithPageable(ctx context.Context, db, collection string, page Pageable, query, result interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	opts := options.Find().SetSort(page.Sort).SetSkip(page.Size * (page.No - 1)).SetLimit(page.Size)
	cur, err := c.client.Database(db).Collection(collection).Find(ctx, query, opts)
	if err != nil {
		return err
	}
	defer func() {
		_ = cur.Close(ctx)
	}()
	if !cur.TryNext(ctx) {
		return ErrNoDocuments
	}
	return cur.All(ctx, result)
}

func (c *Client) SearchPagination(ctx context.Context, db, collection string, pageSize, pageNumber int64, query, result interface{}) error {
	var page = Pageable{Size: pageSize, No: pageNumber, Sort: bson.M{"_id": 1}}
	return c.SearchWithPageable(ctx, db, collection, page, query, result)
}

func (c *Client) Count(ctx context.Context, db, collection string, query interface{}) (int64, error) {
	if c.client == nil {
		return 0, ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	return c.client.Database(db).Collection(collection).CountDocuments(ctx, query)
}

func (c *Client) Insert(ctx context.Context, db, collection string, data interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}

	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	_, err := c.client.Database(db).Collection(collection).InsertOne(ctx, data)
	return err
}

func (c *Client) InsertAll(ctx context.Context, db, collection string, list []interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	var writes []mongo.WriteModel
	for _, ins := range list {
		model := mongo.NewInsertOneModel().SetDocument(ins)
		writes = append(writes, model)
	}
	_, err := c.client.Database(db).Collection(collection).BulkWrite(ctx, writes)
	return err
}

func (c *Client) Remove(ctx context.Context, db, collection string, query interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	result, err := c.client.Database(db).Collection(collection).DeleteOne(ctx, query)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNoDocuments
	}
	return nil
}

func (c *Client) RemoveAll(ctx context.Context, db, collection string, query interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	result, err := c.client.Database(db).Collection(collection).DeleteMany(ctx, query)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNoDocuments
	}
	return nil
}

func (c *Client) Update(ctx context.Context, db, collection string, query, into interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	result, err := c.client.Database(db).Collection(collection).UpdateOne(ctx, query, into)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNoDocuments
	}
	return nil
}

func (c *Client) FindOneAndUpdate(ctx context.Context, db, collection string, query, update, projection, result interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}

	if projection == nil {
		return c.client.Database(db).Collection(collection).FindOneAndUpdate(ctx, query, update).Decode(result)
	}
	return c.client.Database(db).Collection(collection).
		FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetProjection(projection)).
		Decode(result)
}

func (c *Client) Upsert(ctx context.Context, db, collection string, query, into interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	if _, err := c.client.Database(db).Collection(collection).UpdateOne(ctx, query, into, options.Update().SetUpsert(true)); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpsertWithOperationResult(ctx context.Context, db, collection string, query, into interface{}) (int64, int64, error) {
	if c.client == nil {
		return 0, 0, ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	updateResult, err := c.client.Database(db).Collection(collection).UpdateOne(ctx, query, into, options.Update().SetUpsert(true))
	if err != nil {
		return 0, 0, err
	}
	return updateResult.ModifiedCount, updateResult.UpsertedCount, nil
}

func (c *Client) UpdateAll(ctx context.Context, db, collection string, query, into interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	result, err := c.client.Database(db).Collection(collection).UpdateMany(ctx, query, into)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNoDocuments
	}
	return nil
}

func (c *Client) Distinct(ctx context.Context, db, collection, key string, query interface{}) ([]interface{}, error) {
	if c.client == nil {
		return nil, ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}
	if query == nil {
		query = bson.D{}
	}
	return c.client.Database(db).Collection(collection).Distinct(ctx, key, query)
}

func (c *Client) Aggregate(ctx context.Context, db, collection string, pipeline []bson.D, result interface{}) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}

	cursor, err := c.client.Database(db).Collection(collection).Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	return cursor.All(ctx, result)
}

func (c *Client) CreateUniqueIndex(ctx context.Context, db, collection, field string) error {
	if c.client == nil {
		return ErrUninitiated
	}
	if db == "" {
		db = c.cfg.DialInfo.Database
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := c.client.Database(db).Collection(collection).Indexes().CreateOne(ctx, indexModel)

	return err
}

func (c *Client) Close() error {
	return c.client.Disconnect(context.Background())
}

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 || we.Code == 11001 || we.Code == 12582 || we.Code == 16460 && strings.Contains(we.Error(), " E11000 ") {
				return true
			}
		}
	}
	return false
}

func NewObjectID() ObjectID {
	return primitive.NewObjectID()
}
