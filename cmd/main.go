package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/internal/cache/pgcache"
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type CmdLineOpts struct {
	ConfigurationProvider string
	ConfigurationSource   string
}

// type CacheHandler struct {
// 	Cache *pgcache.Cache
// }

func main() {
	opts := parseOptions()
	setupLogger()
	config := getConfiguration(opts)
	pgcache.NewCache(config)

	// todo: start router

	// handler := CacheHandler{
	// 	Cache: cache,
	// }

	// router := httprouter.New()
	// router.GET("/documents/:id", handler.get)
	// router.POST("/documents/:id", handler.create)
	// router.PUT("/documents/:id", handler.update)
	// router.DELETE("/documents/:id", handler.delete)
	// router.DELETE("/documents", handler.clean)

	// server := &http.Server{
	// 	Addr:         ":" + opts.Port,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	Handler:      router,
	// }

	// log.Fatal(server.ListenAndServe())
}

func parseOptions() *CmdLineOpts {
	configProvider := flag.String("configuration-provider", "file", "Configuration provider")
	configSource := flag.String("configuration-source", "", "Configuration source")

	flag.Parse()

	return &CmdLineOpts{
		ConfigurationProvider: *configProvider,
		ConfigurationSource:   *configSource,
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true, TimeFormat: time.RFC3339Nano}).With().Caller().Logger()
}

func getConfiguration(opts *CmdLineOpts) configuration.Configuration {
	var configProvider configuration.ConfigurationSource

	switch opts.ConfigurationProvider {
	case "file":
		configProvider = configuration.NewFileSource()
	default:
		panic("Unsupported configuration provider")
	}

	config := configProvider.Configuration(opts.ConfigurationSource)

	log.Printf("Loaded configuration: %+v", config)

	return config
}

// func (handler *CacheHandler) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Pragma", "no-cache")
// 	w.Header().Set("Cache-Control", "no-store")

// 	document, err := handler.Cache.Get(p.ByName("id"))

// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, document)
// }

// func (handler *CacheHandler) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Pragma", "no-cache")
// 	w.Header().Set("Cache-Control", "no-store")

// 	document, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Cannot not read body: %s\n", err)
// 		return
// 	}

// 	err = handler.Cache.Create(p.ByName("id"), document)

// 	if err != nil {
// 		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
// 			w.WriteHeader(http.StatusConflict)
// 			fmt.Fprint(w, err.Error())
// 			return
// 		}

// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprint(w, "")
// }

// func (handler *CacheHandler) update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Pragma", "no-cache")
// 	w.Header().Set("Cache-Control", "no-store")

// 	document, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Cannot not read body: %s\n", err)
// 		return
// 	}

// 	err = handler.Cache.Update(p.ByName("id"), document)

// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// 	fmt.Fprint(w, "")
// }

// func (handler *CacheHandler) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Pragma", "no-cache")
// 	w.Header().Set("Cache-Control", "no-store")

// 	err := handler.Cache.Delete(p.ByName("id"))

// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// 	fmt.Fprintf(w, "")
// }

// func (handler *CacheHandler) clean(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Pragma", "no-cache")
// 	w.Header().Set("Cache-Control", "no-store")

// 	err := handler.Cache.Clean()

// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// 	fmt.Fprintf(w, "")
// }
