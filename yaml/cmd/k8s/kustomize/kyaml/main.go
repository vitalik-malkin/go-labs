package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/stelligent/config-lint/assertion"
	kio "sigs.k8s.io/kustomize/kyaml/kio"
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
	res, err := loadLinterResource(kyamltext)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res {
		fmt.Printf("ID: %s (%s)\n", v.ID, v.Type)
	}

	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	fmt.Printf("text: %s", sc.Text())
}

type locRes assertion.Resource

func loadLinterResource(manifest string) (resources []assertion.Resource, err error) {
	resReader := &kio.ByteReader{
		Reader:                bytes.NewBufferString(manifest),
		OmitReaderAnnotations: true,
	}
	resNodes, err := resReader.Read()
	if err != nil {
		log.Fatal(err)
	}
	var _ []*kyaml.RNode = resNodes
	for _, resNode := range resNodes {
		resYNode := resNode.YNode()
		resMeta, err := resNode.GetMeta()
		if err != nil && err != kyaml.ErrMissingMetadata {
			return nil, err
		}
		var lntResID, lntResType string
		if err == kyaml.ErrMissingMetadata {
			lntResID, lntResType = "<n/a>", "<n/a>"
			err = nil
		} else {
			lntResID, lntResType = fmt.Sprintf("%s/%s", resMeta.Namespace, resMeta.Name), resMeta.Kind
		}
		lntResProps := map[string]interface{}{}
		err = resYNode.Decode(lntResProps)
		if err != nil {
			return nil, err
		}
		lntRes := assertion.Resource{
			ID:         lntResID,
			Type:       lntResType,
			LineNumber: resYNode.Line,
			Properties: lntResProps,
		}
		resources = append(resources, lntRes)
	}
	return
}

func (l *locRes) m1() string {
	return ""
}

func (l locRes) m1() string {
	return ""
}
