package links_test

import (
	"strings"
	"testing"

	"github.com/ythosa/linkschecker/src/internal/app/apiserver/links"
)

func TestExtract(t *testing.T) {
	testCases := []struct {
		url          string
		expectedURLs []string
	}{
		{
			url: "http://ythosa.github.io",
			expectedURLs: []string{"https://www.youtube.com/channel/UCtzplSh3NNIX1zr63P2c-tQ", "https://vk.com/ythosa",
				"https://www.instagram.com/y1hosa/", "https://github.com/ythosa", "https://t.me/ythosa"},
		},
	}

	for _, testCase := range testCases {
		res, doc, err := links.CheckURL(links.ParsingURL(testCase.url))
		if err != nil {
			t.Error(err)
		}
		extracted := links.Extract(res, doc)
		for _, link := range extracted {
			found := false
			for _, e := range testCase.expectedURLs {
				if strings.EqualFold(string(link), e) {
					found = true
				}
			}
			if !found {
				t.Errorf("%s not extracted", link)
			}
		}
	}
}
