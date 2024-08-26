package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName   string = "localhost:27017"
	dbName     string = "golang-todo-app"
	collection string = "todos"
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

func main() {
	stopChannel := make(chan os.Signal)
	signal.Notify(stopChannel, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeRoutes)
	r.Mount("/todo", todoRoutes())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port: ", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Listen:&s\n", err)
		}
	}()

	<-stopChannel
	log.Println(("Shutting down server..."))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	srv.Shutdown(ctx)

	defer cancel()
	log.Println("Server gracefully shut down")

}
