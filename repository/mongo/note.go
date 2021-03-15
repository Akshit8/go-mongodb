package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Akshit8/go-mongodb/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NoteRepository defines repository actions for resource NOTE.
type NoteRepository interface {
	CreateNote(newNote *entity.Note) error
	GetNoteByID(noteID string) (*entity.Note, error)
	GetNotes() ([]*entity.Note, error)
	UpdateNoteByID(updatedNote *entity.Note) error
	DeleteNoteByID(noteID string) error
}

type noteRepository struct {
	client     *mongo.Client
	database   string
	collection string
	timeout    time.Duration
}

// NewNoteRepository creates new instance of noteRepository
func NewNoteRepository(client *mongo.Client, database, collection string, timeout time.Duration) NoteRepository {
	return &noteRepository{
		client:     client,
		database:   database,
		collection: collection,
		timeout:    timeout,
	}
}

func (n *noteRepository) getNoteCollection() *mongo.Collection {
	collection := n.client.Database(n.database).Collection(n.collection)
	return collection
}

func (n *noteRepository) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), n.timeout)
}

func (n *noteRepository) CreateNote(newNote *entity.Note) error {
	ctx, cancel := n.getContext()
	defer cancel()

	collection := n.getNoteCollection()

	_, err := collection.InsertOne(ctx, newNote)
	if err != nil {
		return err
	}

	return nil
}

// getOne is generic function to get one note.
func (n *noteRepository) getOne(filter bson.M) (*entity.Note, error) {
	ctx, cancel := n.getContext()
	defer cancel()

	collection := n.getNoteCollection()
	var retrievedNote entity.Note
	err := collection.FindOne(ctx, filter).Decode(&retrievedNote)
	if err != nil {
		return nil, err
	}

	return &retrievedNote, nil
}

func (n *noteRepository) GetNoteByID(noteID string) (*entity.Note, error) {
	filter := bson.M{"_id": noteID}
	return n.getOne(filter)
}

func (n *noteRepository) GetNotes() ([]*entity.Note, error) {
	ctx, cancel := n.getContext()
	defer cancel()

	collection := n.getNoteCollection()
	cursor, err := collection.Find(ctx, bson.M{"title": "twitch"})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result []*entity.Note
	for cursor.Next(ctx) {
		var note entity.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, err
		}
		result = append(result, &note)
	}

	return result, nil
}

// updateOne is generic function to update one note
// uses ReplaceOne method to perform the update note.
// for performing update, the function is passed the new state of note and replaces it with existing document.
// this makes update operation consistent over various orm and db drivers(e.g mysql, postgres cassandra)
func (n *noteRepository) updateOne(filter bson.M, updatedNote *entity.Note) error {
	ctx, cancel := n.getContext()
	defer cancel()

	collection := n.getNoteCollection()
	result, err := collection.ReplaceOne(ctx, filter, updatedNote)
	if err != nil {
		return err
	}

	if result.ModifiedCount != 1 {
		return errors.New("error updating note")
	}

	return nil
}

func (n *noteRepository) UpdateNoteByID(updatedNote *entity.Note) error {
	filter := bson.M{"_id": updatedNote.NoteID}
	return n.updateOne(filter, updatedNote)
}

// deleteOne is generic function to delete one note.
func (n *noteRepository) deleteOne(filter bson.M) error {
	ctx, cancel := n.getContext()
	defer cancel()

	collection := n.getNoteCollection()
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("error deleting note")
	}

	return nil
}

func (n *noteRepository) DeleteNoteByID(noteID string) error {
	filter := bson.M{"_id": noteID}
	return n.deleteOne(filter)
}
