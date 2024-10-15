package main

import (
	"database/sql"
	"flag"
	nodeHandler "juno/pkg/api/node/handler"
	nodeMig "juno/pkg/api/node/migration/mysql"
	nodePolicy "juno/pkg/api/node/policy"
	nodeRepo "juno/pkg/api/node/repo/mysql"
	nodeSvc "juno/pkg/api/node/service"

	balancerHandler "juno/pkg/api/balancer/handler"
	balancerMig "juno/pkg/api/balancer/migration/mysql"
	balancerPolicy "juno/pkg/api/balancer/policy"
	balancerRepo "juno/pkg/api/balancer/repo/mysql"
	balancerSvc "juno/pkg/api/balancer/service"

	assignmentHandler "juno/pkg/api/assignment/handler"
	assignmentMig "juno/pkg/api/assignment/migration/mysql"
	assignmentPolicy "juno/pkg/api/assignment/policy"
	assignmentRepo "juno/pkg/api/assignment/repo/mysql"
	assignmentSvc "juno/pkg/api/assignment/service"

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

func setupDatabases() (*sql.DB, *sql.DB, *sql.DB, *sql.DB) {
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

	assignmentDB, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/assignment?parseTime=true")

	if err != nil {
		panic(err)
	}

	err = assignmentMig.ExecuteMigrations(assignmentDB)

	if err != nil {
		panic(err)
	}

	return nodeDB, userDB, balancerDB, assignmentDB
}

func main() {

	var portFlag string
	flag.StringVar(&portFlag, "port", "8080", "port to run the server on")

	flag.Parse()

	nodeDB, userDB, balancerDB, assignmentDB := setupDatabases()

	logger := logrus.New()

	nodeRepo := nodeRepo.New(nodeDB)
	nodeSvc := nodeSvc.New(nodeRepo)
	nodePolicy := nodePolicy.New()
	nodeHandler := nodeHandler.New(logger, nodePolicy, nodeSvc)

	userRepo := userRepo.New(userDB)

	userSvc := userSvc.New(logger, userRepo)
	policy := userPolicy.New()
	userHandler := userHandler.New(logger, policy, userSvc)

	authSvc := authSvc.New(logger, userSvc)
	authHandler := authHandler.New(logger, authSvc)

	balancerRepo := balancerRepo.New(balancerDB)
	balancerSvc := balancerSvc.New(balancerRepo)
	balancerPolicy := balancerPolicy.New()
	balancerHandler := balancerHandler.New(logger, balancerPolicy, balancerSvc)

	assignmentRepo := assignmentRepo.New(assignmentDB)
	assignmentSvc := assignmentSvc.New(assignmentRepo)
	assignmentPolicy := assignmentPolicy.New()
	assignmentHandler := assignmentHandler.New(logger, assignmentPolicy, assignmentSvc)

	r := router.New(
		nodeHandler,
		userHandler,
		authHandler,
		balancerHandler,
		assignmentHandler,
	)

	r.Run(":" + portFlag)
}
