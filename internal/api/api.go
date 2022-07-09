package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"stochastic_indicator/internal/controller"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Server interface {
	GracefullListenAndServe(port int) error
}

type server struct {
	router     *mux.Router
	log        *logrus.Logger
	controller *controller.StochasticController
}

func NewServer(controller *controller.StochasticController) Server {
	server := &server{
		router:     mux.NewRouter(),
		log:        logrus.New(),
		controller: controller,
	}

	server.configureRouter()

	server.log.Info("starting api server")

	return server
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) bindingAddressFromPort(port int) string {
	s.log.Infof("starting on port %d", port)
	return fmt.Sprintf(":%d", port)
}

func (s *server) configureRouter() {
	s.router.Path("/calculate").Handler(s.controller.HandleCalculate()).Methods(http.MethodPost)
}

func (s *server) GracefullListenAndServe(port int) error {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := http.Server{
		Addr:    s.bindingAddressFromPort(port),
		Handler: s,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	return g.Wait()
}
