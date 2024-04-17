package models

type Customer struct {
	Id          string  `json:"id"`
	Mail       string   `json:"mail"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Password    string  `json:"password"`
	Phone       string  `json:"phone"`
	Sex         string  `json:"sex"`
	Active      bool  `json:"active"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type CreateCustomer struct {
	Id          string  `json:"id"`
	Mail       string   `json:"mail"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Password    string  `json:"password"`
	Phone       string  `json:"phone"`
	Sex         string  `json:"sex"`
	Active      bool    `json:"active"`
}

type UpdateCustomer struct {
	Id          string  `json:"id"`
	Mail       string   `json:"mail"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Password    string  `json:"password"`
	Phone       string  `json:"phone"`
	Sex         string  `json:"sex"`
	Active      bool    `json:"active"`
}

type GetAllCustomersResponse struct {
	Customers []Customer `json:"customers"`
	Count int16 `json:"count"`
}

type GetAllCustomersRequest struct {
    Search string `json:"search"`
	Page uint64 `json:"page"`
	Limit uint64 `json:"limit"`
}

type PasswordOfCustomer struct {
	Phone       string  `json:"phone"`
	NewPassword string `json:"new_password"`
	Password string `json:"password"`
}

type ForgetpasswordofEmail struct {
	Mail    string `json:"mail"`
	Optcode string `json:"otp_code"`
	Password string `json:"password"`

}