package storage

import (
	"context"
	"exam/api/models"
	"time"
)




type IStorage interface {
	CloseDB()
	Customer() ICustomerStorage
	Redis()    IRedisStorage
}


type ICustomerStorage interface {
	Create(context.Context,models.CreateCustomer) (string, error)
	UpdateCustomer(context.Context,models.UpdateCustomer) (string, error)
	UpdateCustomerStatus(ctx context.Context, customer models.UpdateCustomer) (string, error)
	GetAllCustomer(ctx context.Context,req models.GetAllCustomersRequest) (models.GetAllCustomersResponse, error)
	GetByID(ctx context.Context,id string) (models.Customer, error)
	Delete(ctx context.Context,id string) error
	UpdateCustomerPassword(context.Context,models.PasswordOfCustomer)(string, error)
	GetPasswordforLogin(ctx context.Context, phone string) (string, error)
	GetByLogin(context.Context, string) (models.CreateCustomer, error)
	GetGmail(ctx context.Context, gmail string) (string, error)
	CustomerRegisterCreate(ctx context.Context, customer models.LoginCustomer) (string, error)
    UpdateForgotPasswordofEmail(ctx context.Context, User models.ForgetpasswordofEmail) (string, error) 

}

type IRedisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (interface{},error)
	Del(ctx context.Context, key string) error
}