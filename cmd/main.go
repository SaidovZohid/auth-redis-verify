package main

import (
	"fmt"
	"log"

	"github.com/SaidovZohid/auth-redis-verify/api"
	"github.com/SaidovZohid/auth-redis-verify/config"
	"github.com/SaidovZohid/auth-redis-verify/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println("Configuration: ", cfg)
	fmt.Println("Connected Succesfully!")

	strg := storage.NewStorage(psqlConn)

	server := api.New(&api.RouteOption{
		Cfg: &cfg,
		Storage: strg,
	})

	err = server.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
