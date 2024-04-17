package service

import (
	"context"
	"exam/api/models"
	"exam/pkg/logger"
	"exam/storage"
)



type customerService struct {
   storage storage.IStorage
   logger logger.ILogger
   redis storage.IRedisStorage
}

func NewCustomerService(storage storage.IStorage,logger logger.ILogger,redis storage.IRedisStorage) customerService {
	return customerService{
		storage: storage,
		logger: logger,
		redis: redis,
	}
}

func (cs customerService) Create(ctx context.Context, customer models.CreateCustomer) (string,error) {
	pkey,err := cs.storage.Customer().Create(ctx,customer)
	if err != nil {
		cs.logger.Error("Failed to create customer", logger.Error(err))
		return "", err
	}
	return pkey,nil
}

func (cs customerService) Update(ctx context.Context, customer models.UpdateCustomer) (string,error) {
	pkey, err := cs.storage.Customer().UpdateCustomer(ctx,customer)
	if err != nil {
		cs.logger.Error("Failed to update customer information", logger.Error(err))
		return "",err
	}

	err = cs.redis.Del(ctx,customer.Id)
	if err != nil {
		cs.logger.Error("Failed to clear customer cache after update", logger.Error(err))
		return "Cache clear error",err
	}

	return pkey,nil
}

func (cs customerService) UpdateStatus(ctx context.Context, customer models.UpdateCustomer) (string,error) {
	pkey, err := cs.storage.Customer().UpdateCustomerStatus(ctx,customer)
	if err != nil {
		cs.logger.Error("Failed to update customer status", logger.Error(err))
		return "",err
	}

	err = cs.redis.Del(ctx,customer.Id)
	if err != nil {
		cs.logger.Error("Failed to clear customer cache after status update", logger.Error(err))
		return "Cache clear error",err
	}

	return pkey,nil
}

func (cs customerService) GetCustomerAll(ctx context.Context,customer models.GetAllCustomersRequest) (models.GetAllCustomersResponse, error) {
	customers, err := cs.storage.Customer().GetAllCustomer(ctx,customer)
	if err != nil {
		cs.logger.Error("Failed to retrieve all customers", logger.Error(err))
		return customers,err
	}
	return customers,nil
}

func (cs customerService) GetByIDCustomer(ctx context.Context, id string) (models.Customer,error) {
	customer, err := cs.storage.Customer().GetByID(ctx,id)
	if err != nil {
		cs.logger.Error("Failed to retrieve customer by ID", logger.Error(err))
		return models.Customer{},err
	}
	return customer,nil
}


func (cs customerService) Delete(ctx context.Context, id string) (error) {
	err := cs.storage.Customer().Delete(ctx,id)
	if err != nil {
		cs.logger.Error("Failed to delete customer", logger.Error(err))
		return err
	}

	err = cs.redis.Del(ctx, id)
	if err != nil {
		cs.logger.Error("Failed to clear customer cache after deletion", logger.Error(err))
		return err
	}
	
	return nil
}

func (cs customerService) UpdatePassword(ctx context.Context, customer models.PasswordOfCustomer) (string, error) {

	pKey, err := cs.storage.Customer().UpdateCustomerPassword(ctx, customer)
	if err != nil {
		cs.logger.Error("Failed to update customer password", logger.Error(err))
		return "", err
	}

	return pKey, nil
}
