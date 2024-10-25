package main

import (
	"database/sql"
	"flag"
	"juno/cmd/api/run/config"
	nodeHandler "juno/pkg/api/node/handler"
	nodeMig "juno/pkg/api/node/migration/mysql"
	nodePolicy "juno/pkg/api/node/policy"
	nodeRepo "juno/pkg/api/node/repo/mysql"
	nodeSvc "juno/pkg/api/node/service"
	"log"

	tranHandler "juno/pkg/api/transaction/handler"
	tranMig "juno/pkg/api/transaction/migration/mysql"
	tranRepo "juno/pkg/api/transaction/repo/mysql"
	tranSvc "juno/pkg/api/transaction/service"

	tokenHandler "juno/pkg/api/token/handler"
	tokenService "juno/pkg/api/token/service"

	extractorJobHandler "juno/pkg/api/extractor/job/handler"
	extractorJobMig "juno/pkg/api/extractor/job/migration/mysql"
	extractorJobPolicy "juno/pkg/api/extractor/job/policy"
	extractorJobRepo "juno/pkg/api/extractor/job/repo/mysql"
	extractorJobSvc "juno/pkg/api/extractor/job/service"

	selectorHandler "juno/pkg/api/extractor/selector/handler"
	selectorMig "juno/pkg/api/extractor/selector/migration/mysql"
	selectorPolicy "juno/pkg/api/extractor/selector/policy"
	selectorRepo "juno/pkg/api/extractor/selector/repo/mysql"
	selectorService "juno/pkg/api/extractor/selector/service"

	filterHandler "juno/pkg/api/extractor/filter/handler"
	filterMig "juno/pkg/api/extractor/filter/migration/mysql"
	filterPolicy "juno/pkg/api/extractor/filter/policy"
	filterRepo "juno/pkg/api/extractor/filter/repo/mysql"
	filterService "juno/pkg/api/extractor/filter/service"

	fieldHandler "juno/pkg/api/extractor/field/handler"
	fieldMig "juno/pkg/api/extractor/field/migration/mysql"
	fieldPolicy "juno/pkg/api/extractor/field/policy"
	fieldRepo "juno/pkg/api/extractor/field/repo/mysql"
	fieldService "juno/pkg/api/extractor/field/service"

	strategyHandler "juno/pkg/api/extractor/strategy/handler"
	strategyMig "juno/pkg/api/extractor/strategy/migration/mysql"
	strategyPolicy "juno/pkg/api/extractor/strategy/policy"
	strategyRepo "juno/pkg/api/extractor/strategy/repo/strategy/mysql"
	strategyService "juno/pkg/api/extractor/strategy/service"

	strategyFieldRepo "juno/pkg/api/extractor/strategy/repo/field/mysql"
	strategyFilterRepo "juno/pkg/api/extractor/strategy/repo/filter/mysql"
	strategySelectorRepo "juno/pkg/api/extractor/strategy/repo/selector/mysql"

	balancerHandler "juno/pkg/api/balancer/handler"
	balancerMig "juno/pkg/api/balancer/migration/mysql"
	balancerPolicy "juno/pkg/api/balancer/policy"
	balancerRepo "juno/pkg/api/balancer/repo/mysql"
	balancerSvc "juno/pkg/api/balancer/service"

	userHandler "juno/pkg/api/user/handler"
	userMig "juno/pkg/api/user/migration/mysql"
	userPolicy "juno/pkg/api/user/policy"
	userRepo "juno/pkg/api/user/repo/mysql"
	userSvc "juno/pkg/api/user/service"

	authHandler "juno/pkg/api/auth/handler"
	authSvc "juno/pkg/api/auth/service"

	"juno/pkg/api/router"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sirupsen/logrus"
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
	nodeDB := setupDatabase(config.NodeDB, nodeMig.ExecuteMigrations)
	userDB := setupDatabase(config.UserDB, userMig.ExecuteMigrations)
	balancerDB := setupDatabase(config.BalancerDB, balancerMig.ExecuteMigrations)
	tranDB := setupDatabase(config.TranDB, tranMig.ExecuteMigrations)
	extractionJobDB := setupDatabase(config.ExtractionJobDB, extractorJobMig.ExecuteMigrations)
	selectorDB := setupDatabase(config.SelectorDB, selectorMig.ExecuteMigrations)
	filterDB := setupDatabase(config.FilterDB, filterMig.ExecuteMigrations)
	fieldDB := setupDatabase(config.FieldDB, fieldMig.ExecuteMigrations)
	strategyDB := setupDatabase(config.StrategyDB, strategyMig.ExecuteMigrations)

	logger := logrus.New()

	nodeRepo := nodeRepo.New(nodeDB)
	nodeSvc := nodeSvc.New(nodeRepo)
	nodePolicy := nodePolicy.New()
	nodeHandler := nodeHandler.New(logger, nodePolicy, nodeSvc)

	tranRepo := tranRepo.New(tranDB)
	tranSvc := tranSvc.New(logger, tranRepo)
	tranHandler := tranHandler.New(tranSvc)

	tokenSvc := tokenService.New(tranSvc)
	tokenHandler := tokenHandler.New(logger, tokenSvc)

	balancerRepo := balancerRepo.New(balancerDB)
	balancerSvc := balancerSvc.New(balancerRepo)
	balancerPolicy := balancerPolicy.New()
	balancerHandler := balancerHandler.New(logger, balancerPolicy, balancerSvc)

	selectorRepo := selectorRepo.New(selectorDB)
	selectorSvc := selectorService.New(selectorRepo)
	selectorPolicy := selectorPolicy.New()
	selectorHandler := selectorHandler.New(selectorPolicy, selectorSvc)

	filterRepo := filterRepo.New(filterDB)
	filterSvc := filterService.New(filterRepo)
	filterPolicy := filterPolicy.New()
	filterHandler := filterHandler.New(filterPolicy, filterSvc)

	fieldRepo := fieldRepo.New(fieldDB)
	fieldSvc := fieldService.New(fieldRepo)
	fieldPolicy := fieldPolicy.New()
	fieldHandler := fieldHandler.New(fieldPolicy, fieldSvc)

	strategyRepo := strategyRepo.New(strategyDB)
	strategySelectorRepo := strategySelectorRepo.New(strategyDB)
	strategyFieldRepo := strategyFieldRepo.New(strategyDB)
	strategyFilterRepo := strategyFilterRepo.New(strategyDB)
	strategySvc := strategyService.New(strategyRepo, strategyFilterRepo, strategyFieldRepo, strategySelectorRepo, filterSvc, fieldSvc, selectorSvc)
	strategyPolicy := strategyPolicy.New()
	strategyHandler := strategyHandler.New(strategyPolicy, strategySvc)

	extractionJobRepo := extractorJobRepo.New(extractionJobDB)
	extractionJobSvc := extractorJobSvc.New(extractionJobRepo, strategySvc)
	extractionJobPolicy := extractorJobPolicy.New()
	extractionJobHandler := extractorJobHandler.New(extractionJobSvc, extractionJobPolicy)

	userRepo := userRepo.New(userDB)

	userSvc := userSvc.New(logger, userRepo)
	policy := userPolicy.New()
	userHandler := userHandler.New(logger, policy, userSvc)

	authSvc := authSvc.New(logger, userSvc)
	authHandler := authHandler.New(logger, authSvc)

	r := router.New(
		nodeHandler,
		balancerHandler,
		tranHandler,
		extractionJobHandler,
		selectorHandler,
		filterHandler,
		fieldHandler,
		strategyHandler,
		tokenHandler,
		userHandler,
		authHandler,
	)

	r.Run(":" + portFlag)
}
