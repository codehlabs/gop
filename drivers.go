package gop

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActionDriver interface {
	Login(username, email, phone string, password string) (id string, err error)
}

type MongoADriver struct {
	client     *mongo.Client
	db         *mongo.Database // database to operate with
	collection *mongo.Collection
}

func (md MongoADriver) Login(username, email, phone string, password string) (id string, err error) {
	filter := bson.D{}
	projection := bson.D{}

	if username != "" {
		filter = bson.D{{"username", username}}
		projection = bson.D{{"username", 1}}
	} else if email != "" {
		filter = bson.D{{"email", email}}
		projection = bson.D{{"email", 1}}
	} else if phone != "" {
		filter = bson.D{{"phone", phone}}
		projection = bson.D{{"phone", 1}}
	}

	filter = append(filter)
	projection = append(projection, bson.E{"password", 1}, bson.E{"salt", 1})
	ctx := context.Background()
	defer md.client.Disconnect(ctx)

	msr := md.collection.FindOne(ctx, filter, options.FindOne().SetProjection(bson.D{}))
	var u User

	err = msr.Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", ErrUnabelToAuthenticate
		}
		return "", err
	}

	hash, err := ValidateHash(u.Salt, password)
	if err != nil {
		return "", err
	}

	if hash != u.Password {
		return "", ErrUnabelToAuthenticate
	}

	return u.Id, nil
}

func NewMongoADriver(conn string, databaseName string, collection string) (MongoADriver, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return MongoADriver{}, err
	}
	d := MongoADriver{}
	d.client = client
	d.db = d.client.Database(databaseName)
	d.collection = d.db.Collection(collection)

	return d, nil
}

type SqlLADriver struct {
}
