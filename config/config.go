package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	User     string
	Host     string
	Port     string
	Password string
	DbName   string
}

func validateEnvsDb(user string, host string, port string, password string, dbName string) error {
	if user == "" && host == "" && port == "" && password == "" && dbName == ""{
		return errors.New("Nenhuma env foi DB foi definida")
	}
	if user == ""{
		return errors.New("user DB não definido")
	}
	if host == ""{
		return errors.New("host DB não definido")
	}
	if port == ""{
		return errors.New("port DB não definido")
	}
	if password == ""{
		return errors.New("password DB não definido")
	}
	if dbName == ""{
		return errors.New("dbName DB não definido")
	}
	return nil
}

func LoadDbConfig() *DbConfig {

	err := godotenv.Load()
	if err != nil {
		log.Println("info: não encotnrado arquivo env (esperado em produção)")
	}

	user := os.Getenv("USER")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	password := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")

	err = validateEnvsDb(user, host, port, password, dbName)
	if err != nil {
		log.Fatal(err)
	}

	return &DbConfig{
		User: user,
		Host: host,
		Port: port,
		Password: password,
		DbName: dbName,
	}
}
