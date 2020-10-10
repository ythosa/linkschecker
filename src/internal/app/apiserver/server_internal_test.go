package apiserver

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestServer_HandleLinkValidation(t *testing.T) {
    s := newServer()

    type testCase struct {
        name             string
        payload          interface{}
        expectedCode     int
        expectedResponse interface{}
    }

    testCases := []testCase{
        {
            name: "invalid request",
            payload: map[string]string{
                "error_json": "",
            },
            expectedCode: 400,
            expectedResponse: map[string]string{
                "error": "json: unknown field \"error_json\"",
            },
        },
        {
            name: "invalid link: response 404",
            payload: map[string]string{
                "link": "https://github.com/yththththththththaaa",
            },
            expectedCode: 200,
            expectedResponse: map[string]string{
                "ok":    "false",
                "error": "https://github.com/yththththththththaaa - bad status code response - 404",
            },
        },
        {
            name: "invalid link: unreachable",
            payload: map[string]string{
                "link": "http://asdasfasasd//asd",
            },
            expectedCode: 200,
            expectedResponse: map[string]string{
                "ok":    "false",
                "error": "http://asdasfasasd//asd - is unreachable",
            },
        },
    }

    for _, tc := range testCases {
        func(tc testCase) {
            t.Run(tc.name, func(t *testing.T) {
                rec := httptest.NewRecorder()
                b := &bytes.Buffer{}
                if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
                    t.Error(err)
                }
                req, _ := http.NewRequest(http.MethodPost, "/validate_link", b)

                s.ServeHTTP(rec, req)

                var resp map[string]string

                if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
                    t.Error(err)
                }

                assert.Equal(
                    t,
                    map[string]interface{}{
                        "code":     tc.expectedCode,
                        "response": tc.expectedResponse,
                    },
                    map[string]interface{}{
                        "code":     rec.Code,
                        "response": resp,
                    },
                )
            })
        }(tc)
    }
}

func TestServer_HandleFindBrokenLinks(t *testing.T) {
    s := newServer()

    type testCase struct {
        name             string
        payload          interface{}
        expectedCode     int
        expectedResponse map[string]interface{}
    }

    testCases := []testCase{
        {
            name: "invalid request",
            payload: map[string]string{
                "error_json": "",
            },
            expectedCode: 400,
            expectedResponse: map[string]interface{}{
                "error": "json: unknown field \"error_json\"",
            },
        },
        {
            name: "all links valid",
            payload: map[string]string{
                "base_url": "https://ythosa.github.io",
            },
            expectedCode: 200,
            expectedResponse: map[string]interface{}{
                "broken_links": map[string]interface{}{},
            },
        },
    }

    for _, tc := range testCases {
        func(tc testCase) {
            t.Run(tc.name, func(t *testing.T) {
                rec := httptest.NewRecorder()
                b := &bytes.Buffer{}
                if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
                    t.Error(err)
                }
                req, _ := http.NewRequest(http.MethodPost, "/get_broken_links", b)

                s.ServeHTTP(rec, req)

                var resp map[string]interface{}
                if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
                    t.Error(err)
                }

                assert.Equal(
                    t,
                    map[string]interface{}{
                        "code":     tc.expectedCode,
                        "response": tc.expectedResponse,
                    },
                    map[string]interface{}{
                        "code":     rec.Code,
                        "response": resp,
                    },
                )
            })
        }(tc)
    }
}

func findURLAndErrorValuesInResponse(url string, error string, values []map[string]string) bool {
    found := false
    for _, v := range values {
        if url == v["url"] && error == v["error"] {
            found = true
            break
        }
    }
    return found
}

func TestServer_HandleLinksValidations(t *testing.T) {
    s := newServer()

    type urlResult struct {
        URL   string `json:"url"`
        Error string `json:"error"`
    }

    type testCase struct {
        name             string
        payload          map[string][]string
        expectedResponse []urlResult
    }

    testCases := []testCase{
        {
            name: "invalid request",
            payload: map[string][]string{
                "links": {
                    "http://youtube.com",
                    "https://vk.com",
                    "https://github.com/nonononovalid",
                    "https://unreachable.unreachable/",
                },
            },
            expectedResponse: []urlResult{
                {
                    URL:   "https://unreachable.unreachable/",
                    Error: "https://unreachable.unreachable/ - is unreachable",
                },
                {
                    URL:   "https://vk.com",
                    Error: "null",
                },
                {
                    URL:   "https://github.com/nonononovalid",
                    Error: "https://github.com/nonononovalid - bad status code response - 404",
                },
                {
                    URL:   "http://youtube.com",
                    Error: "null",
                },
            },
        },
    }

    for _, tc := range testCases {
        func(tc testCase) {
            t.Run(tc.name, func(t *testing.T) {
                rec := httptest.NewRecorder()
                b := &bytes.Buffer{}
                if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
                    t.Error(err)
                }
                req, _ := http.NewRequest(http.MethodPost, "/validate_links", b)
                s.ServeHTTP(rec, req)
                var resp []map[string]string
                if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
                    t.Error(err)
                }

                valid := true
                for _, r := range tc.expectedResponse {
                    valid = findURLAndErrorValuesInResponse(r.URL, r.Error, resp)
                }

                if !valid {
                    assert.Equal(t, tc.expectedResponse, resp)
                }
            })
        }(tc)
    }
}
