package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	fullPath := home + "/.gatorconfig.json"

	jsonData, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, nil
	}

	var config Config
	err = json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fullPath := home + "/.gatorconfig.json"

	c.CurrentUserName = username

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}

	return nil
}
