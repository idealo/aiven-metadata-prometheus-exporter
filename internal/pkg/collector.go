package pkg

import (
	"github.com/aiven/aiven-go-client"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var metrics []prometheus.Metric

type AivenCollector struct {
	Client *aiven.Client
}

func (ac *AivenCollector) CollectAsync() {
	metrics = make([]prometheus.Metric, 0)
	projects := ac.GetProjects()
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

func (ac *AivenCollector) processProjects(projects []*aiven.Project) {
	metrics = append(metrics, prometheus.MustNewConstMetric(projectCount, prometheus.GaugeValue, float64(len(projects))))
	for _, project := range projects {
		log.Debug("Fetching infos for " + project.Name)
		countClustersPerProject(ac.Client, project)
		processServices(ac.Client, project)
	}
}

func processServices(client *aiven.Client, project *aiven.Project) {
	services, _ := client.Services.List(project.Name)
	for _, service := range services {
		log.Debug("Fetching service infos for " + project.Name + " and " + service.Name)
		collectServiceNodeCount(service, project)
		collectServiceNodeStates(service, project)
		collectServiceUsersPerService(service, project)
	}
}

func collectServiceNodeCount(service *aiven.Service, project *aiven.Project) {
	metrics = append(metrics, prometheus.MustNewConstMetric(nodeCount, prometheus.GaugeValue, float64(service.NodeCount), project.Name, service.Name))
}

func collectServiceNodeStates(service *aiven.Service, project *aiven.Project) {
	for _, state := range service.NodeStates {
		metrics = append(metrics, prometheus.MustNewConstMetric(nodeState, prometheus.CounterValue, float64(1), project.Name, service.Name, state.Name, state.State))
	}
}

func countClustersPerProject(client *aiven.Client, project *aiven.Project) {
	services, _ := client.Services.List(project.Name)
	metrics = append(metrics, prometheus.MustNewConstMetric(serviceCount, prometheus.GaugeValue, float64(len(services)), project.Name))
}

func collectServiceUsersPerService(service *aiven.Service, project *aiven.Project) {
	metrics = append(metrics, prometheus.MustNewConstMetric(serviceUserCount, prometheus.GaugeValue, float64(len(service.Users)), project.Name, service.Name))
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
	projectCount = prometheus.NewDesc("aiven_project_count_total", "The number of projects registered in the account", nil, nil)
	serviceCount = prometheus.NewDesc("aiven_service_count", "The number of services per project", []string{"project"}, nil)

	// Service related info
	nodeCount        = prometheus.NewDesc("aiven_service_node_count", "Node Count per Service", []string{"project", "service"}, nil)
	nodeState        = prometheus.NewDesc("aiven_service_node_state", "Node State per Service", []string{"project", "service", "node_name", "state"}, nil)
	serviceUserCount = prometheus.NewDesc("aiven_service_serviceuser_count", "Serviceuser Count per Service", []string{"project", "service"}, nil)
)
