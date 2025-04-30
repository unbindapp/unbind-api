package models

type S3BackendCreateInput struct {
	Endpoint string `json:"endpoint" required:"true"`
}
