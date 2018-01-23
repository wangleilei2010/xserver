package main

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"./httpmodule"
	"./model"
)

func main() {
	go func() {
		for {
			model.ScanOnlineUsers()
			time.Sleep(time.Minute * 5)
		}
	}()

	regexpHandler := httpmodule.RegexpHandler{}

	regexpHandler.AddHandleFunc(regexp.MustCompile("^/configs/?$"), httpmodule.SSConfigApi)

	regexpHandler.AddHandleFunc(regexp.MustCompile("^/servers/?$"), httpmodule.SServerApi)
	regexpHandler.AddHandleFunc(regexp.MustCompile("^/speed/?$"), httpmodule.SpeedApi)

	regexpHandler.AddHandleFunc(regexp.MustCompile("^/admin/?$"), httpmodule.AdminApi)
	regexpHandler.AddHandleFunc(regexp.MustCompile("^/download/?$"), httpmodule.DownloadApi)
	regexpHandler.AddHandleFunc(regexp.MustCompile("^/heartbeat/?$"), httpmodule.HeartBeatApi)

	regexpHandler.AddHandleFunc(regexp.MustCompile("^/messageget/?$"), httpmodule.GetMessageApi)
	regexpHandler.AddHandleFunc(regexp.MustCompile("^/messagepush/?$"), httpmodule.PushMessageApi)

	http.HandleFunc("/", regexpHandler.ServeHTTP)
	err := http.ListenAndServe(":80", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
