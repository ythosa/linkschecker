package links

// ParsingURL is type-wrapper over string for links.
type ParsingURL string

// BrokenURL is type for broken links.
type BrokenURL struct {
    ParsingURL       // Link URL
    Error      error // Error of link
}
