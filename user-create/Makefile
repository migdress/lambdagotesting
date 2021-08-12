.PHONY: build
build: 
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/v1 v1/*.go
 
.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose


