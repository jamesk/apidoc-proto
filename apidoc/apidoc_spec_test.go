package apidoc

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func TestUnmarshal(t *testing.T) {
	testSpecJSON, err := ioutil.ReadFile("../test_service.json")
	if err != nil {
		t.Fatal(err)
	}

	var testSpec Spec
	err = json.Unmarshal(testSpecJSON, &testSpec)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", testSpec)
}