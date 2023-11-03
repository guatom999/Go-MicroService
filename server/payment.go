package server

import (
	"github.com/guatom999/Go-MicroService/modules/payment/paymentHandlers"
	"github.com/guatom999/Go-MicroService/modules/payment/paymentRepositories"
	"github.com/guatom999/Go-MicroService/modules/payment/paymentUseCases"
)

func (s *server) paymentService() {
	paymentRepository := paymentRepositories.NewPaymentRepository(s.db)
	paymentUseCase := paymentUseCases.NewPaymentUseCase(paymentRepository)
	paymentHtppHandler := paymentHandlers.NewPaymentHttpHandler(s.cfg, paymentUseCase)
	paymentGrpcHandler := paymentHandlers.NewPaymentGrpcHandler(paymentUseCase)

	_ = paymentHtppHandler
	_ = paymentGrpcHandler

	payment := s.app.Group("/payment_v1")

	// Health Check
	payment.GET("", s.healthCheckService)
}
