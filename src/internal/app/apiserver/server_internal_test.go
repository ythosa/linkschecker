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

    testCases := []struct {
        name             string
        payload          interface{}
        expectedCode     int
        expectedResponse interface{}
    }{
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
                "ok": "false",
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
                "ok": "false",
                "error": "http://asdasfasasd//asd - is unreachable",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            rec := httptest.NewRecorder()
            b := &bytes.Buffer{}
            json.NewEncoder(b).Encode(tc.payload)
            req, _ := http.NewRequest(http.MethodPost, "/validate_link", b)

            s.ServeHTTP(rec, req)

            var resp map[string]string
            json.NewDecoder(rec.Body).Decode(&resp)

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
    }
}
