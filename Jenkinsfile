/* groovylint-disable CompileStatic */
CONTAINER_NAME = 'journal'

pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                scmSkip(deleteBuild: true, skipPattern: '^Skip CI.*')
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

        stage('Package and Release') {
            when {
                branch 'main'
            }
            steps {
                build job: '/github/journal-folder/release', wait: true
            }
        }

        // stage('Deploy') {
        //     when {
        //         branch 'main'
        //     }
        //     steps {
        //         library identifier: 'jenkins@main'
        //         build job: '/github/journal-folder/deploy', wait: true, parameters: [
        //             string(name: 'targetVersion', value: getVersion(repo: 'jamiefdhurst/journal').full)
        //         ]
        //     }
        // }
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
