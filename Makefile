build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/raspadmin_amd64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o build/raspadmin_arm64 main.go

clean:
	rm -rf build/raspadmin_*

.PHONY: all build clean