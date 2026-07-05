package usecase

import "github.com/bytedance/go-noba-trial/internal/domain"

type HealthUsecase interface {
	Ping() domain.HealthStatus
}

type healthUsecase struct {
	serviceName string
}

func NewHealthUsecase(serviceName string) HealthUsecase {
	return &healthUsecase{
		serviceName: serviceName,
	}
}

func (u *healthUsecase) Ping() domain.HealthStatus {
	return domain.HealthStatus{
		Service: u.serviceName,
		Status:  "ok",
	}
}
