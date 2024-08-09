package commands

type LoginCommand struct {
	Nick     string `json:"nick" binding:"required"`
	Password string `json:"password" binding:"required"`
}
