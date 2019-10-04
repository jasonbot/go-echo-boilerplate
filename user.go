package quoteapi

type userData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Robot    string `json:"robot"`
	Location string `json:"location"`
}

func (user *userData) PopulateFields() {
	user.Robot = "TODO"
}

// User represents a user in the system
type User interface {
	PopulateFields()
}
