package mongo

import (
	"sync"
	"testing"
	"time"

	"github.com/Akshit8/go-mongodb/entity"
	"github.com/Akshit8/go-mongodb/random"
	"github.com/stretchr/testify/require"
)

func createRandomNote(t *testing.T) entity.Note {
	newNote := getRandomNote()

	err := noteRepo.CreateNote(&newNote)

	require.NoError(t, err)

	return newNote
}

func getRandomNote() entity.Note {
	newNote := entity.Note{
		NoteID:      random.GetRandomUUID(),
		Title:       random.GetRandomString(8),
		Description: random.GetRandomString(20),
		Tags:        random.GetRandomStringListOfSizeN(5),
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return newNote
}

func TestCreateNote(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	wg.Add(4)
	// creating 4 random notes concurrently
	for i := 0; i < 4; i++ {
		go func() {
			createRandomNote(t)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGetNoteByID(t *testing.T) {
	t.Parallel()

	var notesList []entity.Note

	for i := 0; i < 4; i++ {
		notesList = append(notesList, createRandomNote(t))
	}

	var wg sync.WaitGroup
	wg.Add(4)
	// getting notes by id concurrently
	for _, note := range notesList {
		go func(note entity.Note) {
			retrievedNote, err := noteRepo.GetNoteByID(note.NoteID)
			require.NoError(t, err)
			require.NotEmpty(t, retrievedNote)

			require.Equal(t, note.NoteID, retrievedNote.NoteID)
			require.Equal(t, note.Title, retrievedNote.Title)
			require.Equal(t, note.Completed, retrievedNote.Completed)

			wg.Done()
		}(note)
	}
	wg.Wait()
}

func TestGetCompleteNotes(t *testing.T) {
	t.Parallel()

	for i := 0; i < 4; i++ {
		newNote := getRandomNote()
		newNote.Completed = true

		err := noteRepo.CreateNote(&newNote)

		require.NoError(t, err)
	}

	var notes []*entity.Note

	notes, err := noteRepo.GetCompletedNotes()

	require.NoError(t, err)

	for _, note := range notes {
		require.Equal(t, note.Completed, true)
	}

}

func TestUpdateNoteByID(t *testing.T) {
	t.Parallel()

	newNote := createRandomNote(t)

	newNote.Completed = true
	newNote.UpdatedAt = time.Now()
	err := noteRepo.UpdateNoteByID(&newNote)

	require.NoError(t, err)

	retrievedNote, err := noteRepo.GetNoteByID(newNote.NoteID)

	require.NoError(t, err)
	require.Equal(t, newNote.Completed, retrievedNote.Completed)
}

func TestDeleteNoteByID(t *testing.T) {
	t.Parallel()

	newNote := createRandomNote(t)
	err := noteRepo.DeleteNoteByID(newNote.NoteID)

	require.NoError(t, err)

	retrievedNote, err := noteRepo.GetNoteByID(newNote.NoteID)

	require.Error(t, err)
	require.Empty(t, retrievedNote)
}
