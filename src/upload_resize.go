package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	s3Client = s3.New(sess)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body struct {
		ImageData string `json:"image_data"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	}
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	imgData, err := os.ReadFile(body.ImageData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	resizedImg := resize(img, body.Width, body.Height)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("your-s3-bucket"),
		Key:    aws.String("resized/" + body.ImageData),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("Image resized and uploaded: %s", body.ImageData),
	}, nil
}

func resize(img image.Image, width, height int) image.Image {
	rect := img.Bounds()
	dx := float64(rect.Dx())
	dy := float64(rect.Dy())
	ratio := dx / dy

	var newWidth, newHeight int
	if dx > dy {
		newWidth = width
		newHeight = int(float64(width) / ratio)
	} else {
		newHeight = height
		newWidth = int(float64(height) * ratio)
	}

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			newImg.Set(x, y, img.At(int(float64(x)/float64(newWidth)*dx), int(float64(y)/float64(newHeight)*dy)))
		}
	}

	return newImg
}

func main() {
	lambda.Start(handler)
}
