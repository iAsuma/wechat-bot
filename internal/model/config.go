package model

type OpenConfig struct {
	BaseUrl string
	ApiKey  string
	Model   string
}

type EmailConfig struct {
	Host     string
	Port     string
	From     string
	To       string
	Password string
}
