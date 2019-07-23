package config

import (
	"github.com/cloudfoundry-community/go-cfenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server server `mapstructure:"server"`
	S3 s3 `mapstructure:"s3"`
}

type s3 struct {
	Protocol string
	Host string
	Port int
	TmpDir string
	Secret string
	AccessKey string
}

type server struct {
	Host string
	Port int
}

func LoadConfig(env string) *viper.Viper {

	log.Debug("load static configs..")

	var confFile string

	switch env {
		case "dev":
			confFile = "config.dev"
		case "cloud":
			confFile = "config.cloud"
		default:
			confFile = "config.dev"
	}

	v := viper.New()
	v.SetConfigName(confFile)
	v.AddConfigPath("config")

	err := v.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Fatal("Fatal error config file: ", err)
	}

	log.Info("static configuration: ", v)

	vcapErr := v.BindEnv("vcap", "VCAP_SERVICES")

	if vcapErr != nil {
		log.Error("Failed to read vcap services")
	}

	if v.GetString("vcap") != "" {

		log.Debug("Content of VCAP_SERVICES variable: ", v.GetString("vcap"))
		appEnv, err := cfenv.Current()

		if err != nil {
			log.Error("failed to read cf env ", err)
		}

		log.Debug("cf environment: ", appEnv)
		log.Debug("cf services: ", appEnv.Services)

		if v.GetString("s3.serviceName") == "" {
			log.Fatal("you have to provide a s3.serviceName in config to read the service details")
		}

		storage, err := appEnv.Services.WithName(v.GetString("s3.serviceName"))

		if err != nil {
			log.Error("failed to read storage service from cf env")
		} else {
			log.Info("cf storage: ", storage)

			accessHost, ok := storage.CredentialString("accessHost")
			accessKey, ok := storage.CredentialString("accessKey")
			sharedSecret, ok := storage.CredentialString("sharedSecret")
			namespace, ok := storage.CredentialString("namespace")
			namespaceHost, ok := storage.CredentialString("namespaceHost")

			if ok {
				v.Set("s3.host", accessHost)
				v.Set("s3.accessKey", accessKey)
				v.Set("s3.secret", sharedSecret)
				v.Set("s3.namespace", namespace)
				v.Set("s3.namespaceHost", namespaceHost)
			}
		}
	}

	return v
}
