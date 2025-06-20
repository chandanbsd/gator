package gator/internal/config

import (
	"os"
	"fmt"
	"json"
	"errors"
)

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const gatorConfigFileName string = ".gatorconfig.json"

func Read() Config, error {

	configFilePath := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Cannot open gator config file")
		return nil, err
	}

	defer file.Close()

	data := []bytes{}
	_, err := file.Read(data)
	if err != nil {
		fmt.Println("Unable to read the contents of the gator config file")
		return nil, err
	}

	var config Config

	err := json.Unmarshall(data, &config)
	if err != nil {
		fmt.Println("Unable to parse the config file")
		return nil, err
	}

	return config, nil
}

func (c *Config) SetUser(userName string) {
	c.UserName = userName

	configFilePath := getConfigFilePath()

	var bytes []bytes
	err := json.Unmarshall(bytes, c)
	if err != nil {
		fmt.Printf("SetUser: Bad data")
	}
	err := write(err)
	if err != nil {
		return
	}

	return nil
}

func write(bytes []byte , filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("write: Failed to open file")
		return err
	}

	_, err := file.Write(bytes)
	if err != nil {
		fmt.Println("write: Failed to write to file")
		return err
	}

	return nil
}

func getConfigFilePath() string, error {
	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to obtain your home directory")
		return "", err
	}


	filePath = fmt.Sprintf("%s/%s", homeDir, gatorConfigFileName)
	configFilePath = fmt.Sprintf("%s/%s", homeDir, gatorConfigFileName)

	return configFilePath, nil
}
