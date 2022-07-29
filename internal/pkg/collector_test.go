package pkg

import (
	"github.com/aiven/aiven-go-client"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
			Name:             "TestProject",
			EstimatedBalance: "42.00",
		},
	}
}

func (m MockAivenClient) GetServicesList(_ string) []*aiven.Service {
	return []*aiven.Service{
		{
			Name:       "TestService",
			State:      "Running",
			Type:       "kafka",
			NodeStates: []*aiven.NodeState{{State: "Running"}},
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

func (m MockAivenClient) GetVpcPeeringConnectionsList(_ string, _ string) []*aiven.VPCPeeringConnection {
	return []*aiven.VPCPeeringConnection{
		{
			//	empty
		},
	}

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
			name: "Should meter also when the type is uppercase or capitalized",
			args: args{
				client:  MockAivenClient{},
				service: &aiven.Service{Name: "TestService", Type: "KaFKa"},
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

func TestAivenCollector_processProjects(t *testing.T) {
	mock := MockAivenClient{}
	ac := AivenCollector{client: mock}
	projects := []*aiven.Project{
		{
			EstimatedBalance: "42.00",
			AccountId:        "TestAccountId",
		},
	}

	t.Run("Happy Path", func(t *testing.T) {
		ac.processProjects(projects)

		wantedMetrics := 10

		if len(metrics) != wantedMetrics {
			t.Error("Wanted", wantedMetrics, "got", len(metrics))
		}
	})
}

func TestAivenCollector_CollectAsync(t *testing.T) {
	mock := MockAivenClient{}
	ac := AivenCollector{client: mock}

	t.Run("Happy Path", func(t *testing.T) {
		ac.CollectAsync()

		wantedMetrics := 12
		if len(metrics) != wantedMetrics {
			for _, metric := range metrics {
				log.Error(metric.Desc())
			}
			t.Error("Wanted", wantedMetrics, "got", len(metrics))
		}
	})
}
