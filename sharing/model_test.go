package main

import "testing"

func TestNotifyUsersByTypeRequest_Validate(t *testing.T) {
	type fields struct {
		Message  string
		UserType string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Validate fail, empty Message & UserType",
			fields:  fields{Message: "", UserType: ""},
			wantErr: true,
		},
		{
			name:    "Validate fail, empty Message",
			fields:  fields{Message: "", UserType: UserTypePremium},
			wantErr: true,
		},
		{
			name:    "Validate fail, empty UserType",
			fields:  fields{Message: "Message", UserType: ""},
			wantErr: true,
		},
		{
			name:    "Validate success",
			fields:  fields{Message: "Message", UserType: UserTypePremium},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := NotifyUsersByTypeRequest{
				Message:  tt.fields.Message,
				UserType: tt.fields.UserType,
			}
			if err := sr.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
