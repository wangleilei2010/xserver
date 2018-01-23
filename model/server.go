package model

import (
	"encoding/json"
	"fmt"
	"time"

	"../utils"
	"github.com/go-redis/redis"
)

var client4server = redis.NewClient(&redis.Options{
	Addr:     RedisAddr,
	Password: "njust2006",
	DB:       8,
})

type SServer struct {
	Name   string
	Speed  float64  `json:"speed"`
	Config SSConfig `json:"config"`
}

type SSConfig struct {
	Server     string `json:"server"`
	ServerPort string `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	Remarks    string `json:"remarks"`
}

func GetServer(servName string) *SServer {
	servers, _ := client4server.Keys("*").Result()
	exists, _ := utils.Contains(servName, servers)

	if exists {
		serverJsonStr, _ := client4server.Get(servName).Result()
		var server SServer
		err := json.Unmarshal([]byte(serverJsonStr), &server)
		if err == nil {
			server.Name = servName
			return &server
		} else {
			fmt.Println("GetServer error")
			return nil
		}

	} else {
		server := &SServer{Name: servName}
		//user.save()
		return server
	}
}
func GetConfigs() []SSConfig {
	var servers []string
	var err interface{}

	servers, err = client4server.Keys("*").Result()

	if err != nil {
		go FetchAll()
		time.Sleep(1)
		servers, _ = client4server.Keys("*").Result()
	}

	var configs = []SSConfig{}

	for _, servName := range servers {
		serverJsonStr, _ := client4server.Get(servName).Result()
		var server SServer
		err := json.Unmarshal([]byte(serverJsonStr), &server)

		if err == nil {
			configs = append(configs, server.Config)
		}
	}
	return configs
}

func GetServers() []string {
	var servers []string
	var err interface{}

	servers, err = client4server.Keys("*").Result()

	if err != nil {
		go FetchAll()
		time.Sleep(1)
		servers, _ = client4server.Keys("*").Result()
	}

	return servers
}

func FlushServersDB() string {
	ret, _ := client4server.FlushDb().Result()
	return ret
}

func GetValidAndFastServers(fastServers []string) []string {
	var validAndFastServers = []string{}

	validServers := GetServers()

	for _, s := range fastServers {
		exists, _ := utils.Contains(s, validServers)
		if exists {
			validAndFastServers = append(validAndFastServers, s)
		}
	}

	return validAndFastServers
}

func GetFastestServer(except string) *SServer {
	var validServers []string
	validServers = GetServers()

	if len(validServers) == 0 {
		go FetchAll()
		time.Sleep(1)
		validServers = GetServers()
	}

	var begin_idx = 0

	if validServers[0] == except && len(validServers) > 1 {
		begin_idx = 1
	}

	var fastest = GetServer(validServers[begin_idx])

	for _, s := range validServers {
		if s == except {
			fmt.Println("continue block:", s, except)
			continue
		}
		curServer := GetServer(s)

		if curServer.Speed < fastest.Speed {
			fastest = curServer
		}
	}
	return fastest
}

func GetRandomServer() *SServer {
	rnd_key := func() string {
		rnd_key, err := client4server.RandomKey().Result()

		if err != nil {
			go FetchAll()
			time.Sleep(1)
			rnd_key, _ = client4server.RandomKey().Result()
		}

		return rnd_key
	}()

	return GetServer(rnd_key)
}

func (server *SServer) Save() {
	if jsonStr, err := json.Marshal(server); err == nil {
		client4server.Set(server.Name, jsonStr, 0)
	} else {
		fmt.Println(err.Error())
	}
}

func (server *SServer) Del() int64 {
	if ret, err := client4server.Del(server.Name).Result(); err != nil {
		fmt.Println(err.Error())
		return -1
	} else {
		return ret
	}
}
