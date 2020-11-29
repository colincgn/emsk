package pkg

import (
	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
)

func NewKafka(awsRegion string, bootstrapServers []string, tlsEnabled bool) (*Kafka, error) {

	scfg := sarama.NewConfig()
	scfg.Net.TLS.Enable = tlsEnabled
	scfg.Consumer.Return.Errors = true

	cfg, err := config.LoadDefaultConfig(config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &Kafka{
		bootstrapServers: bootstrapServers,
		config: scfg,
		awsConfig: cfg,
	}, nil
}

type Kafka struct {
	bootstrapServers []string
	config *sarama.Config
	awsConfig aws.Config
}

func (k *Kafka) ListTopics() ([]string, error) {
	c, err := sarama.NewClient(k.bootstrapServers, k.config)
	if err != nil {
		log.Println("Unable to create a kafka client", err)
		return nil, err
	}
	defer c.Close()
	return c.Topics()
}
