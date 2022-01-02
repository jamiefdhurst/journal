/* groovylint-disable CompileStatic */
CONTAINER_NAME = 'journal'

pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh "docker build -t $CONTAINER_NAME -f Dockerfile.test ."
            }
        }

        stage('Test') {
            steps {
                sh """
                docker run --name $CONTAINER_NAME $CONTAINER_NAME make test > journal-tests.xml
                docker cp $CONTAINER_NAME:/go/src/github.com/jamiefdhurst/journal/coverage.xml journal-coverage.xml
                """
                junit 'journal-tests.xml'
                step([$class: 'CoberturaPublisher', coberturaReportFile: 'journal-coverage.xml'])
            }
        }
    }

    post {
        always {
            sh """
            docker stop $CONTAINER_NAME
            docker rm $CONTAINER_NAME
            """
        }
    }
}
