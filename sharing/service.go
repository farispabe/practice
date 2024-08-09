package main

import (
	"context"
	"log"

	"github.com/practice/sharing/util/custerror"
	"github.com/practice/sharing/util/json"
	"github.com/practice/sharing/util/validator"
)

type UserService struct {
	userRepository  UserRepository
	cacheRepository CacheRepository
	phoneNotifier   Notifier
	emailNotifier   Notifier
}

// NotifyUsersByType notifies a Message to users identified by UserType
func (us *UserService) NotifyUsersByType(ctx context.Context, request NotifyUsersByTypeRequest) (resp NotifyUsersByTypeResponse, err error) {
	// validate request
	if err = validator.Validate(request); err != nil {
		return resp, err
	}

	// get users
	var users []User
	users, err = us.getActiveUsersByType(ctx, request)
	if err != nil {
		return resp, err
	}

	// notify users
	resp = us.notifyUsers(ctx, users, request.Message)

	return resp, nil
}

// getActiveUsersByType gets active users by type from cache or database if not exist in cache
func (us *UserService) getActiveUsersByType(ctx context.Context, request NotifyUsersByTypeRequest) (users []User, err error) {
	// get from cache
	var usersJson string
	cacheKey := getCacheKeyActiveUsersByType(request.UserType)
	usersJson, err = us.cacheRepository.Get(ctx, cacheKey)
	if err == nil {
		err = json.Unmarshal([]byte(usersJson), &users)
		if err != nil {
			return nil, custerror.NewInternal(err.Error())
		} else {
			return users, nil
		}
	}

	// get from database
	getUsersReq := createGetActiveUsersByTypeRequest(request)
	users, err = us.userRepository.GetByTypeAndState(ctx, getUsersReq)
	if err != nil || users == nil {
		if err == nil {
			return nil, custerror.NewNotFound("users not found")
		}
		return nil, custerror.NewInternal(err.Error())
	}

	go func() {
		bytesData, errMarshal := json.Marshal(users)
		if errMarshal != nil {
			log.Println(errMarshal.Error(), "users", users)
			return
		}
		if errSet := us.cacheRepository.Set(ctx, cacheKey, string(bytesData), CacheTtlActiveUserByType); errSet != nil {
			log.Println(errSet.Error(), "key", cacheKey)
		}
	}()

	return users, nil
}

// notifyUsers notifies a message to users by phone or email based on their score
func (us *UserService) notifyUsers(ctx context.Context, users []User, message string) (resp NotifyUsersByTypeResponse) {
	for _, user := range users {
		var err error
		if user.Score > 50 {
			err = us.emailNotifier.Notify(ctx, user.Email, message)
		} else {
			err = us.phoneNotifier.Notify(ctx, user.PhoneNumber, message)
		}
		if err != nil {
			resp.FailedNotifyUsers = append(resp.FailedNotifyUsers, NotifyUserResult{
				UserId:  user.Id,
				Message: err.Error(),
			})
		} else {
			resp.SuccessNotifyUsers = append(resp.SuccessNotifyUsers, NotifyUserResult{
				UserId: user.Id,
			})
		}
	}
	return resp
}

func createGetActiveUsersByTypeRequest(request NotifyUsersByTypeRequest) GetUsersByTypeRequest {
	return GetUsersByTypeRequest{
		UserType:  request.UserType,
		IsDeleted: false,
		IsActive:  true,
	}
}
