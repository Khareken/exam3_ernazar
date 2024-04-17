package service

import (
	"context"
	"errors"
	"exam/api/models"
	"exam/config"
	"exam/pkg"
	"exam/pkg/jwt"
	"exam/pkg/logger"
	"exam/pkg/password"
	"fmt"

	"exam/pkg/smtp"
	"exam/storage"
	"time"
)

type authService struct {
	storage storage.IStorage
	log     logger.ILogger
	redis   storage.IRedisStorage
	

}

func NewAuthService(storage storage.IStorage, log logger.ILogger,redis storage.IRedisStorage) authService {
	return authService{
		storage: storage,
		log:     log,
		redis: redis,
	}
}

func (a authService) CustomerLogin(ctx context.Context, loginRequest models.CustomerLoginRequest) (models.CustomerLoginResponse, error) {
	fmt.Println(" Attempting login with: ", loginRequest.Login)
	customer, err := a.storage.Customer().GetByLogin(ctx, loginRequest.Login)
	if err != nil {
		a.log.Error("Failed to retrieve customer by login", logger.Error(err))
		return models.CustomerLoginResponse{}, err
	}

	if err = password.CompareHashAndPassword(customer.Password, loginRequest.Password); err != nil {
		a.log.Error("Password mismatch error", logger.Error(err))
		return models.CustomerLoginResponse{}, err
	}

	m := make(map[interface{}]interface{})

	m["user_id"] = customer.Id
	m["user_role"] = config.CUSTOMER_ROLE

    accessToken,refreshToken,err :=jwt.GenJWT(m)
	if err != nil {
		a.log.Error("Failed to generate tokens during customer login",logger.Error(err))
	}

	return models.CustomerLoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	},nil
}


func (a authService) CustomerRegister(ctx context.Context, loginRequest models.CustomerRegisterRequest) error {
	_,err := a.storage.Customer().GetGmail(ctx,loginRequest.Mail)
	if err != nil {
		a.log.Error("This email is already associated with an account",logger.Error(err))
		return err
	}
	otpCode := pkg.GenerateOTP()

	msg := fmt.Sprintf("Your otp code is: %v, for registering exam. Don't give it to anyone", otpCode)
	
	fmt.Printf("Your otp code is: %v, for registering exam. Don't give it to anyone", otpCode)
	
	fmt.Println(loginRequest.Mail,"----------",otpCode)

	err = a.redis.SetX(ctx, loginRequest.Mail, otpCode, time.Minute*2)
	if err != nil {
		a.log.Error("error while setting otpCode to redis customer register", logger.Error(err))
		return err
	}
	
    err = smtp.SendMail(loginRequest.Mail, msg)
	if err != nil {
		a.log.Error("Failed to set OTP in Redis for customer registration", logger.Error(err))
		return err
	}
   return nil
}


func (a authService) CustomerRegisterConfirm(ctx context.Context,req models.CustomerRegisterConfRequest) (models.CustomerLoginResponse,error) {
	resp := models.CustomerLoginResponse{}

	otp,err := a.redis.Get(ctx,req.Mail)
	if err != nil {
		a.log.Error("Failed to retrieve OTP for email confirmation",logger.Error(err))
		return resp, errors.New("incorrect otp code")	
	}
   
    if req.Otp != otp {
		a.log.Error("OTP mismatch during registration confirmation", logger.Error(err))
	 return resp,errors.New("incorrect otp code")
	}

	req.Customer.Mail = req.Mail
	id,err := a.storage.Customer().Create(ctx,req.Customer)
	if err != nil {
		a.log.Error("Failed to create customer account",logger.Error(err))
		return resp,err
	}
	m := make(map[interface{}]interface{})
	m["user_id"] = id
	m["user_role"] = config.CUSTOMER_ROLE

	accessToken,refreshToken,err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("Failed to generate tokens after registration confirmation",logger.Error(err))
	return resp,err
	}
	resp.AccessToken = accessToken
	resp.RefreshToken =  refreshToken

	return resp,nil
}

func (a authService) ForgetpasswordofEmail(ctx context.Context, customer models.ForgetpasswordofEmail) (id string,err error) {
	
	resp := models.ForgetpasswordofEmail{}


   otpCode,err := a.storage.Redis().Get(context.Background(),customer.Mail)
   if err != nil {
	a.log.Error(err.Error())
   }

    OTPCODEStr, ok := otpCode.(string)
	if !ok {
		a.log.Error("ERROR in service layer while login user", logger.Error(errors.New(err.Error())))
		return id, errors.New(err.Error())
	}

if OTPCODEStr!=customer.Optcode{
	a.log.Error("ERROR in service layer while login user", logger.Error(err))
	return id,errors.New(err.Error())
}

id, err = a.storage.Customer().UpdateForgotPasswordofEmail(ctx,resp)
if err == nil {
	a.log.Error("ERROR in service layer while login", logger.Error(err))
	return id, err
}
	

	return id, nil
}
