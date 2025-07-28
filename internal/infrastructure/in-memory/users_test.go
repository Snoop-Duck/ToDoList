package inmemory

import (
	"testing"

	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryUsers(t *testing.T) {
	im := NewUsers()

	user1 := users.User{
		UID:   "11111111-1111-1111-1111-111111111111",
		Email: "user1@example.com",
		Name:  "User One",
	}

	user2 := users.User{
		UID:   "22222222-2222-2222-2222-222222222222",
		Email: "user2@example.com",
		Name:  "User Two",
	}

	t.Run("SaveUser - new user", func(t *testing.T) {
		err := im.SaveUser(user1)
		assert.NoError(t, err)
		assert.Equal(t, user1, im.userStorage[user1.UID])
	})

	t.Run("SaveUser - duplicate email", func(t *testing.T) {
		duplicateUser := users.User{
			UID:   uuid.New().String(),
			Email: user1.Email,
			Name:  "Duplicate User",
		}

		err := im.SaveUser(duplicateUser)
		assert.ErrorIs(t, err, users.ErrUserAlredyExists)
	})

	t.Run("GetUser - existing user", func(t *testing.T) {
		err := im.SaveUser(user2)
		assert.NoError(t, err)

		foundUser, err := im.GetUser(user2.Email)
		assert.NoError(t, err)
		assert.Equal(t, user2, foundUser)
	})

	t.Run("GetUser - non-existent user", func(t *testing.T) {
		_, err := im.GetUser("nonexistent@example.com")
		assert.ErrorIs(t, err, users.ErrUserNotFound)
	})

	t.Run("GetUserID - existing user", func(t *testing.T) {
		foundUser, err := im.GetUserID(user1.UID)
		assert.NoError(t, err)
		assert.Equal(t, user1, foundUser)
	})

	t.Run("GetUserID - non-existent user", func(t *testing.T) {
		_, err := im.GetUserID("non-existent-uid")
		assert.ErrorIs(t, err, users.ErrUserNotFound)
	})

	t.Run("DeleteUser - existing user", func(t *testing.T) {
		err := im.DeleteUser(user1.UID)
		assert.NoError(t, err)
		_, exists := im.userStorage[user1.UID]
		assert.False(t, exists)
	})

	t.Run("DeleteUser - non-existent user", func(t *testing.T) {
		err := im.DeleteUser("non-existent-uid")
		assert.ErrorIs(t, err, users.ErrUserNotFound)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		im.userStorage = make(map[string]users.User)
		im.SaveUser(user1)
		im.SaveUser(user2)

		usersList, err := im.GetAllUsers()
		assert.NoError(t, err)
		assert.Len(t, usersList, 2)
		assert.Contains(t, usersList, user1)
		assert.Contains(t, usersList, user2)
	})

	t.Run("GetAllUsers - empty storage", func(t *testing.T) {
		im.userStorage = make(map[string]users.User)
		usersList, err := im.GetAllUsers()
		assert.ErrorIs(t, err, users.ErrNoUsersAvailable)
		assert.Nil(t, usersList)
	})

	t.Run("UpdateUserID - existing user", func(t *testing.T) {
		err := im.SaveUser(user2)
		assert.NoError(t, err)

		updatedUser := user2
		updatedUser.Name = "Updated User Two"

		err = im.UpdateUserID(user2.UID, updatedUser)
		assert.NoError(t, err)

		storedUser := im.userStorage[user2.UID]
		assert.Equal(t, updatedUser, storedUser)
	})

	t.Run("UpdateUserID - non-existent user", func(t *testing.T) {
		err := im.UpdateUserID("non-existent-uid", user1)
		assert.ErrorIs(t, err, users.ErrUserNotFound)
	})

	t.Run("Close", func(t *testing.T) {
		err := im.Close()
		assert.NoError(t, err)
	})
}
