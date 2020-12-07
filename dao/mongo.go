package dao

import (
	"context"
	"sync"

	"github.com/seachenjy/go-comment/config"
	"github.com/seachenjy/go-comment/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	lock   *sync.Mutex = &sync.Mutex{}
	inited bool
	ss     *MongoService
)

//MongoService the db service
type MongoService struct {
	db *mongo.Database
}

//AddComment add comment to mongodb
func (m *MongoService) AddComment(c *Comment) bool {
	ctx := context.Background()
	_, err := m.db.Collection("comments").InsertOne(ctx, c)
	if err != nil {
		return false
	}
	return true
}

//GetComments get comments by sourceid
func (m *MongoService) GetComments(s SourceID, offset, limit int64) []*Comment {
	var cs []*Comment
	ctx := context.Background()
	opts := &options.FindOptions{}
	opts.SetLimit(limit)
	opts.SetSkip(offset)
	opts.SetSort(bson.M{
		"c_time": -1,
	})
	res, err := m.db.Collection("comments").Find(ctx, bson.M{
		"source_id": bson.M{"$eq": s},
	}, opts)
	if err == nil {
		var c Comment
		for res.Next(ctx) {
			if err := res.Decode(&c); err != nil {
				log.GetLogger().Error(err)
			} else {
				cs = append(cs, &c)
			}
		}
	}
	return cs
}

//NewMongo return a new mongodb connect service
func NewMongo(cfg *config.Config) *MongoService {
	lock.Lock()
	defer lock.Unlock()
	if inited {
		return ss
	}
	inited = true
	clientOptions := options.Client().ApplyURI(cfg.Mongo.MongoURL).SetMaxPoolSize(cfg.Mongo.Poolsize)
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	if err := conn.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}
	ss = &MongoService{}
	ss.db = conn.Database(cfg.Mongo.DbName)
	return ss
}
