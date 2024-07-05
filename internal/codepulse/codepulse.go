package codepulse

import (
	"database/sql"
	"log"

	"github.com/SathvikPN/codepulse/internal/utils"
)

type codepulse struct {
	name   string
	logger *log.Logger
	db     *sql.DB
}

func StartApplication() {
	app := codepulse{
		name:   "CodePulse",
		logger: utils.GetLogger(),
	}
	app.logger.Println("starting", app.name, "...")

	// connect db

	// start REST server

}
