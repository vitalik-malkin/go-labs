package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
)

// Builder ...
type Builder interface {
	Build(out io.Writer) error
}

type builder struct {
	options           *krusty.Options
	kustomizationPath string
}

func (b *builder) Build(out io.Writer) error {
	fileSys := filesys.MakeFsOnDisk()
	kustomizer := krusty.MakeKustomizer(fileSys, b.options)
	resMap, err := kustomizer.Run(b.kustomizationPath)
	if err != nil {
		return err
	}
	manifestYAML, err := resMap.AsYaml()
	_, err = out.Write(manifestYAML)
	return err
}

// NewBuilder ...
func NewBuilder(kustomizationPath string) Builder {
	return &builder{
		kustomizationPath: kustomizationPath,
		options:           krusty.MakeDefaultOptions(),
	}
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	kustomizationPath := "../testdata/wb-deploy-service/gateway/overlays/test"
	outputPath := filepath.Join(wd, "/output.yaml")

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	builder := NewBuilder(kustomizationPath)
	err = builder.Build(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
}
