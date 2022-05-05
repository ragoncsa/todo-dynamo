# todo

## Overview

This is a fork from https://github.com/ragoncsa/todo/tree/v0.1.1
See its README https://github.com/ragoncsa/todo/blob/v0.1.1/README.md

This is an example how to port the "todo" demo project from gorm/postgres to Dynamo.

**What's included:**

* removal of gorm dependency
* implementation of the existing methods (Task, Tasks, DeleteTask, DeleteTasks) as is

**What's not included:**

* this only tries out Dynamo API, the implementation is inefficient (for example task ID as primary key is a bad choice, so is full table scan, or deleting all items one by one)

## Create the DynamoDB table "Task" in AWS

```bash
source aws-creds-to-env.sh
cd terraform
terraform init
terraform apply
```

For more see:

* [How to use AWS from Terraform](https://learn.hashicorp.com/tutorials/terraform/aws-build)
* [Resource: aws_dynamodb_table](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table)

## To start and test the service

```bash
go run main.go
```

http://localhost:8080/swagger/index.html

## Clean up

```bash
cd terraform
terraform destroy
unset AWS_ACCESS_KEY_ID && unset AWS_SECRET_ACCESS_KEY && unset AWS_SESSION_TOKEN
```