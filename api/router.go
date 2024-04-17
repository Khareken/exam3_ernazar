package api

import (
	"exam/api/handler"
	"exam/pkg/logger"
	"exam/service"
	

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)
// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceManager,log logger.ILogger) *gin.Engine {
	h :=handler.NewStrg(services,log)
    
	r :=gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/customer/login",h.CustomerLogin)
	r.POST("/customer/register", h.CustomerRegister)
	r.POST("/customer/register-confirm", h.CustomerRegisterConfirm)
   
	r.POST("/customer",h.CreateCustomer)
	r.PUT("/customer/:id",h.UpdateCustomer)
    r.GET("/customers",h.GetAllCustomer)
    r.GET("/customer/:id",h.GetByIDCustomer)
    r.DELETE("/customer/:id",h.DeleteCustomer)

	r.PATCH("/customer/password",h.UpdateCustomerPassword)
	r.PATCH("/customer/:id",h.UpdateCustomerStatus)

	return r
}
