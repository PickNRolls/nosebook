package presenteruser

import presenterdto "nosebook/src/application/presenters/dto"

type User = presenterdto.User

type FindOutUser = presenterdto.FindOut[*User]

type FindByTextInput struct {
	Text string
	Next string
}

func (this *FindByTextInput) BuildFromMap(m map[string]any) FindByTextInput {
	out := FindByTextInput{}

	if next, ok := m["next"].(string); ok {
		out.Next = next
	}

	if text, ok := m["text"].(string); ok {
    out.Text = text
	}

	return out
}
