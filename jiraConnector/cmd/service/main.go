package main

import (
	"jiraAnalyzer/jiraConnector/cmd/service/internal/app"
	"jiraAnalyzer/jiraConnector/cmd/service/internal/config"
	"jiraAnalyzer/jiraConnector/internal/repository/database"
)

func main() {
	cfg, err := config.LoadConfig(*ConfigPathFlag)
	if err != nil {
		panic(err)
	}

	newApp, db, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			panic(err)
		}
	}()
	
	newApp.Run()
	defer newApp.Close()

}
