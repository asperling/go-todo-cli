APP_NAME=todo
OUT_DIR=bin
SRC=.
INSTALL_DIR=$(HOME)/bin

.PHONY: all build test cover clean install \
        build-linux build-mac build-arm build-windows

## 🔧 Standard: Local build
all: build

## 🛠 Local build
build:
	@go build -o $(OUT_DIR)/$(APP_NAME) $(SRC)
	@echo "✅ Built $(OUT_DIR)/$(APP_NAME)"

## ✅ Tests including coverage
test:
	go test -v -cover ./...

## 🚨 Run lint
lint:
	golangci-lint run --timeout 5m

## 📈 Coverage HTML report
cover:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

## 🧽 Remove artifacts
clean:
	rm -rf $(OUT_DIR)/*
	rm -f cover.out

## 💻 Cross builds
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUT_DIR)/$(APP_NAME)-linux $(SRC)

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o $(OUT_DIR)/$(APP_NAME)-mac $(SRC)

build-arm:
	GOOS=linux GOARCH=arm64 go build -o $(OUT_DIR)/$(APP_NAME)-arm $(SRC)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUT_DIR)/$(APP_NAME).exe $(SRC)

## 📦 Install globally
install: build
	@mkdir -p $(INSTALL_DIR)
	@install -m 755 $(OUT_DIR)/$(APP_NAME) $(INSTALL_DIR)/$(APP_NAME)
	@echo "✅ Installed to $(INSTALL_DIR)/$(APP_NAME)"
