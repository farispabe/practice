package main

import "errors"

type NotifyUsersByTypeRequest struct {
	Message  string
	UserType string
}

func (sr NotifyUsersByTypeRequest) Validate() error {
	if sr.Message == "" {
		return errors.New("message should not be empty")
	}

	if sr.UserType == "" {
		return errors.New("user type should not be empty")
	}

	return nil
}

type NotifyUserResult struct {
	UserId  int64
	Message string
}

type NotifyUsersByTypeResponse struct {
	FailedNotifyUsers  []NotifyUserResult
	SuccessNotifyUsers []NotifyUserResult
}

type GetUsersByTypeRequest struct {
	UserType  string
	IsDeleted bool
	IsActive  bool
}

type User struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Score       int    `json:"score"`
}
