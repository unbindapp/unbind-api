package models

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

type S3Response struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Endpoint  string         `json:"endpoint"`
	Region    string         `json:"region"`
	AccessKey string         `json:"access_key"`
	SecretKey string         `json:"secret_key"`
	Buckets   []types.Bucket `json:"buckets" nullable:"false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// TransformS3Entity transforms an ent.S3 entity into a S3Response
func TransformS3Entity(entity *ent.S3, accessKey string, secretKey string, buckets []types.Bucket) *S3Response {
	if buckets == nil {
		buckets = []types.Bucket{}
	}

	response := &S3Response{}
	if entity != nil {
		response = &S3Response{
			ID:        entity.ID,
			Name:      entity.Name,
			Endpoint:  entity.Endpoint,
			Region:    entity.Region,
			AccessKey: accessKey,
			SecretKey: secretKey,
			Buckets:   buckets,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		}
	}
	return response
}

// Transforms a slice of ent.S3 entities into a slice of S3Response
func TransformS3Entities(entities []*ent.S3, accessKeyMap map[uuid.UUID]string, secretKeyMap map[uuid.UUID]string, bucketsMap map[uuid.UUID][]types.Bucket) []*S3Response {
	responses := make([]*S3Response, len(entities))
	for i, entity := range entities {
		responses[i] = TransformS3Entity(entity, accessKeyMap[entity.ID], secretKeyMap[entity.ID], bucketsMap[entity.ID])
	}
	return responses
}
