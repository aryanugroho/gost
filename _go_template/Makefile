APP_EXECUTABLE="bin/"
APP_BOOTSTRAP="./cmd//"

compile:
	mkdir -p bin/
	go build -o $(APP_EXECUTABLE) $(APP_BOOTSTRAP)

test:
	go test -short -count=1 -race ./...

static-check:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

copy-config:
	cp ./config.yml.tmpl ./config.yml