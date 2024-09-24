package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go/aws"
)

// Start the container and return ip
func CreateECSContainer(userId string, projectId string) (string, error) {
	subnetsIDsFile := os.Getenv("SUBNET_IDS_FILE")
	subnetsIDsData, err := os.ReadFile(subnetsIDsFile)
	if err != nil {
		fmt.Printf("failed to read subnets ID files: %v\n", err)
		return "", err
	}
	var subnetsIDs []string
	err = json.Unmarshal(subnetsIDsData, subnetsIDs)
	if err != nil {
		fmt.Printf("Failed to parse subnetes IDs: %v\n", err)
		return "", nil
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}
	client := ecs.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)
	input := &ecs.RunTaskInput{
		Cluster:        aws.String("socket-servers"),
		TaskDefinition: aws.String("socket-server-task"),
		Count:          aws.Int32(1),
		LaunchType:     types.LaunchTypeFargate,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        subnetsIDs,
				AssignPublicIp: types.AssignPublicIpEnabled,
			},
		},
	}

	result, err := client.RunTask(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to run task: %v\n", err)
		return "", err
	}

	if len(result.Tasks) == 0 {
		fmt.Println("No tasks started")
		return "", nil
	}

	taskArn := *result.Tasks[0].TaskArn
	for {
		describeTasksInput := &ecs.DescribeTasksInput{
			Cluster: aws.String("socket-servers"),
			Tasks:   []string{taskArn},
		}
		describeTasksOutput, err := client.DescribeTasks(context.TODO(), describeTasksInput)
		if err != nil {
			fmt.Printf("Failed to describe task: %v\n", err)
			return "", err
		}

		if len(describeTasksOutput.Tasks) == 0 {
			fmt.Println("Task not found")
			return "", err
		}

		task := describeTasksOutput.Tasks[0]
		if task.LastStatus != nil && *task.LastStatus == "RUNNING" {
			if len(task.Attachments) > 0 && len(task.Attachments[0].Details) > 0 {
				var eniID string
				for _, detail := range task.Attachments[0].Details {
					if *detail.Name == "networkInterfaceId" {
						eniID = *detail.Value
						break
					}
				}
				if eniID != "" {
					describeNetworkInterfacesInput := &ec2.DescribeNetworkInterfacesInput{
						NetworkInterfaceIds: []string{eniID},
					}
					describeNetworkInterfaceOutput, err := ec2Client.DescribeNetworkInterfaces(context.TODO(), describeNetworkInterfacesInput)
					if err != nil {
						fmt.Printf("Failed to describe network interface: %v\n", err)
						return "", err
					}

					if len(describeNetworkInterfaceOutput.NetworkInterfaces) > 0 {
						publicIP := describeNetworkInterfaceOutput.NetworkInterfaces[0].Association.PublicIp
						return *publicIP, nil
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

// Delete the container
