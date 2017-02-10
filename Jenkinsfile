node("docker-image") {
  stage("Checkout"){
    properties([
      parameters([string(defaultValue: 'hub.docker.com', description: '', name: 'REGISTRY')]),
      [$class: 'RebuildSettings', autoRebuild: false, rebuildDisabled: false],
      pipelineTriggers([])
    ])

    step([$class: 'StashNotifier'])
    checkout scm
  }
  stage("Build") {
        sh '''
          TAG=`date "+%Y-%m-%d"`.$BUILD_NUMBER
          IMAGE_NAME=journald2graylog
          DOCKERFILE=Dockerfile

          docker build -t $REGISTRY/$IMAGE_NAME:$TAG -f $DOCKERFILE .
          # Because packer is using docker 1.9 we can't use multiple -t
          docker build -t $REGISTRY/$IMAGE_NAME:latest -f $DOCKERFILE .
          docker push $REGISTRY/$IMAGE_NAME:$TAG
          docker push $REGISTRY/$IMAGE_NAME:latest
        '''
  }
  stage('Push') {
    currentBuild.result = 'SUCCESS'
    step([$class: 'StashNotifier'])
  }
}
