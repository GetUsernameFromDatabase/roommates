package locales

type Language string

const (
	ET      Language = "et"
	Default Language = ET
)

func (e Language) Valid() bool {
	switch e {
	case ET:
		return true
	}
	return false
}
