provider "aws" {
  region = "ap-south-1"
}

resource "aws_s3_bucket" "image_bucket" {
  bucket = "prashanna-s3-bucket"
}

resource "aws_iam_role" "lambda_role" {
  name = "lambda_execution_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/AWSLambdaExecute"
}

resource "aws_lambda_function" "upload_resize" {
  filename         = "upload_resize.zip"
  function_name    = "upload_resize"
  role             = aws_iam_role.lambda_role.arn
  handler          = "upload_resize"
  runtime          = "go1.x"
  source_code_hash = filebase64sha256("upload_resize.zip")
  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.image_bucket.bucket
    }
  }
}

resource "aws_lambda_function" "retrieve" {
  filename         = "retrieve.zip"
  function_name    = "retrieve"
  role             = aws_iam_role.lambda_role.arn
  handler          = "retrieve"
  runtime          = "go1.x"
  source_code_hash = filebase64sha256("retrieve.zip")
  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.image_bucket.bucket
    }
  }
}

resource "aws_api_gateway_rest_api" "image_resizer_api" {
  name        = "ImageResizerAPI"
  description = "API for Image Resizer Service"
}

resource "aws_api_gateway_resource" "upload_resource" {
  rest_api_id = aws_api_gateway_rest_api.image_resizer_api.id
  parent_id   = aws_api_gateway_rest_api.image_resizer_api.root_resource_id
  path_part   = "upload"
}

resource "aws_api_gateway_method" "upload_method" {
  rest_api_id   = aws_api_gateway_rest_api.image_resizer_api.id
  resource_id   = aws_api_gateway_resource.upload_resource.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "upload_integration" {
  rest_api_id = aws_api_gateway_rest_api.image_resizer_api.id
  resource_id = aws_api_gateway_resource.upload_resource.id
  http_method = aws_api_gateway_method.upload_method.http_method
  integration_http_method = "POST"
  type        = "AWS_PROXY"
  uri         = aws_lambda_function.upload_resize.invoke_arn
}

resource "aws_api_gateway_resource" "retrieve_resource" {
  rest_api_id = aws_api_gateway_rest_api.image_resizer_api.id
  parent_id   = aws_api_gateway_rest_api.image_resizer_api.root_resource_id
  path_part   = "{key}"
}

resource "aws_api_gateway_method" "retrieve_method" {
  rest_api_id   = aws_api_gateway_rest_api.image_resizer_api.id
  resource_id   = aws_api_gateway_resource.retrieve_resource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "retrieve_integration" {
  rest_api_id = aws_api_gateway_rest_api.image_resizer_api.id
  resource_id = aws_api_gateway_resource.retrieve_resource.id
  http_method = aws_api_gateway_method.retrieve_method.http_method
  integration_http_method = "POST"
  type        = "AWS_PROXY"
  uri         = aws_lambda_function.retrieve.invoke_arn
}

output "api_url" {
  value = aws_api_gateway_rest_api.image_resizer_api.execution_arn
}
