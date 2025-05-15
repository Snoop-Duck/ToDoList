package inmemory

import "github.com/Snoop-Duck/ToDoList/internal/domain/users"

var emtyUser = users.User{}

type InMemory struct {
	userStorage map[string]users.User
}

func New() *InMemory {
	return &InMemory{
		userStorage: make(map[string]users.User),
	}
}

func (im *InMemory) SaveUser(user users.User) error {
	for _, us := range im.userStorage {
		if us.Email == user.Email {
			return users.ErrUserAlredyExists
		}
	}

	im.userStorage[user.UID] = user
	return nil
}

func (im *InMemory) GetUser(login string) (users.User, error) {
	for _, us := range im.userStorage {
		if us.Email == login {
			return us, nil
		}
	}
	return emtyUser, users.ErrUserNotFound
}

func (im *InMemory) DeleteUser(userID string) error {
	if _, ok := im.userStorage[userID]; !ok {
		return users.ErrUserNotFound
	}
	delete(im.userStorage, userID)
	return nil
}

func (im *InMemory) GetAllUsers() (map[string]users.User, error) {
	if len(im.userStorage) == 0 {
		return nil, users.ErrNoUsersAvailable
	}
	return im.userStorage, nil
}

func (im *InMemory) GetUserID(userID string) (users.User, error) {
	user, ok := im.userStorage[userID]
	if !ok {
		return users.User{}, users.ErrUserNotFound
	}
	return user, nil
}

func (im *InMemory) UpdateUserID(userID string, user users.User) error {
	if _, ok := im.userStorage[userID]; !ok {
		return users.ErrUserNotFound
	}
	im.userStorage[userID] = user
	return nil 
}