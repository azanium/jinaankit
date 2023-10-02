package aws

type CloudConfig struct {
	ID       string
	Secret   string
	Address  string
	Region   string
	Profile  string
	Endpoint string
	S3       S3Config
}

type S3Config struct {
	Bucket            string
	PresignedTTLInSec int
	UploadConcurrency int
	UploadPartSize    int
}
