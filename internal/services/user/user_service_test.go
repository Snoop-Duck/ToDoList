package user

import (
	"errors"
	"testing"

	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
	"github.com/Snoop-Duck/ToDoList/internal/services/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginUser(t *testing.T) {
	type want struct {
		userID string
		err    error
	}

	type test struct {
		name    string
		user    users.User
		userReq users.UserRequest
		want    want
	}

	tests := []test{
		{
			name: "test 1: success call",
			userReq: users.UserRequest{
				Email:    "email",
				Password: "password",
			},
			user: users.User{
				UID:      "uuid",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			want: want{
				userID: "uuid",
				err:    nil,
			},
		},
		{
			name: "test 2: fail call",
			userReq: users.UserRequest{
				Email:    "email",
				Password: "password",
			},
			user: users.User{
				UID:      "uuid",
				Name:     "John Doe",
				Email:    "email",
				Password: "password1234",
			},
			want: want{
				userID: "",
				err:    users.ErrInvalidUserCreds,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("GetUser", tc.userReq.Email).Return(tc.user, nil)

			testUserService := New(repoMock)

			userID, err := testUserService.LoginUser(tc.userReq)
			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tc.want.userID, userID)
		})
	}
}

func TestRegisterUser(t *testing.T) {
	type want struct {
		err error
	}

	type test struct {
		name string
		user users.User
		want want
	}

	tests := []test{
		{
			name: "test 1: success call",
			user: users.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "test 2: unique error case",
			user: users.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			want: want{
				err: users.ErrUserAlredyExists,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("SaveUser", mock.MatchedBy(func(user users.User) bool {
				return user.Name == tc.user.Name &&
					user.Email == tc.user.Email &&
					user.Password == tc.user.Password
			})).Return(tc.want.err)

			testUserService := New(repoMock)

			userID, err := testUserService.RegisterUser(tc.user)
			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			assert.NoError(t, err)

			assert.NotEmpty(t, userID)
		})
	}
}

func TestDeleteUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mockErr error
		wantErr error
	}{
		{
			name:    "successful deletion",
			userID:  "user123",
			wantErr: nil,
		},
		{
			name:    "user not found",
			userID:  "nonexistent",
			mockErr: users.ErrUserNotFound,
			wantErr: users.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("DeleteUser", tt.userID).Return(tt.mockErr)

			service := New(repoMock)
			err := service.DeleteUserID(tt.userID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	mockUsers := []users.User{
		{
			UID:      "user1",
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "pass1",
		},
		{
			UID:      "user2",
			Name:     "Jane Doe",
			Email:    "jane@example.com",
			Password: "pass2",
		},
	}

	tests := []struct {
		name     string
		mockResp []users.User
		mockErr  error
		wantErr  error
	}{
		{
			name:     "successful get all users",
			mockResp: mockUsers,
			wantErr:  nil,
		},
		{
			name:     "database error",
			mockResp: nil,
			mockErr:  errors.New("database error"),
			wantErr:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("GetAllUsers").Return(tt.mockResp, tt.mockErr)

			service := New(repoMock)
			users, err := service.GetUsers()

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResp, users)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	mockUser := users.User{
		UID:      "user123",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	tests := []struct {
		name     string
		userID   string
		mockResp users.User
		mockErr  error
		wantErr  error
	}{
		{
			name:     "successful get user",
			userID:   "user123",
			mockResp: mockUser,
			wantErr:  nil,
		},
		{
			name:     "user not found",
			userID:   "nonexistent",
			mockResp: users.User{},
			mockErr:  users.ErrUserNotFound,
			wantErr:  users.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("GetUserID", tt.userID).Return(tt.mockResp, tt.mockErr)

			service := New(repoMock)
			user, err := service.GetUser(tt.userID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, users.User{}, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResp, user)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	updatedUser := users.User{
		UID:      "user123",
		Name:     "New Name",
		Email:    "new@example.com",
		Password: "newpass",
	}

	tests := []struct {
		name    string
		userID  string
		user    users.User
		mockErr error
		wantErr error
	}{
		{
			name:    "successful update",
			userID:  "user123",
			user:    updatedUser,
			wantErr: nil,
		},
		{
			name:    "user not found",
			userID:  "nonexistent",
			user:    updatedUser,
			mockErr: users.ErrUserNotFound,
			wantErr: users.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("UpdateUserID", tt.userID, tt.user).Return(tt.mockErr)

			service := New(repoMock)
			err := service.UpdateUser(tt.userID, tt.user)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
