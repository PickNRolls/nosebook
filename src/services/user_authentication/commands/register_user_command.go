package commands

type RegisterUserCommand struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Nick      string `json:"nick"`
	Password  string `json:"password"`
}
