package server

import (
	"main/ToDoList/internal/domain/notes"
	"main/ToDoList/internal/services/note"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *NotesAPI) createNote(ctx *gin.Context) {
	var nReq notes.Note
	if err := ctx.BindJSON(&nReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	noteService := note.New(s.repoNote)

	noteID, err := noteService.CreateNote(nReq)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusCreated, "Note add: %s", noteID)
}

func (s *NotesAPI) getNotes(ctx *gin.Context) {
	noteService := note.New(s.repoNote)

	allNotes, err := noteService.GetNotes()
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "No tasks"})
		return
	}
	ctx.String(http.StatusAccepted, "Notes get: %s", allNotes)
}

func (s *NotesAPI) getNoteID(ctx *gin.Context) {
	noteID := ctx.Param("id")
	noteService := note.New(s.repoNote)
	note, err := noteService.GetNoteID(noteID)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "No task"})
		return
	}
	ctx.String(http.StatusAccepted, "Note get: %s", note)
}

func (s *NotesAPI) deleteNote(ctx *gin.Context) {
	noteID := ctx.Param("id")
	noteService := note.New(s.repoNote)
	err := noteService.DeleteNoteID(noteID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No task"})
		return
	}
	ctx.String(http.StatusOK, "Note deleted: %s", noteID)
}

func (s *NotesAPI) updateNote(ctx *gin.Context) {
	var nReq notes.Note
	if err := ctx.BindJSON(&nReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	noteID := ctx.Param("id")
	noteService := note.New(s.repoNote)
	err := noteService.UpdateNoteID(noteID, nReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No task"})
		return
	}
	ctx.String(http.StatusOK, "Note update: %s", noteID)
}
