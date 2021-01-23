
def build() {
    sh """
    mkdir -p /go/src/github.com/r3boot
    ln -s `pwd` /go/src/github.com/r3boot/go-paste
    cd /go/src/github.com/r3boot/go-paste && make
    """
}

podTemplate(containers: [
    containerTemplate(name: 'golang-libc', image: 'golang:latest', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'golang-musl', image: 'golang:alpine', ttyEnabled: true, command: 'cat')
  ]) {

  options {
      skipDefaultCheckout true
  }

    node(POD_LABEL) {
        stage('Preparation') {
            steps {
                sh """
                ssh-keyscan -f ~/.ssh/known_hosts -p 2222 gitea-ssh.develop.svc
                """
            }
        }
        stage('Checkout') {
            steps {
                 checkout scm
            }
        }
        stage('Build go-paste') {
            git url: 'ssh://git@gitea-ssh.develop.svc:2222/r3boot/go-paste.git'
            container('golang-libc') {
                stage('Build binary for libc-amd64') {
                    sh """
                    mkdir -p /go/src/github.com/r3boot
                    ln -s `pwd` /go/src/github.com/r3boot/go-paste
                    cd /go/src/github.com/r3boot/go-paste && make
                    """
                }
            }
            container('golang-musl') {
                stage('Build binary for musl-amd64') {
                    sh """
                    mkdir -p /go/src/github.com/r3boot
                    ln -s `pwd` /go/src/github.com/r3boot/go-paste
                    cd /go/src/github.com/r3boot/go-paste && make
                    """
                }
            }
        }

    }
}