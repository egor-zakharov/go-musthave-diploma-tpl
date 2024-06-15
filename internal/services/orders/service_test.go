package orders

import (
	"context"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/orders"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"
)

func Test_service_Add(t *testing.T) {
	type fields struct {
		log     *zap.Logger
		storage func(ctrl *gomock.Controller) orders.Storage
	}
	type args struct {
		ctx     context.Context
		orderID string
		userID  string
	}

	ctx := context.Background()

	successCase := args{
		ctx:     ctx,
		orderID: "12345678903",
		userID:  "1",
	}
	errLuhnCase := args{
		ctx:     ctx,
		orderID: "1",
		userID:  "1",
	}
	dataConflictCase := args{
		ctx:     ctx,
		orderID: "12345678903",
		userID:  "1",
	}
	anotherUserCase := args{
		ctx:     ctx,
		orderID: "12345678903",
		userID:  "1",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successCase",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					mock.EXPECT().Add(gomock.Any(), successCase.orderID, successCase.userID).Return(&models.Order{UserID: successCase.userID}, nil)
					return mock
				},
			},
			args: args{
				ctx:     successCase.ctx,
				orderID: successCase.orderID,
				userID:  successCase.userID,
			},
			wantErr: false,
		},
		{
			name: "err luhn",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					return mock
				},
			},
			args: args{
				ctx:     errLuhnCase.ctx,
				orderID: errLuhnCase.orderID,
				userID:  errLuhnCase.userID,
			},
			wantErr: true,
		},
		{
			name: "data conflict",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					mock.EXPECT().Add(gomock.Any(), dataConflictCase.orderID, dataConflictCase.userID).Return(&models.Order{}, orders.ErrConflict)
					return mock
				},
			},
			args: args{
				ctx:     dataConflictCase.ctx,
				orderID: dataConflictCase.orderID,
				userID:  dataConflictCase.userID,
			},
			wantErr: true,
		},
		{
			name: "another user case",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					mock.EXPECT().Add(gomock.Any(), anotherUserCase.orderID, anotherUserCase.userID).Return(&models.Order{UserID: "2"}, nil)
					return mock
				},
			},
			args: args{
				ctx:     anotherUserCase.ctx,
				orderID: anotherUserCase.orderID,
				userID:  anotherUserCase.userID,
			},
			wantErr: true,
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
			if err := s.Add(tt.args.ctx, tt.args.orderID, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_Get(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	order := models.Order{
		Number:     "1",
		UserID:     "1",
		Status:     "New",
		Accrual:    0,
		UploadedAt: now,
	}

	type fields struct {
		log     *zap.Logger
		storage func(ctrl *gomock.Controller) orders.Storage
	}
	type args struct {
		ctx     context.Context
		orderID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Order
		wantErr bool
	}{
		{
			name: "successCase",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					mock.EXPECT().Get(gomock.Any(), "1").Return(&order, nil)
					return mock
				},
			},
			args: args{
				ctx:     ctx,
				orderID: "1",
			},
			want:    &order,
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					mock.EXPECT().Get(gomock.Any(), "1").Return(nil, orders.ErrNotFound)
					return mock
				},
			},
			args: args{
				ctx:     ctx,
				orderID: "1",
			},
			want:    nil,
			wantErr: true,
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
			got, err := s.Get(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetAllByUser(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	ords := []models.Order{
		{
			Number:     "1",
			UserID:     "1",
			Status:     "NEW",
			Accrual:    0,
			UploadedAt: now,
		},
		{
			Number:     "2",
			UserID:     "1",
			Status:     "NEW",
			Accrual:    0,
			UploadedAt: now,
		},
	}
	type fields struct {
		log     *zap.Logger
		storage func(ctrl *gomock.Controller) orders.Storage
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]models.Order
		wantErr bool
	}{
		{
			name: "successCase",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					mock.EXPECT().GetAllByUser(gomock.Any(), "1").Return(&ords, nil)
					return mock
				},
			},
			args: args{
				ctx:    ctx,
				userID: "1",
			},
			want:    &ords,
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				log: nil,
				storage: func(ctrl *gomock.Controller) orders.Storage {
					mock := orders.NewMockStorage(ctrl)
					mock.EXPECT().GetAllByUser(gomock.Any(), "1").Return(nil, orders.ErrNotFound)
					return mock
				},
			},
			args: args{
				ctx:    ctx,
				userID: "1",
			},
			want:    nil,
			wantErr: true,
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
			got, err := s.GetAllByUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllByUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
