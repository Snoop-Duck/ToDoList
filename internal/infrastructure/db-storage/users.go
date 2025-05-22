package dbstorage

import (
	"context"
	"time"

	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
)

func (db *DBStorage) SaveUser(user users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "INSERT INTO users(id, name, email, password) VALUES ($1, $2, $3, $4)", user.UID, user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBStorage) GetUser(login string) (users.User, error) {
	panic("unimplemented")
}
