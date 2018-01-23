package httpmodule

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"../model"
	"../utils"
)

var fileStoredPath = "/home/ec2-user/"
var rules = map[string]string{}

var client4conf = redis.NewClient(&redis.Options{
	Addr:     model.RedisAddr,
	Password: "njust2006",
	DB:       4,
})

func getUrlParam(r *http.Request, paramKey string) string {
	var paramValue string

	for k, v := range r.Form {
		if k == paramKey {
			for _, value := range v {
				paramValue = value
			}
		}
	}
	return paramValue
}

func getUserID(computer_id string) string {
	clientInfo := strings.Split(computer_id, "-")
	userId := clientInfo[1]
	return userId
}

func getClientVersion(computer_id string) string {
	clientInfo := strings.Split(computer_id, "-")
	cliVer := clientInfo[0]
	return cliVer
}

func findXSocks() string {
	var xsocks_fullname string
	filepath.Walk(fileStoredPath, func(path string, info os.FileInfo, err error) error {
		var e error
		m, _ := filepath.Match("xsocks*.exe.gz", info.Name())

		if m {
			xsocks_fullname = info.Name()
		}

		return e
	})
	return xsocks_fullname
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func SSConfigApi(w http.ResponseWriter, r *http.Request) {
	configs := model.GetConfigs()

	rjson, _ := json.Marshal(configs)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(rjson))
}

func SServerApi(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fastServers, _ := client4conf.LRange("good_servers", 0, -1).Result()
	addr_info := strings.Split(r.RemoteAddr, ":")

	//token := r.Header.Get("access-token")
	//if token != "b25seS1mb3ItZmV3LXBlcnNvbnMtdGhhdC1yZWFsbHktbmVlZA" {
	//	w.Header().Set("Content-Type", "application/json")
	//	fmt.Fprint(w, `[	{"server": "100.100.100.100","server_port": "8888","password": "password","method": "aes-256-cfb","remarks": "请及时删除Xsocks!"}]`)
	//	return
	//}

	computer_id := getUrlParam(r, "computerid")
	userId := getUserID(computer_id)

	user := model.GetUser(userId)

	keyParam := getUrlParam(r, "key")

	if keyParam == "client-close" {
		user.Online = 0
		user.Save()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"client-close\":\"OK\"}")
		return
	} else if keyParam != "" {
		ret, _ := utils.Contains(keyParam, fastServers)
		if !ret {
			serv := model.GetServer(keyParam)
			serv.Del()
		}
	}

	var validAndFastServers = model.GetValidAndFastServers(fastServers)

	var rnd_key string
	var sserver *model.SServer

	var vipUsers = []string{"wll17331", "llwang"}

	isVip, _ := utils.Contains(userId, vipUsers)

	fmt.Println(userId, isVip)

	if sserverName, ok := rules[addr_info[0]]; ok {
		rnd_key = sserverName
	} else if len(validAndFastServers) > 0 {
		triedNum := 0
		for {
			if !isVip {
				randIndex := rand.Intn(len(validAndFastServers))
				rnd_key = validAndFastServers[randIndex]
			} else {
				fastest := model.GetFastestServer(user.SServer)
				rnd_key = fastest.Name
				fmt.Println(user.SServer, rnd_key)
			}

			if rnd_key != user.SServer || triedNum == len(validAndFastServers) {
				break
			}
			triedNum++
		}
	}

	if rnd_key == "" {
		if isVip {
			//VIP user, give him the fast server
			sserver = model.GetFastestServer("")
		} else {
			sserver = model.GetRandomServer()
		}
	} else {
		sserver = model.GetServer(rnd_key)
	}

	rjson, _ := json.Marshal(sserver.Config)

	user.SServer = sserver.Name
	user.Save()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "[%s]", rjson)
}

func GetMessageApi(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//token := r.Header.Get("access-token")
	//if token != "b25seS1mb3ItZmV3LXBlcnNvbnMtdGhhdC1yZWFsbHktbmVlZA" {
	//	fmt.Fprint(w, "[]")
	//	return
	//}

	computer_id := getUrlParam(r, "computerid")
	userId := getUserID(computer_id)

	user := model.GetUser(userId)

	messages := user.ConsumeMessages()

	if messages != nil {
		fmt.Fprintf(w, "MESSAGE:%s", strings.Join(messages, ", "))
	} else {
		fmt.Fprint(w, "[]")
	}
}

func PushMessageApi(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	message := getUrlParam(r, "message")

	for _, user := range model.GetUsers() {
		user.PushMessages(message)
	}

	fmt.Fprintf(w, "[%s]", "success")
}

func SpeedApi(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	info := getUrlParam(r, "info")
	speedInfo := strings.Split(info, "$")
	servName := speedInfo[0]
	speed, _ := strconv.ParseFloat(speedInfo[1], 64)

	server := model.GetServer(servName)
	server.Speed = speed
	server.Save()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{\"set speed\":\"ok\"}")
}

func AdminApi(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var actions = []string{}
	var resp = []string{}

	for k, v := range r.Form {
		if k == "action" {
			for _, action := range v {
				actions = append(actions, action)
			}
		}
	}

	for _, action := range actions {
		if action == "flushdb" {
			ret := model.FlushServersDB()
			resp = append(resp, fmt.Sprintf("{\"flushdb\": {\"serversDB\": \"%s\"}}", ret))

		} else if action == "getall" {
			servers := model.GetServers()
			if servers == nil {
				resp = append(resp, "{\"getall\": \"error\"}")
			} else {
				new_keys := utils.Map(servers, func(s string) string {
					sserver := model.GetServer(s)
					if sserver.Speed != 60 {
						return fmt.Sprintf("\"%s--%.3f\"", s, sserver.Speed)
					} else {
						return fmt.Sprintf("\"%s\"", s)
					}

				})
				resp = append(resp, fmt.Sprintf("{\"getall\": {\"count\": %d, \"keys\": [%s]}}",
					len(servers),
					strings.Join(new_keys, ",")))
			}

		} else if strings.HasPrefix(action, "get->") {
			parts := strings.Split(action, "->")
			sserver := model.GetServer(parts[1])
			if sserver == nil {
				resp = append(resp, fmt.Sprintf("{\"%s\": \"error\"}", action))
			} else {
				jsonBytes, _ := json.Marshal(sserver.Config)
				resp = append(resp, fmt.Sprintf("{\"%s\": %s}",
					action,
					string(jsonBytes)))
			}

		} else if strings.HasPrefix(action, "del->") {
			parts := strings.Split(action, "->")
			sserver := model.GetServer(parts[1])

			if sserver == nil {
				resp = append(resp, fmt.Sprintf("{\"%s\": \"error\"}", action))
			} else {
				ret := sserver.Del()
				resp = append(resp, fmt.Sprintf("{\"%s\": %d}",
					action,
					ret))
			}

		} else if strings.HasPrefix(action, "ar->") {
			parts := strings.Split(action, "->")
			addr_info := strings.Split(r.RemoteAddr, ":")

			servers := model.GetServers()
			exists, _ := utils.Contains(parts[1], servers)
			if exists {
				rules[addr_info[0]] = parts[1]
				resp = append(resp, fmt.Sprintf("{\"%s\": \"success\"}", action))
			} else {
				resp = append(resp, fmt.Sprintf("{\"%s\": \"failed: no such server: %s\"}", action, parts[1]))
			}

		} else if strings.HasPrefix(action, "dr") {
			addr_info := strings.Split(r.RemoteAddr, ":")
			delete(rules, addr_info[0])
			resp = append(resp, fmt.Sprintf("{\"%s\": \"success\"}", action))

		} else if strings.HasPrefix(action, "usage") {
			var usage = map[string][]string{}
			aliveUserNum := 0

			users := model.GetUsers()

			for _, user := range users {
				if user.Online == 1 {
					usage[user.SServer] = append(usage[user.SServer], user.ClientIP+"#"+user.ClientVersion+"-"+user.UserID)
					aliveUserNum++
				}
			}

			usage_json, _ := json.Marshal(usage)
			resp = append(resp, fmt.Sprintf("{\"%s\": %s}", action+"("+strconv.Itoa(aliveUserNum)+")", usage_json))

		} else {
			resp = append(resp, "{\"msg\": \"unsupported action\"}")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "[%s]", strings.Join(resp, ","))
}

func DownloadApi(w http.ResponseWriter, r *http.Request) {
	var action string
	r.ParseForm()
	for k, v := range r.Form {
		if k == "action" {
			action = v[0]
		}
	}
	xsocksFullName := findXSocks()

	reg := regexp.MustCompile(`xsocks-(\d+\.\d+\.\d+).exe.gz`)
	results := reg.FindAllStringSubmatch(string(xsocksFullName), -1)

	if action == "getversion" {
		fmt.Fprint(w, results[0][1])
	} else if action == "getfile" {
		file := fileStoredPath + xsocksFullName
		if exist := isExist(file); !exist {
			http.NotFound(w, r)
		}
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", xsocksFullName))
		w.Header().Add("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, file)
	} else {
		fmt.Fprint(w, "unsupported action")
	}
}

func HeartBeatApi(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr_info := strings.Split(r.RemoteAddr, ":")

	//token := r.Header.Get("access-token")
	//if token != "b25seS1mb3ItZmV3LXBlcnNvbnMtdGhhdC1yZWFsbHktbmVlZA" {
	//	fmt.Fprint(w, "[]")
	//	return
	//}

	computer_id := getUrlParam(r, "computerid")
	sserver := getUrlParam(r, "server")

	userId := getUserID(computer_id)
	cliVer := getClientVersion(computer_id)

	user := model.GetUser(userId)

	if sserver != "" {
		user.SServer = sserver
	}

	user.ClientIP = addr_info[0]
	user.Online = 1
	//timeLocal := time.FixedZone("CST", 3600*8)
	//time.Local = timeLocal
	user.LastAliveTime = time.Now().Format("2006-01-02 15:04:05")
	user.ClientVersion = cliVer

	user.Save()

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprint(w, "{\"receive hb\":\"ok\"}")
}
