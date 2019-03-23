package commonmap

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cgi"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Serve() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	var listenAddr = "localhost:7070"

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			os.Exit(1)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/wms", render)
	r.Handle("/{path:.*}", http.StripPrefix("/", http.FileServer(http.Dir(appDir))))

	log.Printf("listening on http://%s", listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, handlers.LoggingHandler(os.Stdout, r)))
}

func render(w http.ResponseWriter, r *http.Request) {

	//w.Header().Add("Content-Type", "image/png")
	err := MapRender(w, r)

	if err != nil {
		internalError(w, r, err)
		return
	}
}

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Print(err)
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal error"))
}

func MapRender(dst io.Writer, r *http.Request /*mapReq Request2*/) error {
	//fmt.Println("Map = " + MapfilePath)

	wd := filepath.Dir(MapfilePath)
	handler := cgi.Handler{
		Path: mapservPath,
		Dir:  wd,
	}

	w := &responseRecorder{
		Body: dst,
	}

	query := "/?MAP=" + url.QueryEscape(MapfilePath) + "&" + r.URL.RawQuery
	//	if !strings.Contains(query, "GetCapabilities") {
	//		query = query + "&LAYERS=map"
	//	}

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return err
	}

	handler.ServeHTTP(w, req)

	if w.Code != 200 {
		return fmt.Errorf("error while calling mapserv CGI (status %d)", w.Code)
	}
	//if ct := w.Header().Get("Content-type"); ct != "" && !strings.HasPrefix(ct, "image") {
	//	return fmt.Errorf(" mapserv CGI did not return image (%v)", w.Header())
	//}
	return nil
}

// responseRecorder from net/http/httptest
// copied here to work around global -httptest.server flag from httptest package

// responseRecorder is an implementation of http.ResponseWriter that
// records its mutations for later inspection in tests.
type responseRecorder struct {
	Code      int         // the HTTP response code from WriteHeader
	HeaderMap http.Header // the HTTP response headers
	Body      io.Writer   // if non-nil, the io.Writer to append written data to
	Flushed   bool

	wroteHeader bool
}

func (rw *responseRecorder) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

func (rw *responseRecorder) Write(buf []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

func (rw *responseRecorder) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.Code = code
	}
	rw.wroteHeader = true
}

func (rw *responseRecorder) Flush() {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	rw.Flushed = true
}
