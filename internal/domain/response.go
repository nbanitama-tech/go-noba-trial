package domain

type GeneralResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type HealthStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}
