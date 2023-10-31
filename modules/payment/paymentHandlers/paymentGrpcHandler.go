package paymentHandlers

import "github.com/guatom999/Go-MicroService/modules/payment/paymentUseCases"

type (
	paymentGrpcHandler struct {
		paymentUseCase paymentUseCases.IPaymentUseCaseService
	}
)

func NewpaymentGrpcHandler(paymentUseCase paymentUseCases.IPaymentUseCaseService) *paymentGrpcHandler {
	return &paymentGrpcHandler{paymentUseCase: paymentUseCase}
}
