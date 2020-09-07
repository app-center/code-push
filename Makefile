Project := Code Push
BuildDist = build
ReleaseDist = release

Platforms = linux darwin windows
CmdTypes = svr cli
GOOS = $(shell go env GOOS)

Svrs := $(foreach n,$(shell go list ./cmd/svr/*),$(notdir $(n)))
#Clis := $(foreach n,$(shell go list ./cmd/cli/*),$(notdir $(n)))

Version := $(shell git describe --tags --dirty --match="v*" 2> /dev/null || echo v0.0.0-dev)
Date := $(shell date -u '+%Y-%m-%d-%H%M UTC')

SvrReleaseDistribution = $(foreach p,$(Platforms),$(foreach c,$(Svrs),$(ReleaseDist)/svr/$(Version)/$(p)-amd64/$(c)))
SvrBuildDistribution = $(foreach c,$(Svrs),$(BuildDist)/svr/$(c))

CliReleaseDistribution = $(foreach p,$(Platforms),$(foreach c,$(Clis),$(ReleaseDist)/cli/$(Version)/$(p)-amd64/$(c)))
CliBuildDistribution = $(foreach c,$(Clis),$(BuildDist)/cli/$(c))

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

$(SvrReleaseDistribution):
	$(call go_build,$@)
.PHONY: $(SvrReleaseDistribution)

$(SvrBuildDistribution):
	$(call go_build,$@)
.PHONY: $(SvrBuildDistribution)

$(CliReleaseDistribution):
	$(call go_build,$@)
.PHONY: $(CliReleaseDistribution)

$(CliBuildDistribution):
	$(call go_build,$@)
.PHONY: $(CliBuildDistribution)

#build: $(CliBuildDistribution)
build: $(SvrBuildDistribution)
.PHONY: build

#release: $(CliReleaseDistribution)
release: $(SvrReleaseDistribution)
.PHONY: release

define go_build
	@-rm $1
	$(eval buildPlatform = $(shell $(foreach p,$(Platforms),echo $1 | grep -owh $(p);)))
	$(eval buildCmdType = $(shell $(foreach p,$(CmdTypes),echo $1 | grep -owh $(p);)))
	$(eval buildCmd := $(notdir $1))
	$(eval buildGOOS := $(if $(buildPlatform),$(buildPlatform),$(GOOS)))
	CGO_ENABLED=0 GOOS=$(buildGOOS) GOARCH=amd64 \
		go build \
			-ldflags="-X 'github.com/funnyecho/code-push/pkg/svrkit.BuildPlatform=$(buildGOOS)-amd64' -X 'github.com/funnyecho/code-push/pkg/svrkit.Version=$(Version)' -X 'github.com/funnyecho/code-push/pkg/svrkit.BuildTime=$(Date)'" \
			-o ./$1 \
			./cmd/$(buildCmdType)/$(buildCmd);
endef