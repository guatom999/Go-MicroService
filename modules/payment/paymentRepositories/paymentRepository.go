package paymentRepositories

import "go.mongodb.org/mongo-driver/mongo"

type (
	IPaymentRepositoryService interface {
	}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPaymentRepository(db *mongo.Client) IPaymentRepositoryService {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) paymentDbConn() *mongo.Database {
	return r.db.Database("payment_db")
}
