package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/config"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/helper"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
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
	switch r.Method {
	case "GET":
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
	default:
		if h.logger != nil {
			h.logger.Error("Only GET allowed")
		}

		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}

func (h *sHandler) root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("Корень"))
		if err != nil {
			h.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	default:
		h.logger.Error("Only GET and POST are allowed")
		http.Error(w, "Only GET and POST are allowed", http.StatusMethodNotAllowed)
	}
}

func NewServer(cfg config.Config, logger logger.Logger, app Application) *Server {
	handler := &sHandler{
		logger: logger,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", loggingMiddleware(logger, handler.root))
	mux.HandleFunc("/hello", loggingMiddleware(logger, handler.hello))

	server := &http.Server{
		Addr:         helper.ConCat(cfg.GetString("server.host"), cfg.GetString("server.port")),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.GetInt("servet.timeout.read") * int(time.Second)),
		WriteTimeout: time.Duration(cfg.GetInt("servet.timeout.write") * int(time.Second)),
	}

	return &Server{
		app: app,
		srv: server,
		//		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		if err := s.Stop(ctx); err != nil {
			return err
		}
	default:
		if err := s.srv.ListenAndServe(); err != nil {
			return errors.Wrap(err, "cannot listen http")
		}
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Close(); err != nil {
		return errors.Wrap(err, "cannot close http")
	}

	return nil
}
