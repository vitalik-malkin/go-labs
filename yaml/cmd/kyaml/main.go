package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/stelligent/config-lint/assertion"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	kyamltext = `---
---
---
---
apiVersion: v1
kind: Service
metadata:
  namespace: deploy-service
  name: vault
  labels:
    app: vault
spec:
  ports:
    - name: vault
      port: 8200
  clusterIP: None
  selector:
    app: vault
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  namespace: deploy-service
  name: vault
spec:
  selector:
    matchLabels:
      app: vault
  serviceName: "vault"
  replicas: 1
  template:
    metadata:
      labels:
        app: vault
    spec:
      containers:
        - name: vault
          image: git.wildberries.ru:4567/infrastructure/deploy-service/deploy-service/vault:v0.0.1-test-284-23
          ports:
            - containerPort: 8200
              name: vault
          securityContext:
            capabilities:
              add: ["IPC_LOCK"]
      imagePullSecrets:
        - name: gitlab-registry-secret
---
---
---
a: n
  `
)

func main() {

	// resTextSet := strings.Split(kyamltext, "\n---\n")
	// for _, resText := range resTextSet {
	// 	node := &y.Node{}
	// 	y.Unmarshal([]byte(resText), node)

	// 	switch node.Kind {
	// 	case y.DocumentNode:
	// 		fmt.Printf("document, %v\n", node.Content[0].Tag)
	// 	default:
	// 		fmt.Printf("unknown\n")
	// 	}
	// }

	resReader := &kio.ByteReader{
		Reader:                bytes.NewBufferString(kyamltext),
		OmitReaderAnnotations: true,
	}
	resNodes, err := resReader.Read()
	if err != nil {
		log.Fatal(err)
	}
	var _ []*yaml.RNode = resNodes
	for i, resNode := range resNodes {
		resYNode := resNode.YNode()
		lineNum := resNode.YNode().Line
		resMeta, err := resNode.GetMeta()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d, %d: name=%s; ns=%s; kind=%s\n", i, lineNum, resMeta.Name, resMeta.Namespace, resMeta.Kind)
	}

	// var _ []string = values

	// err := kio.Pipeline{
	// 	Inputs: []kio.Reader{&kio.ByteReader{Reader: bytes.NewBufferString(kyamltext)}},
	// }.Execute()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// h := make([]interface{}, 0, 16)
	// err := y.NewDecoder(bytes.NewReader([]byte(kyamltext))).Decode(h)
	// if err != nil {
	// 	fmt.Printf("error: %v'n", err)
	// }

	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	fmt.Printf("text: %s", sc.Text())
}

func loadLinterResource(manifest string) (resources []assertion.Resource, err error) {
	resReader := &kio.ByteReader{
		Reader:                bytes.NewBufferString(manifest),
		OmitReaderAnnotations: true,
	}
	resNodes, err := resReader.Read()
	if err != nil {
		log.Fatal(err)
	}
	var _ []*yaml.RNode = resNodes
	for i, resNode := range resNodes {
		resYNode := resNode.YNode()
		lineNum := resNode.YNode().Line
		resMeta, err := resNode.GetMeta()
		if err != nil && err != kyaml.ErrMissingMetadata {
			return nil, err
		}
		var lntResID, lntResType string
		if err == kyaml.ErrMissingMetadata {
			lntResID, lntResType = "<n/a>", "<n/a>"
		} else {
			lntResID = fmt.Sprintf("%s.%s", resMeta.Namespace, resMeta.Name)
		}

		lntRes := assertion.Resource{
			ID:         lntResID,
			Type:       lntResType,
			LineNumber: lineNum,
		}

		fmt.Printf("%d, %d: name=%s; ns=%s; kind=%s\n", i, lineNum, resMeta.Name, resMeta.Namespace, resMeta.Kind)
	}

}
