package httpKit

import (
	"net/http"
	"testing"
)

var (
	server_test Server
)

func TestNew(t *testing.T) {
	server_test = New()
	if server_test == nil {
		t.Fatal("Server is nil!")
	} else {
		t.Log("Server created ..OK!")
	}
}

func TestRunWithoutConfig(t *testing.T) {
	err := server_test.Run()
	if err.Error() != "Server not configured!" {
		t.Fatal("Server.Run() not return error ..FAIL!")
	} else {
		t.Log("Server.Run return error:", err.Error(), "..OK!")
	}
}

func TestRunWithConfig(t *testing.T) {
	test_log := make(chan bool, 10)
	log_func := func(handler http.Handler) http.Handler {
		test_log <- true
		return DefaultLog
	}
	server_test.Configure("localhost", "9999", log_func)
	err := server_test.Run()
	if err.Error() != "" {
		t.Fatal("Server.Run() return error:", err.Error(), " ..FAIL!")
	} else {
		t.Log("Server.Run() ..OK!")
	}
	test, ok := <-test_log
	if (!test) || (!ok) {
		t.Fatal("Server.SetLogFunc() not work ..FAIL!")
	} else {
		t.Log("Server.SetLogFunc() work ..OK!")
	}
}

func TestDoubleRun(t *testing.T) {
	err := server_test.Run()
	if err.Error() != "Server already running" {
		t.Fatal("Server.Run() not return err ..FAIL!")
	} else {
		t.Log("Server.Run() return error:", err.Error(), "..OK!")
	}
}

func TestRoute(t *testing.T) {
	test_chan := make(chan bool, 10)
	test_handler := func(w http.ResponseWriter, r *http.Request) {
		test_chan <- true
	}
	server_test.AddRoute("/test_route", "GET", "TestRoute", test_handler)
	http.Get("http://localhost:9999/test_route")
	test, ok := <-test_chan
	if (!test) || (!ok) {
		t.Fatal("Server.AddRoute() not work ..FAIL!")
	} else {
		t.Log("Server.AddRoute() work ..OK!")
	}
}

func TestGetRouter(t *testing.T) {
	router := server_test.GetRouter()
	if router == nil {
		t.Fatal("Server.GetRouter() not return router ..FAIL!")
	} else {
		t.Log("Server.GetRouter() return router ..OK!")
	}
}

func Test404Route(t *testing.T) {
	test_chan := make(chan bool, 2)
	test_handler := func(w http.ResponseWriter, r *http.Request) {
		test_chan <- true
	}
	server_test.Set404Route(test_handler)
	http.Get("http://localhost:9999/404")
	test, ok := <-test_chan
	if (!test) || (!ok) {
		t.Fatal("Server.Set404Route() not work ..FAIL!")
	} else {
		t.Log("Server.Set404Route work! ..OK!")
	}
}

func TestStatus(t *testing.T) {
	conf, run := server_test.Status()
	if !conf {
		t.Fatal("Incorrect configure status! ..FAIL!")
	}
	if !run {
		t.Fatal("Incorrect run status! ..FAIL!")
	}
	t.Log("Server.Status() work! ..OK!")
}
