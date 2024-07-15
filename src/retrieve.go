package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"os"
)

var s3Client *s3.S3

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	s3Client = s3.New(sess)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key := request.PathParameters["key"]

	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("your-s3-bucket"),
		Key:    aws.String("resized/" + key),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(buf.Bytes()),
		Headers:    map[string]string{"Content-Type": "image/jpeg"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
