package main

import (
	"time"

	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName   string = "localhost:27017"
	dbName     string = "golang-todo-app"
	collection string = "todo"
	port       string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"createdAt"`
	}

	todo struct {
		ID        bson.ObjectId `json:"id"`
		Title     string        `json:"title"`
		Completed bool          `json:"completed"`
		CreatedAt time.Time     `json:"created_at"`
	}
)

func init() {
	rnd = renderer.New()
	sess, err := mgo.Dial(hostName)
	checkError(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}
