node ('master') {
    def root = tool name: 'Go 1.14', type: 'go'
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        stage 'Clean up'
        cleanWs()
        stage 'Checkout'
        checkout scm

        stage 'Run Tests'
        sh 'make install'
        sh 'make test'

        stage "Release build"
        sh 'make release'

        stage 'Clean up'
        cleanWs()
    }
}
