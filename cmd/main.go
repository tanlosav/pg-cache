package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tanlosav/pg-cache/internal/cache/pgcache"
)

type Opts struct {
	Port     string
	User     string
	Password string
	Name     string
	Host     string
}

type CacheHandler struct {
	Cache *pgcache.Cache
}

func main() {
	opts := parseOptions()
	cache := pgcache.NewCache(opts.User, opts.Password, opts.Name, opts.Host)
	cache.Connect()

	handler := CacheHandler{
		Cache: cache,
	}

	router := httprouter.New()
	router.GET("/documents/:id", handler.get)
	router.POST("/documents/:id", handler.create)
	router.PUT("/documents/:id", handler.update)
	router.DELETE("/documents/:id", handler.delete)
	router.DELETE("/documents", handler.clean)

	server := &http.Server{
		Addr:         ":" + opts.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	log.Fatal(server.ListenAndServe())
}

func (handler *CacheHandler) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")

	document, err := handler.Cache.Get(p.ByName("id"))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, document)
}

func (handler *CacheHandler) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")

	document, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Cannot not read body: %s\n", err)
		return
	}

	err = handler.Cache.Create(p.ByName("id"), document)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "")
}

func (handler *CacheHandler) update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")

	document, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Cannot not read body: %s\n", err)
		return
	}

	err = handler.Cache.Update(p.ByName("id"), document)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "")
}

func (handler *CacheHandler) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")

	err := handler.Cache.Delete(p.ByName("id"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "")
}

func (handler *CacheHandler) clean(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")

	err := handler.Cache.Clean()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "")
}

func parseOptions() *Opts {
	port := flag.String("port", "8080", "Port")
	user := flag.String("user", "", "DB user")
	password := flag.String("password", "", "DB password")
	name := flag.String("db", "", "DB name")
	host := flag.String("host", "", "DB host")

	flag.Parse()

	if *port == "" || *user == "" || *password == "" || *name == "" || *host == "" {
		fmt.Printf("Usage: %v [--port <server port>:8080] --user <DB user> --password <DB password> --db <DB name> --host <DB host>\n", os.Args[0])
		os.Exit(1)
	}

	return &Opts{
		Port:     *port,
		User:     *user,
		Password: *password,
		Name:     *name,
		Host:     *host,
	}
}
