# Options.
#
GOOS ?= linux

# Targets.
#
BUILD_TARGETS := build
DEPLOY_TARGETS := devel
PHONY_TARGETS := $(BUILD_TARGETS) $(DEPLOY_TARGETS)
.PHONY: $(PHONY_TARGETS)

build_tlshttp_selfhost: Dockerfile.tlshttp.selfhost
	docker build 								\
		--build-arg "GOOS=$(GOOS)" 				\
		-t tlshttp_selfhost .