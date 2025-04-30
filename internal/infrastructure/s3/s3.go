package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// S3Client provides methods to interact with S3-compatible storage.
type S3Client struct {
	client *s3.Client
}

// NewS3Client creates a new S3Client with the provided credentials.
func NewS3Client(ctx context.Context, endpoint, region, accessKeyID, secretKey string) (*S3Client, error) {
	// Validate/parse the custom endpoint once up-front.
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid endpoint URL")
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, ""),
		),
		// From R2 docs
		config.WithRequestChecksumCalculation(0),
		config.WithResponseChecksumValidation(0),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &S3Client{client: client}, nil
}

// ListBuckets lists all S3 buckets in the configured account.
func (c *S3Client) ListBuckets(ctx context.Context) ([]*models.S3Bucket, error) {
	buckets, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}
	return models.TransformBucketEntities(buckets.Buckets), nil
}

// ProbeAnyBucketRW tries to write + read a probe object in the first bucket the
// credentials can see.  If every step succeeds it returns nil.  All SDK errors
// are normalised via mapS3Error so callers get the domain-specific errors they
// expect.
func (c *S3Client) ProbeAnyBucketRW(ctx context.Context) error {
	lbOut, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return mapS3Error(err)
	}
	if len(lbOut.Buckets) == 0 {
		return errdefs.NewCustomError(errdefs.ErrTypeNotFound,
			"credentials are valid but no buckets are visible")
	}

	for _, b := range lbOut.Buckets {
		if err := c.ProbeBucketRW(ctx, aws.ToString(b.Name)); err == nil {
			return nil // success on this bucket
		}
	}
	return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
		"credentials can't read & write any visible bucket")
}

// probeBucketRW writes → heads → deletes a tiny object in one bucket.
func (c *S3Client) ProbeBucketRW(ctx context.Context, bucket string) error {
	key := fmt.Sprintf(".probe-%s", uuid.NewString())

	// Put (write)
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket, Key: &key,
		Body: bytes.NewReader([]byte("probe")),
	})
	if err != nil {
		log.Error(err)
		return mapS3Error(err)
	}

	// Head (read)
	_, err = c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket, Key: &key,
	})
	if err != nil {
		_ = c.deleteSilent(ctx, bucket, key)
		return mapS3Error(err)
	}

	// Delete (cleanup) – ignore classification result
	_ = c.deleteSilent(ctx, bucket, key)
	return nil
}

// deleteSilent is best-effort cleanup.
func (c *S3Client) deleteSilent(ctx context.Context, bucket, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket, Key: &key,
	})
	return err
}

// Translate S3 errors to our error types
func mapS3Error(err error) error {
	var respErr *smithyhttp.ResponseError
	if errors.As(err, &respErr) {
		switch respErr.Response.StatusCode {
		case 403:
			// Forbidden – caller treats this like an invalid endpoint for R2/MinIO
			return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"invalid endpoint URL or access forbidden (HTTP 403)")
		case 404:
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound,
				"bucket or object not found (HTTP 404)")
		case 301, 307:
			return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"wrong region for bucket (redirect)")
		}
	}

	// Network/URL parse issues
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
			fmt.Sprintf("invalid endpoint URL: %v", urlErr))
	}

	// Fallback – surface the original error.
	return err
}
