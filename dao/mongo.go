package dao

import (
	"context"
	"fmt"
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
	db         *mongo.Database
	collection *mongo.Collection
}

//AddComment add comment to mongodb
func (m *MongoService) AddComment(c *Comment) bool {
	ctx := context.Background()
	_, err := m.collection.InsertOne(ctx, c)
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
	res, err := m.collection.Find(ctx, bson.M{
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

//Aggregate all comments count and average score
func (m *MongoService) Aggregate(s SourceID) *CommentStatistics {
	query := []bson.M{
		{"$match": bson.M{
			"source_id": bson.M{"$eq": s},
		}},
		{"$group": bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1},
			"grade": bson.M{"$avg": "$grade"},
		}},
	}
	var showsWithInfo CommentStatistics
	ctx := context.Background()
	res, err := m.collection.Aggregate(ctx, query)
	if err != nil {
		log.GetLogger().Error(err)
	} else {
		err = res.All(ctx, &showsWithInfo)
		if err != nil {
			log.GetLogger().Error(err)
		}
	}
	return &showsWithInfo
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
		log.GetLogger().Error(err)
		return nil
	}
	if err := conn.Ping(context.TODO(), nil); err != nil {
		log.GetLogger().Error(err)
		return nil
	}
	ss = &MongoService{}

	ss.db = conn.Database(cfg.Mongo.DbName)

	//test table exists
	ctx := context.Background()
	result, err := ss.db.ListCollectionNames(ctx, bson.M{"name": cfg.Mongo.TableName})
	if err != nil {
		panic(err)
	}
	if len(result) <= 0 {
		fmt.Printf("create document: '%s' in mongo database: '%s'", cfg.Mongo.TableName, cfg.Mongo.DbName)
		sourceIndex := mongo.IndexModel{
			Keys:    bson.M{"source_id": -1},
			Options: nil,
		}
		ctimeIndex := mongo.IndexModel{
			Keys:    bson.M{"c_time": -1},
			Options: nil,
		}
		parentIndex := mongo.IndexModel{
			Keys:    bson.M{"parent": -1},
			Options: nil,
		}
		_, err := ss.db.Collection(cfg.Mongo.TableName).Indexes().CreateMany(ctx, []mongo.IndexModel{sourceIndex, ctimeIndex, parentIndex})
		if err != nil {
			panic(err)
		}
	}
	ss.collection = ss.db.Collection(cfg.Mongo.TableName)
	return ss
}
