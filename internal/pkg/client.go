package pkg

import "github.com/aiven/aiven-go-client"

type Client interface {
	GetAccountsList() []aiven.Account
	GetAccountTeamsList(accountId string) *[]aiven.AccountTeam
	GetAccountTeamMembersList(accountId string, teamId string) *[]aiven.AccountTeamMember
	GetKafkaTopicsList(projectName string, serviceName string) []*aiven.KafkaListTopic
	GetProjectsList() []*aiven.Project
	GetServicesList(projectName string) []*aiven.Service
	GetVpcsList(projectName string) []*aiven.VPC
	GetVpcPeeringConnectionsList(projectName string, vpcId string) []*aiven.VPCPeeringConnection
}

type AivenClient struct {
	client *aiven.Client
}

func (c AivenClient) GetAccountTeamMembersList(accountId string, teamId string) *[]aiven.AccountTeamMember {
	list, err := c.client.AccountTeamMembers.List(accountId, teamId)
	handle(err)
	return &list.Members
}

func (c AivenClient) GetKafkaTopicsList(projectName string, serviceName string) []*aiven.KafkaListTopic {
	topics, err := c.client.KafkaTopics.List(projectName, serviceName)
	if aivenError, ok := err.(aiven.Error); ok {
		if aivenError.Status == 501 {
			// Aiven returns 501 if a Kafka Cluster is powered off
			return nil
		}
		handle(err)
	}
	return topics
}

func (c AivenClient) GetServicesList(projectName string) []*aiven.Service {
	services, err := c.client.Services.List(projectName)
	handle(err)
	return services
}

func (c AivenClient) GetProjectsList() []*aiven.Project {
	projects, err := c.client.Projects.List()
	handle(err)
	return projects
}

func (c AivenClient) GetAccountTeamsList(accountId string) *[]aiven.AccountTeam {
	response, err := c.client.AccountTeams.List(accountId)
	handle(err)
	return &response.Teams
}

func (c AivenClient) GetVpcsList(projectName string) []*aiven.VPC {
	vpcs, err := c.client.VPCs.List(projectName)
	handle(err)
	return vpcs
}

func (c AivenClient) GetVpcPeeringConnectionsList(projectName string, vpcId string) []*aiven.VPCPeeringConnection {
	connections, err := c.client.VPCPeeringConnections.List(projectName, vpcId)
	handle(err)
	return connections
}

func (c AivenClient) Client() *aiven.Client {
	return c.client
}

func (c AivenClient) GetAccountsList() []aiven.Account {
	accountsResponse, err := c.Client().Accounts.List()
	handle(err)
	return accountsResponse.Accounts
}
