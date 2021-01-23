
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
