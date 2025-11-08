package postgress

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/celio001/prodify/config"
)

func NewInstance(dbConfig *config.DbConfig) (*sql.DB, error){
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
	dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName)

	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	log.Print("Conected DB")

	return db, nil
}