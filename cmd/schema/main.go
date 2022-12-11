package main

import (
	"github.com/tanlosav/pg-cache/internal/cmd"
	"github.com/tanlosav/pg-cache/internal/configuration"
	"github.com/tanlosav/pg-cache/internal/db"
	"github.com/tanlosav/pg-cache/internal/logger"
)

func main() {
	opts := cmd.ParseOptions()
	logger.SetupLogger()
	config := configuration.NewConfiguration(opts)
	driver := db.NewDriver(config)
	driver.Connect()
	schema := db.NewSchema(config, driver)
	schema.Init()
}
