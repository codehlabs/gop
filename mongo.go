package gop

import (
	"context"
	"database/sql"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoADriver struct {
	client     *mongo.Client
	db         *mongo.Database // database to operate with
	collection *mongo.Collection
	connstr    string
}

func (md *MongoADriver) connect() error {
	if md.connstr == "" {
		return errors.New("connetion string is not set")
	}
	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(md.connstr))
	if err != nil {
		return err
	}
	md.client = c
	return nil
}

func (md MongoADriver) Db() *sql.DB {
	return nil
}

// Authenticates [User] in the system using username, email of phone
func (md MongoADriver) Login(username, email, phone string, password string) (id string, err error) {
	if e := md.connect(); e != nil {
		return "", e
	}

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

func (md MongoADriver) Save(u User) error {
	if e := md.connect(); e != nil {
		return e
	}
	defer md.client.Disconnect(context.Background())

	matching := []bson.M{}

	if driverbehavior.UniqueEmail {
		matching = append(matching, bson.M{"email": u.Email})
	}

	if driverbehavior.UniquePhone {
		matching = append(matching, bson.M{"phone": u.Phone})
	}

	if driverbehavior.UniqueUsername {
		matching = append(matching, bson.M{"username": u.Username})
	}

	filter := bson.M{
		"$or": matching,
	}

	if len(matching) != 0 {
		cursor, err := md.collection.Find(context.TODO(), filter)
		if err != nil {
			return err
		}
		defer cursor.Close(context.TODO())

		if cursor.TryNext(context.TODO()) != false {
			return errors.New("the email, phone or username you provided already exits in the database")
		}
	}

	_, err := md.collection.InsertOne(context.Background(), u)
	if err != nil {
		return err
	}

	return nil
}

func (md MongoADriver) Update(u User) error {
	return nil
}

func (md MongoADriver) Delete(id string) error {
	return nil
}

func (md MongoADriver) Read(id string) (User, error) {
	return User{}, nil
}

func NewMongoADriver(conn string, databaseName string, collection string) (*MongoADriver, error) {
	d := &MongoADriver{connstr: conn}
	e := d.connect()
	if e != nil {
		return nil, e
	}
	d.db = d.client.Database(databaseName)
	d.collection = d.db.Collection(collection)
	return d, nil
}
