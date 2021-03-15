package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Akshit8/go-mongodb/config"
	"github.com/Akshit8/go-mongodb/entity"
	"github.com/Akshit8/go-mongodb/repository/mongo"
	"github.com/google/uuid"
)

func main() {
	var appConfig config.AppConfig
	err := config.LoadConfig("config", &appConfig)
	if err != nil {
		log.Fatalln("error loading config: ", err)
	}

	timeout := time.Duration(10) * time.Second
	client, err := mongo.NewMongoClient(appConfig.MongoURI, timeout)
	if err != nil {
		log.Fatalln("error connecting to mongodb: ", err)
	}

	noteRepo := mongo.NewNoteRepository(client, appConfig.DBName, appConfig.NotesTable, timeout)

	// creating a test note.
	testID := uuid.New().String()
	newNote := entity.Note{
		NoteID:      testID,
		Title:       "GO MONGODB",
		Description: "a note to complete this repository",
		Completed:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// saving note in db.
	err = noteRepo.CreateNote(&newNote)
	if err != nil {
		log.Fatalln("CreateNote error: ", err)
	}

	// retrieving note from db.
	note, err := noteRepo.GetNoteByID(testID)
	if err != nil {
		log.Fatalln("GetNoteByID error: ", err)
	}
	fmt.Println("retrieved note: ", note)

	// updating title of the note.
	note.Title = "UPDATED GO TITLE"
	err = noteRepo.UpdateNoteByID(note)
	if err != nil {
		log.Fatalln("UpdateNoteByID error: ", err)
	}

	// deleting the above created note.
	err = noteRepo.DeleteNoteByID(note.NoteID)
	if err != nil {
		log.Fatalln("DeleteNoteByID error: ", err)
	}
}
