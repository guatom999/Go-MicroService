package paymentUseCases

import "github.com/guatom999/Go-MicroService/modules/payment/paymentRepositories"

type (
	IPaymentUseCaseService interface {
	}

	paymentUseCase struct {
		paymentRepo paymentRepositories.IPaymentRepositoryService
	}
)

func NewPaymentUseCase(paymentRepo paymentRepositories.IPaymentRepositoryService) IPaymentUseCaseService {
	return &paymentUseCase{paymentRepo: paymentRepo}
}
