
def build() {
    sh """
    mkdir -p /go/src/github.com/r3boot
    ln -s `pwd` /go/src/github.com/r3boot/go-paste
    cd /go/src/github.com/r3boot/go-paste && make
    """
}

pipeline {
    agent {
        kubernetes {
            yaml """
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: golang-libc
                image: golang:latest
                imagePullPolicy: Always
                command:
                - cat
                tty: true
              - name: golang-musl
                image: golang:alpine
                imagePullPolicy: Always
                command:
                - cat
                tty: true
            """
        }
    }

    stages {
        stage('Build') {
            steps {
                container('golang-libc') {
                    sh """
                    mkdir -p /go/src/github.com/r3boot
                    ln -s `pwd` /go/src/github.com/r3boot/go-paste
                    cd /go/src/github.com/r3boot/go-paste && make
                    """
                }
                container('golang-musl') {
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
