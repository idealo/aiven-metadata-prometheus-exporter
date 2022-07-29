package pkg

import (
	"github.com/aiven/aiven-go-client"
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

type MockAivenClient struct{}

func (m MockAivenClient) GetAccountsList() []aiven.Account {
	return []aiven.Account{
		{
			Id:   "TestAccountId",
			Name: "TestAccount"},
	}
}

func (m MockAivenClient) GetAccountTeamsList(accountId string) *[]aiven.AccountTeam {
	return &[]aiven.AccountTeam{
		{
			Name:      "TestTeam",
			AccountId: accountId,
		},
	}

}

func (m MockAivenClient) GetAccountTeamMembersList(_ string, teamId string) *[]aiven.AccountTeamMember {
	return &[]aiven.AccountTeamMember{
		{
			UserId: "TestUser",
			TeamId: teamId,
		},
	}
}

func (m MockAivenClient) GetKafkaTopicsList(_ string, _ string) []*aiven.KafkaListTopic {
	return []*aiven.KafkaListTopic{
		{TopicName: "TestTopic"},
	}
}

func (m MockAivenClient) GetProjectsList() []*aiven.Project {
	return []*aiven.Project{
		{
			Name: "TestProject",
		},
	}
}

func (m MockAivenClient) GetServicesList(_ string) []*aiven.Service {
	return []*aiven.Service{
		{
			Name: "TestService",
		},
	}
}

func (m MockAivenClient) GetVpcsList(_ string) []*aiven.VPC {
	return []*aiven.VPC{
		{
			ProjectVPCID: "TestVPCId",
		},
	}
}

func (m MockAivenClient) GetVpcPeeringConnectionsList(projectName string, vpcId string) []*aiven.VPCPeeringConnection {
	return []*aiven.VPCPeeringConnection{
		{
			//	empty
		},
	}

}

func TestAivenCollector_processAccountInfo(t *testing.T) {
	mock := MockAivenClient{}
	ac := AivenCollector{client: mock}

	t.Run("Happy Path", func(t *testing.T) {
		ac.processAccountInfo()
		if len(metrics) == 0 {
			t.Fail()
		}
	})
}

func Test_collectServiceTopicCount(t *testing.T) {
	type args struct {
		client  Client
		service *aiven.Service
		project *aiven.Project
	}

	tests := []struct {
		name          string
		args          args
		wantedMetrics int
	}{
		{
			name: "Happy Path",
			args: args{
				client:  MockAivenClient{},
				service: &aiven.Service{Name: "TestService", Type: "kafka"},
				project: &aiven.Project{Name: "TestProject"},
			},
			wantedMetrics: 1,
		},
		{
			name: "Should not meter when service is not of type kafka",
			args: args{
				client:  MockAivenClient{},
				service: &aiven.Service{Name: "TestService", Type: "not-kafka"},
				project: &aiven.Project{Name: "TestProject"},
			},
			wantedMetrics: 0,
		},
	}

	for _, tt := range tests {
		metrics = make([]prometheus.Metric, 0)
		t.Run(tt.name, func(t *testing.T) {
			collectServiceTopicCount(tt.args.client, tt.args.service, tt.args.project)
			if len(metrics) != tt.wantedMetrics {
				t.Fail()
			}
		})
	}
}
