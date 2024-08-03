package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// LogReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	// GET etc.
	method  string
	uri     string
	referer string
	ipaddr  string
	// response code, like 200, 404
	code int
	// number of bytes of the response sent
	size int64
	// how long did it take to
	duration  time.Duration
	userAgent string
}

// func logRequestHandler(h http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		ri := &HTTPReqInfo{
// 			method: r.Method,
// 			uri: r.URL.String(),
// 			referer: r.Header.Get("Referer"),
// 			userAgent: r.Header.Get("User-Agent"),
// 		}

// 		ri.ipaddr = requestGetRemoteAddress(r)

// 		// this runs handler h and captures information about
// 		// HTTP request
// 		m := httpsnoop.CaptureMetrics(h, w, r)

// 		ri.code = m.Code
// 		ri.size = m.BytesWritten
// 		ri.duration = m.Duration
// 		logHTTPReq(ri)
// 	}
// 	return http.HandlerFunc(fn)
// }

// Create a request logging middleware handler called Logger
type Logger struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

// NewLogger constructs a new Logger middleware handler
func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

func main() {

	mux := http.NewServeMux()

	fileHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	mux.Handle("/static/", fileHandler)

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	//wrap entire mux with logger middleware
	wrappedMux := NewLogger(mux)

	fmt.Println("Starting front end service on port 80")
	log.Fatal(http.ListenAndServe(":80", wrappedMux))
}

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string

	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))
	templateSlice = append(templateSlice, partials...)

	type _viewState struct {
		Name string
		Year int
	}

	v := _viewState{Name: "Jaan", Year: 2024}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, &v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
