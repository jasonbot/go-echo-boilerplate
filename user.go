package quoteapi

type userData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Robot    string `json:"robot"`
}

// User represents a user in the system
type User interface {
}
