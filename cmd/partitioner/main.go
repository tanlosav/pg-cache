package main

import (
	"fmt"

	"github.com/tanlosav/pg-cache/internal/api/rest/partitioner"
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
	schema := db.NewSchema(config, driver)
	router := partitioner.NewRouter(config)
	driver.Connect()
	schema.Init()
	fmt.Printf("Service terminated with status: %s", router.Run())
}
