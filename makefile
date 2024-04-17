swag_init:
	swag init -g api/router.go -o api/docs
c_m:
	#creates a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)