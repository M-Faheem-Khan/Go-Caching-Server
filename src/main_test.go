package main

import "testing"

func TestGetUserID(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "no_userId_given_without_trailing_slash", args: args{path: "/users"}, want: 0, wantErr: true},
		{name: "no_userId_given_with_trailing_slash", args: args{path: "/users/"}, want: 0, wantErr: true},

		{name: "userId_given_without_trailing_slash", args: args{path: "/users/1234"}, want: 1234, wantErr: false},
		{name: "userId_given_with_trailing_slash", args: args{path: "/users/1234/"}, want: 1234, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserID(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
