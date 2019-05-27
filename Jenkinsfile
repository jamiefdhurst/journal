CONTAINER_NAME="journal"

node {

    stage('Checkout') {
        checkout scm
    }

    stage('Build') {
        sh "docker build -t $CONTAINER_NAME ."
    }

    stage('Test') {
        sh "docker run --name $CONTAINER_NAME -it $CONTAINER_NAME go test ./..."
    }

    stage('Cleanup') {
        sh "docker rm -f $CONTAINER_NAME"
    }
}