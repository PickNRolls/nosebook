package presenterpostlike

type postResource struct {
}

func (this *postResource) IDColumn() string {
	return "post_id"
}

func (this *postResource) Table() string {
	return "post_likes"
}
