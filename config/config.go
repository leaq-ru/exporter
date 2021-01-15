package config

type Config struct {
	ServiceName string
	Grpc        grpc
	MongoDB     mongodb
	S3          s3
	STAN        stan
	NATS        nats
	Service     service
	LogLevel    string `envconfig:"LOGLEVEL"`
}

type grpc struct {
	Port string `envconfig:"GRPC_PORT"`
}

type mongodb struct {
	URL string `envconfig:"MONGODB_URL"`
}

type s3 struct {
	ExporterBucketName string `envconfig:"S3_EXPORTERBUCKETNAME"`
	Endpoint           string `envconfig:"S3_ENDPOINT"`
	AccessKeyID        string `envconfig:"S3_ACCESSKEYID"`
	SecretAccessKey    string `envconfig:"S3_SECRETACCESSKEY"`
	Secure             string `envconfig:"S3_SECURE"`
	Region             string `envconfig:"S3_REGION"`
}

type service struct {
	Parser   string `envconfig:"SERVICE_PARSER"`
	City     string `envconfig:"SERVICE_CITY"`
	Category string `envconfig:"SERVICE_CATEGORY"`
}

type stan struct {
	ClusterID string `envconfig:"STAN_CLUSTERID"`
}

type nats struct {
	URL string `envconfig:"NATS_URL"`
}
