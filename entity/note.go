package entity

import "time"

// Note struct defines entity of resource NOTE.
// bson tags are used by mongo-driver.
// for NoteID uuid as string is used to prevent dependency on mongo object id.
type Note struct {
	NoteID      string    `json:"noteId" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Completed   bool      `json:"completed" bson:"completed"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
}
