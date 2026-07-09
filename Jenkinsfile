pipeline {
    agent none

    environment {
        IMAGE    = "ghcr.io/yoshikrit/inventory"
        GOPATH   = "/tmp/go"
        GOCACHE  = "/tmp/go-cache"
    }

    stages {
        stage('Lint') {
            agent { docker { image 'golang:1.26-alpine' } }
            steps {
                sh 'go vet ./...'
                sh 'go install honnef.co/go/tools/cmd/staticcheck@latest'
                sh '$GOPATH/bin/staticcheck ./...'
            }
        }

        stage('Test') {
            agent { docker { image 'golang:1.26-alpine' } }
            steps {
                sh 'go mod verify'
                sh 'go test ./... -cover'
            }
        }

        stage('Docker Build & Push') {
            agent any
            when { branch 'main' }
            steps {
                script {
                    docker.withRegistry('https://ghcr.io', 'ghcr-credentials') {
                        def img = docker.build("${IMAGE}:${env.GIT_COMMIT.take(7)}")
                        img.push()
                        img.push('latest')
                    }
                }
            }
        }
    }

    post {
        success { echo '✅ Pipeline passed' }
        failure  { echo '❌ Pipeline failed' }
    }
}
