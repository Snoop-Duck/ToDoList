package server

import (
	"net/http"

	"github.com/Snoop-Duck/ToDoList/internal/services/user"

	"github.com/Snoop-Duck/ToDoList/internal/domain/users"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags users
// @Accept  json
// @Produce  json
// @Param   credentials body users.UserRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Successfully logged in"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/login [post]
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

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user body users.User true "User registration data"
// @Success 200 {object} map[string]interface{} "Successfully registered"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/register [post]
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

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user account by ID
// @Tags users
// @Produce  json
// @Param   id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User successfully deleted"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/del/{id} [delete]
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

// GetUsers godoc
// @Summary Get all users
// @Description Get list of all registered users
// @Tags users
// @Produce  json
// @Success 200 {array} users.User "List of users"
// @Failure 202 {object} map[string]interface{} "No users found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/profile [get]
func (s *NotesAPI) getUsers(ctx *gin.Context) {
	userService := user.New(s.repo)
	allUsers, err := userService.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusAccepted, gin.H{"error": "No users"})
		return
	}
	ctx.String(http.StatusOK, "Users get: %s", allUsers)
}

// GetUserID godoc
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags users
// @Produce  json
// @Param   id path string true "User ID"
// @Success 200 {object} users.User "User details"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/profile/{id} [get]
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

// UpdateUserID godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id path string true "User ID"
// @Param   user body users.User true "User data to update"
// @Success 200 {object} map[string]interface{} "User successfully updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/upd/{id} [put]
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
