# Serverless Image Resizer

## Overview
Create a serverless image resizer application using Go. The application will use AWS Lambda for serverless computing, S3 for storing the images, and API Gateway to handle HTTP requests. Docker will be used for local development and testing, while Terraform will manage the infrastructure as code.

## Features
- Upload Image: Accepts an image file and uploads it to an S3 bucket.
- Resize Image: On request, resizes the image to the specified dimensions.
- Retrieve Image: Retrieves the resized image from the S3 bucket.
Steps to Create the Project
## Repository Setup
- Create a new GitHub repository named serverless-image-resizer.
- Initialize the repository with a README.md and .gitignore for Go.
## Local Development with Docker
- Set up a Dockerfile to create a development environment for the Go code.
- Create a docker-compose.yml to manage local development dependencies.
## Go Lambda Function
- Develop the main logic of the image resizer using Go.
- Create two Lambda functions: one for uploading and resizing images, and another for retrieving images.
- Use the AWS SDK to interact with S3.
## API Gateway Integration
- Set up an API Gateway to handle HTTP requests and trigger the Lambda functions.
## Infrastructure as Code with Terraform
Write Terraform scripts to provision the necessary AWS resources:
- API Gateway
- Lambda Functions
- S3 Bucket
- IAM Roles and Policies
## Deployment
- Use Terraform to deploy the infrastructure to AWS.
- Use AWS CLI or other deployment tools to upload your Lambda code.