package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
	"github.com/jamesk/apidoc-proto/aproto"
)

const usageMessage = `apidoc-proto is a tool to convert apidoc.me to .proto format

Usage:
	apidoc-proto INPUT.JSON OUTPUT.proto
`

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Printf("Need 1 or 2 arguments:\n%v", usageMessage)
		return
	}

	fmt.Printf("Reading apidoc spec: [%v]\n", os.Args[1])
	aSpec, err := apidoc.GetSpecFromFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	pFile := aproto.GetProtoFromAPISpec(aSpec)

	//TODO: safe name?
	outputFilepath := aSpec.Name + ".proto"
	if len(os.Args) > 2 {
		outputFilepath = os.Args[2]
	}
	fmt.Printf("Outputting .proto file: [%v]\n", outputFilepath)
	f, err := os.Create(outputFilepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	proto.NewFormatter(w, "  ").Format(&pFile)
	w.Flush()
}
