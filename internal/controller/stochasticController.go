package controller

import (
	"encoding/json"
	"net/http"
	"stochastic_indicator/internal/service"

	"github.com/sirupsen/logrus"
)

type StochasticController struct {
	stochasticSerivce service.StochasticService
	log               logrus.Logger
}

func NewStochasticController(stochasticService service.StochasticService) *StochasticController {
	return &StochasticController{
		stochasticSerivce: stochasticService,
		log:               *logrus.New(),
	}
}

func (c *StochasticController) HandleCalculate() http.HandlerFunc {
	type request struct {
		Value      float64   `json:"value"`
		Period     int       `json:"period"`
		PrevKPoint *float64  `json:"prevKPoint"`
		ListPoints []float64 `json:"listPoints"`
	}
	type response struct {
		KPoint     float64   `json:"kPoint"`
		DPoint     float64   `json:"dPoint"`
		LastPoints []float64 `json:"lastPoints"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body := &request{}
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			// TODO send response
			c.log.Errorln(err)
			return
		}

		kPoint, dPoint, lastPoints, err := c.stochasticSerivce.Calculate(
			body.Value,
			body.PrevKPoint,
			body.Period,
			body.ListPoints,
		)
		if err != nil {
			// TODO send response
			c.log.Fatalln(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response{
			KPoint:     kPoint,
			DPoint:     dPoint,
			LastPoints: lastPoints,
		}); err != nil {
			// TODO send response
			c.log.Fatalln(err)
		}
	}
}
