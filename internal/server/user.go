package server

import (
	"net/http"

	"github.com/Snoop-Duck/ToDoList/internal/services/user"

	"github.com/Snoop-Duck/ToDoList/internal/domain/users"

	"github.com/gin-gonic/gin"
)

func (s *NotesAPI) login(ctx *gin.Context) {
	var uReq users.UserRequest

	if err := ctx.ShouldBindJSON(&uReq); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userService := user.New(s.repo)

	userID, err := userService.LoginUser(uReq)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := jwtToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("Authorization", token)
	ctx.String(http.StatusOK, "user logined: %s", userID)
}

func (s *NotesAPI) register(ctx *gin.Context) {
	var uReq users.User

	if err := ctx.ShouldBindJSON(&uReq); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userService := user.New(s.repo)

	userID, err := userService.RegisterUser(uReq)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	token, err := jwtToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("Authorization", token)
	ctx.String(http.StatusOK, "user registered: %s", userID)
}

func (s *NotesAPI) deleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	userService := user.New(s.repo)
	err := userService.DeleteUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user"})
		return
	}
	ctx.String(http.StatusOK, "User deleted: %s", userID)
}

func (s *NotesAPI) getUsers(ctx *gin.Context) {
	userService := user.New(s.repo)
	allUsers, err := userService.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusAccepted, gin.H{"error": "No users"})
		return
	}
	ctx.String(http.StatusOK, "Users get: %s", allUsers)
}

func (s *NotesAPI) getUserID(ctx *gin.Context) {
	userID := ctx.Param("id")
	userService := user.New(s.repo)
	getUser, err := userService.GetUser(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No users"})
		return
	}
	ctx.String(http.StatusOK, "User get: %s", getUser)
}

func (s *NotesAPI) updateUserID(ctx *gin.Context) {
	var uReq users.User
	userID := ctx.Param("id")
	userService := user.New(s.repo)
	err := userService.UpdateUser(userID, uReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user"})
		return
	}
	ctx.String(http.StatusOK, "User update: %s", userID)
}
