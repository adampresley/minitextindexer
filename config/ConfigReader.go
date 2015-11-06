package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
)

/*
LoadConfiguration reads data from a Reader into a new Configuration structure.
*/
func LoadConfiguration(reader io.Reader) (*Configuration, error) {
	var err error
	var contents bytes.Buffer
	var buffer = make([]byte, 4096)
	var bytesRead int

	result := &Configuration{}
	bufferedReader := bufio.NewReader(reader)

	for {
		bytesRead, err = bufferedReader.Read(buffer)
		if err != nil && err != io.EOF {
			return result, err
		}

		if bytesRead == 0 {
			break
		}

		if _, err := contents.Write(buffer[:bytesRead]); err != nil {
			return result, err
		}
	}

	err = json.Unmarshal(contents.Bytes(), result)
	if err != nil {
		return result, err
	}

	return result, nil
}

/*
LoadConfigurationFromFile reads data from a file into a Configuration object. Makes use of
LoadConfiguration().
*/
func LoadConfigurationFromFile(fileName string) (*Configuration, error) {
	result := &Configuration{}

	configFileHandle, err := os.Open(fileName)
	if err != nil {
		return result, err
	}

	result, err = LoadConfiguration(configFileHandle)
	if err != nil {
		return result, err
	}

	return result, nil
}
