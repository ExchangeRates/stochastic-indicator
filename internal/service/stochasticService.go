package service

import (
	"math"
	"stochastic_indicator/internal/feign"

	"github.com/sirupsen/logrus"
)

type StochasticService interface {
	Calculate(value float64, prevKPoint *float64, period int, lastPoints []float64) (kPoint, dPoint float64, lastResPoints []float64, err error)
}

type stochasticServiceImpl struct {
	epsilon   float64
	emaClient feign.EmaFeignClient
}

func NewKStochasticService(emaClient feign.EmaFeignClient) StochasticService {
	return &stochasticServiceImpl{
		epsilon:   0.000001,
		emaClient: emaClient,
	}
}

func (s *stochasticServiceImpl) Calculate(value float64, prevKPoint *float64, period int, lastPoints []float64) (kPoint, dPoint float64, lastResPoints []float64, err error) {
	if len(lastPoints) < period {
		lastPoints = append(lastPoints, value)
	} else {
		lastPoints = append(lastPoints, value)
		lastPoints = append(lastPoints[:0], lastPoints[1:]...)
	}

	highest, lowest := s.highestLowest(lastPoints)
	logrus.Infof("H: %d, L: %d", highest, lowest)
	kPoint = s.calcualteKPoint(value, highest, lowest)
	dPoint, err = s.emaClient.Calculate(prevKPoint, kPoint, period)
	if err != nil {
		return 0, 0, []float64{}, err
	}

	return kPoint, dPoint, lastPoints, nil

}

func (s *stochasticServiceImpl) calcualteKPoint(value, highest, lowest float64) float64 {
	divider := highest - lowest
	if divider < s.epsilon {
		return 0
	}
	return ((value - lowest) / divider) * 100
}

func (s *stochasticServiceImpl) highestLowest(points []float64) (float64, float64) {
	highest := 0.
	lowest := math.MaxFloat64
	for _, point := range points {
		if point > highest {
			highest = point
		}
		if point < lowest {
			lowest = point
		}
	}
	return highest, lowest
}
