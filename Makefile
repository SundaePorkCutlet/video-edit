app_name := app

all: clean build

build:
	$(info # Build binary file. Output path is bin)
	go build -o ./bin/$(app_name) cmd/$(app_name)/main.go
	
run:
	go run ./cmd/$(app_name)/main.go

test:
	$(info Running test!)
	go test -v -cover ./internal/testing

clean:
	rm -f ./bin/$(app_name)


