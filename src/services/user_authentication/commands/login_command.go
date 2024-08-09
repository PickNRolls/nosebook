package commands

type LoginCommand struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}
