package pkg

import (
	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
)

type Client interface {
	ListTopics() ([]string, error)
}

func NewKlient(awsRegion string, bootstrapServers []string, tlsEnabled bool) (*Klient, error) {

	scfg := sarama.NewConfig()
	scfg.Net.TLS.Enable = tlsEnabled
	scfg.Consumer.Return.Errors = true

	cfg, err := config.LoadDefaultConfig(config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &Klient{
		bootstrapServers: bootstrapServers,
		config: scfg,
		awsConfig: cfg,
	}, nil
}

type Klient struct {
	bootstrapServers []string
	config *sarama.Config
	awsConfig aws.Config
}

func (k Klient) ListTopics() ([]string, error) {
	c, err := sarama.NewClient(k.bootstrapServers, k.config)
	if err != nil {
		log.Println("Unable to create a kafka client", err)
		return nil, err
	}
	defer c.Close()
	return c.Topics()
}
