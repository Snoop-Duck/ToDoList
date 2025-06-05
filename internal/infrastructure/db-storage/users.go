package dbstorage

import (
	"context"
	"time"

	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
)

func (db *DBStorage) SaveUser(user users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "INSERT INTO users(uid, name, email, password) VALUES ($1, $2, $3, $4)", user.UID, user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBStorage) GetUser(login string) (users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user users.User

	row := db.db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", login)
	err := row.Scan(&user.UID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return users.User{}, err
	}

	return user, nil
}

func (db *DBStorage) DeleteUser(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStorage) GetAllUsers() ([]users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.db.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var usersSlice []users.User
	for rows.Next() {
		var user users.User
		if err := rows.Scan(&user.UID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		usersSlice = append(usersSlice, user)
	}
	return usersSlice, nil
}

func (db *DBStorage) GetUserID(userID string) (users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user users.User
	row := db.db.QueryRow(ctx, "SELECT * FROM users WHERE uid = $1", userID)
	err := row.Scan(&user.UID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (db *DBStorage) UpdateUserID(userID string, user users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "UPDATE users SET name = $1, email = $2 WHERE uid = $3", user.Name, user.Email, userID)
	if err != nil {
		return err
	}
	return nil
}
