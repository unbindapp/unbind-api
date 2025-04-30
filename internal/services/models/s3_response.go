package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

type S3Response struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Endpoint  string    `json:"endpoint"`
	Region    string    `json:"region"`
	AccessKey string    `json:"access_key"`
	SecretKey string    `json:"secret_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransformS3Entity transforms an ent.S3 entity into a S3Response
func TransformS3Entity(entity *ent.S3, accessKey string, secretKey string) *S3Response {
	response := &S3Response{}
	if entity != nil {
		response = &S3Response{
			ID:        entity.ID,
			Name:      entity.Name,
			Endpoint:  entity.Endpoint,
			Region:    entity.Region,
			AccessKey: accessKey,
			SecretKey: secretKey,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		}
	}
	return response
}

// Transforms a slice of ent.S3 entities into a slice of S3Response
func TransformS3Entities(entities []*ent.S3, accessKeyMap map[uuid.UUID]string, secretKeyMap map[uuid.UUID]string) []*S3Response {
	responses := make([]*S3Response, len(entities))
	for i, entity := range entities {
		responses[i] = TransformS3Entity(entity, accessKeyMap[entity.ID], secretKeyMap[entity.ID])
	}
	return responses
}
