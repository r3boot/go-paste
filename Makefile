TARGET = go-paste

BUILD_DIR = ./build
PREFIX = /usr/local

all: ${TARGET}

${TARGET}:
	[[ -d "${BUILD_DIR}" ]] || mkdir -vp ${BUILD_DIR}
	go build -v -o ${BUILD_DIR}/${@} ${@}.go

install:
	install -o root -g root -m 0755 ${BUILD_DIR}/${TARGET} \
		${PREFIX}/bin/${TARGET}
	install -o root -g root -m 0644 templates/index.html \
		/usr/share/go-paste.html

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rvf ${BUILD_DIR} || true
