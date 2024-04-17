package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"exam/api/models"
	"exam/config"
	"exam/pkg"
    
	"exam/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)



type customerRepo struct {
	db *pgxpool.Pool
	logger logger.ILogger
}


func NewCustomer(db *pgxpool.Pool,log logger.ILogger) customerRepo {
	return customerRepo{
		db: db,
		logger: log,
	}
}

func (c *customerRepo) Create(ctx context.Context, customer models.CreateCustomer) (string, error) {
	id := uuid.New()

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return "error while hashing password here", err
	}
	customer.Password = string(hashedpassword)

	query := `insert into customers(id,mail,first_name,last_name,password,phone,sex,active) values($1,$2,$3,$4,$5,$6,$7,$8)`

	ctx, cancel := context.WithTimeout(ctx, config.TimewithContex)
	defer cancel()

	_, err = c.db.Exec(ctx, query, id.String(),customer.Mail, customer.FirstName, customer.LastName, hashedpassword, customer.Phone,customer.Sex,customer.Active)
	if err != nil {
		return "error:", err
	}
	return id.String(), nil
}

func (c *customerRepo) UpdateCustomer(ctx context.Context, customer models.UpdateCustomer) (string, error) {
	query := `update customers set 
	first_name=$1,
	last_name=$2,
	mail=$3,
	phone=$4,
	active=$5,
	sex=$6,
	updated_at=CURRENT_TIMESTAMP
	WHERE id = $7`

	ctx, cancel := context.WithTimeout(ctx, config.TimewithContex)
	defer cancel()

	_, err := c.db.Exec(ctx, query,
		customer.FirstName,
		customer.LastName,
		customer.Mail,
		customer.Phone,
		customer.Active,
		customer.Sex,
		customer.Id)
		if err != nil {
		return "", err
	}
	return customer.Id, nil
}
func (c *customerRepo) UpdateCustomerStatus(ctx context.Context, customer models.UpdateCustomer) (string, error) {
	query := `update customers set 
	active=$1,
	updated_at=CURRENT_TIMESTAMP
	WHERE id = $2`

	ctx, cancel := context.WithTimeout(ctx, config.TimewithContex)
	defer cancel()

	_, err := c.db.Exec(ctx, query,
		customer.Active,
		customer.Id)
		if err != nil {
		return "", err
	}
	return customer.Id, nil
}

func (c *customerRepo) GetAllCustomer(ctx context.Context, req models.GetAllCustomersRequest) (models.GetAllCustomersResponse, error) {
	var (
		resp   = models.GetAllCustomersResponse{}
		filter = ""
	)

	offset := (req.Page - 1) * req.Limit
	if req.Search != "" {
		filter += fmt.Sprintf(`and first_name ILIKE '%%%v%%'`, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter:", filter)

	query := `select id,mail,first_name,last_name,password,phone,sex,active,created_at,updated_at from customers`
	ctx, cancel := context.WithTimeout(ctx, config.TimewithContex)
	defer cancel()

	rows, err := c.db.Query(ctx, query+filter+``)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			customer = models.Customer{}
			updateAt sql.NullString
			createdAt sql.NullString
		)
		if err := rows.Scan(
			&customer.Id,
			&customer.Mail,
			&customer.FirstName,
			&customer.LastName,
			&customer.Password,
			&customer.Phone,
			&customer.Sex,
			&customer.Active,
			&createdAt,
			&updateAt); err != nil {
			return resp, err
		}
		customer.CreatedAt = pkg.NullStringToString(createdAt)
		customer.UpdatedAt = pkg.NullStringToString(updateAt)
		resp.Customers = append(resp.Customers, customer)
	}
	if err = rows.Err(); err != nil {
		return resp, err
	}
	countQuery := `Select count(*) from customers`
	err = c.db.QueryRow(ctx, countQuery).Scan(&resp.Count)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *customerRepo) GetByID(ctx context.Context, id string) (models.Customer, error) {
	customer := models.Customer{}

     var (
		updateAt sql.NullString
		createdAt sql.NullString
	 )
	 
	 
	ctx, cancel := context.WithTimeout(ctx, config.TimewithContex)
	defer cancel()
	if err := c.db.QueryRow(ctx,`select id,first_name,last_name,mail,phone,active,password,sex,created_at,updated_at from customers where id = $1`, id).Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.LastName,
		&customer.Mail,
		&customer.Phone,
		&customer.Active,
		&customer.Password,
		&customer.Sex,
		&createdAt,
		&updateAt); err != nil {
		return models.Customer{}, err
	}
	customer.CreatedAt = pkg.NullStringToString(createdAt)
	customer.UpdatedAt = pkg.NullStringToString(updateAt)
	return customer, nil
}

func (c *customerRepo) Delete(ctx context.Context, id string) error {
	queary := `delete from customers where id = $1`

	ctx, cancel := context.WithTimeout(ctx, config.TimewithContex)
	defer cancel()
	_, err := c.db.Exec(ctx, queary, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *customerRepo) UpdateCustomerPassword(ctx context.Context, customer models.PasswordOfCustomer) (string, error) {
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(customer.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error while hashing new password")
	}

	var Pass_cur string
	query := `select password from customers where phone = $1`

	err = c.db.QueryRow(ctx, query, customer.Phone).Scan(&Pass_cur)
	if err != nil {
		return "error:", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(Pass_cur), []byte(hashedNewPassword))
	if err != nil {
		log.Println("error while comparing new and old passwords !!!")
	}

	ctx, cancel := context.WithTimeout(ctx, config.TimewithContex)
	defer cancel()

	_, err = c.db.Exec(ctx, `update customers set password = $1 where phone = $2`, hashedNewPassword, customer.Phone)
	if err != nil {
		return "error:", err
	}

	return "OK", nil
}

func (c *customerRepo) GetPasswordforLogin(ctx context.Context, phone string) (string, error) {
	var hashedPasswordforLogin string

	query := `SELECT password
	FROM customers
	WHERE phone = $1`

	err := c.db.QueryRow(ctx, query, phone).Scan(&hashedPasswordforLogin)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("phone is incorrect here")
		} else {
			return "", err
		}
	}

	return hashedPasswordforLogin, nil
}

func (c *customerRepo) GetByLogin(ctx context.Context, login string)(models.CreateCustomer, error) {
	var (
		firstname sql.NullString
		lastname  sql.NullString
		phone     sql.NullString
		mail      sql.NullString
		sex       sql.NullString
	)

	query := `SELECT 
		id, 
		first_name, 
		last_name, 
		phone,
		mail,
		sex,
		active,
		password
		FROM customers WHERE phone = $1`

	row := c.db.QueryRow(ctx, query, login)

	customer := models.CreateCustomer{}

	err := row.Scan(
		&customer.Id,
		&firstname,
		&lastname,
		&phone,
		&mail,
		&sex,
		&customer.Active,
		&customer.Password,
	)

	if err != nil {
		return models.CreateCustomer{}, err
	}

	customer.FirstName = firstname.String
	customer.LastName = lastname.String
	customer.Phone = phone.String
	customer.Mail = mail.String
	customer.Sex = sex.String


	return customer, nil
}


func (c *customerRepo) UpdateForgotPasswordofEmail(ctx context.Context, User models.ForgetpasswordofEmail) (string, error) {

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error hashing new password")
	}

	_, err = c.db.Exec(ctx, `UPDATE Users SET password = $1 WHERE mail = $2`, hashedNewPassword, User.Mail)
	if err != nil {
		return "", errors.New("error updating password")
	}

	return "Password updated successfully", nil
}

