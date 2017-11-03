package xsocks

import (
	"fmt"
	"net/http"
	"github.com/go-redis/redis"
	"time"
	"reflect"
	"errors"
	"math/rand"
	"strings"
	"path/filepath"
	"os"
	"regexp"
	"encoding/json"
	"encoding/base64"
	"io/ioutil"
	"strconv"
)

//var file_stored_path = "D:\\"
var file_stored_path = "/home/ec2-user/"
var redis_Addr = "localhost:6379"

var client = redis.NewClient(&redis.Options{
	Addr:     redis_Addr,
	Password: "njust2006",
	DB:       0,
})

var client2 = redis.NewClient(&redis.Options{
	Addr:     redis_Addr,
	Password: "njust2006",
	DB:       1,
})

var client3 = redis.NewClient(&redis.Options{
	Addr:     redis_Addr,
	Password: "njust2006",
	DB:       2,
})

var msg_client = redis.NewClient(&redis.Options{
	Addr:     redis_Addr,
	Password: "njust2006",
	DB:       3,
})

var config_client = redis.NewClient(&redis.Options{
	Addr:     redis_Addr,
	Password: "njust2006",
	DB:       4,
})

var doub_client = redis.NewClient(&redis.Options{
	Addr:     redis_Addr,
	Password: "njust2006",
	DB:       5,
})

var rules = map[string]string{}

func Contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in array")
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func random_key() string {
	rnd_key, err := client.RandomKey().Result()

	if err != nil {
		c := make(chan bool)
		go HttpGet(c)
		time.Sleep(1)
		rnd_key, _ = client.RandomKey().Result()
	}

	return rnd_key
}

func HttpServe(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	good_servers, _ := config_client.LRange("good_servers", 0, -1).Result()
	fmt.Println(good_servers)
	addr_info := strings.Split(r.RemoteAddr, ":")

	token := r.Header.Get("access-token")
	if token != "b25seS1mb3ItZmV3LXBlcnNvbnMtdGhhdC1yZWFsbHktbmVlZA" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `[	{"server": "100.100.100.100","server_port": "8888","password": "password","method": "aes-256-cfb","remarks": "请及时删除Xsocks!"}]`)
		return
	}

	var computer_id string

	for k, v := range r.Form {
		if k == "computerid" {
			for _, value := range v {
				computer_id = value
			}
		}
	}

	for k, v := range r.Form {
		if k == "key" {
			for _, value := range v {
				if value == "client-close" {
					client2.Del(addr_info[0] + "#" + computer_id).Result()
					client3.Del(addr_info[0] + "#" + computer_id).Result()
					resp := "{\"client-close\":\"OK\"}"
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintf(w, resp)
					return
				} else {
					ret, _ := Contains(value, good_servers)
					if !ret {
						client.Del(value).Result()
					}
				}
			}
		}
	}

	var keys []string
	var keys_err interface{}
	keys, keys_err = client.Keys("*").Result()

	if keys_err != nil {
		c := make(chan bool)
		go HttpGet(c)
		time.Sleep(1)
		keys, _ = client.Keys("*").Result()
	}

	var real_good_servers = []string{}
	var rnd_key string
	//var err interface{}

	for _, s := range good_servers {
		exists, _ := Contains(s, keys)
		if exists {
			real_good_servers = append(real_good_servers, s)
		}
	}

	if ip, ok := rules[addr_info[0]]; ok {
		rnd_key = ip
	} else if len(real_good_servers) > 0 {
		randIndex := rand.Intn(len(real_good_servers))
		rnd_key = real_good_servers[randIndex]
	} else {
		//rnd_key, err = client.RandomKey().Result()
		//
		//if err != nil {
		//	c := make(chan bool)
		//	go HttpGet(c)
		//	time.Sleep(1)
		//	rnd_key, _ = client.RandomKey().Result()
		//}
		rnd_key = random_key()
	}

	rjson, _ := client.Get(rnd_key).Result()
	if rjson == "" {
		rnd_key = random_key()
		rjson, _ = client.Get(rnd_key).Result()
	}

	client2.Set(addr_info[0]+"#"+computer_id, rnd_key, 0).Result()

	resp := fmt.Sprintf("[%s]", rjson)
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, resp)
}

func HttpMessageGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token := r.Header.Get("access-token")
	if token != "b25seS1mb3ItZmV3LXBlcnNvbnMtdGhhdC1yZWFsbHktbmVlZA" {
		fmt.Fprintf(w, "[]")
		return
	}

	addr_info := strings.Split(r.RemoteAddr, ":")
	var computer_id string
	var message string

	for k, v := range r.Form {
		if k == "computerid" {
			for _, value := range v {
				computer_id = value
			}
		}
	}

	message, _ = msg_client.Get(addr_info[0] + "#" + computer_id).Result()

	resp := fmt.Sprintf("[%s]", message)

	_, err := msg_client.Del(addr_info[0] + "#" + computer_id).Result()

	if err != nil {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, resp)
	}
}

func HttpMessagePut(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var message string

	for k, v := range r.Form {
		if k == "message" {
			for _, value := range v {
				message = value
			}
		}
	}

	db2_keys, _ := client3.Keys("*").Result()

	alive_keys := Filter(db2_keys, func(k string) bool {
		v, _ := client3.Get(k).Result()
		last_hb, _ := time.Parse("2006-01-02 15:04:05", v)
		subM := time.Now().UTC().Sub(last_hb)
		return subM.Seconds() < 360
	})

	for _, r2key := range alive_keys {
		msg_client.Set(r2key, message, 0).Result()
	}

	resp := fmt.Sprintf("[%s]", "success")

	fmt.Fprintf(w, resp)
}

func HttpAdmin(w http.ResponseWriter, r *http.Request) {
	var actions = []string{}
	r.ParseForm()
	for k, v := range r.Form {
		if k == "action" {
			for _, action := range v {
				actions = append(actions, action)
			}
		}
	}
	var resp = []string{}
	for _, action := range actions {
		if action == "flushdb" {
			ret, _ := client.FlushDb().Result()
			//ret2, _ := client2.FlushDb().Result()
			resp = append(resp, fmt.Sprintf("{\"flushdb\": {\"db0\": \"%s\"}}", ret))

		} else if action == "getall" {
			keys, keys_err := client.Keys("*").Result()
			if keys_err != nil {
				resp = append(resp, "{\"getall\": \"error\"}")
			} else {
				new_keys := Map(keys, func(s string) string {
					return fmt.Sprintf("\"%s\"", s)
				})
				resp = append(resp, fmt.Sprintf("{\"getall\": {\"count\": %d, \"keys\": [%s]}}",
					len(keys),
					strings.Join(new_keys, ",")))
			}

		} else if strings.HasPrefix(action, "get->") {
			parts := strings.Split(action, "->")
			sserver, err := client.Get(parts[1]).Result()
			if err != nil {
				resp = append(resp, fmt.Sprintf("{\"%s\": \"error\"}", action))
			} else {
				resp = append(resp, fmt.Sprintf("{\"%s\": %s}",
					action,
					sserver))
			}

		} else if strings.HasPrefix(action, "del->") {
			parts := strings.Split(action, "->")
			ret, err := client.Del(parts[1]).Result()
			if err != nil {
				resp = append(resp, fmt.Sprintf("{\"%s\": \"error\"}", action))
			} else {
				resp = append(resp, fmt.Sprintf("{\"%s\": %d}",
					action,
					ret))
			}

		} else if strings.HasPrefix(action, "add_rule->") {
			parts := strings.Split(action, "->")
			addr_info := strings.Split(r.RemoteAddr, ":")
			keys, _ := client.Keys("*").Result()
			exists, _ := Contains(parts[1], keys)
			if exists {
				rules[addr_info[0]] = parts[1]
				resp = append(resp, fmt.Sprintf("{\"%s\": \"success\"}", action))
			} else {
				resp = append(resp, fmt.Sprintf("{\"%s\": \"failed: no such server: %s\"}", action, parts[1]))
			}

		} else if strings.HasPrefix(action, "del_rule") {
			addr_info := strings.Split(r.RemoteAddr, ":")
			delete(rules, addr_info[0])
			resp = append(resp, fmt.Sprintf("{\"%s\": \"success\"}", action))

		} else if strings.HasPrefix(action, "usage") {
			db2_keys, _ := client3.Keys("*").Result()

			alive_keys := Filter(db2_keys, func(k string) bool {
				v, _ := client3.Get(k).Result()
				last_hb, _ := time.Parse("2006-01-02 15:04:05", v)
				subM := time.Now().UTC().Sub(last_hb)
				return subM.Seconds() < 360
			})
			var usage = map[string][]string{}

			for _, r2key := range alive_keys {
				r2value, _ := client2.Get(r2key).Result()
				if strings.HasPrefix(r2key, r2value) {
					usage[r2value] = append(usage[r2value], r2key+"(全局)")
				} else {
					usage[r2value] = append(usage[r2value], r2key)
				}

			}
			usage_json, _ := json.Marshal(usage)
			resp = append(resp, fmt.Sprintf("{\"%s\": %s}", action+"("+strconv.Itoa(len(alive_keys))+")", usage_json))

		} else {
			resp = append(resp, "{\"msg\": \"unsupported action\"}")
		}
	}

	resp_json := fmt.Sprintf("[%s]", strings.Join(resp, ","))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, resp_json)
}

func HttpDownload(w http.ResponseWriter, r *http.Request) {
	var action string
	r.ParseForm()
	for k, v := range r.Form {
		if k == "action" {
			action = v[0]
		}
	}
	xsocks_fullname := findXSocks()

	reg := regexp.MustCompile(`xsocks-(\d+\.\d+\.\d+).exe.gz`)
	results := reg.FindAllStringSubmatch(string(xsocks_fullname), -1)

	if action == "getversion" {
		fmt.Fprintf(w, results[0][1])
	} else if action == "getfile" {
		file := file_stored_path + xsocks_fullname
		if exist := isExist(file); !exist {
			http.NotFound(w, r)
		}
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", xsocks_fullname))
		w.Header().Add("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, file)
	} else {
		fmt.Fprintf(w, "unsupported action")
	}
}

func HttpHeartBeat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr_info := strings.Split(r.RemoteAddr, ":")

	token := r.Header.Get("access-token")
	if token != "b25seS1mb3ItZmV3LXBlcnNvbnMtdGhhdC1yZWFsbHktbmVlZA" {
		fmt.Fprintf(w, "[]")
		return
	}

	var computer_id string
	var server string = ""

	for k, v := range r.Form {
		if k == "computerid" {
			for _, value := range v {
				computer_id = value
			}
		}
	}

	for k, v := range r.Form {
		if k == "server" {
			for _, value := range v {
				server = value
			}
		}
	}

	if server != "" {
		client2.Set(addr_info[0]+"#"+computer_id, server, 0).Result()
	}

	ret, _ := client3.Set(addr_info[0]+"#"+computer_id, time.Now().UTC().Format("2006-01-02 15:04:05"), 0).Result()

	resp := fmt.Sprintf("{\"receive hb\":\"%s\"}", ret)
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, resp)
}

func HttpShell(w http.ResponseWriter, r *http.Request) {
	file := file_stored_path + "getss.sh"
	if exist := isExist(file); !exist {
		http.NotFound(w, r)
	}
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "getss.sh"))
	w.Header().Add("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, file)
}

func HttpRouter(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, _ := http.NewRequest("GET", "https://app.arukas.io/api/containers", nil)
	req.Header.Add("Authorization", "Basic "+basicAuth("784360d3-42b1-4bae-88c6-5e1336a853bc", "LHybOCmtEW1YwyRRLzW1EUMmQdRab5SUkACGdxL5ZBDX9n65cqI2YjugsWMhf9pP"))

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	body_as_bytes, _ := ioutil.ReadAll(resp.Body)

	body_as_str := string(body_as_bytes)
	reg := regexp.MustCompile(`"port_mappings":\s*\[\s*\[\s*\{([^\{\}]+)\}\s*\]`)

	row := reg.FindString(body_as_str)
	subreg := regexp.MustCompile(`\d+-\d+-\d+-\d+`)
	ip := strings.Replace(subreg.FindString(row), "-", ".", -1)
	subreg2 := regexp.MustCompile(`"service_port":\s*(\d+)`)
	port := subreg2.FindAllStringSubmatch(row, -1)[0][1]

	fmt.Fprintf(w, ip+" \""+port+"\"")
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return nil
}

func findXSocks() string {
	var xsocks_fullname string
	filepath.Walk(file_stored_path, func(path string, info os.FileInfo, err error) error {
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
