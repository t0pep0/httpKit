//Litle http server package
package httpKit

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

//HTTP Server interface
type Server interface {
	//Set server host, port and function for logging
	Configure(host, port string, logFunc func(http.Handler) http.Handler)
	//Launch server
	Run() error
	//Get router (*gorilla/mux.Router)
	GetRouter() *mux.Router
	//Add route
	AddRoute(pattern, method, name string, handler func(http.ResponseWriter, *http.Request))
	//404 error route
	Set404Route(handler func(http.ResponseWriter, *http.Request))
	//Get server status
	Status() (configured, run bool)
}

type sError string

func (err *sError) Error() string {
	return string(*err)
}

type err404Handler struct {
	fn404 func(http.ResponseWriter, *http.Request)
}

func (e *err404Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.fn404(w, r)
}

type server struct {
	server     *http.Server
	host, port string
	router     *mux.Router
	logFunc    func(http.Handler) http.Handler
	configured bool
	run        bool
	err404     *err404Handler
}

func (s *server) log(handler http.Handler) (ret http.Handler) {
	if s.logFunc != nil {
		return s.logFunc(handler)
	}
	return handler
}

//Create new server
func New() (serv Server) {
	servp := new(server)
	servp.router = mux.NewRouter()
	servp.configured = false
	servp.run = false
	servp.err404 = new(err404Handler)
	return servp
}

func DefaultLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		log.Println(r.Proto, r.Method, r.RemoteAddr, r.URL, r.UserAgent())
	})
}

func (s *server) Configure(host, port string, logFunc func(http.Handler) http.Handler) {
	s.host = host
	s.port = port
	s.logFunc = logFunc
	s.configured = true
}

func (s *server) Run() error {
	err := new(sError)
	if s.run {
		*err = sError("Server already running")
		return err
	}
	if s.configured {
		http.Handle("/", s.router)
		s.server = &http.Server{
			Addr:           s.host + ":" + s.port,
			Handler:        s.log(http.DefaultServeMux),
			ReadTimeout:    100 * time.Second,
			WriteTimeout:   100 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go s.server.ListenAndServe()
		s.run = true
	} else {
		*err = sError("Server not configured!")
		return err
	}
	return err
}

func (s *server) GetRouter() (router *mux.Router) {
	return s.router
}

func (s *server) AddRoute(pattern, method, name string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler).Methods(method).Name(name)
}

func (s *server) Set404Route(handler func(http.ResponseWriter, *http.Request)) {
	s.err404.fn404 = handler
	s.router.NotFoundHandler = s.err404
}

func (s *server) Status() (configured, run bool) {
	return s.configured, s.run
}
