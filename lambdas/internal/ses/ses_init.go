package ses

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func Ses_Init() *sesv2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("SES_REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return sesv2.NewFromConfig(cfg)
}
