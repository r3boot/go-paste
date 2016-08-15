TARGET = go-paste
UTIL = gp

VERSION = latest
ACI = ${TARGET}-${VERSION}-amd64.aci

BUILD_DIR = ./build
INSTALL_DIR = ./installed
PREFIX = /usr/local


all: ${TARGET} ${UTIL}

${TARGET}:
	[[ -d "${BUILD_DIR}" ]] || mkdir -vp ${BUILD_DIR}
	[[ -d "${INSTALL_DIR}" ]] || mkdir -vp ${INSTALL_DIR}
	mkdir -p ${BUILD_DIR}/src/github.com/r3boot/go-paste
	cp -Rp lib templates go-paste.go ${BUILD_DIR}/src/github.com/r3boot/go-paste
	install -o root -g root -m 0755 files/build.sh ${BUILD_DIR}/build.sh
	rkt-builder

${UTIL}:
	[[ -d "${BUILD_DIR}" ]] || mkdir -vp ${BUILD_DIR}
	go build -v -o ${BUILD_DIR}/gp gp.go

${ACI}:
	./scripts/build_aci.sh ${VERSION}

install:
	rkt fetch --insecure-options=image ${BUILD_DIR}/${ACI}
	install -o root -g root -m 0755 ${BUILD_DIR}/${UTIL} /usr/local/bin/${UTIL}

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rvf ${BUILD_DIR} || true
	[[ -d "${INSTALL_DIR}" ]] && rm -rvf ${INSTALL_DIR} || true
