package handler

import (
	"fmt"
	"net/http"
	_ "exam/api/docs"
	"exam/api/models"
	"exam/pkg/check"
	"github.com/gin-gonic/gin"
)

// CustomerLogin godoc
// @Router       /customer/login [POST]
// @Summary      Customer login
// @Description  Customer login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.CustomerLoginRequest true "login"
// @Success      201  {object}  models.CustomerLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerLogin(c *gin.Context)  {
	loginReq := models.CustomerLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq);err != nil {
		handlerResponseLog(c, h.Log, "Failed to process login data", http.StatusBadRequest, err)
		return
	}

	if err := check.ValidatePhoneNumberOfCustomer(loginReq.Login); !true {
		handlerResponseLog(c,h.Log,"Provided email is invalid: "+loginReq.Login, http.StatusBadRequest, err)
		return
	}

	if err := check.ValidatePassword(loginReq.Password); !true {
		handlerResponseLog(c,h.Log,"Password does not meet security requirements:"+loginReq.Password, http.StatusBadRequest, err)
		return
	}

	loginResp,err := h.Services.Auth().CustomerLogin(c.Request.Context(),loginReq)
	if err != nil {
		handlerResponseLog(c,h.Log,"Login attempt failed",http.StatusUnauthorized,err)
		return
	}
	handlerResponseLog(c,h.Log,"Login successful",http.StatusOK,loginResp)
}


// CustomerRegister godoc
// @Router       /customer/register [POST]
// @Summary      Customer register
// @Description  Customer register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.CustomerRegisterRequest true "register"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerRegister(c *gin.Context)  {
	loginReq := models.CustomerRegisterRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handlerResponseLog(c, h.Log, "Failed to read registration data", http.StatusBadRequest, err)
		return
	}
	
	if err := check.CheckEmail(loginReq.Mail); err != nil {
		handlerResponseLog(c,h.Log,"Email validation failed for: "+loginReq.Mail, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Services.Auth().CustomerRegister(c.Request.Context(), loginReq)
	if err != nil {
		handlerResponseLog(c, h.Log, "Registration failed", http.StatusInternalServerError,err)
		return
	}
	handlerResponseLog(c, h.Log, "Verification email sent successfully", http.StatusOK, "okey")
}



// CustomerRegister godoc
// @Router       /customer/register-confirm [POST]
// @Summary      Customer register
// @Description  Customer register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.CustomerRegisterConfRequest true "register"
// @Success      201  {object}  models.CustomerLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerRegisterConfirm(c *gin.Context) {
	req := models.CustomerRegisterConfRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		handlerResponseLog(c, h.Log, "Failed to process confirmation data", http.StatusBadRequest, err)
		return
	}
	fmt.Println("req: ", req)

    if err := check.CheckEmail(req.Mail); err != nil {
		handlerResponseLog(c, h.Log, "Email validation failed during confirmation for:", http.StatusBadRequest, err.Error())
	  }

	confResp, err := h.Services.Auth().CustomerRegisterConfirm(c.Request.Context(), req)
	if err != nil {
		handlerResponseLog(c, h.Log, "Registration confirmation failed", http.StatusUnauthorized, err.Error())
		return
	}

	handlerResponseLog(c, h.Log, "Registration confirmed successfully", http.StatusOK, confResp)

}


