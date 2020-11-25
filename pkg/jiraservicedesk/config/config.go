package config

import (
	"os"

	util "github.com/stakater/operator-utils/util"
	secretsUtil "github.com/stakater/operator-utils/util/secrets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("config")

const (
	JiraServiceDeskDefaultSecretName   string = "jira-service-desk-config"
	JiraServiceDeskAPITokenSecretKey   string = "JIRA_SERVICE_DESK_API_TOKEN"
	JiraServiceDeskAPIBaseURLSecretKey string = "JIRA_SERVICE_DESK_API_BASE_URL"
	JiraServiceDeskEmailSecretKey      string = "JIRA_SERVICE_DESK_EMAIL"
)

var (
	JiraServiceDeskSecretName = getConfigSecretName()
)

type ControllerConfig struct {
	ApiToken   string
	ApiBaseUrl string
	Email      string
}

func getConfigSecretName() string {
	configSecretName, _ := os.LookupEnv("CONFIG_SECRET_NAME")
	if len(configSecretName) == 0 {
		configSecretName = JiraServiceDeskDefaultSecretName
		log.Info("CONFIG_SECRET_NAME is unset, using default value: " + JiraServiceDeskDefaultSecretName)
	}
	return configSecretName
}

func LoadControllerConfig(apiReader client.Reader) (ControllerConfig, error) {
	log.Info("Loading Configuration from secret")

	// Retrieve operator namespace
	operatorNamespace, _ := os.LookupEnv("OPERATOR_NAMESPACE")
	if len(operatorNamespace) == 0 {
		operatorNamespaceTemp, err := util.GetOperatorNamespace()
		if err != nil {
			if err.Error() == "namespace not found for current environment" {
				log.Info("Skipping leader election; not running in a cluster.")
			}
			log.Error(err, "Unable to get operator namespace")
		}
		operatorNamespace = operatorNamespaceTemp
	}

	apiToken, err := secretsUtil.LoadSecretData(apiReader, JiraServiceDeskSecretName, operatorNamespace, JiraServiceDeskAPITokenSecretKey)
	if err != nil {
		log.Error(err, "Unable to fetch apiToken from secret")
	}

	apiBaseUrl, err := secretsUtil.LoadSecretData(apiReader, JiraServiceDeskSecretName, operatorNamespace, JiraServiceDeskAPIBaseURLSecretKey)
	if err != nil {
		log.Error(err, "Unable to fetch apiBaseUrl from secret")
	}

	email, err := secretsUtil.LoadSecretData(apiReader, JiraServiceDeskSecretName, operatorNamespace, JiraServiceDeskEmailSecretKey)
	if err != nil {
		log.Error(err, "Unable to fetch email from secret")
	}

	controllerConfig := ControllerConfig{ApiToken: apiToken, ApiBaseUrl: apiBaseUrl, Email: email}

	return controllerConfig, err
}
