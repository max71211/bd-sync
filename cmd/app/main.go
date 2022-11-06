package main

import (
	"aopt-db-sync/datastore/aopt_mysql"
	"aopt-db-sync/datastore/auto_mysql"
	"aopt-db-sync/usecases"
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

	// AUTO Repos
	autoCarMarkRepo := auto_mysql.NewCarMarRepository(autoDB)
	autoCarModelRepo := auto_mysql.NewCarModelRepository(autoDB)
	autoCarModificationRepo := auto_mysql.NewCarModificationRepository(autoDB)
	characteristicRepo := auto_mysql.NewCarCharacteristicRepository(autoDB)

	// AOPT repos
	aoptBrandRepo := aopt_mysql.NewBrandsRepository(aoptDB)
	aoptVehicleRepo := aopt_mysql.NewVehiclesRepository(aoptDB)
	aoptModificationsRepo := aopt_mysql.NewModificationsRepository(aoptDB)

	// endregion datastore

	// region useCases

	aoptUseCase := usecases.NewAoptUseCase(aoptBrandRepo, aoptVehicleRepo, aoptModificationsRepo)
	autoUseCase := usecases.NewAutoUseCase(autoCarMarkRepo, autoCarModelRepo, autoCarModificationRepo, characteristicRepo)

	syncUseCase := usecases.NewSyncUseCase(aoptUseCase, autoUseCase)

	syncUseCase.SyncData(mainCtx, cfg.Main.UpdateBrands, cfg.Main.UpdateVehicles, cfg.Main.UpdateModifications)

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
