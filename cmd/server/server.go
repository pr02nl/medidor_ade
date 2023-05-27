package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pr02nl/medidor_ade/configs"
	"github.com/pr02nl/medidor_ade/internal/infra/database"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	println("Starting...")
	t := time.Now()
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		configs.DBUser,
		configs.DBPassword,
		configs.DBHost,
		configs.DBPort,
		configs.DBName))
	if err != nil {
		panic(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	medidorRepository := database.NewMedidorRepository(db)
	err = medidorRepository.InitTable()
	if err != nil {
		log.Fatal(err)
	}
	medicaoRepository := database.NewMedicaoRepository(db)
	err = medicaoRepository.InitTable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Time Elapsed", time.Since(t).Milliseconds())
}
