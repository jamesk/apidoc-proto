package apidoc

import (
	"encoding/json"
	"io/ioutil"
)

func GetSpecFromFile(filePath string) (Spec, error) {
	testSpecJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Spec{}, err
	}

	var spec Spec
	err = json.Unmarshal(testSpecJSON, &spec)
	if err != nil {
		return Spec{}, err
	}

	return spec, nil
}