.PHONY: default auth build push clean
.DEFAULT_GOAL := default

FUNC := aws-cloudtrail-watcher
NAME := aws-cop

default: test

test:
	go test -v ./...

config:
	@aws lambda update-function-configuration \
	--function-name ${FUNC} \
	--environment "Variables={`cat config.txt | tr '\n' ','`}"

deploy:
	- rm -rf build/linux
	GOOS=linux GOARCH=amd64 go build -o build/linux/${NAME} .
	cd build/linux; zip ${NAME}.zip ${NAME}
	aws lambda update-function-code \
	--function-name ${FUNC} \
	--zip-file "fileb://build/linux/${NAME}.zip"