Project := Code Push
BuildDist = build
ReleaseDist = release

Platforms = linux darwin windows

Cmds := $(foreach n,$(shell go list ./cmd/*),$(notdir $(n)))

Version := $(shell git describe --tags --dirty --match="v*" 2> /dev/null || echo v0.0.0-dev)
Date := $(shell date -u '+%Y-%m-%d-%H%M UTC')

ReleaseDistribution = $(foreach p,$(Platforms),$(foreach c,$(Cmds),$(ReleaseDist)/$(Version)/$(p)-amd64/$(c)))
BuildDistribution = $(foreach c,$(Cmds),$(BuildDist)/$(c))

go-clean:
	go clean ./cmd/...
.PHONY: go-clean

go-get:
	go get
	go mod download
.PHONY: go-get

clean: go-clean
	rm -rf $(ReleaseDist)/$(Version)
.PHONY: clean

$(ReleaseDistribution): platform = $(shell $(foreach p,$(Platforms),echo $@ | grep -oh $(p);))
$(ReleaseDistribution): cmd = $(notdir $@)
$(ReleaseDistribution):
	@-rm $@
	@cd cmd/$(cmd); CGO_ENABLED=0 GOOS=$(platform) GOARCH=amd64 go build -ldflags="-X 'main.Version=$(Version)' -X 'main.BuildTime=$(Date)'" -o ../../$@;
.PHONY: $(ReleaseDistribution)

$(BuildDistribution): cmd = $(notdir $@)
$(BuildDistribution):
	@-rm $@
	cd cmd/$(cmd); CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=amd64 go build -ldflags="-X 'main.Version=$(Version)' -X 'main.BuildTime=$(Date)'" -o ../../$@;
.PHONY: $(BuildDistribution)

build: $(BuildDistribution)
.PHONY: build

release: $(ReleaseDistribution)
.PHONY: release

test:
	echo $(Cmds)
.PHONY: test