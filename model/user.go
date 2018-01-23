package model

import (
	"encoding/json"
	"fmt"
	"time"

	"../utils"
	"github.com/go-redis/redis"
)

//106.14.224.61
var RedisAddr = "localhost:6379"

var client4user = redis.NewClient(&redis.Options{
	Addr:     RedisAddr,
	Password: "njust2006",
	DB:       7,
})

type User struct {
	UserID        string
	Online        int    `json:"online"`
	SServer       string `json:"sserver"`
	ClientVersion string `json:"clientversion"`
	ClientIP      string `json:"clientip"`
	LastAliveTime string `json:"lastalivetime"`

	Messages []string `json:"messages"`
}

func GetUser(userid string) *User {
	users, _ := client4user.Keys("*").Result()
	exists, _ := utils.Contains(userid, users)

	if exists {
		userJsonStr, _ := client4user.Get(userid).Result()
		var user User
		err := json.Unmarshal([]byte(userJsonStr), &user)
		if err == nil {
			user.UserID = userid
			return &user
		} else {
			fmt.Println("GetUser error")
			return nil
		}

	} else {
		user := &User{UserID: userid}
		//user.save()
		return user
	}
}

func ScanOnlineUsers() {
	users := GetUsers()

	for _, user := range users {
		go func(user *User) {
			lastAliveTime, _ := time.Parse("2006-01-02 15:04:05", user.LastAliveTime)
			subM := time.Now().Sub(lastAliveTime)

			fmt.Println("user:", user.UserID, subM.Seconds())

			if subM.Seconds() > 360 {
				user.Online = 0
				user.Save()
			}
		}(user)
	}
}

func GetUsers() []*User {
	userNames, err := client4user.Keys("*").Result()

	if err != nil {
		return nil
	} else {
		var users = []*User{}

		for _, id := range userNames {
			users = append(users, GetUser(id))
		}
		return users
	}
}

func (user *User) Save() {
	//timeLocal := time.FixedZone("CST", 3600*8)
	//time.Local = timeLocal

	//lastAliveTime := time.Now().Format("2006-01-02 15:04:05")
	//user.LastAliveTime = lastAliveTime

	if jsonStr, err := json.Marshal(user); err == nil {
		client4user.Set(user.UserID, jsonStr, 0)
	} else {
		fmt.Println(err.Error())
	}
}

func (user *User) Update(userInfo map[string]string) {
	var ss_server, cli_Version, cli_IP string

	if sserver, ok := userInfo["sserver"]; ok && sserver != "" {
		ss_server = sserver
	}

	if cliVersion, ok := userInfo["clientversion"]; ok && cliVersion != "" {
		cli_Version = cliVersion
	}

	if cliIP, ok := userInfo["userid"]; ok && cliIP != "" {
		cli_IP = cliIP
	}

	user.SServer = ss_server
	user.ClientVersion = cli_Version
	user.ClientIP = cli_IP

	user.Save()
}

func (user *User) PushMessages(message string) {
	user.Messages = append(user.Messages, message)
	user.Save()
}

func (user *User) ConsumeMessages() []string {
	messages := user.Messages
	user.Messages = nil
	user.Save()
	return messages
}
