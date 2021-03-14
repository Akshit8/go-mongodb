package mongo

import (
	"github.com/Akshit8/go-mongodb/config"
	"log"
	"os"
	"testing"
	"time"
)

var noteRepo NoteRepository

func TestMain(m *testing.M) {
	var appConfig config.AppConfig
	err := config.LoadConfig("../../../config", &appConfig)
	if err != nil {
		log.Fatalln("error loading config: ", err)
	}

	timeout := time.Duration(10) * time.Second
	client, err := NewMongoClient(appConfig.MongoURI, timeout)
	if err != nil {
		log.Fatalln("error connecting to mongodb: ", err)
	}

	noteRepo = NewNoteRepository(client, appConfig.DBName, appConfig.NotesTable, timeout)

	os.Exit(m.Run())
}
