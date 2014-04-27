package http

import (
	"regexp"
	"net/http"
	"strings"
	"github.com/ian-kent/MailHog/mailhog/config"
	"github.com/ian-kent/MailHog/mailhog/http/api"
	"github.com/ian-kent/MailHog/mailhog/http/handler"
)

var exitChannel chan int
var cfg *config.Config

func web_exit(w http.ResponseWriter, r *http.Request, route *handler.Route) {
	web_headers(w)
	w.Write([]byte("Exiting MailHog!"))
	exitChannel <- 1
}

func web_index(w http.ResponseWriter, r *http.Request, route *handler.Route) {
	web_headers(w)
	data, _ := cfg.Assets("assets/templates/index.html")
	w.Write([]byte(web_render(string(data))))
}

func web_jscontroller(w http.ResponseWriter, r *http.Request, route *handler.Route) {
	w.Header().Set("Content-Type", "text/javascript")
	data, _ := cfg.Assets("assets/js/controllers.js")
	w.Write(data)
}

func web_imgcontroller(w http.ResponseWriter, r *http.Request, route *handler.Route) {
	w.Header().Set("Content-Type", "image/png")
	data, _ := cfg.Assets("assets/images/hog.png")
	w.Write(data)
}

func web_render(content string) string {
	data, _ := cfg.Assets("assets/templates/layout.html")
	layout := string(data)
	html := strings.Replace(layout, "<%= content %>", content, -1)
	return html
}

func web_headers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func Start(exitCh chan int, conf *config.Config) {
	exitChannel = exitCh
	cfg = conf

	server := &http.Server{
		Addr: conf.HTTPBindAddr,
		Handler: &handler.RegexpHandler{},
	}

	server.Handler.(*handler.RegexpHandler).HandleFunc(regexp.MustCompile("^/exit/?$"), web_exit)
	server.Handler.(*handler.RegexpHandler).HandleFunc(regexp.MustCompile("^/js/controllers.js$"), web_jscontroller)
	server.Handler.(*handler.RegexpHandler).HandleFunc(regexp.MustCompile("^/images/hog.png$"), web_imgcontroller)
	server.Handler.(*handler.RegexpHandler).HandleFunc(regexp.MustCompile("^/$"), web_index)

	api.CreateAPIv1(exitCh, conf, server)

	server.ListenAndServe()
}