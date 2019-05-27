CONTAINER_NAME="journal"

node {

    stage('Checkout') {
        checkout scm
    }

    stage('Build') {
        sh "docker build -t $CONTAINER_NAME ."
    }

    stage('Test') {
        sh "docker run --rm $CONTAINER_NAME go test ./..."
    }
}