package main

import (
	"./xsocks"
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/servers", xsocks.HttpServe)
	http.HandleFunc("/admin", xsocks.HttpAdmin)
	http.HandleFunc("/download/", xsocks.HttpDownload)
	http.HandleFunc("/heartbeat", xsocks.HttpHeartBeat)
	http.HandleFunc("/router", xsocks.HttpRouter)
	http.HandleFunc("/shell", xsocks.HttpShell)
	http.HandleFunc("/messageget", xsocks.HttpMessageGet)
	http.HandleFunc("/messageput", xsocks.HttpMessagePut)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
