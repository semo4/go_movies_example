go mod init ${name}

psql -f ./postgreSQL/shcema.sql -d movies_db

run: 
	$(shell cat .env | tr "\n" " ") go run .
 
go run cmd/cli.go -movieSecretKey 93f1029151c799988599e489eb7443f2