package server

import (
	"net/http"

	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/services/note"

	"github.com/gin-gonic/gin"
)

// CreateNote godoc
// @Summary Create new note
// @Description Create a new note for authenticated user
// @Tags notes
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   note body notes.Note true "Note data"
// @Success 201 {object} map[string]interface{} "Note successfully created"
// @Failure 400 {object} map[string]interface{} "Invalid JSON body"
// @Failure 409 {object} map[string]interface{} "Note creation conflict"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notes/add [post]
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

// GetNotes godoc
// @Summary Get all notes
// @Description Get list of all notes for authenticated user
// @Tags notes
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} notes.Note "List of notes"
// @Failure 204 {object} map[string]interface{} "No notes found"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notes/list [get]
func (s *NotesAPI) getNotes(ctx *gin.Context) {
	s.log.Debug().Str("uid", ctx.GetString("uid")).Msg("user id from gin context")
	noteService := note.New(s.repoNote)

	notes, err := noteService.GetNotes()
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "No tasks"})
		return
	}
	ctx.String(http.StatusAccepted, "Notes get: %v", notes)
}

// GetNoteID godoc
// @Summary Get note by ID
// @Description Get note details by note ID
// @Tags notes
// @Produce  json
// @Param   id path string true "Note ID"
// @Success 200 {object} notes.Note "Note details"
// @Failure 204 {object} map[string]interface{} "Note not found"
// @Failure 400 {object} map[string]interface{} "Invalid note ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notes/list/{id} [get]
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

// DeleteNote godoc
// @Summary Delete note
// @Description Delete note by ID
// @Tags notes
// @Produce  json
// @Param   id path string true "Note ID"
// @Success 200 {object} map[string]interface{} "Note successfully deleted"
// @Failure 400 {object} map[string]interface{} "Invalid note ID or note not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notes/del/{id} [delete]
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

// UpdateNote godoc
// @Summary Update note
// @Description Update existing note by ID
// @Tags notes
// @Accept  json
// @Produce  json
// @Param   id path string true "Note ID"
// @Param   note body notes.Note true "Note data to update"
// @Success 200 {object} map[string]interface{} "Note successfully updated"
// @Failure 400 {object} map[string]interface{} "Invalid JSON body or note ID"
// @Failure 404 {object} map[string]interface{} "Note not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notes/upd/{id} [put]
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
