package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
)

type MClient struct {
	Client     *minio.Client
	BucketName string
}

func NewMClient(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MClient{
		Client:     client,
		BucketName: bucketName,
	}, nil
}

func (mc *MClient) GetImage(objectName string) ([]byte, error) {
	object, err := mc.Client.GetObject(context.Background(), mc.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("err1: %w", err)
	}

	imageBytes, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("err2: %w", err)
	}

	return imageBytes, nil
}

func (mc *MClient) UploadImage(objectName string, file *multipart.FileHeader) error {
	src, err := file.Open()
	_, err = mc.Client.PutObject(
		context.Background(),
		mc.BucketName,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	return err
}
