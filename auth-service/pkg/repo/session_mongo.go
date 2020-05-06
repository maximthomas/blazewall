package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/maximthomas/blazewall/auth-service/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSessionRepository struct {
	client     *mongo.Client
	db         string
	collection string
}

type mongoRepoSession struct {
	models.Session `bson:",inline"`
}

func NewMongoSessionRepository(uri, db, c string) (*MongoSessionRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("connecting to mongo, uri: %v", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoSessionRepository{
		client:     client,
		db:         db,
		collection: c,
	}, nil

}

func (sr *MongoSessionRepository) CreateSession(session models.Session) (models.Session, error) {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	repoSession := mongoRepoSession{
		Session: session,
	}
	collection := sr.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, &repoSession)
	if err != nil {
		return session, err
	}

	return session, nil
}

func (sr *MongoSessionRepository) DeleteSession(id string) error {
	collection := sr.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	err := collection.FindOneAndDelete(ctx, filter).Err()
	if err != nil {
		return err
	}
	return nil
}

func (sr *MongoSessionRepository) GetSession(id string) (models.Session, error) {
	var session models.Session
	collection := sr.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	filter := bson.M{"id": id}

	var repoSession mongoRepoSession
	err := collection.FindOne(ctx, filter).Decode(&repoSession)

	if err != nil {
		return session, err
	}

	return repoSession.Session, nil
}

func (sr *MongoSessionRepository) UpdateSession(session models.Session) error {
	collection := sr.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.M{"id": session.ID}
	var repoSession mongoRepoSession
	err := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": session}).Decode(&repoSession)

	if err != nil {
		return err
	}
	return nil
}

func (sr *MongoSessionRepository) getCollection() *mongo.Collection {
	return sr.client.Database(sr.db).Collection(sr.collection)
}
