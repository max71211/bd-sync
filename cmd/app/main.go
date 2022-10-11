package main

import (
	"aopt-db-sync/datastore/auto_mysql"
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cfg Config

func main() {
	mainCtx := context.Background()

	// region configuration

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Load .env error:", zap.Error(err))
	}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal("Parse .env error:", zap.Error(err))
	}

	// endregion configuration

	// region datastore

	aoptDBURL := fmt.Sprintf("%s:%s@(%s:%s)/%s", cfg.AoptDB.User, cfg.AoptDB.Password, cfg.AoptDB.Host, cfg.AoptDB.Port, cfg.AoptDB.Database)
	aoptDB, err := sqlx.Connect("mysql", aoptDBURL)
	if err != nil {
		log.Fatal("connect to aopt error:", zap.Error(err))
	}
	defer aoptDB.Close()

	autoDBURL := fmt.Sprintf("%s:%s@(%s:%s)/%s", cfg.AutoDB.User, cfg.AutoDB.Password, cfg.AutoDB.Host, cfg.AutoDB.Port, cfg.AutoDB.Database)
	autoDB, err := sqlx.Connect("mysql", autoDBURL)
	if err != nil {
		log.Fatal("connect to auto error:", zap.Error(err))
	}
	defer autoDB.Close()

	autoCarMarkRepo := auto_mysql.NewCarMarRepository(autoDB)

	carMarks, err := autoCarMarkRepo.GetAll(mainCtx)
	if err != nil {
		log.Println("Get car marks error:", zap.Error(err))
	}

	log.Println(fmt.Sprintf("LEN %d", len(carMarks)))
	for _, cm := range carMarks {
		log.Println(cm.Name)
	}

	// endregion datastore

	// region useCases

	// endregion useCases

	// region waiter

	log.Println("the service is started")
	defer log.Println("the service is stopped")

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sig

	// endregion waiter
}
