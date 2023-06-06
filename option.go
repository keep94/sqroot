package sqroot

// Option represents an option for the Print, Fprint, and Sprint methods of
// Number and Digits.
type Option interface {
	mutate(p *printerSettings)
}

// DigitsPerRow sets the number of digits per row. Zero or negative means no
// separate rows.
func DigitsPerRow(count int) Option {
	return optionFunc(func(p *printerSettings) {
		p.digitsPerRow = count
	})
}

// DigitsPerColumn sets the number of digits per column. Zero or negative
// means no separate columns.
func DigitsPerColumn(count int) Option {
	return optionFunc(func(p *printerSettings) {
		p.digitsPerColumn = count
	})
}

// ShowCount shows the digit count in the left margin if on is true.
func ShowCount(on bool) Option {
	return optionFunc(func(p *printerSettings) {
		p.showCount = on
	})
}

// MissingDigit sets the character to represent a missing digit.
func MissingDigit(missingDigit rune) Option {
	return optionFunc(func(p *printerSettings) {
		p.missingDigit = missingDigit
	})
}

type optionFunc func(p *printerSettings)

func (o optionFunc) mutate(p *printerSettings) {
	o(p)
}

func mutateSettings(
	options []Option, settings *printerSettings) *printerSettings {
	for _, option := range options {
		option.mutate(settings)
	}
	return settings
}
