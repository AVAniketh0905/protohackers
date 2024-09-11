# Protohackers Problems

## Deploy to AWS ECS

- [Resources](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/create-container-image.html)

### IMPORTANT

**DELETE** After the test is done

`aws ecr delete-repository --repository-name protohackers --region us-east-2 --force`

## 1. Smoke Test

[Code](./questions/smoke_test)

### Problem

- TCP Echo Server
- At max 5 connections
- Echo back the message
