package entity

type SendRequest struct {
	QueueURL   string
	Body       string
	Attributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
	Type  string
}

type Message struct {
	ID            string
	ReceiptHandle string
	Body          string
	Attributes    map[string]string
}

type S3Bucket struct {
	Name string `json:"name"`
}
type S3Object struct {
	Key  string `json:"key"`
	ETag string `json:"eTag"`
}
type S3 struct {
	Bucket S3Bucket `json:"bucket"`
	Object S3Object `json:"object"`
}
type S3Record struct {
	S3 S3 `json:"s3"`
}
type S3Records struct {
	Records []S3Record `json:"records"`
}
