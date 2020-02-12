TARGETS   = go-paste gp
VERSION   = latest
NAMESPACE = as65342

BUILD_DIR = ./build
PREFIX    = /usr/local

all: $(TARGETS)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

dependencies:
	go get -v ./...

$(TARGETS): $(BUILD_DIR) dependencies
	go build -v -o ${BUILD_DIR}/$@ ./cmd/$@/main.go

install:
	install -o root -g root -m 0755 ${BUILD_DIR}/${TARGET} /usr/local/bin/${TARGET}
	install -o root -g root -m 0755 ${BUILD_DIR}/${UTIL} /usr/local/bin/${UTIL}

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rvf ${BUILD_DIR} || true
