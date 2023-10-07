.PHONY: all

OS = linux
AGENT_BUILD_NAME = agent
PATH_MAIN_AGENT_GO = ./cmd/agent/main.go
SERVER_BUILD_NAME = server
PATH_MAIN_SERVER_GO = ./cmd/server/main.go


all: clean get build-agent build-server

build-agent:
	@echo "  >  Building agent"
	@CGO_ENABLED=0 GOOS=$(OS) go build -ldflags "-w" -a -o $(AGENT_BUILD_NAME) $(PATH_MAIN_AGENT_GO)

build-server:
	@echo "  >  Building server"
	@CGO_ENABLED=0 GOOS=$(OS) go build -ldflags "-w" -a -o $(SERVER_BUILD_NAME) $(PATH_MAIN_SERVER_GO)

gotest:
	go test `go list ./... | grep -v test` -count 1
	
gotestcover:
	go test `go list ./... | grep -v test` -count 1 -cover

get:
	@echo "  >  Checking dependencies"
	@go mod download
	@go install $(PATH_MAIN_AGENT_GO)
	@go install $(PATH_MAIN_SERVER_GO)

clean:
	@echo "  >  Clearing folder"
	@rm -f ./$(AGENT_BUILD_NAME)
	@rm -f ./$(SERVER_BUILD_NAME)