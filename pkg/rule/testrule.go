package rule

// NewTestRule returns a new Rule meant to be used only in tests
func NewTestRule() *Rule {
	r := Rule{
		Name:         "test-rule",
		Terms:        []string{"testrule", "test-rule"},
		Alternatives: []string{"better-rule"},
		Severity:     SevError,
		Options: Options{
			WordBoundary: false,
		},
	}
	r.SetRegexp()
	return &r
}