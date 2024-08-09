package commands

type RegisterUserCommand struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Nick      string `json:"nick" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
