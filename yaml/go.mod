module github.com/vitalik-mironov/go-labs/yaml

go 1.14

require sigs.k8s.io/kustomize/kyaml v0.4.1

require (
	github.com/stelligent/config-lint v1.6.0
	sigs.k8s.io/kustomize/api v0.5.1
)

replace (
	github.com/stelligent/config-lint => git.wildberries.ru/mironov.vitaliy3/config-lint v1.6.2
	github.com/valyala/fasthttp => git.wildberries.ru/mironov.vitaliy3/fasthttp v1.14.1
)
