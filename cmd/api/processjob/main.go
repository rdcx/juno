package main

import (
	"database/sql"
	"flag"
	"juno/cmd/api/run/config"
	"log"

	extractorJobMig "juno/pkg/api/extractor/job/migration/mysql"
	extractorJobRepo "juno/pkg/api/extractor/job/repo/mysql"
	extractorJobSvc "juno/pkg/api/extractor/job/service"

	selectorMig "juno/pkg/api/extractor/selector/migration/mysql"
	selectorRepo "juno/pkg/api/extractor/selector/repo/mysql"
	selectorService "juno/pkg/api/extractor/selector/service"

	filterMig "juno/pkg/api/extractor/filter/migration/mysql"
	filterRepo "juno/pkg/api/extractor/filter/repo/mysql"
	filterService "juno/pkg/api/extractor/filter/service"

	fieldMig "juno/pkg/api/extractor/field/migration/mysql"
	fieldRepo "juno/pkg/api/extractor/field/repo/mysql"
	fieldService "juno/pkg/api/extractor/field/service"

	strategyMig "juno/pkg/api/extractor/strategy/migration/mysql"
	strategyRepo "juno/pkg/api/extractor/strategy/repo/strategy/mysql"
	strategyService "juno/pkg/api/extractor/strategy/service"

	strategyFieldRepo "juno/pkg/api/extractor/strategy/repo/field/mysql"
	strategyFilterRepo "juno/pkg/api/extractor/strategy/repo/filter/mysql"
	strategySelectorRepo "juno/pkg/api/extractor/strategy/repo/selector/mysql"

	ranagMig "juno/pkg/api/ranag/migration/mysql"
	ranagRepo "juno/pkg/api/ranag/repo/mysql"
	ranagSvc "juno/pkg/api/ranag/service"

	_ "github.com/go-sql-driver/mysql"
)

func setupDatabase(connectionString string, migrations func(*sql.DB) error) *sql.DB {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := migrations(db); err != nil {
		log.Fatalf("failed to execute migrations: %v", err)
	}

	return db
}

func main() {

	var portFlag string
	flag.StringVar(&portFlag, "port", "8080", "port to run the server on")

	flag.Parse()

	config := config.LoadConfig()
	extractionJobDB := setupDatabase(config.ExtractionJobDB, extractorJobMig.ExecuteMigrations)
	selectorDB := setupDatabase(config.SelectorDB, selectorMig.ExecuteMigrations)
	filterDB := setupDatabase(config.FilterDB, filterMig.ExecuteMigrations)
	fieldDB := setupDatabase(config.FieldDB, fieldMig.ExecuteMigrations)
	strategyDB := setupDatabase(config.StrategyDB, strategyMig.ExecuteMigrations)
	ranagDB := setupDatabase(config.RanagDB, ranagMig.ExecuteMigrations)

	ranagRepo := ranagRepo.New(ranagDB)
	ranagSvc := ranagSvc.New(ranagRepo)

	selectorRepo := selectorRepo.New(selectorDB)
	selectorSvc := selectorService.New(selectorRepo)

	filterRepo := filterRepo.New(filterDB)
	filterSvc := filterService.New(filterRepo)

	fieldRepo := fieldRepo.New(fieldDB)
	fieldSvc := fieldService.New(fieldRepo)

	strategyRepo := strategyRepo.New(strategyDB)
	strategySelectorRepo := strategySelectorRepo.New(strategyDB)
	strategyFieldRepo := strategyFieldRepo.New(strategyDB)
	strategyFilterRepo := strategyFilterRepo.New(strategyDB)
	strategySvc := strategyService.New(strategyRepo, strategyFilterRepo, strategyFieldRepo, strategySelectorRepo, filterSvc, fieldSvc, selectorSvc)

	extractionJobRepo := extractorJobRepo.New(extractionJobDB)
	extractionJobSvc := extractorJobSvc.New(extractionJobRepo, strategySvc, ranagSvc)

	extractionJobSvc.ProcessPending()
}
