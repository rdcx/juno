package main

import (
	"database/sql"
	"flag"
	nodeHandler "juno/pkg/api/node/handler"
	nodeMig "juno/pkg/api/node/migration/mysql"
	nodeRepo "juno/pkg/api/node/repo/mysql"
	nodeSvc "juno/pkg/api/node/service"

	userHandler "juno/pkg/api/user/handler"
	userMig "juno/pkg/api/user/migration/mysql"
	userPolicy "juno/pkg/api/user/policy"
	userRepo "juno/pkg/api/user/repo/mysql"
	userSvc "juno/pkg/api/user/service"

	"juno/pkg/api/router"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sirupsen/logrus"
)

func main() {

	var portFlag string
	flag.StringVar(&portFlag, "port", "8080", "port to run the server on")

	flag.Parse()

	nodeDB, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/node?parseTime=true")

	if err != nil {
		panic(err)
	}

	err = nodeMig.ExecuteMigrations(nodeDB)

	if err != nil {
		panic(err)
	}

	logger := logrus.New()

	nodeRepo := nodeRepo.New(nodeDB)
	nodeSvc := nodeSvc.New(nodeRepo)
	nodeHandler := nodeHandler.New(logger, nodeSvc)

	userRepo := userRepo.New(nodeDB)

	err = userMig.ExecuteMigrations(nodeDB)

	if err != nil {
		panic(err)
	}

	userSvc := userSvc.New(logger, userRepo)
	policy := userPolicy.New()
	userHandler := userHandler.New(logger, policy, userSvc)

	r := router.New(
		nodeHandler,
		userHandler,
	)

	r.Run(":" + portFlag)
}
