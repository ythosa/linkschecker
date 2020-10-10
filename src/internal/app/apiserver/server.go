package apiserver

import (
    "context"
    "encoding/json"
    "net/http"
    "sync"
    "time"

    "github.com/google/uuid"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"

    "github.com/ythosa/linkschecker/src/internal/app/apiserver/links"
)

const ctxKeyRequestID ctxKey = iota

type ctxKey int8

type server struct {
    router *mux.Router
    logger *logrus.Logger
}

func newServer() *server {
    s := &server{
        router: mux.NewRouter(),
        logger: logrus.New(),
    }

    s.configureRouter()

    return s
}

func (s *server) configureRouter() {
    s.router.Use(s.setRequestID)
    s.router.Use(s.logRequest)
    s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
    s.router.HandleFunc("/get_broken_links", s.HandleFindBrokenLinks()).Methods("POST")
    s.router.HandleFunc("/validate_link", s.HandleLinkValidation()).Methods("POST")
    s.router.HandleFunc("/validate_links", s.HandleLinksValidations()).Methods("POST")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := uuid.New().String()
        w.Header().Set("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
    })
}

func (s *server) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger := s.logger.WithField("request", logrus.Fields{
            "remote_addr": r.RemoteAddr,
            "request_id":  r.Context().Value(ctxKeyRequestID),
        })
        logger.Infof("started %s %s", r.Method, r.RequestURI)

        start := time.Now()

        rw := &responseWriter{w, http.StatusOK}
        next.ServeHTTP(rw, r)

        logger.Infof(
            "completed with %d %s in %v\n",
            rw.code,
            http.StatusText(rw.code),
            time.Since(start),
        )
    })
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.router.ServeHTTP(w, r)
}

func (s *server) HandleFindBrokenLinks() http.HandlerFunc {
    type request struct {
        BaseURL string `json:"base_url"`
    }

    return func(w http.ResponseWriter, r *http.Request) {
        req := &request{}
        decoder := json.NewDecoder(r.Body)
        decoder.DisallowUnknownFields()
        if err := decoder.Decode(req); err != nil {
            s.error(w, http.StatusBadRequest, err)
            return
        }

        foundBrokenLinks := links.FindBrokenLinks(links.ParsingURL(req.BaseURL))
        badLinks := make(map[string]string)

        for _, l := range foundBrokenLinks {
            badLinks[string(l.ParsingURL)] = l.Error.Error()
        }

        s.respond(w, http.StatusOK, map[string]interface{}{"broken_links": badLinks})
    }
}

func (s *server) HandleLinkValidation() func(http.ResponseWriter, *http.Request) {
    type request struct {
        Link string `json:"link"`
    }

    type response struct {
        OK    string `json:"ok"`
        Error string `json:"error"`
    }

    return func(w http.ResponseWriter, r *http.Request) {
        req := &request{}
        decoder := json.NewDecoder(r.Body)
        decoder.DisallowUnknownFields()
        if err := decoder.Decode(req); err != nil {
            s.error(w, http.StatusBadRequest, err)
            return
        }

        res := response{}
        _, _, isLinkValid := links.CheckURL(links.ParsingURL(req.Link)) // nolint:bodyclose
        if isLinkValid != nil {
            res.OK = "false"
            res.Error = isLinkValid.Error()
        } else {
            res.OK = "true"
        }
        s.respond(w, http.StatusOK, res)
    }
}

func (s *server) HandleLinksValidations() func(http.ResponseWriter, *http.Request) {
    type request struct {
        Links []links.ParsingURL `json:"links"`
    }

    type responseElement struct {
        URL   string `json:"url"`
        Error string `json:"error"`
    }

    return func(w http.ResponseWriter, r *http.Request) {
        req := &request{}
        decoder := json.NewDecoder(r.Body)
        decoder.DisallowUnknownFields()
        if err := decoder.Decode(req); err != nil {
            s.error(w, http.StatusBadRequest, err)
            return
        }

        var wg sync.WaitGroup
        var mu sync.Mutex
        var response []responseElement

        for _, l := range req.Links {
            wg.Add(1)
            go func(l links.ParsingURL) {
                defer wg.Done()
                _, _, err := links.CheckURL(l) // nolint:bodyclose
                mu.Lock()
                if err != nil {
                    response = append(response, responseElement{
                        URL:   string(l),
                        Error: err.Error(),
                    })
                } else {
                    response = append(response, responseElement{
                        URL:   string(l),
                        Error: "null",
                    })
                }
                mu.Unlock()
            }(l)
        }

        wg.Wait()
        s.respond(w, http.StatusOK, response)
    }
}

func (s *server) error(w http.ResponseWriter, code int, err error) {
    s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, code int, data interface{}) {
    w.WriteHeader(code)
    if data != nil {
        _ = json.NewEncoder(w).Encode(data)
    }
}
