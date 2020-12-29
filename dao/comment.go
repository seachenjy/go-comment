package dao

import (
	"html"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Dao comment db
type Dao interface {
	AddComment(*Comment) bool
	GetComments(s SourceID, offset, limit int64) []*Comment
	Aggregate(s SourceID) *CommentStatistics
}

//SourceID comment source
//source id can build from page params or url
//it's must be unique identification
type SourceID string

//CommentStatistics statistics comments
type CommentStatistics struct {
	Count int     `bson:"count"`
	Grade float64 `bson:"grade"`
}

//Comment mongo fields
type Comment struct {
	ID       primitive.ObjectID `bson:"_id"`                        //comment id
	Parent   primitive.ObjectID `bson:"parent" json:"parent"`       //parent comment id, reply comment
	SourceID SourceID           `bson:"source_id" json:"source_id"` //comment source id
	Content  string             `bson:"content" json:"content"`     //comment content
	Grade    float32            `bson:"grade" json:"grade"`         //source grade .5 - 5
	Ctime    time.Time          `bson:"c_time" json:"c_time"`       //create time
	Utime    time.Time          `bson:"u_time" json:"u_time"`       //update time
	User     `bson:",inline" json:",inline"`
}

//User comment user
type User struct {
	Avatar    string `bson:"avatar" json:"avatar"`
	NickName  string `bson:"nick_name" json:"nick_name"`
	IPAddress string `bson:"ip_address" json:"ip_address"`
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
	c.Avatar = html.EscapeString(c.Avatar)
	c.Content = html.EscapeString(c.Content)
	return d.AddComment(c)
}

//Get get comments by source id
func Get(s SourceID, d Dao, offset, limit int64) ([]*Comment, *CommentStatistics) {
	comments := d.GetComments(s, offset, limit)
	aggregate := d.Aggregate(s)
	return comments, aggregate
}
