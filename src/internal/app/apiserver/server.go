package apiserver

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
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
    s.router.HandleFunc("/get_broken", s.handleFindBrokenLinks()).Methods("POST")
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
            time.Now().Sub(start),
        )
    })
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.router.ServeHTTP(w, r)
}

func (s *server) handleFindBrokenLinks() http.HandlerFunc {
    type request struct {
        BaseURL string `json:"base_url"`
    }

    if len(os.Args) == 1 {
        fmt.Println("Please, pass something in arguments :(")
        os.Exit(1)
    }

    baseURL := links.ParsingURL(os.Args[1])

    for _, errLink := range links.FindBrokenLinks(baseURL) {
        fmt.Println(errLink.Err)
    }

    return func(w http.ResponseWriter, r *http.Request) {
        req := &request{}
        if err := json.NewDecoder(r.Body).Decode(req); err != nil {
            s.error(w, r, http.StatusBadRequest, err)
            return
        }

        foundBrokenLinks := links.FindBrokenLinks(links.ParsingURL(req.BaseURL))

        s.respond(w, r, http.StatusOK, map[string][]links.BrokenURL{"links": foundBrokenLinks})
    }
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
    s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
    w.WriteHeader(code)
    if data != nil {
        json.NewEncoder(w).Encode(data)
    }
}
