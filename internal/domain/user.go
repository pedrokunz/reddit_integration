package domain

type User struct {
	ID         string            `json:"id"`
	Email      string            `json:"email"`
	Password   string            `json:"password"`
	Roles      []string          `json:"roles"`
	Attributes map[string]string `json:"attributes"`
	Customer   Customer          `json:"customer"`
}
