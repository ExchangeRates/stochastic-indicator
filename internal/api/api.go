package api

import (
	"fmt"
	"net/http"
	"stochastic_indicator/internal/controller"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router     *mux.Router
	log        *logrus.Logger
	controller *controller.StochasticController
}

func NewServer(controller *controller.StochasticController) *server {
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

func (s *server) BindingAddressFromPort(port int) string {
	return fmt.Sprintf(":%d", port)
}

func (s *server) configureRouter() {
	s.router.Path("/calcualte").Handler(s.controller.HandleCalculate()).Methods(http.MethodPost)
}
