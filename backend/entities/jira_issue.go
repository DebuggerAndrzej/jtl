package entities

type Issue struct {
	Key              string
	Status           string
	ShortDescription string
	OriginalEstimate string
	LoggedTime       string
}

func (t Issue) FilterValue() string {
	return t.ShortDescription
}
