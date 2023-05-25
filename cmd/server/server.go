package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pr02nl/medidor_ade/configs"
	"github.com/pr02nl/medidor_ade/internal/infra/database"
	_ "github.com/sijms/go-ora/v2"
)

//(description= (retry_count=20)
// (retry_delay=3)(address=(protocol=tcps)(port=1522)
// (host=adb.us-ashburn-1.oraclecloud.com))
// (connect_data=(service_name=gff6197df21179c_owk5v277wxb677sd_high.adb.oraclecloud.com))(security=(ssl_server_dn_match=yes)))

func main() {
	println("Starting...")
	t := time.Now()
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		configs.DBDriver,
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
