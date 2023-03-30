# install packages
install:
	go mod download

# build bin file prod
build:
	go mod download && go build -o ./.bin/app ./src/cmd/app/main.go