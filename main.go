package main

import (
	"log"
	"time"

	"github.com/Akshit8/go-mongodb/repository/mongo"

	"github.com/Akshit8/go-mongodb/config"
)

var noteRepo mongo.NoteRepository

func main() {
	var appConfig config.AppConfig
	err := config.LoadConfig("cmd/config", &appConfig)
	if err != nil {
		log.Fatalln("error loading config: ", err)
	}

	timeout := time.Duration(10) * time.Second
	client, err := mongo.NewMongoClient(appConfig.MongoURI, timeout)
	if err != nil {
		log.Fatalln("error connecting to mongodb: ", err)
	}

	noteRepo = mongo.NewNoteRepository(client, appConfig.DBName, appConfig.NotesTable, timeout)
}

func createNote() {}

func getNote() {}
