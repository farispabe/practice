package main

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"

	"github.com/practice/sharing/util/json"
	"github.com/practice/sharing/util/validator"
)

type userServiceMocks struct {
	userRepository  *MockUserRepository
	cacheRepository *MockCacheRepository
	phoneNotifier   *MockNotifier
	emailNotifier   *MockNotifier

	jsonHandler      *json.MockHandler
	validatorHandler *validator.MockHandler
}

type NotifyUsersByTypeTestParam struct {
	ctx     context.Context
	request NotifyUsersByTypeRequest
	mocks   userServiceMocks
}

type NotifyUsersByTypeTestResult struct {
	expectedResp NotifyUsersByTypeResponse
	expectedErr  error
	shouldWait   bool
	cleanupFunc  func()
}

func NotifyUsersByType_fail_errValidate(req NotifyUsersByTypeTestParam) (resp NotifyUsersByTypeTestResult) {
	validateErr := errors.New("failed")

	validator.SetHandler(req.mocks.validatorHandler)
	req.mocks.validatorHandler.EXPECT().Validate(req.request).
		Return(validateErr)

	resp.expectedResp = NotifyUsersByTypeResponse{}
	resp.expectedErr = validateErr
	resp.shouldWait = false
	resp.cleanupFunc = func() {
		validator.SetHandler(validator.Default())
	}
	return resp
}

func NotifyUsersByType_fail_errgetActiveUsersByType(req NotifyUsersByTypeTestParam) (resp NotifyUsersByTypeTestResult) {
	validator.SetHandler(req.mocks.validatorHandler)
	req.mocks.validatorHandler.EXPECT().Validate(req.request).
		Return(nil)
	getUserCaseResp := getActiveUsersByType_fail_errUnmarshal(getActiveUsersByTypeTestParam{
		ctx:     req.ctx,
		request: req.request,
		mocks:   req.mocks,
	})

	resp.expectedResp = NotifyUsersByTypeResponse{}
	resp.expectedErr = getUserCaseResp.expectedErr
	resp.shouldWait = getUserCaseResp.shouldWait
	resp.cleanupFunc = func() {
		getUserCaseResp.cleanupFunc()
		validator.SetHandler(validator.Default())
	}
	return resp
}

func NotifyUsersByType_success(req NotifyUsersByTypeTestParam) (resp NotifyUsersByTypeTestResult) {
	validator.SetHandler(req.mocks.validatorHandler)
	req.mocks.validatorHandler.EXPECT().Validate(req.request).
		Return(nil)
	getUserCaseResp := getActiveUsersByType_succ_noErr(getActiveUsersByTypeTestParam{
		ctx:     req.ctx,
		request: req.request,
		mocks:   req.mocks,
	})
	notifyUserCaseResp := notifyUsers_succEmailNotifier(notifyUsersTestParam{
		ctx:     req.ctx,
		users:   getUserCaseResp.expectedRes,
		message: req.request.Message,
		mocks:   req.mocks,
	})

	resp.expectedResp = notifyUserCaseResp.expectedRes
	resp.expectedErr = nil
	resp.shouldWait = getUserCaseResp.shouldWait
	resp.cleanupFunc = func() {
		getUserCaseResp.cleanupFunc()
		validator.SetHandler(validator.Default())
	}
	return resp
}

func TestUserService_NotifyUsersByType(t *testing.T) {
	ctx := context.Background()
	request := NotifyUsersByTypeRequest{
		Message:  "message",
		UserType: UserTypePremium,
	}

	type args struct {
		ctx     context.Context
		request NotifyUsersByTypeRequest
	}
	tests := []struct {
		name         string
		args         args
		testCaseFunc func(req NotifyUsersByTypeTestParam) (resp NotifyUsersByTypeTestResult)
	}{
		{
			name:         "NotifyUsersByType fail, error validator.Validate",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: NotifyUsersByType_fail_errValidate,
		},
		{
			name:         "NotifyUsersByType fail, error getActiveUsersByType",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: NotifyUsersByType_fail_errgetActiveUsersByType,
		},
		{
			name:         "NotifyUsersByType success",
			args:         args{ctx: ctx, request: request},
			testCaseFunc: NotifyUsersByType_success,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := userServiceMocks{
				userRepository:   NewMockUserRepository(ctrl),
				cacheRepository:  NewMockCacheRepository(ctrl),
				emailNotifier:    NewMockNotifier(ctrl),
				phoneNotifier:    NewMockNotifier(ctrl),
				jsonHandler:      json.NewMockHandler(ctrl),
				validatorHandler: validator.NewMockHandler(ctrl),
			}
			testCaseResp := tt.testCaseFunc(NotifyUsersByTypeTestParam{
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
				phoneNotifier:   mocks.phoneNotifier,
				emailNotifier:   mocks.emailNotifier,
			}
			gotResp, err := us.NotifyUsersByType(tt.args.ctx, tt.args.request)
			if !assertErr(err, testCaseResp.expectedErr) {
				t.Errorf("NotifyUsersByType() error = %v, wantErr %v", err, testCaseResp.expectedErr)
				return
			}
			if !reflect.DeepEqual(gotResp, testCaseResp.expectedResp) {
				t.Errorf("NotifyUsersByType() gotResp = %v, want %v", gotResp, testCaseResp.expectedResp)
			}
		})
	}
}
