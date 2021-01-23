
def build() {
    sh """
    mkdir -p /go/src/github.com/r3boot
    ln -s `pwd` /go/src/github.com/r3boot/go-paste
    cd /go/src/github.com/r3boot/go-paste && make
    """
}

podTemplate(containers: [
    containerTemplate(name: 'golang_libc', image: 'golang:latest', ttyEnabled: true, command: 'cat')
    containerTemplate(name: 'golang_musl', image: 'golang:alpine', ttyEnabled: true, command: 'cat')
  ]) {

    node(POD_LABEL) {
        stage('Build go-paste') {
            git url: 'ssh://git@gitea.as65342.net:2222/r3boot/go-paste.git'
            container('golang_libc') {
                stage('Build binary for libc-amd64') {
                    build()
                }
            }
            container('golang_musl') {
                stage('Build binary for musl-amd64') {
                    build()
                }
            }
        }

    }
}