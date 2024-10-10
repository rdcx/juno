package main

import (
	"database/sql"
	"flag"
	nodeHandler "juno/pkg/api/node/handler"
	nodeMig "juno/pkg/api/node/migration/mysql"
	nodeRepo "juno/pkg/api/node/repo/mysql"
	nodeSvc "juno/pkg/api/node/service"
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

	r := router.New(
		nodeHandler,
	)

	r.Run(":" + portFlag)
}
