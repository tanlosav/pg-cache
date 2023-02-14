package partitioner

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type Router struct {
	config *configuration.Configuration
	router *httprouter.Router
}

func NewRouter(config *configuration.Configuration) *Router {
	return &Router{
		config: config,
		router: httprouter.New(),
	}
}

func (r *Router) Run() error {
	port := strconv.Itoa(r.config.Server.Port)

	r.router.GET("/", r.Ping)
	r.router.GET("/status", r.Status)

	log.Info().Msg("Run server on port " + port)

	return http.ListenAndServe(":"+port, r.router)
}

func (r *Router) Ping(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w = setNoCacheHeaders(w)
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "")
}

func (r *Router) Status(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w = setNoCacheHeaders(w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "up")
}

func setNoCacheHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")

	return w
}
