package model

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"sync"
	"time"

	"../utils"
)

var waitgroup sync.WaitGroup

func save_ss(server string, port string, password string, method string, remarks string) {
	var ssConf = &SSConfig{
		Server:     server,
		ServerPort: port,
		Password:   password,
		Method:     method,
		Remarks:    remarks,
	}

	var sServer = &SServer{
		Name:   server,
		Speed:  60,
		Config: *ssConf,
	}

	sServer.Save()
}

func fetchByUrl(useProxy bool, url string, regexpStr string, remarks string) {
	waitgroup.Add(1)
	defer waitgroup.Done()

	var resp *http.Response
	var err error

	if useProxy {
		socks5client, err := utils.Socks5Client("127.0.0.1:1080")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		resp, err = socks5client.Get(url)
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		fmt.Println("error for httpGet:", err.Error())
		return
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := strings.Replace(string(body), "\r\n", "", -1)

		reg := regexp.MustCompile(regexpStr)

		for _, matches := range reg.FindAllStringSubmatch(bodyStr, -1) {
			save_ss(matches[1], matches[2], matches[3], strings.ToLower(matches[4]), remarks)
		}
	}
}

func getDoub() {
	waitgroup.Add(1)

	doub_html := func() string {
		client := &http.Client{
			Timeout: time.Duration(10 * time.Second),
		}

		header := map[string]string{"User-Agent":
		`Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36`,
			"Cache-Control": "max-age=0",
			"Upgrade-Insecure-Requests": "1"}

		request, _ := http.NewRequest("GET", "https://doub.bid/sszhfx/", nil)

		for k, v := range header {
			request.Header.Add(k, v)
		}

		if response, err := client.Do(request); err == nil {
			body, _ := ioutil.ReadAll(response.Body)
			return string(body)
		} else {
			return "error"
		}
	}()

	reg := regexp.MustCompile(`ss://([^"\s&#]+)`)
	for _, row := range reg.FindAllString(doub_html, -1) {
		s := strings.Split(row, "//")
		decodeInput := strings.Replace(s[1], "!", "", -1)
		data, err := base64.StdEncoding.DecodeString(decodeInput)
		if err != nil {
			sub_data, sub_err := base64.StdEncoding.DecodeString(decodeInput + "=")
			if sub_err != nil {
				continue
			} else {
				data = sub_data
			}
			//continue
		}

		subReg := regexp.MustCompile(`([^:@]+)`)
		filteredStr := strings.Replace(string(data), "\t", "", -1)

		items := subReg.FindAllString(filteredStr, -1)

		if len(items) == 4 {
			save_ss(items[2], items[3], items[1], items[0], "doub.io")
		}
	}
	waitgroup.Done()
}

func getFromGoogleGroup(useProxy bool, url string) {
	waitgroup.Add(1)
	defer waitgroup.Done()

	var resp *http.Response
	var err error

	if useProxy {
		socks5client, err := utils.Socks5Client("127.0.0.1:1080")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		resp, err = socks5client.Get(url)
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		fmt.Println("error for httpGet:", err.Error())
		return
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := strings.Replace(string(body), "\r\n", "", -1)

		reg := regexp.MustCompile(`ss://([^"\s&#%]+)`)
		for _, row := range reg.FindAllString(bodyStr, -1) {
			s := strings.Split(row, "//")
			decodeInput := strings.Replace(s[1], "!", "", -1)
			data, err := base64.StdEncoding.DecodeString(decodeInput)
			if err != nil {
				sub_data, sub_err := base64.StdEncoding.DecodeString(decodeInput + "=")
				if sub_err != nil {
					continue
				} else {
					data = sub_data
				}
				//continue
			}

			subReg := regexp.MustCompile(`([^:@]+)`)
			filteredStr := strings.Replace(string(data), "\t", "", -1)

			items := subReg.FindAllString(filteredStr, -1)

			if len(items) == 4 {
				fmt.Println(items)
				save_ss(items[2], items[3], items[1], items[0], "google+")
			}
		}
	}
}

func getManualSS() {
	waitgroup.Add(1)
	defer waitgroup.Done()

	save_ss("45.77.190.39", "8080", "FaLunDaFaHao@513", "aes-256-cfb", "github")

	save_ss("172.196.14.158", "6688", "ntdtv.com", "aes-256-cfb", "github0")
	save_ss("159.89.135.7", "999", "ntdtv.com", "aes-256-cfb", "github0")

	save_ss("mirror2.wordao.net", "8443", "fmqV4sCYoykN", "chacha20-ietf-poly1305", "google+")
	save_ss("180.188.196.157", "11622", "https://goo.gl/CGxNQN", "chacha20-ietf-poly1305", "google+")

	//save_ss("52.52.170.10", "8388", "leilei2010", "rc4-md5", "amazon")
}

func parseQRCode(url string, id string) {
	jar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar: jar,
	}

	reg := regexp.MustCompile(`//([^"]+)`)

	resp, _ := httpClient.Post("http://tool.oschina.net/action/qrcode/decode",
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

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			for _, subItems := range reg.FindAllStringSubmatch(string(body), -1) {
				if data, err1 := base64.StdEncoding.DecodeString(subItems[1]); err1 == nil {
					dataAsStr := strings.Replace(string(data), "\n", "", -1)

					s1 := strings.Split(dataAsStr, ":")
					s2 := strings.Split(s1[1], "@")

					save_ss(s2[1], s1[2], s2[0], s1[0], id)
				}
			}
		}
	}
}

func getSS8() {
	waitgroup.Add(1)

	//https://en.ss8.fun/images/server01.png
	//http://s8.6gg6.net/images/server01.png
	parseQRCode("https://en.ss8.fun/images/server01.png", "ss8")
	parseQRCode("https://en.ss8.fun/images/server02.png", "ss8")
	parseQRCode("https://en.ss8.fun/images/server03.png", "ss8")

	waitgroup.Done()
}

func FetchAll() {
	go getSS8()

	go fetchByUrl(false,
		"https://global.ishadowx.net/",
		`<h4>IP Address:<span id="ip[^"]*">([^<]+)<.*?Port:<span id="port[^"]*">(\d+)\s*<.*?Password:<span id="[^"]*">([^<\s]+)\s*<.*?Method:([^<]+)<`,
		"isdx")

	go getDoub()

	go getManualSS()

	go getFromGoogleGroup(true, "https://plus.google.com/communities/104092405342699579599/stream/8a593591-2091-4096-bb00-7d9c5659db93")

	go getFromGoogleGroup(true, "https://plus.google.com/u/0/communities/110896176697381748150")

	waitgroup.Wait()
}
