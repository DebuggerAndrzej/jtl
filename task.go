package main

type Issue struct {
	title             string
	status            string
	short_description string
}

func (t Issue) FilterValue() string {
	return t.title
}

func (t Issue) Title() string {
	return t.title
}

func (t Issue) Description() string {
	return t.short_description
}
