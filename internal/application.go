package internal

import (
	"net/http"
	"stochastic_indicator/internal/api"
	"stochastic_indicator/internal/config"
	"stochastic_indicator/internal/controller"
	"stochastic_indicator/internal/feign"
	"stochastic_indicator/internal/service"
)

func Start(config *config.Config) error {

	emaFeignClient := feign.NewEmaFeignClient(config.EmaClientURL)
	stochasticService := service.NewKStochasticService(emaFeignClient)
	stochasticController := controller.NewStochasticController(stochasticService)
	srv := api.NewServer(stochasticController)
	bindingAddress := srv.BindingAddressFromPort(config.Port)

	return http.ListenAndServe(bindingAddress, srv)
}
