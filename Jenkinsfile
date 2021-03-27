/* groovylint-disable CompileStatic */
CONTAINER_NAME = 'journal'

node {
    stage('Checkout') {
        checkout scm
    }

    stage('Build') {
        sh "docker build -t $CONTAINER_NAME -f Dockerfile.test ."
    }

    stage('Test - Latest Go') {
        sh """
        docker run --name $CONTAINER_NAME $CONTAINER_NAME make test > journal-test.xml
        docker cp $CONTAINER_NAME:/go/src/github.com/jamiefdhurst/journal/coverage.xml journal-coverage.xml
        docker stop $CONTAINER_NAME
        docker rm $CONTAINER_NAME
        """
        junit 'journal-test.xml'
        step([$class: 'CoberturaPublisher', coberturaReportFile: 'journal-coverage.xml'])
    }
}
