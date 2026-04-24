package profiles

const DefaultLocale = "en"

type Profile struct {
	Username    string
	DisplayName string
	Locale      string
}

func (p Profile) String() string {
	return p.Username
}
