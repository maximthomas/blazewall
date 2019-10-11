package repo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/google/uuid"
	"github.com/maximthomas/blazewall/user-service/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepoUser struct {
	models.User `bson:",inline"`
	Password    string `json:"password,omitempty"`
}

type UserRepositoryMongoDB struct {
	client     *mongo.Client
	db         string
	collection string
}

/*
	GetUser(realm, userID string) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(realm, userID string) error
	SetPassword(realm, userID, password string) error
	ValidatePassword(realm, userID, password string) (bool, error)
*/

func (ur *UserRepositoryMongoDB) GetUser(realm, userID string) (models.User, error) {
	var user models.User
	collection := ur.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	filter := bson.M{"id": userID, "realm": realm}

	var repoUser mongoRepoUser
	err := collection.FindOne(ctx, filter).Decode(&repoUser)

	if err != nil {
		return user, err
	}

	return repoUser.User, nil
}

func (ur *UserRepositoryMongoDB) CreateUser(user models.User) (models.User, error) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	repoUser := mongoRepoUser{
		User:     user,
		Password: "",
	}
	collection := ur.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, &repoUser)
	if err != nil {
		return user, err
	}

	return user, nil

}

func (ur *UserRepositoryMongoDB) UpdateUser(user models.User) (models.User, error) {

	collection := ur.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.M{"id": user.ID, "realm": user.Realm}
	var updatedUser mongoRepoUser
	err := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": user}).Decode(&updatedUser)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserRepositoryMongoDB) DeleteUser(realm, userID string) error {
	collection := ur.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.M{"id": userID, "realm": realm}
	err := collection.FindOneAndDelete(ctx, filter).Err()

	if err != nil {
		return err
	}
	return nil

}

func (ur *UserRepositoryMongoDB) SetPassword(realm, userID, password string) error {
	collection := ur.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	filter := bson.M{"id": userID, "realm": realm}
	var updatedUser mongoRepoUser
	err := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": bson.M{"password": password}}).Decode(&updatedUser)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepositoryMongoDB) ValidatePassword(realm, userID, password string) (bool, error) {
	collection := ur.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	filter := bson.M{"id": userID, "realm": realm}
	var repoUser mongoRepoUser
	err := collection.FindOne(ctx, filter).Decode(&repoUser)
	var valid bool
	if err != nil {
		return valid, err
	}
	valid = repoUser.Password == password

	return valid, nil
}

func (ur *UserRepositoryMongoDB) getCollection() *mongo.Collection {
	return ur.client.Database(ur.db).Collection(ur.collection)
}

func NewUserRepositoryMongoDB(uri, db, collection string) UserRepositoryMongoDB {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("connecting to mongo, uri: %v", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	return UserRepositoryMongoDB{
		client:     client,
		db:         db,
		collection: collection,
	}
}
