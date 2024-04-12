package main

import (
	"flag"
	"fmt"
	"github.com/Vadim992/avito/internal"
	"github.com/Vadim992/avito/internal/cfg"
	"github.com/Vadim992/avito/internal/mws"
	"github.com/Vadim992/avito/internal/postgres"
	"github.com/Vadim992/avito/internal/storage"
	"github.com/Vadim992/avito/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		logger.ErrLog.Fatal(err)
	}

	cfgDB := cfg.NewCfgDB()
	cfgDB.SetFromEnv()

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		" dbname=%s sslmode=disable",
		cfgDB.HostDB, cfgDB.PortDB, cfgDB.UsernameDB, cfgDB.PasswordDB, cfgDB.NameDB)

	addr := flag.String("addr", ":3000", "HTTP network address")
	flag.Parse()

	db, err := postgres.InitDB(conn)

	if err != nil {
		logger.ErrLog.Fatalf("cannot connect to postgres: %v", err)
	}

	cfgTokens := cfg.NewCfgTokens()
	cfgTokens.SetFromEnv()

	if cfgTokens.AdminToken == "" || cfgTokens.UserToken == "" {
		logger.ErrLog.Fatalf("dont have tokens in .env file: %v", err)
	}

	tokenMap := map[string]int{
		cfgTokens.AdminToken: mws.ADMIN,
		cfgTokens.UserToken:  mws.USER,
	}

	DB := postgres.NewDB(db)

	//if err := DB.FillDb(); err != nil {
	//	logger.ErrLog.Fatalf("cannot fill DB: %v", err)
	//}

	inMemory := storage.NewStorage()
	err = DB.FillStorage(inMemory)

	if err != nil {
		logger.ErrLog.Fatalf("cant save data to inMemory storage %v:", err)
	}
	fmt.Println(inMemory.Get(1000, 1000))

	app := internal.NewApp(DB, inMemory, tokenMap)

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.Routes(),
		ErrorLog: logger.ErrLog,
	}

	logger.InfoLog.Printf("Starting server on port %s\n", *addr)

	logger.ErrLog.Fatal(srv.ListenAndServe())
}
