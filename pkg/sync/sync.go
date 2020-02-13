package sync

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-logr/logr"
)

// Sync ...
type Sync struct {
	Log    logr.Logger
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
	Region string `json:"region"`
}

// AwsSync ...
func (f *Sync) AwsSync() error {
	log := f.Log.WithName("AwsSync")
	if err := f.isValid(); err != nil {
		return err
	}

	di, err := NewDirectoryIterator(f.Bucket, f.Path)
	if err != nil {
		return err
	}

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(f.Region)},
	)
	if err != nil {
		return fmt.Errorf("failed to create a session %v", err)
	}

	uploader := s3manager.NewUploader(sess)
	if err := uploader.UploadWithIterator(aws.BackgroundContext(), di); err != nil {
		return fmt.Errorf("failed to upload %v", err)
	}
	log.Info("Sync done")
	return nil
}

// List
func (f *Sync) List() ([]string, error) {
	di, err := NewDirectoryIterator(f.Bucket, f.Path)
	if err != nil {
		return nil, err
	}
	return di.Path(), nil
}

func (f *Sync) isValid() error {
	if len(f.Path) <= 0 {
		return fmt.Errorf("Invalid Path")
	}

	if len(f.Bucket) <= 0 {
		return fmt.Errorf("Invalid Bucket")
	}

	if len(f.Region) <= 0 {
		return fmt.Errorf("Invalid Region")
	}
	return nil
}
