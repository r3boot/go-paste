TARGETS   = go-paste gp
VERSION   = latest
NAMESPACE = as65342

BUILD_DIR = ./build
PREFIX    = /usr/local

OS_ID = $(shell awk -F= '/^ID/{ print $2 }' /etc/os-release)

all: $(TARGETS)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

dependencies:
	go get -v ./...

$(TARGETS): $(BUILD_DIR) dependencies
	if [[ "$(OS_ID)" == "debian" ]]; then
		go build -v -o $(BUILD_DIR)/$@-libc-amd64 ./cmd/$@/main.go
	elif [[ "$(OS_ID)" == "alpine" ]]; then
		go build -v -o $(BUILD_DIR)/$@-musl-amd64 ./cmd/$@/main.go
	else
	  	go build -v -o $(BUILD_DIR)/$@ ./cmd/$@/main.go
	fi

install:
	install -o root -g root -m 0755 $(BUILD_DIR)/${TARGET} /usr/local/bin/${TARGET}
	install -o root -g root -m 0755 $(BUILD_DIR)/${UTIL} /usr/local/bin/${UTIL}

clean:
	[[ -d "$(BUILD_DIR)" ]] && rm -rvf $(BUILD_DIR) || true
