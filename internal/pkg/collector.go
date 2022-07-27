package pkg

import (
	"github.com/aiven/aiven-go-client"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	metrics     []prometheus.Metric
	accountInfo = make(map[string]string)

	// Account related
	teamCount       = prometheus.NewDesc("aiven_account_team_count_total", "The number of teams per account", []string{"account"}, nil)
	teamMemberCount = prometheus.NewDesc("aiven_account_team_member_count_total", "The number of members per team for an account", []string{"account", "team"}, nil)

	// Project related
	projectCount = prometheus.NewDesc("aiven_project_count_total", "The number of projects registered in the account", []string{"account"}, nil)
	serviceCount = prometheus.NewDesc("aiven_service_count_total", "The number of services per project", []string{"account", "project"}, nil)

	// Service related info
	nodeCount        = prometheus.NewDesc("aiven_service_node_count_total", "Node count per service", []string{"account", "project", "service"}, nil)
	nodeState        = prometheus.NewDesc("aiven_service_node_state_info", "Node state per service", []string{"account", "project", "service", "node_name", "state"}, nil)
	serviceUserCount = prometheus.NewDesc("aiven_service_serviceuser_coun_total", "Serviceuser count per service", []string{"account", "project", "service"}, nil)
)

type AivenCollector struct {
	AivenClient *aiven.Client
}

func (ac AivenCollector) CollectAsync() {
	metrics = make([]prometheus.Metric, 0)
	ac.processAccountInfo()
	projects := ac.getProjects()
	ac.processProjects(projects)
}

func (ac AivenCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ac, descs)
}

func (ac AivenCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range metrics {
		ch <- metric
	}
}

func (ac AivenCollector) processAccountInfo() {
	accountsResponse, err := ac.AivenClient.Accounts.List()
	handle(err)
	log.Debug("Fetching account infos")
	for _, acc := range accountsResponse.Accounts {

		accountInfo[acc.Id] = acc.Name

		teamsResponse, err := ac.AivenClient.AccountTeams.List(acc.Id)
		handle(err)
		metrics = append(metrics, prometheus.MustNewConstMetric(teamCount, prometheus.GaugeValue, float64(len(teamsResponse.Teams)), acc.Name))
		ac.collectAccountTeams(teamsResponse, acc)
	}
}

func (ac AivenCollector) collectAccountTeams(teamsResponse *aiven.AccountTeamsResponse, acc aiven.Account) {
	log.Debug("Collecting account team infos for account: " + acc.Name)
	for _, team := range teamsResponse.Teams {
		membersResponse, _ := ac.AivenClient.AccountTeamMembers.List(acc.Id, team.Id)
		metrics = append(metrics, prometheus.MustNewConstMetric(teamMemberCount, prometheus.GaugeValue, float64(len(membersResponse.Members)), acc.Name, team.Name))
	}
}

func (ac AivenCollector) processProjects(projects []*aiven.Project) {
	projectCountPerAccount := make(map[string]int)

	for _, project := range projects {
		log.Debug("Fetching infos for " + project.Name)
		ac.countClustersPerProject(project)
		ac.processServices(project)
		projectCountPerAccount[project.AccountId]++
	}

	for key, count := range projectCountPerAccount {
		metrics = append(metrics, prometheus.MustNewConstMetric(projectCount, prometheus.GaugeValue, float64(count), accountInfo[key]))
	}
}

func (ac AivenCollector) processServices(project *aiven.Project) {
	services, err := ac.AivenClient.Services.List(project.Name)
	handle(err)
	for _, service := range services {
		log.Debug("Fetching service infos for " + project.Name + " and " + service.Name)
		collectServiceNodeCount(service, project)
		collectServiceNodeStates(service, project)
		collectServiceUsersPerService(service, project)
	}
}

func collectServiceNodeCount(service *aiven.Service, project *aiven.Project) {
	metrics = append(metrics, prometheus.MustNewConstMetric(nodeCount, prometheus.GaugeValue, float64(service.NodeCount), accountInfo[project.AccountId], project.Name, service.Name))
}

func collectServiceNodeStates(service *aiven.Service, project *aiven.Project) {
	for _, state := range service.NodeStates {
		metrics = append(metrics, prometheus.MustNewConstMetric(nodeState, prometheus.CounterValue, float64(1), accountInfo[project.AccountId], project.Name, service.Name, state.Name, state.State))
	}
}

func (ac AivenCollector) countClustersPerProject(project *aiven.Project) {
	services, err := ac.AivenClient.Services.List(project.Name)
	handle(err)
	metrics = append(metrics, prometheus.MustNewConstMetric(serviceCount, prometheus.GaugeValue, float64(len(services)), accountInfo[project.AccountId], project.Name))
}

func collectServiceUsersPerService(service *aiven.Service, project *aiven.Project) {
	metrics = append(metrics, prometheus.MustNewConstMetric(serviceUserCount, prometheus.GaugeValue, float64(len(service.Users)), accountInfo[project.AccountId], project.Name, service.Name))
}

func (ac AivenCollector) getProjects() []*aiven.Project {
	log.Debug("Start fetching all projects")
	list, err := ac.AivenClient.Projects.List()
	handle(err)
	return list
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
