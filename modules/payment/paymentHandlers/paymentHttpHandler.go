package paymentHandlers

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/payment/paymentUseCases"
)

type (
	IPaymentHttpHandlerService interface {
	}

	paymentHttpHandler struct {
		config         *config.Config
		paymentUseCase paymentUseCases.IPaymentUseCaseService
	}
)

func NewPaymentHttpHandler(config *config.Config, paymentUseCase paymentUseCases.IPaymentUseCaseService) IPaymentHttpHandlerService {
	return &paymentHttpHandler{
		config:         config,
		paymentUseCase: paymentUseCase,
	}
}
