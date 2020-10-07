package links

type ParsingURL string

type BadURL struct {
    ParsingURL
    Err error
}
