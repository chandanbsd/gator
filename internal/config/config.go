package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const gatorConfigFileName string = ".gatorconfig.json"

func Read() (Config, error) {

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Cannot open gator config file")
		return Config{}, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Unable to read the contents of the gator config file")
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Unable to parse the config file")
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(userName string) {
	c.CurrentUserName = userName

	configFilePath, err:= getConfigFilePath()
	if err != nil {
		return
	}
	var bytes []byte
	bytes, err = json.Marshal(c)
	if err != nil {
		fmt.Println("SetUser: Bad data")
	}
	err = write(bytes, configFilePath)
	if err != nil {
		return
	}
}

func write(bytes []byte , filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("write: Failed to create file")
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println("write: Failed to write to file")
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to obtain your home directory")
		return "", err
	}


	configFilePath := fmt.Sprintf("%s/%s", homeDir, gatorConfigFileName)

	return configFilePath, nil
}
