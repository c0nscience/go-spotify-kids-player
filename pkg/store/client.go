package store

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Store interface {
	Disconnect(ctx context.Context) error
	//Create(ctx context.Context, userId string, d interface{}) (interface{}, error)
	Save(ctx context.Context, d HasObjectId) error
	FindOne(ctx context.Context, filter interface{}, rec interface{}) error
	Find(ctx context.Context, filter interface{}, sort interface{}, rec interface{}) error
	//DeleteAll(ctx context.Context, userId string) (int64, error)
	//DropCollection(ctx context.Context) error
	Delete(ctx context.Context, filter interface{}, rec interface{}) error
	//Distinct(ctx context.Context, userId string, field string) ([]interface{}, error)
}

type CollectionName string

type mongoDbStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func New(uri, db string, collection CollectionName) (Store, error) {
	opts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cl, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal().Err(err).
			Msg("could not create mongodb client")
		return nil, err
	}

	err = cl.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err).
			Msg("cloud not ping database")
		return nil, err
	}

	database := cl.Database(db)
	return &mongoDbStore{
		client: cl,
		coll:   database.Collection(string(collection)),
	}, err
}

func (me *mongoDbStore) Save(ctx context.Context, d HasObjectId) error {
	opts := options.
		FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)
	var objectId = d.ObjectId()
	if objectId == primitive.NilObjectID {
		objectId = primitive.NewObjectID()
		d.SetObjectId(objectId)
	}
	res := me.coll.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.D{{Key: "$set", Value: d}}, opts)

	return res.Err()
}

func (me *mongoDbStore) Find(ctx context.Context, filter interface{}, sort interface{}, rec interface{}) error {
	opts := options.Find()
	if sort != nil {
		opts.SetSort(sort)
	}
	cur, err := me.coll.Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	return cur.All(ctx, rec)
}

func (me *mongoDbStore) FindOne(ctx context.Context, filter interface{}, rec interface{}) error {
	res := me.coll.FindOne(ctx, filter)

	if rec != nil {
		return res.Decode(rec)
	}

	return res.Err()
}

func (me *mongoDbStore) Delete(ctx context.Context, filter interface{}, rec interface{}) error {
	res := me.coll.FindOneAndDelete(ctx, filter)

	if rec != nil {
		return res.Decode(rec)
	}

	return res.Err()
}

func (me *mongoDbStore) Disconnect(ctx context.Context) error {
	return me.client.Disconnect(ctx)
}

type HasObjectId interface {
	ObjectId() primitive.ObjectID
	SetObjectId(id primitive.ObjectID)
}
