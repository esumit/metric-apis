.PHONY: build clean run

APP_NAME=metric-apis
export CUR_DIR=$(shell pwd)

build:
	rm -rf ./dist
	mkdir dist
	CGO_ENABLED=0 go build -o dist/$(APP_NAME)
	cp .env dist/.env

run:
	./dist/$(APP_NAME)

docker:
	docker build -f Dockerfile -t esumit/$(APP_NAME) .
clean:
	rm -rf ./dist/*



