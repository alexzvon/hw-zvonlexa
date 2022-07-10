package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/alexzvon/hw12_13_14_15_calendar/internal/config"
	"github.com/alexzvon/hw12_13_14_15_calendar/internal/logger"
	"github.com/alexzvon/hw12_13_14_15_calendar/internal/myutils"
	"github.com/pkg/errors"
)

type Server struct {
	srv *http.Server
	app Application
}

type Application interface{}

type sHandler struct {
	logger logger.Logger
}

func (h *sHandler) hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)

	strResponse := "hello-world"
	_, err := w.Write([]byte(strResponse))
	if err != nil {
		if h.logger != nil {
			h.logger.Error(err.Error())
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *sHandler) root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Корень"))
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *sHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch myutils.ConCat(r.Method, " ", r.URL.Path) {
	case "GET /":
		loggingMiddleware(h.logger, h.root)(w, r)
	case "GET /hello":
		loggingMiddleware(h.logger, h.hello)(w, r)
	default:
		h.logger.Error("Not Found")
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func NewServer(cfg config.Config, logger logger.Logger, app Application) *Server {
	handlerHTTP := &sHandler{
		logger: logger,
	}

	server := &http.Server{
		Addr:         myutils.ConCat(cfg.GetString("server.host"), cfg.GetString("server.port")),
		Handler:      handlerHTTP,
		ReadTimeout:  time.Duration(cfg.GetInt("servet.timeout.read") * int(time.Second)),
		WriteTimeout: time.Duration(cfg.GetInt("servet.timeout.write") * int(time.Second)),
	}

	return &Server{
		app: app,
		srv: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil {
		return errors.Wrap(err, "cannot listen http")
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Close(); err != nil {
		return errors.Wrap(err, "cannot close http")
	}

	return nil
}
