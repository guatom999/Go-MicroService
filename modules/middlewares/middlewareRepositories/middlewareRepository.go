package middlewareRepositories

type (
	IMiddlewareRepositoryService interface {
	}

	middlewareRepository struct {
	}
)

func NewMiddlewareRepository() IMiddlewareRepositoryService {
	return &middlewareRepository{}
}
