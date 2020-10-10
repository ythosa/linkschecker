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
