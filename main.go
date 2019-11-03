package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"

	// "github.com/ogurtsov/deploy5p/sshclient"

	"github.com/ogurtsov/deploy5p/ssh"
	"github.com/sirupsen/logrus"
)

var command string

type RealmConfig struct {
	Name     string
	Host     string
	Port     int
	Username string
	Password string
	Commands []string
}

type Config struct {
	TelegramAPIKey string
	TelegramChatID string
	Realms         []RealmConfig
}

var config Config

var path = ".deploy5p.json"

func getFilePath() string {
	usr, err := user.Current()
	if err != nil {
		logrus.Fatal(err)
	}
	return usr.HomeDir + "/" + path
}

func loadConfig() {
	filePath := getFilePath()
	configFile, _ := os.Open(filePath)
	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		// panic("Unable to load config")
	}
	json.Unmarshal(data, &config)
}

func TelegramSend(text string) {
	logger := logrus.New()
	api_key := config.TelegramAPIKey
	chat_id := config.TelegramChatID

	logger.Info(fmt.Sprintf("Trying to send: %s\n", text))

	output := make(map[string]interface{})

	response, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", api_key, chat_id, text))
	if err != nil {
		panic(err)
	}

	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Sending succeeded: %s\n", output))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func save() {
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(b))
	// d1 := []byte("hello\ngo\n")

	if fileExists(getFilePath()) {
		panic("Default config file already created")
	}

	f, err2 := os.Create(getFilePath())
	if err2 != nil {
		fmt.Println(err2)
		panic("Unable to create file")
	}
	defer f.Close()
	f.Write(b)

}

func list() {
	loadConfig()
	fmt.Println(config.Realms)
}

func setdefaults() {
	config = Config{
		TelegramAPIKey: "",
		TelegramChatID: "",
		Realms: []RealmConfig{
			RealmConfig{
				Name:     "default",
				Host:     "127.0.0.1",
				Port:     22,
				Username: "username",
				Password: "password",
				Commands: []string{"whoami", "pwd"},
			},
		},
	}
	save()
}

func initDeploy(Realm RealmConfig) {
	logrus.Info("Deploy for <" + Realm.Name + "> started...")
	// TelegramSend("Starting deploy for realm <" + Realm.Name + ">")
	host := Realm.Host + ":" + strconv.Itoa(Realm.Port)

	usr, err := user.Current()
	if err != nil {
		logrus.Fatal(err)
	}

	connection, err := ssh.ConnectWithKey(host, Realm.Username, Realm.Password, usr.HomeDir+"/.ssh/id_rsa")
	if err != nil {
		panic(err)
	}

	commands := strings.Join(Realm.Commands, " && ")
	fmt.Println(commands)
	output, err := connection.SendCommands(commands)
	if err != nil {
		panic(err)
	}
	logrus.Info(string(output))

	logrus.Info("Deploy for <" + Realm.Name + "> finished!")
}

func deploy(RealmName string) {
	loadConfig()
	for i := range config.Realms {
		if config.Realms[i].Name == RealmName {
			initDeploy(config.Realms[i])
			return
		}
	}
	panic("Wrong realm name")
}

func main() {

	flag.Parse()
	command = flag.Arg(0)

	switch command {
	case "list":
		list()
	case "setdefaults":
		setdefaults()
	case "deploy":
		deploy(flag.Arg(1))
	default:
		fmt.Println("Available commands: list, setdefaults, deploy")
	}
}
