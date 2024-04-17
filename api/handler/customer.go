package handler

import (
	"context"
	"errors"
	_ "exam/api/docs"
	"exam/api/models"
	"exam/config"
	"exam/pkg/check"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Security ApiKeyAuth
// @Router       /customer [POST]
// @Summary      Creates a new customers
// @Description  create a new customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        customer body models.CreateCustomer false "customer"
// @Success      201 {object} models.Customer
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) CreateCustomer(c *gin.Context) {
	customer := models.CreateCustomer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		handlerResponseLog(c, h.Log, "Failed to parse customer data", http.StatusBadRequest, err.Error())
		return
	}

	if err := check.ValidatePassword(customer.Password); err != nil {
		handlerResponseLog(c, h.Log, "Invalid password provided", http.StatusBadRequest, err.Error())
	}

	if err := check.ValidateGmailCustomer(customer.Mail); !true {
		handlerResponseLog(c, h.Log, "Failed email verification for"+customer.Mail, http.StatusBadRequest, err)
		return
	}

	if err := check.ValidatePhoneNumberOfCustomer(customer.Phone); !true {
		handlerResponseLog(c, h.Log, "Failed phone number verification for"+customer.Phone, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(c, config.TimewithContex)
	defer cancel()


	id, err := h.Services.Customer().Create(ctx, customer)
	if err != nil {
		handlerResponseLog(c, h.Log, "Customer registration failed", http.StatusInternalServerError, err.Error())
		return
	}
	handlerResponseLog(c, h.Log, "Customer created successfully", http.StatusOK, id)
}



// @Security ApiKeyAuth
// @Router       /customer/{id} [PUT]
// @Summary      Update customer
// @Description  Update customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        id path string true "customer_id"
// @Param        car body models.UpdateCustomer true "customer"
// @Success      201 {object} models.Customer
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) UpdateCustomer(c *gin.Context) {
    customer := models.UpdateCustomer{}

    if err := c.ShouldBindJSON(&customer); err != nil {
        handlerResponseLog(c, h.Log, "Failed to parse customer update data", http.StatusBadRequest, err.Error())
        return
    }

    if err := check.ValidatePassword(customer.Password); !true {
        handlerResponseLog(c, h.Log, "Invalid password provided", http.StatusBadRequest, err.Error())
    }

    if err := check.ValidateGmailCustomer(customer.Mail); !true {
        handlerResponseLog(c, h.Log, "Failed email verification for "+customer.Mail, http.StatusBadRequest,err)
        return
    }

    if err := check.ValidatePhoneNumberOfCustomer(customer.Phone); !true {
        handlerResponseLog(c, h.Log, "Failed phone number verification for "+customer.Phone, http.StatusBadRequest,err)
        return
    }

    fmt.Println("Incoming customer ID:", c.Query("id"))

    customer.Id = c.Param("id")
    err := uuid.Validate(customer.Id)
    if err != nil {
        handlerResponseLog(c, h.Log, "Invalid Customer ID", http.StatusBadRequest, err.Error())
        return
    }

    ctx, cancel := context.WithTimeout(c, config.TimewithContex)
    defer cancel()
    id, err := h.Services.Customer().Update(ctx, customer)

    if err != nil {
        handlerResponseLog(c, h.Log, "Failed to update customer details", http.StatusInternalServerError, err.Error())
        return
    }

    handlerResponseLog(c, h.Log, "Customer updated successfully", http.StatusOK, id)
}



// @Security ApiKeyAuth
// @Router       /customer/{id} [PATCH]
// @Summary      Update customer status
// @Description  Update customer status
// @Tags         Password and Status
// @Accept       json
// @Produce      json
// @Param        id path string true "customer_id"
// @Param        car body models.UpdateCustomer true "customer"
// @Success      201 {object} models.Customer
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) UpdateCustomerStatus(c *gin.Context) {
    customer := models.UpdateCustomer{}

    if err := c.ShouldBindJSON(&customer); err != nil {
        handlerResponseLog(c, h.Log, "Failed to parse customer status data", http.StatusBadRequest, err.Error())
        return
    }

    fmt.Println("Incoming customer ID:", c.Param("id"))

    customer.Id = c.Param("id")
    err := uuid.Validate(customer.Id)
    if err != nil {
        handlerResponseLog(c, h.Log, "Invalid customer ID", http.StatusBadRequest, err.Error())
        return
    }

    ctx, cancel := context.WithTimeout(c, config.TimewithContex)
    defer cancel()
    id, err := h.Services.Customer().UpdateStatus(ctx, customer)

    if err != nil {
        handlerResponseLog(c, h.Log, "Failed to update customer status", http.StatusInternalServerError, err.Error())
        return
    }

    handlerResponseLog(c, h.Log, "CUstomer status updated succesfully", http.StatusOK, id)
}


// @Security ApiKeyAuth
// @Router       /customers [GET]
// @Summary      Get customer list
// @Description  get customer list
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      201 {object} models.GetAllCustomersRequest
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetAllCustomer(c *gin.Context) {
	var (
		request = models.GetAllCustomersRequest{}
	)
	request.Search = c.Query("search")

	page, err := ParsePageQueryParam(c)
	if err != nil {
		handlerResponseLog(c, h.Log, "Failed to parse page parametr", http.StatusInternalServerError, err.Error())
		return
	}
	limit, err := ParseLimitQueryParam(c)
	if err != nil {
		handlerResponseLog(c, h.Log, "Failed to parse limit parametr", http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("page: ", page)
	fmt.Println("limit: ", limit)

	request.Page = page
	request.Limit = limit

	ctx, cancel := context.WithTimeout(c, config.TimewithContex)
	defer cancel()
	
	customers, err := h.Services.Customer().GetCustomerAll(ctx, request)
	if err != nil {
		handlerResponseLog(c, h.Log, "Failed to retrieve customer list", http.StatusInternalServerError, err.Error())
		return
	}

	handlerResponseLog(c, h.Log, "Customer list retrived successfully", http.StatusOK, customers)
}


// @Security ApiKeyAuth
// @Router       /customer/{id} [GET]
// @Summary      Gets customer
// @Description  get customer by ID
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        id path string true "customer"
// @Success      201 {object} models.Customer
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) GetByIDCustomer(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id:", id)

	ctx, cancel := context.WithTimeout(c, config.TimewithContex)
	defer cancel()

	customer, err := h.Services.Customer().GetByIDCustomer(ctx, id)
	if err != nil {
		handlerResponseLog(c, h.Log, "Failed to retrieve customer by ID", http.StatusInternalServerError, err.Error())
		return
	}
	handlerResponseLog(c, h.Log, "Customer retrieved successfully", http.StatusOK, customer)
}

// @Security ApiKeyAuth
// @Router       /customer/{id} [DELETE]
// @Summary      Delete customer
// @Description  Delete customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        id path string true "customer_id"
// @Success      201 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) DeleteCustomer(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id:", id)

	err := uuid.Validate(id)
	if err != nil {
		handlerResponseLog(c, h.Log, "Invalid customer ID", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c, config.TimewithContex)
	defer cancel()

	err = h.Services.Customer().Delete(ctx, id)
	if err != nil {
		handlerResponseLog(c, h.Log, "Failed to delete customer", http.StatusInternalServerError, err.Error())
		return
	}
	handlerResponseLog(c, h.Log, "Customer deleted successfully", http.StatusOK, id)
}



// @Security ApiKeyAuth
// UpdatePassword godoc
// @Router 		/customer/password [PATCH]
// @Summary 	updating password
// @Description updating password
// @Tags 		Password and Status
// @Accept		json
// @Produce		json
// @Param		customer body models.PasswordOfCustomer true "customer"
// @Success		200  {object}  models.PasswordOfCustomer
// @Failure		400  {object}  models.Response
// @Failure		404  {object}  models.Response
// @Failure		500  {object}  models.Response
func (h Handler) UpdateCustomerPassword(c *gin.Context) {

	customer := models.PasswordOfCustomer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		handlerResponseLog(c, h.Log, "Failed to parse password data", http.StatusBadRequest, err.Error())
		return
	}

	if customer.NewPassword == customer.Password {
		handlerResponseLog(c,h.Log,"Please provde a new password different from the old one "+customer.NewPassword,http.StatusBadRequest,errors.New("error is here in customer"))
	}

	if err := check.ValidatePassword(customer.Password); err != nil {
		handlerResponseLog(c, h.Log, "error while validating password", http.StatusBadRequest, err.Error())
	}

	
	ctx, cancel := context.WithTimeout(c, config.TimewithContex)
	defer cancel()
	_, err := h.Services.Customer().UpdatePassword(ctx, customer)
	if err != nil {
		handlerResponseLog(c, h.Log, "Failed to update password", http.StatusInternalServerError, err)
		return
	}
	handlerResponseLog(c, h.Log, "Password updated succesfully", http.StatusOK,customer)
}
