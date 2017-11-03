package main

import (
	"./xsocks"
	"net/http"
	"log"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) addHandler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) addHandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func main() {
	regexpHandler:= RegexpHandler{}
	regexpHandler.addHandleFunc(regexp.MustCompile("^/servers/?$"), xsocks.HttpServe)
	regexpHandler.addHandleFunc(regexp.MustCompile("^/admin/?$"), xsocks.HttpAdmin)
	regexpHandler.addHandleFunc(regexp.MustCompile("^/download/?$"), xsocks.HttpDownload)
	regexpHandler.addHandleFunc(regexp.MustCompile("^/heartbeat/?$"), xsocks.HttpHeartBeat)
	regexpHandler.addHandleFunc(regexp.MustCompile("^/messageget/?$"), xsocks.HttpMessageGet)
	regexpHandler.addHandleFunc(regexp.MustCompile("^/messageput/?$"), xsocks.HttpMessagePut)

	http.HandleFunc("/", regexpHandler.ServeHTTP)

	//http.HandleFunc("/servers", xsocks.HttpServe)
	//http.HandleFunc("/admin", xsocks.HttpAdmin)
	//http.HandleFunc("/download/", xsocks.HttpDownload)
	//http.HandleFunc("/heartbeat", xsocks.HttpHeartBeat)
	//http.HandleFunc("/router", xsocks.HttpRouter)
	//http.HandleFunc("/shell", xsocks.HttpShell)
	//http.HandleFunc("/messageget", xsocks.HttpMessageGet)
	//http.HandleFunc("/messageput", xsocks.HttpMessagePut)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
