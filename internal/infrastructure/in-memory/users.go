package inmemory

import (
	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
)

func (im *InMemoryUsers) SaveUser(user users.User) error {
	for _, us := range im.userStorage {
		if us.Email == user.Email {
			return users.ErrUserAlredyExists
		}
	}

	im.userStorage[user.UID] = user
	return nil
}

func (im *InMemoryUsers) GetUser(login string) (users.User, error) {
	for _, us := range im.userStorage {
		if us.Email == login {
			return us, nil
		}
	}
	return emtyUser, users.ErrUserNotFound
}

func (im *InMemoryUsers) DeleteUser(userID string) error {
	if _, ok := im.userStorage[userID]; !ok {
		return users.ErrUserNotFound
	}
	delete(im.userStorage, userID)
	return nil
}

func (im *InMemoryUsers) GetAllUsers() ([]users.User, error) {
	if len(im.userStorage) == 0 {
		return nil, users.ErrNoUsersAvailable
	}

	usersSlice := make([]users.User, 0, len(im.userStorage))
	for _, note := range im.userStorage {
		usersSlice = append(usersSlice, note)
	}

	return usersSlice, nil
}

func (im *InMemoryUsers) GetUserID(userID string) (users.User, error) {
	user, ok := im.userStorage[userID]
	if !ok {
		return users.User{}, users.ErrUserNotFound
	}
	return user, nil
}

func (im *InMemoryUsers) UpdateUserID(userID string, user users.User) error {
	if _, ok := im.userStorage[userID]; !ok {
		return users.ErrUserNotFound
	}
	im.userStorage[userID] = user
	return nil
}

func (im *InMemoryUsers) Close() error { return nil }
