package s3

import (
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3 struct {
	client *minio.Client
	once   sync.Once
}

var instance *s3

func getClient() (*s3, error) {
	if instance != nil {
		return instance, nil
	}

	instance = &s3{}
	var err error
	instance.once.Do(func() {
		instance.client, err = minio.New("s3.amazonaws.com", &minio.Options{
			Creds: credentials.NewChainCredentials([]credentials.Provider{
				&credentials.EnvAWS{},
				&credentials.FileAWSCredentials{},
				&credentials.IAM{Client: nil},
			}),
			Region: "eu-west-2",
			Secure: true,
		})
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
