package postgres

import (
	"context"
	"errors"
	"fmt"
	"exam/api/models"
	"github.com/google/uuid"
)


func (c *customerRepo) CustomerRegisterCreate(ctx context.Context, customer models.LoginCustomer) (string, error) {
	id := uuid.New()
	fmt.Println("PASSWORD-----------",customer.Password)

	query := ` INSERT INTO customers (
		id ,       
		first_name,
		last_name ,
		mail,    
		phone,
		sex,
		password     )
		VALUES($1,$2,$3,$4,$5,$6) `

	_,err := c.db.Exec(ctx,query,id.String(),customer.First_name, customer.Last_name,
	customer.Gmail, customer.Phone, customer.Sex,customer.Password)	
	if err != nil {
		return "",err
	}
	return id.String(),nil
}

func (c *customerRepo) GetGmail(ctx context.Context, gmail string) (string, error) {
	var id string

	query :=`SELECT id
	FROM customers
	WHERE mail = $1`


	err :=c.db.QueryRow(ctx,query,gmail).Scan(&id)
	if err != nil {
		return id,nil
	}
	return id, errors.New("The Gmail address provided is already associated with an existing account")

}


