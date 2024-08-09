package main

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/practice/sharing/util/custerror"
	"github.com/practice/sharing/util/json"
)

type getActiveUsersByTypeTestParam struct {
	ctx     context.Context
	request NotifyUsersByTypeRequest
	mocks   userServiceMocks
}

type getActiveUsersByTypeTestResult struct {
	expectedRes []User
	expectedErr error
	shouldWait  bool
	cleanupFunc func()
}

// getActiveUsersByType_fail_errUnmarshal defines failure, caused by error json.Unmarshal
// (when trying to get from cache)
func getActiveUsersByType_fail_errUnmarshal(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult) {
	cacheKey := getCacheKeyActiveUsersByType(req.request.UserType)
	cacheGetResp := ""
	unmarshalErr := errors.New("failed")

	req.mocks.cacheRepository.EXPECT().Get(req.ctx, cacheKey).
		Return(cacheGetResp, nil)
	json.SetHandler(req.mocks.jsonHandler)
	req.mocks.jsonHandler.EXPECT().Unmarshal([]byte(cacheGetResp), &result.expectedRes).
		Return(unmarshalErr)

	result.expectedRes = nil
	result.expectedErr = custerror.NewInternal(unmarshalErr.Error())
	result.shouldWait = false
	result.cleanupFunc = func() {
		json.SetHandler(json.Default())
	}
	return result
}

// getActiveUsersByType_succ_noErr defines success without failure
// (when trying to get from cache)
func getActiveUsersByType_succ_noErr(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult) {
	cacheKey := getCacheKeyActiveUsersByType(req.request.UserType)
	cacheGetResp := `[{"id":1, "name": "name", "type": "premium", "phone_number": "088888888", "email": "email@test.mail", "score": 60}]`
	resp := []User{
		{
			Id:          1,
			Name:        "name",
			Type:        UserTypePremium,
			PhoneNumber: "088888888",
			Email:       "email@test.mail",
			Score:       60,
		},
	}
	var users []User

	req.mocks.cacheRepository.EXPECT().Get(req.ctx, cacheKey).
		Return(cacheGetResp, nil)
	json.SetHandler(req.mocks.jsonHandler)
	req.mocks.jsonHandler.EXPECT().Unmarshal([]byte(cacheGetResp), &users).
		DoAndReturn(func(d []byte, r *[]User) error {
			*r = resp
			return nil
		})

	result.expectedRes = resp
	result.expectedErr = nil
	result.shouldWait = false
	result.cleanupFunc = func() {
		json.SetHandler(json.Default())
	}
	return result
}

// getActiveUsersByType_fail_errGetByTypeAndState defines failure, caused by error userRepository.GetByTypeAndState
// (when trying to get from database)
func getActiveUsersByType_fail_errGetByTypeAndState(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult) {
	cacheKey := getCacheKeyActiveUsersByType(req.request.UserType)
	getUsersReq := createGetActiveUsersByTypeRequest(req.request)
	errGetUsers := errors.New("failed")

	req.mocks.cacheRepository.EXPECT().Get(req.ctx, cacheKey).
		Return("", errors.New("failed"))
	req.mocks.userRepository.EXPECT().GetByTypeAndState(req.ctx, getUsersReq).
		Return(nil, errGetUsers)

	result.expectedRes = nil
	result.expectedErr = custerror.NewInternal(errGetUsers.Error())
	result.shouldWait = false
	result.cleanupFunc = nil
	return result
}

// getActiveUsersByType_fail_resultEmptyGetByTypeAndState defines failure, caused by empty result from userRepository.GetByTypeAndState
// (when trying to get from database)
func getActiveUsersByType_fail_resultEmptyGetByTypeAndState(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult) {
	cacheKey := getCacheKeyActiveUsersByType(req.request.UserType)
	getUsersReq := createGetActiveUsersByTypeRequest(req.request)

	req.mocks.cacheRepository.EXPECT().Get(req.ctx, cacheKey).
		Return("", errors.New("failed"))
	req.mocks.userRepository.EXPECT().GetByTypeAndState(req.ctx, getUsersReq).
		Return(nil, nil)

	result.expectedRes = nil
	result.expectedErr = custerror.NewNotFound("user not found")
	result.shouldWait = false
	result.cleanupFunc = nil
	return result
}

// getActiveUsersByType_succ_errCacheGet defines success with an error from cacheRepository.Get
// (when trying to get from database)
func getActiveUsersByType_succ_errCacheGet(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult) {
	cacheKey := getCacheKeyActiveUsersByType(req.request.UserType)
	getUsersReq := createGetActiveUsersByTypeRequest(req.request)
	resp := []User{
		{
			Id:          1,
			Name:        "name",
			Type:        UserTypePremium,
			PhoneNumber: "088888888",
			Email:       "email@test.mail",
			Score:       60,
		},
	}
	respString := `[{"id":1, "name": "name", "type": "premium", "phone_number": "088888888", "email": "email@test.mail", "score": 60}]`

	req.mocks.cacheRepository.EXPECT().Get(req.ctx, cacheKey).
		Return("", errors.New("failed"))
	req.mocks.userRepository.EXPECT().GetByTypeAndState(req.ctx, getUsersReq).
		Return(resp, nil)
	json.SetHandler(req.mocks.jsonHandler)
	req.mocks.jsonHandler.EXPECT().Marshal(resp).
		Return([]byte(respString), nil)
	req.mocks.cacheRepository.EXPECT().Set(req.ctx, cacheKey, respString, CacheTtlActiveUserByType).
		Return(nil)

	result.expectedRes = resp
	result.expectedErr = nil
	result.shouldWait = true
	result.cleanupFunc = func() {
		json.SetHandler(json.Default())
	}
	return result
}

// getActiveUsersByType_succ_errCacheGetAndMarshal defines success with error from cacheRepository.Get and json.Marshal
// (when trying to get from database)
func getActiveUsersByType_succ_errCacheGetAndMarshal(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult) {
	cacheKey := getCacheKeyActiveUsersByType(req.request.UserType)
	getUsersReq := createGetActiveUsersByTypeRequest(req.request)
	resp := []User{
		{
			Id:          1,
			Name:        "name",
			Type:        UserTypePremium,
			PhoneNumber: "088888888",
			Email:       "email@test.mail",
			Score:       60,
		},
	}

	req.mocks.cacheRepository.EXPECT().Get(req.ctx, cacheKey).
		Return("", errors.New("failed"))
	req.mocks.userRepository.EXPECT().GetByTypeAndState(req.ctx, getUsersReq).
		Return(resp, nil)
	json.SetHandler(req.mocks.jsonHandler)
	req.mocks.jsonHandler.EXPECT().Marshal(resp).
		Return(nil, errors.New("failed"))

	result.expectedRes = resp
	result.expectedErr = nil
	result.shouldWait = true
	result.cleanupFunc = func() {
		json.SetHandler(json.Default())
	}
	return result
}

// getActiveUsersByType_succ_errCacheGetAndCacheSet defines success with error from cacheRepository.Get and cacheRepository.Set
// // (when trying to get from database)
func getActiveUsersByType_succ_errCacheGetAndCacheSet(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult) {
	cacheKey := getCacheKeyActiveUsersByType(req.request.UserType)
	getUsersReq := createGetActiveUsersByTypeRequest(req.request)
	resp := []User{
		{
			Id:          1,
			Name:        "name",
			Type:        UserTypePremium,
			PhoneNumber: "088888888",
			Email:       "email@test.mail",
			Score:       60,
		},
	}
	respString := `[{"id":1, "name": "name", "type": "premium", "phone_number": "088888888", "email": "email@test.mail", "score": 60}]`

	req.mocks.cacheRepository.EXPECT().Get(req.ctx, cacheKey).
		Return("", errors.New("failed"))
	req.mocks.userRepository.EXPECT().GetByTypeAndState(req.ctx, getUsersReq).
		Return(resp, nil)
	json.SetHandler(req.mocks.jsonHandler)
	req.mocks.jsonHandler.EXPECT().Marshal(resp).
		Return([]byte(respString), nil)
	req.mocks.cacheRepository.EXPECT().Set(req.ctx, cacheKey, respString, CacheTtlActiveUserByType).
		Return(errors.New("failed"))

	result.expectedRes = resp
	result.expectedErr = nil
	result.shouldWait = true
	result.cleanupFunc = func() {
		json.SetHandler(json.Default())
	}
	return result
}

func TestUserService_getActiveUsersByType(t *testing.T) {
	ctx := context.Background()
	request := NotifyUsersByTypeRequest{
		Message:  "test",
		UserType: UserTypePremium,
	}

	type args struct {
		ctx     context.Context
		request NotifyUsersByTypeRequest
	}
	tests := []struct {
		name         string
		args         args
		testCaseFunc func(req getActiveUsersByTypeTestParam) (result getActiveUsersByTypeTestResult)
	}{
		{
			name:         "getActiveUsersByType fail, error json.Unmarshal",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: getActiveUsersByType_fail_errUnmarshal,
		},
		{
			name:         "getActiveUsersByType success no error",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: getActiveUsersByType_succ_noErr,
		},
		{
			name:         "getActiveUsersByType fail, error userRepository.GetByTypeAndState",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: getActiveUsersByType_fail_errGetByTypeAndState,
		},
		{
			name:         "getActiveUsersByType fail, empty result userRepository.GetByTypeAndState",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: getActiveUsersByType_fail_resultEmptyGetByTypeAndState,
		},
		{
			name:         "getActiveUsersByType success, error cacheRepository.Get",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: getActiveUsersByType_succ_errCacheGet,
		},
		{
			name:         "getActiveUsersByType success, error cacheRepository.Get and json.Marshal",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: getActiveUsersByType_succ_errCacheGetAndMarshal,
		},
		{
			name:         "getActiveUsersByType success, error cacheRepository.Get and cacheRepository.Set",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: getActiveUsersByType_succ_errCacheGetAndCacheSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := userServiceMocks{
				userRepository:  NewMockUserRepository(ctrl),
				cacheRepository: NewMockCacheRepository(ctrl),
				jsonHandler:     json.NewMockHandler(ctrl),
			}
			testCaseResp := tt.testCaseFunc(getActiveUsersByTypeTestParam{
				ctx:     tt.args.ctx,
				request: tt.args.request,
				mocks:   mocks,
			})
			if testCaseResp.cleanupFunc != nil {
				defer testCaseResp.cleanupFunc()
			}

			us := &UserService{
				userRepository:  mocks.userRepository,
				cacheRepository: mocks.cacheRepository,
			}
			gotUsers, err := us.getActiveUsersByType(tt.args.ctx, tt.args.request)
			if testCaseResp.shouldWait {
				time.Sleep(1 * time.Second)
			}
			if !assertErr(err, testCaseResp.expectedErr) {
				t.Errorf("getActiveUsersByType() error = %v, wantErr %v", err, testCaseResp.expectedErr)
				return
			}
			if !reflect.DeepEqual(gotUsers, testCaseResp.expectedRes) {
				t.Errorf("getActiveUsersByType() gotUsers = %v, want %v", gotUsers, testCaseResp.expectedRes)
			}
		})
	}
}

func assertErr(want error, expected error) bool {
	if want == nil && expected == nil {
		return true
	} else if want == nil {
		return false
	} else if expected == nil {
		return false
	}
	var badRequest *custerror.BadRequest
	var notFound *custerror.NotFound
	var internal *custerror.Internal

	if errors.As(want, &badRequest) {
		if errors.As(expected, &badRequest) {
			return true
		} else {
			return false
		}
	} else if errors.As(want, &notFound) {
		if errors.As(expected, &notFound) {
			return true
		} else {
			return false
		}
	} else if errors.As(want, &internal) {
		if errors.As(expected, &internal) {
			return true
		} else {
			return false
		}
	}

	if want.Error() == expected.Error() {
		return true
	}
	return false
}
