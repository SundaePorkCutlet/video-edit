pipeline {
    agent {
        dockerfile {
            filename 'Dockerfile'
            additionalBuildArgs '--build-arg YOUR_ARG=your_value'
        }
    }
    environment {
        AWS_DEFAULT_REGION = 'ap-northeast-2'  // AWS 지역 설정
        AWS_ACCOUNT_ID = '068494925351'
        ECR_REPO_URI = '068494925351.dkr.ecr.ap-northeast-2.amazonaws.com/my-ecr-backend' // ECR 리포지토리 URI
        IMAGE_TAG = "latest"
        SSH_KEY_ID = 'your-ssh-key-id' // SSH 키 자격 증명 ID
        SERVICE_SERVER = 'your-service-server' // 서비스 서버 주소
        GITHUB_CREDENTIALS_ID = 'github-credentials'
        AWS_CREDENTIALS_ID = 'aws-credentials-id'
    }
    stages {
        stage('Checkout') {
            steps {
                script {
                    checkout([$class: 'GitSCM',
                              branches: [[name: '*/main']],
                              doGenerateSubmoduleConfigurations: false,
                              extensions: [],
                              userRemoteConfigs: [[
                                  url: 'https://github.com/SundaePorkCutlet/video-edit.git',
                                  credentialsId: "${GITHUB_CREDENTIALS_ID}"
                              ]]
                    ])
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Docker 이미지 빌드
                    sh 'docker build -t my-app .'
                }
            }
        }

        stage('Login to AWS ECR') {
            steps {
                script {
                    withAWS(credentials: "${AWS_CREDENTIALS_ID}", region: "${AWS_DEFAULT_REGION}") {
                        sh 'aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com'
                    }
                }
            }
        }

        stage('Push Docker Image to ECR') {
            steps {
                script {
                    sh "docker tag my-app:latest ${ECR_REPO_URI}:${IMAGE_TAG}"
                    sh "docker push ${ECR_REPO_URI}:${IMAGE_TAG}"
                }
            }
        }

        stage('Deploy to Service Server') {
            steps {
                script {
                    sshagent(credentials: ["${SSH_KEY_ID}"]) {
                        sh """
                        ssh -o StrictHostKeyChecking=no ${SERVICE_SERVER} '
                        aws ecr get-login-password --region ${AWS_DEFAULT_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com &&
                        cd /path/to/deployment/directory &&
                        git pull origin main &&
                        docker-compose down &&
                        docker-compose pull app &&
                        docker-compose up -d
                        '
                        """
                    }
                }
            }
        }
    }
}
