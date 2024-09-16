package config

var Default = &Config{
	Server: &Server{
		Port:     8080,
		ImageDir: "./images",
	},
	Worker: &Worker{
		ImageDir:       "./images",
		FrequencyS:     600,
		ExecAfterStart: true,
	},
}
