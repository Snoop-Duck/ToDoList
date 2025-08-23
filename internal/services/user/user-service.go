package user

import (
	"github.com/Snoop-Duck/ToDoList/internal/domain/users"

	"github.com/google/uuid"
)

type Repository interface {
	SaveUser(user users.User) error
	GetUser(login string) (users.User, error)
	DeleteUser(userID string) error
	GetAllUsers() ([]users.User, error)
	GetUserID(userID string) (users.User, error)
	UpdateUserID(userID string, user users.User) error
	Close() error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (us *Service) RegisterUser(user users.User) (string, error) {
	user.UID = uuid.New().String()

	err := us.repo.SaveUser(user)
	if err != nil {
		return ``, err
	}
	return user.UID, nil
}

func (us *Service) LoginUser(userCreds users.UserRequest) (string, error) {
	dbUser, err := us.repo.GetUser(userCreds.Email)
	if err != nil {
		return ``, err
	}

	if dbUser.Password != userCreds.Password {
		return ``, users.ErrInvalidUserCreds
	}

	return dbUser.UID, nil
}

func (us *Service) DeleteUserID(userID string) error {
	err := us.repo.DeleteUser(userID)
	if err != nil {
		return err
	}
	return nil
}

func (us *Service) GetUsers() ([]users.User, error) {
	return us.repo.GetAllUsers()
}

func (us *Service) GetUser(userID string) (users.User, error) {
	user, err := us.repo.GetUserID(userID)
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (us *Service) UpdateUser(userID string, user users.User) error {
	err := us.repo.UpdateUserID(userID, user)
	if err != nil {
		return err
	}
	return nil
}
