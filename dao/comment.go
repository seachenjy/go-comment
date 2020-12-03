package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Dao comment db
type Dao interface {
	AddComment(*Comment) bool
	GetComments(s SourceID, offset, limit int64) []*Comment
}

//SourceID comment source
//source id can build from page params or url
//it's must be unique identification
type SourceID string

//Comment mongo fields
type Comment struct {
	ID       primitive.ObjectID `bson:"_id"`       //comment id
	Parent   primitive.ObjectID `bson:"parent"`    //parent comment id, reply comment
	SourceID SourceID           `bson:"source_id"` //comment source id
	Content  []byte             `bson:"content"`   //comment content
	Ctime    time.Time          `bson:"c_time"`    //create time
	Utime    time.Time          `bson:"u_time"`    //update time
	User     `bson:",inline"`
}

//User comment user
type User struct {
	Avatar    string `bson:"avatar"`
	NickName  string `bson:"nick_name"`
	IPAddress string `bson:"ip_address"`
}

//New build new comment
func New() *Comment {
	c := &Comment{}
	c.Ctime = time.Now()
	c.ID = primitive.NewObjectID()
	return c
}

//Save save comment to mongo
func (c *Comment) Save(d Dao) bool {
	return d.AddComment(c)
}

//Get get comments by source id
func Get(s SourceID, d Dao, offset, limit int64) []*Comment {
	return d.GetComments(s, offset, limit)
}
