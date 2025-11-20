.PHONY: build run clean

build:
	go mod tidy
	go build -o build/trendradar cmd/trendradar/main.go

run: build
	./trendradar

clean:
	rm -f trendradar

