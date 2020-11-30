package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/google/uuid"
	"log"
)

type Kafka struct {
	bootstrapServers []string
	config           *sarama.Config
	awsConfig        aws.Config
}

type ConsumerGroup struct {
	Id              string
	ActiveMembers   int
	Members         []Member
	LastKnownStatus string
	FunctionArn     *string
}

type Member struct {
	ClientId string
	Topics   []string
}

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
		config:           scfg,
		awsConfig:        cfg,
	}, nil
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

func (k *Kafka) ListConsumerGroups() (*[]ConsumerGroup, error) {
	c, err := sarama.NewClient(k.bootstrapServers, k.config)
	if err != nil {
		log.Println("Unable to create a kafka client", err)
		return nil, err
	}
	admin, err := sarama.NewClusterAdminFromClient(c)

	if err != nil {
		log.Println("Unable to create an admin kafka client", err)
		return nil, err
	}
	defer admin.Close()

	consumerGroups, err := admin.ListConsumerGroups()
	if err != nil {
		log.Println("could not get a list of consumer groups:", err)
		return nil, err
	}

	svc := lambda.NewFromConfig(k.awsConfig)

	groupNames := make([]string, len(consumerGroups))
	for groupName := range consumerGroups {
		groupNames = append(groupNames, groupName)
	}

	consumerGroupDescriptions, err := admin.DescribeConsumerGroups(groupNames)

	if err != nil {
		fmt.Println("could not describe consumer groups: ", err)
		return nil, err
	}

	var cgroups []ConsumerGroup
	for _, groupDescription := range consumerGroupDescriptions {
		members := groupDescription.Members
		cgroup := ConsumerGroup{
			Id:            groupDescription.GroupId,
			ActiveMembers: len(members),
		}

		for _, groupMemberDescription := range members {
			memberMetadata, err := groupMemberDescription.GetMemberMetadata()
			if err != nil {
				cgroup.LastKnownStatus = "Error getting member metadata"
			} else {
				cgroup.LastKnownStatus = "OK"
				topics := memberMetadata.Topics
				member := Member{
					ClientId: groupMemberDescription.ClientId,
					Topics:   topics,
				}
				cgroup.Members = append(cgroup.Members, member)
			}
		}

		// If UUID, then check which lambda function this belongs to.
		if ok := isValidUUID(groupDescription.GroupId); ok {
			input := &lambda.GetEventSourceMappingInput{
				UUID: aws.String(groupDescription.GroupId),
			}
			req, err := svc.GetEventSourceMapping(context.Background(), input)
			if err != nil {
				if errors.Is(err, &types.ResourceNotFoundException{}) != false {
					log.Println("unable to get event source mapping", err)
				}
			} else {
				cgroup.FunctionArn = req.FunctionArn
			}
		}
		cgroups = append(cgroups, cgroup)
	}

	return &cgroups, nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
