#!/usr/bin/env groovy

def shortCommit
def env
def namespace

node {
  checkout scm
  shortCommit = sh(returnStdout: true, script: "git log -n 1 --pretty=format:'%h'").trim()
  env = params.Environment
  namespace = "artifactory-report-generator"
  stage('Coverage') {
    outFile = "coverage.out"
    pubFile = "coverage.html"
    sh """COVFILE=${outFile} PUBFILE=${pubFile} make docker_test"""
  }

stage('Build') {
    shortCommit = sh(returnStdout: true, script: "git log -n 1 --pretty=format:'%h'").trim()
    target = 'docker_build'
    artifact = 'report-generator'
    sh """BUILD_ARTIFACT=${artifact} RELEASE_VERSION=${shortCommit} make ${target}"""
 }

 stage('Release') {
    withCredentials([usernamePassword(credentialsId: 'CREDENTIALS_ID', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
        sh """TAG=${shortCommit} make publish"""
      }
  }

stage('configure and install') {
   node {
        withCredentials([ file(credentialsId: 'KUBECONFIG', variable: 'KUBECONFIG'),
                          file(credentialsId: 'report-generator', variable: 'VALUES')]) {
            sh "helm template --name=report-generator --values=${VALUES} ./charts/ | kubectl --context=${env} -n ${namespace} apply -f -"
       }
   }    
}
}
