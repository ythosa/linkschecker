package links

type ParsingURL string

type BrokenURL struct {
    ParsingURL
    Err error
}
