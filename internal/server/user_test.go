package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
	"github.com/Snoop-Duck/ToDoList/internal/server/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	var srv NotesAPI

	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.POST("/login", srv.login)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	type want struct {
		resultMsg string
		status    int
	}

	type test struct {
		name    string
		request string
		method  string
		uReq    users.UserRequest
		dbUser  users.User
		repoErr error
		want    want
	}

	tests := []test{
		{
			name: "test 1: success call",
			uReq: users.UserRequest{
				Email:    "email",
				Password: "password",
			},
			dbUser: users.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			repoErr: nil,
			request: "/login",
			method:  http.MethodPost,
			want: want{
				resultMsg: "user logined: uuid-1234-55rr",
				status:    200,
			},
		},
		{
			name: "test 2: invalid creds call",
			uReq: users.UserRequest{
				Email:    "email",
				Password: "password",
			},
			dbUser: users.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "1234567",
			},
			repoErr: users.ErrInvalidUserCreds,
			request: "/login",
			method:  http.MethodPost,
			want: want{
				resultMsg: `{"error":"invalid creds"}`,
				status:    401,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockRepo.On("GetUser", tc.uReq.Email).Return(tc.dbUser, tc.repoErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpTest.URL + tc.request

			body, err := json.Marshal(tc.uReq)
			assert.NoError(t, err)
			req.Body = body

			resp, err := req.Send()
			assert.NoError(t, err)

			respBody := string(resp.Body())

			assert.Equal(t, tc.want.status, resp.StatusCode())
			assert.Equal(t, tc.want.resultMsg, respBody)
		})
	}
}

func TestReqister(t *testing.T) {
	var srv NotesAPI

	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.POST("/register", srv.register)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	type want struct {
		resultMsg string
		status    int
	}

	type test struct {
		name    string
		request string
		method  string
		uReq    users.User
		repoErr error
		want    want
	}

	tests := []test{
		{
			name:    "test 1: success call",
			request: "/register",
			method:  http.MethodPost,
			uReq: users.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			repoErr: nil,
			want: want{
				resultMsg: "user registered:",
				status:    200,
			},
		},

		{
			name:    "test 2: conflict call",
			request: "/register",
			method:  http.MethodPost,
			uReq: users.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			repoErr: users.ErrUserAlredyExists,
			want: want{
				resultMsg: `{"error":"user alredy exists"}`,
				status:    409,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockRepo.On("SaveUser", mock.MatchedBy(func(user users.User) bool {
				return user.Name == tc.uReq.Name &&
					user.Email == tc.uReq.Email &&
					user.Password == tc.uReq.Password
			})).Return(tc.repoErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpTest.URL + tc.request

			body, err := json.Marshal(tc.uReq)
			assert.NoError(t, err)
			req.Body = body

			resp, err := req.Send()
			assert.NoError(t, err)

			respBody := string(resp.Body())

			assert.Equal(t, tc.want.status, resp.StatusCode())
			assert.Contains(t, respBody, tc.want.resultMsg)
		})
	}
}

func TestGetUserIDHandler(t *testing.T) {
	testRouter := gin.New()
	srv := NotesAPI{}
	testRouter.GET("/profile/:id", srv.getUserID)
	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	testUser := users.User{
		UID:   "123",
		Name:  "Test User",
		Email: "test@example.com",
	}

	tests := []struct {
		name     string
		userID   string
		mockUser users.User
		mockErr  error
		wantCode int
	}{
		{
			name:     "successful get",
			userID:   "123",
			mockUser: testUser,
			wantCode: http.StatusOK,
		},
		{
			name:     "user not found",
			userID:   "456",
			mockErr:  users.ErrUserNotFound,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetUserID", tt.userID).Return(tt.mockUser, tt.mockErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = httpTest.URL + "/profile/" + tt.userID

			resp, _ := req.Send()
			assert.Equal(t, tt.wantCode, resp.StatusCode())
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	testRouter := gin.New()
	srv := NotesAPI{}
	testRouter.DELETE("/del/:id", srv.deleteUser)
	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	tests := []struct {
		name     string
		userID   string
		mockErr  error
		wantCode int
	}{
		{
			name:     "successful delete",
			userID:   "123",
			wantCode: http.StatusOK,
		},
		{
			name:     "user not found",
			userID:   "456",
			mockErr:  users.ErrUserNotFound,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("DeleteUser", tt.userID).Return(tt.mockErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = http.MethodDelete
			req.URL = httpTest.URL + "/del/" + tt.userID

			resp, _ := req.Send()
			assert.Equal(t, tt.wantCode, resp.StatusCode())
		})
	}
}

func TestGetUsersHandler(t *testing.T) {
	testRouter := gin.New()
	srv := NotesAPI{}
	testRouter.GET("/profile", srv.getUsers)
	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	testUsers := []users.User{
		{
			UID:   "1",
			Name:  "User 1",
			Email: "user1@example.com",
		},
		{
			UID:   "2",
			Name:  "User 2",
			Email: "user2@example.com",
		},
	}

	tests := []struct {
		name      string
		mockUsers []users.User
		mockErr   error
		wantCode  int
	}{
		{
			name:      "successful get all",
			mockUsers: testUsers,
			wantCode:  http.StatusOK,
		},
		{
			name:     "no users",
			mockErr:  users.ErrNoUsersAvailable,
			wantCode: http.StatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetAllUsers").Return(tt.mockUsers, tt.mockErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = httpTest.URL + "/profile"

			resp, _ := req.Send()
			assert.Equal(t, tt.wantCode, resp.StatusCode())
		})
	}
}

func TestGetUserID(t *testing.T) {
	testUser := users.User{
		UID:      "123",
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}

	tests := []struct {
		name     string
		userID   string
		mockUser users.User
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "successful get by ID",
			userID:   "123",
			mockUser: testUser,
			wantErr:  false,
		},
		{
			name:    "user not found",
			userID:  "456",
			mockErr: users.ErrUserNotFound,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetUserID", tt.userID).Return(tt.mockUser, tt.mockErr)

			user, err := mockRepo.GetUserID(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.mockErr != nil {
					assert.ErrorIs(t, err, tt.mockErr)
				}
				assert.Equal(t, users.User{}, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUser, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUserIDHandler(t *testing.T) {
	testRouter := gin.New()
	srv := NotesAPI{}
	testRouter.PUT("/upd/:id", srv.updateUserID)
	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	tests := []struct {
		name     string
		userID   string
		user     users.User
		mockErr  error
		wantCode int
	}{
		{
			name:   "successful update",
			userID: "123",
			user: users.User{
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "user not found",
			userID:   "456",
			mockErr:  users.ErrUserNotFound,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("UpdateUserID", tt.userID, mock.Anything).Return(tt.mockErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = http.MethodPut
			req.URL = httpTest.URL + "/upd/" + tt.userID

			if tt.user.Name != "" {
				body, _ := json.Marshal(tt.user)
				req.Body = body
			}

			resp, _ := req.Send()
			assert.Equal(t, tt.wantCode, resp.StatusCode())
		})
	}
}

func BenchmarkLogin(b *testing.B) {
	var srv NotesAPI

	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()
	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.POST("/login", srv.login)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	uReq := users.UserRequest{
		Email:    "email",
		Password: "password",
	}
	dbUser := users.User{
		UID:      "uuid-1234-55rr",
		Name:     "John Doe",
		Email:    "email",
		Password: "password",
	}

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("GetUser", uReq.Email).Return(dbUser, nil)
	srv.repo = mockRepo

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = httpTest.URL + "/login"

	body, err := json.Marshal(uReq)
	assert.NoError(b, err)
	req.Body = body

	b.ResetTimer()
	for range b.N {
		req.Send()
	}
}

func BenchmarkRegister(b *testing.B) {
	var srv NotesAPI

	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()
	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.POST("/register", srv.register)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	uReq := users.User{
		UID:      "uuid-1234-55rr",
		Name:     "John Doe",
		Email:    "email",
		Password: "password",
	}

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("SaveUser", mock.Anything).Return(nil)
	srv.repo = mockRepo

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = httpTest.URL + "/register"

	body, err := json.Marshal(uReq)
	assert.NoError(b, err)
	req.Body = body

	b.ResetTimer()
	for range b.N {
		req.Send()
	}
}

func BenchmarkGetUserID(b *testing.B) {
	var srv NotesAPI

	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()
	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.GET("/profile/:id", srv.getUserID)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	testUser := users.User{
		UID:   "123",
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("GetUserID", "123").Return(testUser, nil)
	srv.repo = mockRepo

	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = httpTest.URL + "/profile/123"

	b.ResetTimer()
	for range b.N {
		req.Send()
	}
}

func BenchmarkDeleteUser(b *testing.B) {
	var srv NotesAPI

	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()
	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.DELETE("/del/:id", srv.deleteUser)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("DeleteUser", "123").Return(nil)
	srv.repo = mockRepo

	req := resty.New().R()
	req.Method = http.MethodDelete
	req.URL = httpTest.URL + "/del/123"

	b.ResetTimer()
	for range b.N {
		req.Send()
	}
}

func BenchmarkGetUsers(b *testing.B) {
	var srv NotesAPI

	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()
	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.GET("/profile", srv.getUsers)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	testUsers := []users.User{
		{
			UID:   "1",
			Name:  "User 1",
			Email: "user1@example.com",
		},
		{
			UID:   "2",
			Name:  "User 2",
			Email: "user2@example.com",
		},
	}

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("GetAllUsers").Return(testUsers, nil)
	srv.repo = mockRepo

	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = httpTest.URL + "/profile"

	b.ResetTimer()
	for range b.N {
		req.Send()
	}
}

func BenchmarkUpdateUserID(b *testing.B) {
	var srv NotesAPI

	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()
	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.PUT("/upd/:id", srv.updateUserID)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	uReq := users.User{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("UpdateUserID", "123", mock.Anything).Return(nil)
	srv.repo = mockRepo

	req := resty.New().R()
	req.Method = http.MethodPut
	req.URL = httpTest.URL + "/upd/123"

	body, err := json.Marshal(uReq)
	assert.NoError(b, err)
	req.Body = body

	b.ResetTimer()
	for range b.N {
		req.Send()
	}
}
