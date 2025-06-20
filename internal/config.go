package "internal/config"

import (
	"os",
	"fmt",
	"json",
	"errors"
)

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const gatorConfigFileName string = ".gatorconfig.json"

func Read() Config {
	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to obtain your home directory")
		return
	}

	configFilePath = fmt.Sprintf("%s/%s", homeDir, gatorConfigFileName)

	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Cannot open gator config file")
		return
	}
	
	data := []bytes{}
	_, err := file.Read(data)
	if err != nil {
		fmt.Println("Unable to read the contents of the gator config file")
	}

	var config Config
	
	err := json.Unmarshall(data, &config)
	if err != nil {
		fmt.Println("Unable to parse the config file")
	}

	return config
}

func SetUser(userName string) {
	
}

func getConfigFilePath() string, error {
	


	filePath = fmt.Sprintf("%s/%s", homeDir, gatorConfigFileName)
}
