package xsocks

import (
	"fmt"
	"regexp"
	"strings"
	"encoding/base64"
	"net/http"
	"io/ioutil"
	"net/http/cookiejar"
	"net"
	"time"
	"golang.org/x/net/proxy"
	"sync"
	"io"
)

var waitgroup sync.WaitGroup

func HttpGet(c chan bool) {
	//cookie, err := doub_client.Get("doub_cookie").Result()
	//if err == nil {
	//	fmt.Println(cookie)
	//	go getSS1(cookie)
	//} else {
	//	fmt.Println("get doub cookie error: ", err)
	//}

	html, err := doub_client.Get("doub_html").Result()
	if err == nil {
		//fmt.Println(html)
		go getSS0(html)
	} else {
		fmt.Println("get doub html error: ", err)
	}
	go getSS2()
	go getSS3()
	go getSS4()
	go getSS5()
	go getSS6()
	go getSS7()
	go getSS8()
	go getSS9()
	go getSS10()

	waitgroup.Wait()
	c <- true
}

func update_r_item(key string, value string) {
	if exists, err := client.Exists(key).Result(); err == nil {
		if exists == 0 {
			client.Set(key, value, 0).Result()
		}
	}
}

func getSS0(doub_html string) {
	waitgroup.Add(1)

	reg := regexp.MustCompile(`ss://([^"\s&#]+)`)

	for _, row := range reg.FindAllString(doub_html, -1) {
		s := strings.Split(row, "//")
		data, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			sub_data, sub_err := base64.StdEncoding.DecodeString(s[1] + "==")
			if sub_err != nil {
				continue
			} else {
				data = sub_data
			}
			//continue
		}

		subReg := regexp.MustCompile(`([^:@]+)`)
		items := subReg.FindAllString(string(data), -1)

		if len(items) == 4 {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[2], items[3], items[1], items[0], "doub.io")
			//client.Set(items[2], json, 0).Result()
			update_r_item(items[2], json)
		}
	}
	waitgroup.Done()
}

func getSS1(cookie string) {
	waitgroup.Add(1)

	doub_html := getHtml(cookie)
	fmt.Println(doub_html)
	reg := regexp.MustCompile(`(ss|ssr)://([^"\s&#]+)`)

	for _, row := range reg.FindAllString(doub_html, -1) {
		s := strings.Split(row, "//")
		data, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			sub_data, sub_err := base64.StdEncoding.DecodeString(s[1] + "==")
			if sub_err != nil {
				continue
			} else {
				data = sub_data
			}
			//continue
		}

		subReg := regexp.MustCompile(`([^:@/]+)`)
		items := subReg.FindAllString(string(data), -1)
		if len(items) == 7 {
			password, _ := base64.StdEncoding.DecodeString(items[5] + "==")
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[0], items[1], password, items[3], "doub.io")
			//client.Set(items[0], json, 0).Result()
			update_r_item(items[0], json)
		} else if len(items) == 4 {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[2], items[3], items[1], items[0], "doub.io")
			update_r_item(items[2], json)
		}

	}
	waitgroup.Done()
}

func getSS2() {
	waitgroup.Add(1)
	defer waitgroup.Done()

	set_ss("https://coding.net/u/Alvin9999/p/ip/git/raw/master/ssconfig.txt")
	//set_ss("https://raw.githubusercontent.com/Alvin9999/pac2/master/ssconfig.txt")
	//waitgroup.Done()
}

func getSS3() {
	waitgroup.Add(1)
	var url2 = "http://ss.ishadowx.com/"

	if resp, err := http.Get(url2); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		body_as_str := strings.Replace(string(body), "\r\n", "", -1)
		reg := regexp.MustCompile(`<h4>IP Address:<span id="ip[^"]*">([^<]+)<.*?Port：(\d+)<.*?Password:<span id="[^"]*">([^<]+)<.*?Method:([^<]+)<`)

		for _, items := range reg.FindAllStringSubmatch(body_as_str, -1) {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[1], items[2], items[3], strings.ToLower(items[4]), "ss3")
			//client.Set(items[1], json, 0).Result()
			update_r_item(items[1], json)
		}
	}
	waitgroup.Done()
}

func getSS4() {
	waitgroup.Add(1)

	socks5client, err0 := Socks5Client("127.0.0.1:1080")
	if err0 != nil {
		fmt.Println(err0)
		return
	}

	//var url2 = "https://get.freevpnss.me/"
	var url2 = "http://jn9.org/"

	if resp, err := socks5client.Get(url2); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		body_as_str := strings.Replace(string(body), "\n", "", -1)

		reg := regexp.MustCompile(`>Adress\(IP\)：([^<]+)</p>\s*<p>Port：(\d+)<.*?<p>Password<span class="hidden">[d|e|f]</span>：([^<\s]+)<.*?Method：([^<]+)<`)

		for _, items := range reg.FindAllStringSubmatch(body_as_str, -1) {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[1], items[2], items[3], strings.ToLower(items[4]), "ss4")
			//client.Set(items[1], json, 0).Result()
			update_r_item(items[1], json)
		}
	}
	waitgroup.Done()
}

func getSS5() {
	waitgroup.Add(1)
	var url2 = "http://www.shadowsock.net/index.php"

	if resp, err := http.Get(url2); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		body_as_str := strings.Replace(string(body), "\r\n", "", -1)
		reg := regexp.MustCompile(`<h4>IP地址:<span id="ip[^"]*">([^<]+)<.*?服务器端口：(\d+)<.*?密码:<span id="[^"]*">([^<\s]+)\s*<.*?加密: ([^<]+)<`)

		for _, items := range reg.FindAllStringSubmatch(body_as_str, -1) {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[1], items[2], items[3], strings.ToLower(items[4]), "ss5")
			//client.Set(items[1], json, 0).Result()
			update_r_item(items[1], json)
		}
	}
	waitgroup.Done()
}

func getSS6() {
	waitgroup.Add(1)
	var url2 = "https://freessr.xyz/"

	if resp, err := http.Get(url2); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		body_as_str := strings.Replace(string(body), "\r\n", "", -1)

		reg := regexp.MustCompile(`<h4>服务器地址:([^<]+)</h4>\s*<h4>端口:(\d+)</h4>\s*<h4>密码:([^<\s]+)</h4>\s*<h4>加密方式:([^<]+)</h4>`)

		for _, items := range reg.FindAllStringSubmatch(body_as_str, -1) {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[1], items[2], items[3], strings.ToLower(items[4]), "ss6")
			//client.Set(items[1], json, 0).Result()
			update_r_item(items[1], json)
		}
	}
	waitgroup.Done()
}

func getSS7() {
	waitgroup.Add(1)
	json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
		"ss1.ubox.co", "8388", "password", "aes-256-cfb", "ss7")
	//client.Set("ss1.ubox.co", json, 0).Result()
	update_r_item("ss1.ubox.co", json)

	json2 := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
		"45.76.216.97", "8989", "wf123", "aes-256-cfb", "ss77")
	//client.Set("45.76.216.97", json2, 0).Result()
	update_r_item("45.76.216.97", json2)
	waitgroup.Done()
}

func getSS8() {
	waitgroup.Add(1)

	parse_qrcode("http://s8.6gg6.net/images/server01.png", "ss8")
	parse_qrcode("http://s8.6gg6.net/images/server02.png", "ss8")
	parse_qrcode("http://s8.6gg6.net/images/server03.png", "ss8")

	waitgroup.Done()
}

func getSS9() {
	waitgroup.Add(1)
	var url2 = "https://ss.potvpn.com/"
	resp, err := http.Get(url2)

	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		body_as_str := strings.Replace(string(body), "\r\n", "", -1)

		reg := regexp.MustCompile(`(/img/qrcode_image/\d+/[^"]+)`)

		for _, items := range reg.FindAllStringSubmatch(body_as_str, -1) {
			parse_qrcode("https://ss.potvpn.com"+items[1], "ss9")
		}
	}

	waitgroup.Done()
}

func getSS10() {
	waitgroup.Add(1)
	var url2 = "https://91vps.us/2017/05/30/shadowsocks_share/"
	resp, err := http.Get(url2)

	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		body_as_str := strings.Replace(string(body), "\r\n", "", -1)

		reg := regexp.MustCompile(`<td>(洛杉矶|日本|新加坡|荷兰)</td>\s*<td>([^<]+)</td>\s*<td>([^<]+)</td>\s*<td>([^<]+)</td>\s*<td>([^<]+)</td>`)

		for _, items := range reg.FindAllStringSubmatch(body_as_str, -1) {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[2], items[3], items[4], strings.ToLower(items[5]), "ss10-"+items[1])
			//client.Set(items[2], json, 0).Result()
			update_r_item(items[2], json)
		}
	}

	waitgroup.Done()
}

func parse_qrcode(url string, id string) {
	jar, _ := cookiejar.New(nil)

	http_client := &http.Client{
		Jar: jar,
	}

	reg := regexp.MustCompile(`//([^"]+)`)

	resp, _ := http_client.Post("http://www.sojson.com/deqr.html",
		"multipart/form-data; boundary=----WebKitFormBoundaryeVDH8UqDVqKni9d8",
		strings.NewReader("------WebKitFormBoundaryeVDH8UqDVqKni9d8\r\n"+
			"Content-Disposition: form-data; name=\"upload_ctn\"\r\n"+
			"\r\n"+
			"on\r\n"+
			"------WebKitFormBoundaryeVDH8UqDVqKni9d8\r\n"+
			"Content-Disposition: form-data; name=\"url\"\r\n"+
			"\r\n"+
			url+ "\r\n"+
			"------WebKitFormBoundaryeVDH8UqDVqKni9d8--"))

	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		for _, subItems := range reg.FindAllStringSubmatch(string(body), -1) {
			if data, err1 := base64.StdEncoding.DecodeString(subItems[1]); err1 == nil {
				data_as_str := strings.Replace(string(data), "\n", "", -1)

				s1 := strings.Split(data_as_str, ":")
				s2 := strings.Split(s1[1], "@")

				json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
					s2[1], s1[2], s2[0], s1[0], id)
				//client.Set(s2[1], json, 0).Result()
				update_r_item(s2[1], json)
			}
		}
	}
}

func set_ss(url string) {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		body_base64, _ := ioutil.ReadAll(resp.Body)
		body, _ := base64.StdEncoding.DecodeString(string(body_base64))
		body_as_str := strings.Replace(string(body), "\r\n", "", -1)

		reg := regexp.MustCompile(`"server"\s*:\s*"([^"]+)"\s*,\s*"server_port"\s*:\s*(\d+),\s*.*?"password"\s*:\s*"([^"]+)",\s*"method"\s*:\s*"([^"]+)"`)

		for _, items := range reg.FindAllStringSubmatch(body_as_str, -1) {
			json := fmt.Sprintf("{\"server\": \"%s\", \"server_port\": \"%s\", \"password\": \"%s\",\"method\": \"%s\", \"remarks\":\"%s\"}",
				items[1], items[2], items[3], strings.ToLower(items[4]), "Alvin9999")
			//client.Set(items[1], json, 0).Result()
			update_r_item(items[1], json)
		}
	}
}

func getHtml(cookie string) string {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar:     jar,
		Timeout: time.Duration(10 * time.Second),
	}

	header := map[string]string{"User-Agent":
	`Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36`,
		"Cache-Control": "max-age=0",
		"Upgrade-Insecure-Requests": "1"}
	header["Cookie"] = cookie

	resp := httpRequest(client, "GET", "https://doub.bid/sszhfx/", nil, header)
	return resp
}

func httpRequest(client *http.Client, method string, url string, body io.Reader, header map[string]string) (string) {
	request, _ := http.NewRequest(method, url, body)
	if header != nil && len(header) > 0 {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}

	if response, err := client.Do(request); err == nil {
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)
	} else {
		return "error"
	}
}

func Socks5Client(addr string, auth ...*proxy.Auth) (client *http.Client, err error) {
	dialer, err := proxy.SOCKS5("tcp", addr,
		nil,
		&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	)
	if err != nil {
		return
	}

	transport := &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client = &http.Client{Transport: transport}

	return
}
