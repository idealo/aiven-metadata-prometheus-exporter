package pkg

import (
	"github.com/aiven/aiven-go-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

type AivenCollector struct {
	Client *aiven.Client
}

func (ac AivenCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ac, descs)
}

func (ac AivenCollector) Collect(ch chan<- prometheus.Metric) {

	//https://dev.to/metonymicsmokey/custom-prometheus-metrics-with-go-520n

	projects := ac.GetProjects()
	numProjects.Set(float64(len(projects)))
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
		log.Debug("Fetching infos for " + project.Name + " and " + service.Name)
		collectServiceNodeCount(ch, service, project)
	}
}

func collectServiceNodeCount(ch chan<- prometheus.Metric, service *aiven.Service, project *aiven.Project) {
	ch <- prometheus.MustNewConstMetric(nodeCount, prometheus.GaugeValue, float64(service.NodeCount), project.Name, service.Name)
}

func countClustersPerProject(client *aiven.Client, project *aiven.Project, ch chan<- prometheus.Metric) {
	services, _ := client.Services.List(project.Name)
	//servicesPerProject.WithLabelValues(project.Name).Set(float64(len(services)))
	ch <- prometheus.MustNewConstMetric(serviceInfo, prometheus.GaugeValue, float64(len(services)), project.Name)
}

func countServiceUsersPerService(client *aiven.Client, service string) {

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
	numProjects = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "aiven_projects_count_total",
		Help: "The total number of registered Aiven projects",
	})

	numServiceUsers = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aiven_service_user_count",
	}, []string{"project", "service"})

	servicesPerProject = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aiven_services_count",
	}, []string{"project"})

	nodeCount = prometheus.NewDesc("aiven_service_node_count", "Node Count per Service", []string{"project", "service"}, nil)

	projectInfo = prometheus.NewDesc(prometheus.BuildFQName("aiven", "project", "count"),
		"The number of projects registered in the account",
		nil, nil)

	serviceInfo = prometheus.NewDesc(
		prometheus.BuildFQName("aiven", "service", "count"),
		"The number of services per project",
		[]string{"project"},
		nil)
)
