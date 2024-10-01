package repo

import (
	"context"
	"time"
	. "untitled/internal/musicstorage/mdl"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	defailtCtxTimeoutInSecs = 5
)

type MongoRepo struct {
	client    *mongo.Client
	trackColl *mongo.Collection
	defaultCtxTimeout time.Duration
}

func NewMongoRepo(client *mongo.Client, dbName, collectionName string) *MongoRepo {
	return &MongoRepo{
		client:    client,
		trackColl: client.Database(dbName).Collection(collectionName),
		defaultCtxTimeout: defailtCtxTimeoutInSecs * time.Second, 
	}
}

func NewMongoRepoWithTimeout(client *mongo.Client, dbName, collectionName string, timeout time.Duration) *MongoRepo {
	return &MongoRepo{
		client:    client,
		trackColl: client.Database(dbName).Collection(collectionName),
		defaultCtxTimeout: timeout,
	}
}

func (r *MongoRepo) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.defaultCtxTimeout)
}

func (r *MongoRepo) CreateTrack(track Track_MONGO) error {
	ctx, cancel := r.createContext()
	defer cancel()

	_, err := r.trackColl.InsertOne(ctx, track)
	return err
}

func (r *MongoRepo) UpdateTrackPtr(trackUUID string) error {
	ctx, cancel := r.createContext()
	defer cancel()

	filter := bson.M{"_id": trackUUID}
	update := bson.M{"$set": bson.M{"ptr": trackUUID}}

	_, err := r.trackColl.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepo) DeleteTrackByID(trackID uuid.UUID) error {
	ctx, cancel := r.createContext()
	defer cancel()

	_, err := r.trackColl.DeleteOne(ctx, bson.M{"_id": trackID})
	return err
}

func (r *MongoRepo) UpdateTrackMetadata(trackID uuid.UUID, req TrackRequest) error {
	ctx, cancel := r.createContext()
	defer cancel()

	filter := bson.M{"_id": trackID}
	update := bson.M{
		"$set": bson.M{
			"title":    req.Title,
			"artist":   req.Artist,
			"album":    req.Album,
			"genre":    req.Genre,
			"format":   req.Format,
			"duration": req.Duration,
		},
	}

	_, err := r.trackColl.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepo) ListTracks() ([]Track_MONGO, error) {
	ctx, cancel := r.createContext()
	defer cancel()

	cursor, err := r.trackColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tracks []Track_MONGO
	if err = cursor.All(ctx, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *MongoRepo) GetTrackByID(trackID uuid.UUID) (*Track_MONGO, error) {
	ctx, cancel := r.createContext()
	defer cancel()

	var trackMongo Track_MONGO
	err := r.trackColl.FindOne(ctx, bson.M{"_id": trackID}).Decode(&trackMongo)
	if err != nil {
		return nil, err
	}

	return &trackMongo, nil
}
