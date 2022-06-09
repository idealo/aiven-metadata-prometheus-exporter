package pkg

import (
	"github.com/aiven/aiven-go-client"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type AivenCollector struct {
	Client *aiven.Client
}

func (ac AivenCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ac, descs)
}

func (ac AivenCollector) Collect(ch chan<- prometheus.Metric) {

	projects := ac.GetProjects()
	ch <- prometheus.MustNewConstMetric(projectInfo, prometheus.GaugeValue, float64(len(projects)))

	ac.processProjects(projects, ch)
}

func (ac *AivenCollector) processProjects(projects []*aiven.Project, ch chan<- prometheus.Metric) {
	for _, project := range projects {
		log.Debug("Fetching infos for " + project.Name)
		countClustersPerProject(ac.Client, project, ch)
		processServices(ac.Client, project, ch)
	}
}

func processServices(client *aiven.Client, project *aiven.Project, ch chan<- prometheus.Metric) {
	services, _ := client.Services.List(project.Name)
	for _, service := range services {
		log.Debug("Fetching service infos for " + project.Name + " and " + service.Name)
		collectServiceNodeCount(ch, service, project)
		collectServiceNodeStates(ch, service, project)
		collectServiceUsersPerService(ch, service, project)
	}
}

func collectServiceNodeCount(ch chan<- prometheus.Metric, service *aiven.Service, project *aiven.Project) {
	ch <- prometheus.MustNewConstMetric(nodeCount, prometheus.GaugeValue, float64(service.NodeCount), project.Name, service.Name)
}

func collectServiceNodeStates(ch chan<- prometheus.Metric, service *aiven.Service, project *aiven.Project) {
	for _, state := range service.NodeStates {
		ch <- prometheus.MustNewConstMetric(nodeState, prometheus.CounterValue, float64(1), project.Name, service.Name, state.Name, state.State)
	}
}

func countClustersPerProject(client *aiven.Client, project *aiven.Project, ch chan<- prometheus.Metric) {
	services, _ := client.Services.List(project.Name)
	ch <- prometheus.MustNewConstMetric(serviceInfo, prometheus.GaugeValue, float64(len(services)), project.Name)
}

func collectServiceUsersPerService(ch chan<- prometheus.Metric, service *aiven.Service, project *aiven.Project) {
	ch <- prometheus.MustNewConstMetric(serviceUserCount, prometheus.GaugeValue, float64(len(service.Users)), project.Name, service.Name)
}

func (ac *AivenCollector) GetProjects() []*aiven.Project {
	log.Debug("Start fetching all projects")
	list, err := ac.Client.Projects.List()
	if err != nil {
		log.Fatalln(err)
	}
	return list
}

var (
	// Basic Info
	projectInfo = prometheus.NewDesc("aiven_project_count_total", "The number of projects registered in the account", nil, nil)
	serviceInfo = prometheus.NewDesc("aiven_service_count", "The number of services per project", []string{"project"}, nil)

	// Service related info
	nodeCount        = prometheus.NewDesc("aiven_service_node_count", "Node Count per Service", []string{"project", "service"}, nil)
	nodeState        = prometheus.NewDesc("aiven_service_node_state", "Node State per Service", []string{"project", "service", "node_name", "state"}, nil)
	serviceUserCount = prometheus.NewDesc("aiven_service_serviceuser_count", "Serviceuser Count per Service", []string{"project", "service"}, nil)
)
