package presenterdto

type Likes struct {
	Count            int     `json:"count"`
	RandomFiveLikers []*User `json:"randomFiveLikers"`
	Liked            bool    `json:"liked"`
}
