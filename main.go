package main

import (
	"bufio"
	"os"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
)

func main() {
	aSpec, err := apidoc.GetSpecFromFile("test_service.json")
	if err != nil {
		panic(err)
	}

	pFile := getProtoFromAPISpec(aSpec)

	//TODO: safe name?
	f, err := os.Create(aSpec.Name + ".proto")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	proto.NewFormatter(w, "  ").Format(&pFile)
	w.Flush()
}

const (
	syntaxVersion = "proto3"
)

func getProtoFromAPISpec(spec apidoc.Spec) proto.Proto {
	pFile := proto.Proto{}

	pFile.Elements = append(
		pFile.Elements,
		&proto.Syntax{Value: syntaxVersion},
	)

	for _, aEnum := range spec.Enums {
		pEnum := proto.Enum{Name: aEnum.Name}
		for i, aValue := range aEnum.Values {

			pEnum.Elements = append(
				pEnum.Elements,
				&proto.EnumField{Name:aValue.Name, Integer:i},
			)
		}
		pFile.Elements = append(pFile.Elements, &pEnum)
	}

	return pFile
}

