build:
	dep ensure -v
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/deviceAdd deviceAdd/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/deviceView deviceView/main.go
	build-lambda-zip -o bin/deviceAdd.zip bin/deviceAdd
	build-lambda-zip -o bin/deviceView.zip bin/deviceView

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
