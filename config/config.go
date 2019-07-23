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

	var confFile string

	switch env {
	case "dev":
		confFile = "config.dev"
	case "stag":
		confFile = "config.stag"
	case "prod":
		confFile = "config.prod"
	default:
		confFile = "config.dev"
	}

	// Create new Viper Config Struct
	v := viper.New()
	v.SetConfigName(confFile) // name of config file (without extension)
	v.AddConfigPath("config")

	vcaperr := v.BindEnv("vcap", "VCAP_SERVICES")

	if vcaperr != nil {
		log.Error("Failed to read vcap services")
	}

	if v.GetString("vcap") != "" {

		log.Debug("Available cf env: ", v.GetString("vcap"))
		appEnv, err := cfenv.Current()

		if err != nil {
			log.Error("failed to read cf env ", err)
		}

		log.Debug("cf environment: ", appEnv)
		log.Debug("_________________________________________")
		log.Debug("cf services: ", appEnv.Services)

		storage, err := appEnv.Services.WithName("cd-s3")

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
				v.Set("s3.port", 443)
				v.Set("s3.protocol", "https")
				v.Set("s3.accessKey", accessKey)
				v.Set("s3.secret", sharedSecret)
				v.Set("s3.namespace", namespace)
				v.Set("s3.namespaceHost", namespaceHost)
			}
		}
	}

	err := v.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Fatal("Fatal error config file: ", err)
	}

	return v
}
