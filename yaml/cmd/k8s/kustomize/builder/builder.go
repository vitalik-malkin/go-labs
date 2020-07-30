package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

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
	var s1 []string
	var s2 []string = []string{"A", "B"}
	copy(s2, s1)

	var vi1 i1
	var vi2 i2 = &t1{}

	vi1 = vi2

	var _ = vi1

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// kustomizationPath := "../testdata/wb-deploy-service/gateway/overlays/test"
	// outputPath := filepath.Join(wd, "/output.yaml")

	// outputFile, err := os.Create(outputPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// builder := NewBuilder(kustomizationPath)
	// err = builder.Build(outputFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
}

func super3(v interface{}) {
	switch x := v.(type) {
	case t1:
		fmt.Printf("%v", "t1")
	case *t1:
		fmt.Printf("%v", x == nil)
	default:
		fmt.Printf("unknown")
	}
}

func super1(v i1) string { return v.M1() }

func super2(v i2) string { return v.M1() }

type i1 interface {
	M1() string
	M2() int

	P1() i3
}

type i2 = interface {
	M1() string
	M2() int

	P1() i4
}

type i3 = interface {
	M5() int
}

type i4 = interface {
	M5() int
}

type t1 struct{}

func (t *t1) M1() string { return "hello" }

func (t *t1) M2() int { return 0 }

func (t *t1) P1() i4 { return nil }

func super4() {
	for i := 1; i <= 5; i++ {
		x := i
		defer func() { fmt.Printf("defer %v\n", x) }()
		fmt.Printf("body %v\n", i)
	}
}
