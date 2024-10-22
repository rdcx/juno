package main

import (
	"database/sql"
	"flag"
	nodeHandler "juno/pkg/api/node/handler"
	nodeMig "juno/pkg/api/node/migration/mysql"
	nodePolicy "juno/pkg/api/node/policy"
	nodeRepo "juno/pkg/api/node/repo/mysql"
	nodeSvc "juno/pkg/api/node/service"

	tranHandler "juno/pkg/api/transaction/handler"
	tranMig "juno/pkg/api/transaction/migration/mysql"
	tranRepo "juno/pkg/api/transaction/repo/mysql"
	tranSvc "juno/pkg/api/transaction/service"

	tokenHandler "juno/pkg/api/token/handler"
	tokenService "juno/pkg/api/token/service"

	extractorService "juno/pkg/api/extraction/extractor/service"

	extractionJobHandler "juno/pkg/api/extraction/job/handler"
	extractionJobMig "juno/pkg/api/extraction/job/migration/mysql"
	extractionJobPolicy "juno/pkg/api/extraction/job/policy"
	extractionJobRepo "juno/pkg/api/extraction/job/repo/mysql"
	extractionJobSvc "juno/pkg/api/extraction/job/service"

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

func setupDatabases() (
	*sql.DB, *sql.DB, *sql.DB, *sql.DB, *sql.DB) {
	nodeDB, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/node?parseTime=true")

	if err != nil {
		panic(err)
	}

	err = nodeMig.ExecuteMigrations(nodeDB)

	if err != nil {
		panic(err)
	}

	userDB, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user?parseTime=true")

	if err != nil {
		panic(err)
	}

	err = userMig.ExecuteMigrations(userDB)

	if err != nil {
		panic(err)
	}

	balancerDB, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/balancer?parseTime=true")

	if err != nil {
		panic(err)
	}

	err = balancerMig.ExecuteMigrations(balancerDB)

	if err != nil {
		panic(err)
	}

	tranDB, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/transaction?parseTime=true")

	if err != nil {
		panic(err)
	}

	err = tranMig.ExecuteMigrations(tranDB)

	if err != nil {
		panic(err)
	}

	extractionJobDB, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/extraction_job?parseTime=true")

	if err != nil {
		panic(err)
	}

	err = extractionJobMig.ExecuteMigrations(extractionJobDB)

	if err != nil {
		panic(err)
	}

	return nodeDB, userDB, balancerDB, tranDB, extractionJobDB
}

func main() {

	var portFlag string
	flag.StringVar(&portFlag, "port", "8080", "port to run the server on")

	flag.Parse()

	nodeDB, userDB, balancerDB, tranDB, extractionJobDB := setupDatabases()

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

	extractorSvc := extractorService.New(nil)

	extractionJobRepo := extractionJobRepo.New(extractionJobDB)
	extractionJobSvc := extractionJobSvc.New(extractionJobRepo, extractorSvc)
	extractionJobPolicy := extractionJobPolicy.New()
	extractionJobHandler := extractionJobHandler.New(extractionJobSvc, extractionJobPolicy)

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
		tokenHandler,
		userHandler,
		authHandler,
	)

	r.Run(":" + portFlag)
}
