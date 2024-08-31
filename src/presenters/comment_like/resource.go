package presentercommentlike

type commentResource struct {
}

func (this *commentResource) IDColumn() string {
	return "comment_id"
}

func (this *commentResource) Table() string {
	return "comment_likes"
}
