package convert

// ConverterOption configures a Converter
type ConverterOption func(*Converter)

// Converter converts GitHub data to Markdown
type Converter struct {
	userLinks        bool
	includeReactions bool
}

// WithUserLinks enables rendering usernames as GitHub links
func WithUserLinks(enabled bool) ConverterOption {
	return func(c *Converter) {
		c.userLinks = enabled
	}
}

// WithIncludeReactions includes reaction statistics in output
func WithIncludeReactions(enabled bool) ConverterOption {
	return func(c *Converter) {
		c.includeReactions = enabled
	}
}

// New creates a new Converter with options
func New(opts ...ConverterOption) *Converter {
	c := &Converter{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
