package apidoc2proto

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"fmt"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
	"github.com/jamesk/apidoc-proto/aproto"
)

const (
	testFolder         = "files_for_test"
	testSpecExtension  = ".json"
	testProtoExtension = ".proto"

	testSpecGlob = "*" + testSpecExtension
)

func TestEndToEnd(t *testing.T) {
	fileNames, err := filepath.Glob(filepath.Join(testFolder, testSpecGlob))
	if err != nil {
		t.Fatalf("Error getting test fileNames directory: %v", err)
	}

	for _, fileName := range fileNames {
		//test the file
		aSpec, err := apidoc.GetSpecFromFile(fileName)
		if err != nil {
			t.Error(err)
		}

		pFile := aproto.GetProtoFromAPISpec(aSpec)

		protoFileName := fmt.Sprintf("%v%v", fileName[:len(fileName)-len(testSpecExtension)], testProtoExtension)
		err = outputProtoFile(protoFileName, pFile)
		if err != nil {
			t.Error(err)
		}
	}
}

func outputProtoFile(outputPath string, pFile proto.Proto) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	proto.NewFormatter(w, "  ").Format(&pFile)

	return w.Flush()
}
