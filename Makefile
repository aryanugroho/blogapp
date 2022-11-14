APP_EXECUTABLE="bin/fraudapp"

compile:
	mkdir -p bin/
	go build -o $(APP_EXECUTABLE)

test:
	go test -short -count=1 -race ./...

static-check:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

copy-config:
	cp ./application.yml.sample ./application.yml