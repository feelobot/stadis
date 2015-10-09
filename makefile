build:
	GOOS=linux GOARCH=amd64 go build -o bin/stadis-linux-amd64 && GOOS=darwin GOARH=amd64 go build -o bin/stadis-darwin-amd64
