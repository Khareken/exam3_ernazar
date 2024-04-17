package service

import (
	"exam/pkg/logger"
	"exam/storage"
)


type IServiceManager interface {
	Customer() customerService
	Auth()  authService

}

type Service struct {
 customerService customerService
 logger logger.ILogger
 auth   authService
}

func New(storage storage.IStorage,log logger.ILogger,redis storage.IRedisStorage) Service  {
	services := Service{}
	services.customerService = NewCustomerService(storage,log,redis)
	services.auth = NewAuthService(storage,log,redis)
	services.logger = log

	return services
}


func (s Service) Customer() customerService {
	return s.customerService
}

func (s Service) Auth() authService {
	return s.auth
}