# Continuous Delivery using Docker and Ansible Source Code

This folder contains all of the source code and repositories for this course including:

- API - the sample application repository, including the continuous delivery workflow
- base - Docker base image of the todobackend development and release images
- specs - Node.js test runner that runs acceptance tests against the todobackend sample application
- client - fork of the Todobackend Client application (see https://github.com/TodoBackend/todo-backend-client)
- deploy - Ansible deployment playbook and AWS CloudFormation Stack
- docker-ansible - Ansible playbook runner 
- docker-jenkins - Jenkins Continuous Delivery system Docker image and AWS CloudFormation Stack


## System Requirements

- Mac OS X or Linux (if using Windows, a Linux VM is required)
- Docker 1.10+
- Docker Compose 1.6+
- Docker Machine 0.6+
- VMWare Workstation/Fusion or VirtualBox
- Git
- AWS CLI

## Checking out the Source Code

Each repository includes full git history and is tagged with the following convention:

- module-n-before - source code at the beginning of Module n
- module-n-after - source code at the end of Module n

### Examples

This checks out the API_MAIN repository as at the beginning of Module 3:

```
$ cd API_MAIN
$ git checkout module-3-before
```

This checks out the base repository as at the end of Module 5:

```
$ cd base
$ git checkout module-5-after
```

This checks out the most recent commit on the master branch:

```
$ cd base
$ git checkout master
```