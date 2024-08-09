package main

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

type notifyUsersTestParam struct {
	ctx     context.Context
	users   []User
	message string
	mocks   userServiceMocks
}

type notifyUsersTestResult struct {
	expectedRes NotifyUsersByTypeResponse
}

var (
	user_scoreGreater50_fail = User{
		Id:    1,
		Name:  "user_scoreGreater50_fail",
		Score: 60,
		Email: "user_scoreGreater50_fail",
	}
	user_scoreGreater50_succ = User{
		Id:    2,
		Name:  "user_scoreGreater50_succ",
		Score: 60,
		Email: "user_scoreGreater50_succ",
	}
	user_score50_fail = User{
		Id:          3,
		Name:        "user_score50_fail",
		Score:       50,
		PhoneNumber: "user_score50_fail",
	}
	user_score50_succ = User{
		Id:          4,
		Name:        "user_score50_succ",
		Score:       50,
		PhoneNumber: "user_score50_succ",
	}
	user_scoreLesser50_fail = User{
		Id:          5,
		Name:        "user_scoreLesser50_fail",
		Score:       40,
		PhoneNumber: "user_scoreLesser50_fail",
	}
	user_scoreLesser50_succ = User{
		Id:          6,
		Name:        "user_scoreLesser50_succ",
		Score:       40,
		PhoneNumber: "user_scoreLesser50_succ",
	}
)

// notifyUsers_1scoreGreater50Fail defines resp with 1 failure calling emailNotifier when score > 50
func notifyUsers_1scoreGreater50Fail(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	emailNotifyErr := errors.New("failed")

	req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user_scoreGreater50_fail.Email, req.message).
		Return(emailNotifyErr)

	resp.expectedRes.FailedNotifyUsers = []NotifyUserResult{
		{
			UserId:  user_scoreGreater50_fail.Id,
			Message: emailNotifyErr.Error(),
		},
	}
	return resp
}

// notifyUsers_1scoreGreater50Succ defines resp with 1 success calling emailNotifier when score > 50
func notifyUsers_1scoreGreater50Succ(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user_scoreGreater50_succ.Email, req.message).
		Return(nil)

	resp.expectedRes.SuccessNotifyUsers = []NotifyUserResult{
		{
			UserId: user_scoreGreater50_succ.Id,
		},
	}
	return resp
}

// notifyUsers_1score50Fail defines resp with 1 failure calling phoneNotifier when score = 50
func notifyUsers_1score50Fail(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	phoneNotifierErr := errors.New("failed")

	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_score50_fail.PhoneNumber, req.message).
		Return(phoneNotifierErr)

	resp.expectedRes.FailedNotifyUsers = []NotifyUserResult{
		{
			UserId:  user_score50_fail.Id,
			Message: phoneNotifierErr.Error(),
		},
	}
	return resp
}

// notifyUsers_1score50Succ defines resp with 1 success calling phoneNotifier when score = 50
func notifyUsers_1score50Succ(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_score50_succ.PhoneNumber, req.message).
		Return(nil)

	resp.expectedRes.SuccessNotifyUsers = []NotifyUserResult{
		{
			UserId: user_score50_succ.Id,
		},
	}
	return resp
}

// notifyUsers_1scoreLesser50Fail defines resp with 1 failure calling phoneNotifier when score < 50
func notifyUsers_1scoreLesser50Fail(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	phoneNotifierErr := errors.New("failed")

	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_scoreLesser50_fail.PhoneNumber, req.message).
		Return(phoneNotifierErr)

	resp.expectedRes.FailedNotifyUsers = []NotifyUserResult{
		{
			UserId:  user_scoreLesser50_fail.Id,
			Message: phoneNotifierErr.Error(),
		},
	}
	return resp
}

// notifyUsers_1scoreLesser50Succ defines resp with 1 success calling phoneNotifier when score < 50
func notifyUsers_1scoreLesser50Succ(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_scoreLesser50_succ.PhoneNumber, req.message).
		Return(nil)

	resp.expectedRes.SuccessNotifyUsers = []NotifyUserResult{
		{
			UserId: user_scoreLesser50_succ.Id,
		},
	}
	return resp
}

// notifyUsers_allFail_1scoreGreater50_1score50_1scoreLesser50 defines resp with all failure on 1 score > 50, 1 score = 50, & 1 score < 50
func notifyUsers_allFail_1scoreGreater50_1score50_1scoreLesser50(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	emailNotifierErr := errors.New("emailNotifier failed")
	phoneNotifierErr := errors.New("phoneNotifier failed")

	req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user_scoreGreater50_fail.Email, req.message).
		Return(emailNotifierErr)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_score50_fail.PhoneNumber, req.message).
		Return(phoneNotifierErr)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_scoreLesser50_fail.PhoneNumber, req.message).
		Return(phoneNotifierErr)

	resp.expectedRes.FailedNotifyUsers = []NotifyUserResult{
		{
			UserId:  user_scoreGreater50_fail.Id,
			Message: emailNotifierErr.Error(),
		},
		{
			UserId:  user_score50_fail.Id,
			Message: phoneNotifierErr.Error(),
		},
		{
			UserId:  user_scoreLesser50_fail.Id,
			Message: phoneNotifierErr.Error(),
		},
	}
	return resp
}

// notifyUsers_allSucc_1scoreGreater50_1score50_1scoreLesser50 defines resp with all success on 1 score > 50, 1 score = 50, & 1 score < 50
func notifyUsers_allSucc_1scoreGreater50_1score50_1scoreLesser50(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user_scoreGreater50_succ.Email, req.message).
		Return(nil)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_score50_succ.PhoneNumber, req.message).
		Return(nil)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_scoreLesser50_succ.PhoneNumber, req.message).
		Return(nil)

	resp.expectedRes.SuccessNotifyUsers = []NotifyUserResult{
		{
			UserId: user_scoreGreater50_succ.Id,
		},
		{
			UserId: user_score50_succ.Id,
		},
		{
			UserId: user_scoreLesser50_succ.Id,
		},
	}
	return resp
}

// notifyUsers_fail_1Greater50_1score50_1Lesser50_succ_1Greater50_1score50_1Lesser50 defines resp with failure on 1 score > 50, 1 score = 50, & 1 score < 50
// and success on 1 score > 50, 1 score = 50, & 1 score < 50
func notifyUsers_fail_1Greater50_1score50_1Lesser50_succ_1Greater50_1score50_1Lesser50(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	emailNotifierErr := errors.New("emailNotifier failed")
	phoneNotifierErr := errors.New("phoneNotifier failed")

	req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user_scoreGreater50_fail.Email, req.message).
		Return(emailNotifierErr)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_score50_fail.PhoneNumber, req.message).
		Return(phoneNotifierErr)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_scoreLesser50_fail.PhoneNumber, req.message).
		Return(phoneNotifierErr)
	req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user_scoreGreater50_succ.Email, req.message).
		Return(nil)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_score50_succ.PhoneNumber, req.message).
		Return(nil)
	req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user_scoreLesser50_succ.PhoneNumber, req.message).
		Return(nil)

	resp.expectedRes.FailedNotifyUsers = []NotifyUserResult{
		{
			UserId:  user_scoreGreater50_fail.Id,
			Message: emailNotifierErr.Error(),
		},
		{
			UserId:  user_score50_fail.Id,
			Message: phoneNotifierErr.Error(),
		},
		{
			UserId:  user_scoreLesser50_fail.Id,
			Message: phoneNotifierErr.Error(),
		},
	}
	resp.expectedRes.SuccessNotifyUsers = []NotifyUserResult{
		{
			UserId: user_scoreGreater50_succ.Id,
		},
		{
			UserId: user_score50_succ.Id,
		},
		{
			UserId: user_scoreLesser50_succ.Id,
		},
	}
	return resp
}

func notifyUsers_emptyUsers(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	resp.expectedRes = NotifyUsersByTypeResponse{}
	return resp
}

func notifyUsers_errEmailNotifier(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	emailNotifErr := errors.New("emailNotifier failed")

	for _, user := range req.users {
		req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user.Email, req.message).
			Return(emailNotifErr)

		resp.expectedRes.FailedNotifyUsers = append(resp.expectedRes.FailedNotifyUsers, NotifyUserResult{
			UserId:  user.Id,
			Message: emailNotifErr.Error(),
		})
	}

	return resp
}

func notifyUsers_succEmailNotifier(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	for _, user := range req.users {
		req.mocks.emailNotifier.EXPECT().Notify(req.ctx, user.Email, req.message).
			Return(nil)

		resp.expectedRes.SuccessNotifyUsers = append(resp.expectedRes.SuccessNotifyUsers, NotifyUserResult{
			UserId: user.Id,
		})
	}

	return resp
}

func notifyUsers_errPhoneNotifier(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	phoneNotifErr := errors.New("phoneNotifier failed")

	for _, user := range req.users {
		req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user.PhoneNumber, req.message).
			Return(phoneNotifErr)

		resp.expectedRes.FailedNotifyUsers = append(resp.expectedRes.FailedNotifyUsers, NotifyUserResult{
			UserId:  user.Id,
			Message: phoneNotifErr.Error(),
		})
	}

	return resp
}

func notifyUsers_succPhoneNotifier(req notifyUsersTestParam) (resp notifyUsersTestResult) {
	for _, user := range req.users {
		req.mocks.phoneNotifier.EXPECT().Notify(req.ctx, user.PhoneNumber, req.message).
			Return(nil)

		resp.expectedRes.SuccessNotifyUsers = append(resp.expectedRes.SuccessNotifyUsers, NotifyUserResult{
			UserId: user.Id,
		})
	}

	return resp
}

func TestUserService_notifyUsers(t *testing.T) {
	ctx := context.Background()
	message := "message"

	type args struct {
		ctx     context.Context
		users   []User
		message string
	}
	tests := []struct {
		name         string
		args         args
		testCaseFunc func(req notifyUsersTestParam) (resp notifyUsersTestResult)
	}{
		{
			name:         "notifyUsers on empty users",
			args:         args{ctx: ctx, users: []User{}, message: message},
			testCaseFunc: notifyUsers_emptyUsers,
		},
		{
			name:         "notifyUsers results 1 failure on score greater than 50",
			args:         args{ctx: ctx, users: []User{user_scoreGreater50_fail}, message: message},
			testCaseFunc: notifyUsers_errEmailNotifier,
		},
		{
			name:         "notifyUsers results 1 success on score greater than 50",
			args:         args{ctx: ctx, users: []User{user_scoreGreater50_succ}, message: message},
			testCaseFunc: notifyUsers_succEmailNotifier,
		},
		{
			name:         "notifyUsers results 1 failure on score equal 50",
			args:         args{ctx: ctx, users: []User{user_score50_fail}, message: message},
			testCaseFunc: notifyUsers_errPhoneNotifier,
		},
		{
			name:         "notifyUsers results 1 success on score equal 50",
			args:         args{ctx: ctx, users: []User{user_score50_succ}, message: message},
			testCaseFunc: notifyUsers_succPhoneNotifier,
		},
		{
			name:         "notifyUsers results 1 failure on score lesser than 50",
			args:         args{ctx: ctx, users: []User{user_scoreLesser50_fail}, message: message},
			testCaseFunc: notifyUsers_errPhoneNotifier,
		},
		{
			name:         "notifyUsers results 1 success on score lesser than 50",
			args:         args{ctx: ctx, users: []User{user_scoreLesser50_succ}, message: message},
			testCaseFunc: notifyUsers_succPhoneNotifier,
		},
		{
			name:         "notifyUsers results all failures on 1 score greater than 50, 1 score equal 50, & 1 score lesser than 50",
			args:         args{ctx: ctx, users: []User{user_scoreGreater50_fail, user_score50_fail, user_scoreLesser50_fail}, message: message},
			testCaseFunc: notifyUsers_allFail_1scoreGreater50_1score50_1scoreLesser50,
		},
		{
			name:         "notifyUsers results all successes on 1 score greater than 50, 1 score equal 50, & 1 score lesser than 50",
			args:         args{ctx: ctx, users: []User{user_scoreGreater50_succ, user_score50_succ, user_scoreLesser50_succ}, message: message},
			testCaseFunc: notifyUsers_allSucc_1scoreGreater50_1score50_1scoreLesser50,
		},
		{
			name:         "notifyUsers results successes on 1 score > 50, 1 score = 50, & 1 score < 50; failures on 1 score > 50, 1 score = 50, & 1 score < 50",
			args:         args{ctx: ctx, users: []User{user_scoreGreater50_fail, user_score50_fail, user_scoreLesser50_fail, user_scoreGreater50_succ, user_score50_succ, user_scoreLesser50_succ}, message: message},
			testCaseFunc: notifyUsers_fail_1Greater50_1score50_1Lesser50_succ_1Greater50_1score50_1Lesser50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := userServiceMocks{
				emailNotifier: NewMockNotifier(ctrl),
				phoneNotifier: NewMockNotifier(ctrl),
			}
			testCaseResp := tt.testCaseFunc(notifyUsersTestParam{
				ctx:     tt.args.ctx,
				users:   tt.args.users,
				message: tt.args.message,
				mocks:   mocks,
			})

			us := &UserService{
				phoneNotifier: mocks.phoneNotifier,
				emailNotifier: mocks.emailNotifier,
			}
			if gotResp := us.notifyUsers(tt.args.ctx, tt.args.users, tt.args.message); !reflect.DeepEqual(gotResp, testCaseResp.expectedRes) {
				t.Errorf("notifyUsers() = %v, want %v", gotResp, testCaseResp.expectedRes)
			}
		})
	}
}
