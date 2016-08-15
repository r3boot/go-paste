TARGET = go-paste

BUILD_DIR = ./build

all: ${TARGET}

${TARGET}:
	[[ -d "${BUILD_DIR}" ]] || mkdir -vp ${BUILD_DIR}
	go build -v -o ${BUILD_DIR}/${@} ${@}.go

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rvf ${BUILD_DIR} || true
