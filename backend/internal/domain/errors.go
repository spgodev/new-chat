package domain

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameAlreadyExist  = errors.New("username already exist")
	ErrEmailAlreadyUsed      = errors.New("email already used")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrInvalidUsernameFormat = errors.New("invalid username format")
	ErrShortPassword         = errors.New("minimal size of password = 8")
	ErrTokenNotFound         = errors.New("token not found")
	ErrSomethingWrong        = errors.New("something went wrong")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrAccessDenied          = errors.New("access denied")
	ErrCantCreateChat        = errors.New("can't create chat")
	ErrChatAlreadyExist      = errors.New("chat already exist")
	ErrChatNotFound          = errors.New("chat not found")
	ErrEmptyMessage          = errors.New("you can't send empty message")
)
