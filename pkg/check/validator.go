package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
)



func ValidateGmailCustomer(e string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,3}$`)
    return emailRegex.MatchString(e)
}

func ValidateEmail(address string) (error) {
  _, err := mail.ParseAddress(address)
  if err != nil {
    return  errors.New("email is not valid")
    
  }
  return nil
}

type EmailVerificationResponse struct {
	Data struct {
	  Result string `json:"result"`
	} `json:"data"`
  }
  
  func CheckEmail(email string) (error) {
	apiKey := "a78afa97d76af0e3364a3eb68ed12aae83e247a0"       
  
	url := fmt.Sprintf("https://api.hunter.io/v2/email-verifier?email=%s&api_key=%s", email, apiKey)
	resp, err := http.Get(url)
	if err != nil {
	  fmt.Println("Network error during email verification:", err)
	  return err
	}
	defer resp.Body.Close()
  
	var verificationResponse EmailVerificationResponse
	err = json.NewDecoder(resp.Body).Decode(&verificationResponse)
	if err != nil {
	  fmt.Println("Failed to decode verification response:", err)
	  return err
	}
  
  
	if verificationResponse.Data.Result == "undeliverable" {
	  fmt.Println("Provided email is undeliverable")
	  
	  return errors.New("Email address does not exist or is undeliverable")
	} else if verificationResponse.Data.Result == "deliverable" {
	  fmt.Println("Email address is valid")
	  return nil
	} else {
	  fmt.Println("Email verification inconclusive")
  
	  return errors.New("Unable to verify email address")
	}
  
	return errors.New("Email address does not exist or is undeliverable")
  }
  


func ValidatePhoneNumberOfCustomer(phone string) bool {
	if 12 < len(phone) && len(phone) <= 13{
		phoneregex:= regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
		return phoneregex.MatchString(phone)
	}
	return false
}

func ValidatePassword(password string) error {
  lowercaseRegex := `[a-z]`
  hasLowercase, _ := regexp.MatchString(lowercaseRegex, password)
  uppercaseRegex := `[A-Z]`
  hasUppercase, _ := regexp.MatchString(uppercaseRegex, password)
  digitRegex := `[0-9]`
  hasDigit, _ := regexp.MatchString(digitRegex, password)
  symbolRegex := `[!@#$%^&*()-_+=~\[\]{}|\\:;"'<>,.?\/]`
  hasSymbol, _ := regexp.MatchString(symbolRegex, password)

  if hasLowercase && hasUppercase && hasDigit && hasSymbol && len(password) >= 8 {
    return nil
  }

  return errors.New("Password must include at least one lowercase, one uppercase, a digit, a symbol and be at least 8 characters long")
  }