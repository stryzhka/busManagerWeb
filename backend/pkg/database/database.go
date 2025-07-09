package database

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
)

func NewPostgresDatabase() (*sql.DB, error) {
	viper.SetConfigName("db")
	viper.SetConfigType("toml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	info := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.Get("host"), viper.Get("port"), viper.Get("user"), viper.Get("password"), viper.Get("dbname"))
	db, err := sql.Open("postgres", info)

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("success connection")
	return db, err
}
