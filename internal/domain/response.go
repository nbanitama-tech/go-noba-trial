package domain

type GeneralResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ListUsersResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []User `json:"data"`
	Total   int    `json:"total"`
}

type HealthStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type User struct {
	UUID        string `json:"uuid"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

type CreateUserInput struct {
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	Description string `json:"description"`
}
