package apiserver

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"

    "github.com/ythosa/linkschecker/src/internal/app/apiserver/links"
)

type server struct {
    router       *mux.Router
    logger       *logrus.Logger
}

func newServer() *server {
    s := &server{
        router:       mux.NewRouter(),
        logger:       logrus.New(),
    }

    s.configureRouter()

    return s
}

func (s *server) configureRouter() {
    // s.router.Use(s.setRequestID)
    // s.router.Use(s.logRequest)
    // s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.router.ServeHTTP(w, r)
}

func (s *server) handleFindBrokenLinks() http.HandlerFunc {
    if len(os.Args) == 1 {
        fmt.Println("Please, pass something in arguments :(")
        os.Exit(1)
    }

    baseURL := links.ParsingURL(os.Args[1])

    for _, errLink := range links.FindBrokenLinks(baseURL) {
        fmt.Println(errLink.Err)
    }

    return func(w http.ResponseWriter, r *http.Request) {
        s.respond(w, r, http.StatusOK, "UAU")
    }
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
    w.WriteHeader(code)
    if data != nil {
        json.NewEncoder(w).Encode(data)
    }
}
