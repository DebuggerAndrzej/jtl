package entities

type Issue struct {
	Key              string
	Status           string
	ShortDescription string
	Description      string
	OriginalEstimate string
	LoggedTime       string
}

func (t Issue) FilterValue() string {
	return t.ShortDescription
}
