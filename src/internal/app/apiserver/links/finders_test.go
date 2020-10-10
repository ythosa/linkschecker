package links_test

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/ythosa/linkschecker/src/internal/app/apiserver/links"
)

func TestCheckURL(t *testing.T) {
    testCases := []struct {
        name    string
        url     string
        isValid bool
    }{
        {
            name:    "unreachable",
            url:     "https://fhasjkdhasjdhakdhadkjhahsjd/da/s/d/asd/a/d/",
            isValid: false,
        },
        {
            name:    "bad status code",
            url:     "https://github.com/ythotohotohtohtohotoho",
            isValid: false,
        },
        {
            name:    "valid",
            url:     "https://ythosa.github.io",
            isValid: true,
        },
    }

    for _, tc := range testCases {
        _, _, err := links.CheckURL(links.ParsingURL(tc.url))
        assert.Equal(t, tc.isValid, err == nil)
    }
}
