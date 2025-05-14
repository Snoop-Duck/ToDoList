package users

import "errors"

var (
	ErrInvalidUserCreds = errors.New("invalid creds")
	ErrUserAlredyExists = errors.New("user alredy exists")
	ErrUserNotFound = errors.New("user not found")
	ErrNoUsersAvailable = errors.New("no users avaible")
)
