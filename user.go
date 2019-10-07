package boilerplateapi

// UserInfo represents all public info about a user
type UserInfo struct {
	Username string `json:"username"`
	Location string `json:"location"`
}

type userData struct {
	UserInfo
	Location string
	Password string `json:"password"`
}

func (user *userData) PublicData() UserInfo {
	return UserInfo{
		Username: user.Username,
		Location: user.Location,
	}
}

func (user *userData) PopulateFields() {
	// TODO: Populate other fields
}

// User represents a user in the system
type User interface {
	PublicData() UserInfo
	PopulateFields()
}
