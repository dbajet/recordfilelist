package raicredentials

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type MysqlUser struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func Credentials(credentialsFilepath string) (MysqlUser, error) {
	result := MysqlUser{}
	file, err := os.Open(credentialsFilepath)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	return result, err
}
