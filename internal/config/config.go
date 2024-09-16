package config

type Config struct {
	Server *Server
	Worker *Worker
}

type Server struct {
	ImageDir string
	Port     int
}

type Worker struct {
	ImageDir       string
	FrequencyS     int
	ExecAfterStart bool
}
