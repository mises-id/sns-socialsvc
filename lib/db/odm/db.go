package odm

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/inflection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	db *mongo.Database
}

type collectioner interface {
	CollectionName() string
}

func NewClient(conn *mongo.Database) *Client {
	return &Client{db: conn}
}

func (c *Client) NewSession(ctx context.Context) *DB {
	return &DB{
		ctx:     ctx,
		db:      c.db,
		options: &options.FindOptions{},
	}
}

type DB struct {
	db             *mongo.Database
	ctx            context.Context
	Error          error
	collectionName string
	out            interface{}
	options        *options.FindOptions
	sort           interface{}
	condition      bson.M
}

func (db *DB) Collection(collectionName string) *DB {
	db.collectionName = collectionName
	return db
}

func (db *DB) Model(model interface{}) *DB {
	db.collectionName = ""
	db.reflectCollectionName(model)
	return db
}

func (db *DB) Delete(out interface{}, pk interface{}) *DB {
	db.out = out
	_, err := db.db.Collection(db.reflectCollectionName()).DeleteOne(db.ctx, bson.M{"_id": pk})
	if err != nil {
		db.Error = err
		return db
	}
	return db
}

func (db *DB) Save(out interface{}, pk interface{}) *DB {
	db.out = out
	beforeUpdate(db.out)
	_, err := db.db.Collection(db.reflectCollectionName()).UpdateOne(db.ctx, bson.M{"_id": pk}, out)
	if err != nil {
		db.Error = err
		return db
	}
	return db
}

func (db *DB) Create(out interface{}) *DB {
	db.out = out
	beforeCreate(db.out)
	result, err := db.db.Collection(db.reflectCollectionName()).InsertOne(db.ctx, out)
	if err != nil {
		db.Error = err
		return db
	}
	value := reflect.ValueOf(out).Elem()
	idValue := value.FieldByName("ID")
	resultID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		db.Error = errors.New("invalid inserted id")
		return db
	}
	if idValue.IsValid() {
		if idValue.CanSet() {
			idValue.Set(reflect.ValueOf(resultID))
		}
	}
	return db
}

func (db *DB) Where(condition bson.M) *DB {
	if db.condition == nil {
		db.condition = condition
		return db
	}
	for key, value := range condition {
		db.condition[key] = value

	}
	return db
}

func (db *DB) Count(c *int64) *DB {
	countDocs, err := db.db.Collection(db.reflectCollectionName()).CountDocuments(db.ctx, db.condition)
	db.Error = err
	*c = countDocs
	return db
}

func (db *DB) Sort(sort interface{}) *DB {
	if db.sort == nil {
		db.sort = sort
	}
	return db
}

func (db *DB) Limit(limit int64) *DB {
	db.options.Limit = &limit
	return db
}

func (db *DB) Skip(skip int64) *DB {
	db.options.Skip = &skip
	return db
}

func (db *DB) First(out interface{}, conditions ...bson.M) *DB {
	db.out = out
	for _, condition := range conditions {
		db = db.Where(condition)
	}
	result := db.db.Collection(db.reflectCollectionName()).FindOne(db.ctx, db.condition, &options.FindOneOptions{Sort: bson.M{"_id": 1}})
	db.Error = result.Err()
	if db.Error != nil {
		return db
	}
	db.Error = result.Decode(out)
	return db
}

func (db *DB) Last(out interface{}, conditions ...bson.M) *DB {
	db.out = out
	for _, condition := range conditions {
		db = db.Where(condition)
	}
	result := db.db.Collection(db.reflectCollectionName()).FindOne(db.ctx, db.condition, &options.FindOneOptions{Sort: bson.M{"_id": -1}})
	db.Error = result.Err()
	if db.Error != nil {
		return db
	}
	db.Error = result.Decode(out)
	return db
}

func (db *DB) Find(out interface{}, conditions ...bson.M) *DB {
	db.out = out
	for _, condition := range conditions {
		db = db.Where(condition)
	}
	if db.sort != nil {
		db.options.Sort = db.sort
	}
	var result *mongo.Cursor
	if db.options == nil {
		result, db.Error = db.db.Collection(db.reflectCollectionName()).Find(db.ctx, db.condition)
	} else {
		result, db.Error = db.db.Collection(db.reflectCollectionName()).Find(db.ctx, db.condition, db.options)
	}
	if db.Error != nil {
		return db
	}
	db.Error = result.All(db.ctx, out)
	return db
}

func (db *DB) reflectCollectionName(outs ...interface{}) string {
	if db.collectionName != "" {
		return db.collectionName
	}
	out := db.out
	if len(outs) > 0 {
		out = outs[0]
	}
	t := reflect.TypeOf(out)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		log.Fatal("input must be ptr")
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	value := reflect.New(t)
	item, ok := value.Interface().(collectioner)
	if ok {
		db.collectionName = item.CollectionName()
	} else {
		db.collectionName = strings.ToLower(inflection.Plural(t.Name()))
	}
	return db.collectionName
}

func beforeCreate(out interface{}) {
	value := reflect.ValueOf(out)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	now := time.Now()
	createTimeFiled := value.FieldByName("CreatedAt")
	if createTimeFiled.CanSet() {
		createTimeFiled.Set(reflect.ValueOf(now))
	}
	updateTimeFiled := value.FieldByName("UpdatedAt")
	if updateTimeFiled.CanSet() {
		updateTimeFiled.Set(reflect.ValueOf(now))
	}
}

func beforeUpdate(out interface{}) {
	value := reflect.ValueOf(out)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	now := time.Now()
	updateTimeFiled := value.FieldByName("UpdatedAt")
	if updateTimeFiled.CanSet() {
		updateTimeFiled.Set(reflect.ValueOf(now))
	}
}
