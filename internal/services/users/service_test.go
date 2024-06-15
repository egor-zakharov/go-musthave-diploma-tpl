package users

import (
	"context"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/users"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func Test_service_Login(t *testing.T) {
	ctx := context.Background()

	userIn := models.User{
		Login:    "login",
		Password: "password",
	}

	userOut := models.User{
		UserID:   "1",
		Login:    "login",
		Password: "$2a$10$mFTV7pqNmJC1VWdTtVi2geNGLLlK7Xo7NwjrZDqBrOF1WX.8kMgoC",
	}

	type fields struct {
		log     *zap.Logger
		storage func(ctrl *gomock.Controller) users.Storage
	}
	type args struct {
		ctx    context.Context
		userIn models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) users.Storage {
					mock := users.NewMockStorage(ctrl)
					mock.EXPECT().Login(gomock.Any(), userIn.Login).Return(&userOut, nil)
					return mock
				},
			},
			args: args{
				ctx:    ctx,
				userIn: userIn,
			},
			want:    &userOut,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := &service{
				log:     tt.fields.log,
				storage: tt.fields.storage(ctrl),
			}
			got, err := s.Login(tt.args.ctx, tt.args.userIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO как правильно неэкспортируемые методы самого сервиса замокать?
func Test_service_Register(t *testing.T) {
	ctx := context.Background()

	userIn := models.User{
		Login:    "login",
		Password: "password",
	}

	userOut := models.User{
		UserID:   "1",
		Login:    "login",
		Password: "$2a$10$mFTV7pqNmJC1VWdTtVi2geNGLLlK7Xo7NwjrZDqBrOF1WX.8kMgoC",
	}

	type fields struct {
		log     *zap.Logger
		storage func(ctrl *gomock.Controller) users.Storage
	}
	type args struct {
		ctx    context.Context
		userIn models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) users.Storage {
					mock := users.NewMockStorage(ctrl)
					mock.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&userOut, nil).Times(1)
					return mock
				},
			},
			args: args{
				ctx:    ctx,
				userIn: userIn,
			},
			want:    &userOut,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := &service{
				log:     tt.fields.log,
				storage: tt.fields.storage(ctrl),
			}
			got, err := s.Register(tt.args.ctx, tt.args.userIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}
