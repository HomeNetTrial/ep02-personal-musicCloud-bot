test:
	go test ./... -v

build: get
	go build -ldflags "-X 'musicCloud-bot/config.commit=`git rev-parse --short HEAD`' -X 'musicCloud-bot/config.date=`date`'" -o musiccloud-bot cmd/main.go

get:
	go mod download

run:
	go run .