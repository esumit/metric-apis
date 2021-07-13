package config

type MetricServerConfig struct {
	Port         string
	IPAddress    string
	WriteTimeout int
	ReadTimeout  int
	IdleTimeout  int
}

type MetricCollectionConfig struct {
	CollectionTime int
}