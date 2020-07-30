package config

import (
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/stakater/jira-service-desk-operator/util"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("config")

const (
	JiraServiceDeskDefaultSecretName   string = "jira-service-desk-config"
	JiraServiceDeskAPITokenSecretKey   string = "JIRA_SERVICE_DESK_API_TOKEN"
	JiraServiceDeskAPIBaseURLSecretKey string = "JIRA_SERVICE_DESK_API_BASE_URL"
)

var (
	JiraServiceDeskSecretName = getConfigSecretName()
)

func getConfigSecretName() string {
	configSecretName, _ := os.LookupEnv("CONFIG_SECRET_NAME")
	if len(configSecretName) == 0 {
		configSecretName = JiraServiceDeskDefaultSecretName
		log.Info("CONFIG_SECRET_NAME is unset, using default value: " + JiraServiceDeskDefaultSecretName)
	}
	return configSecretName
}

func LoadControllerConfig(apiReader client.Reader) (string, string, error) {
	log.Info("Loading Configuration from secret")

	// Retrieve operator namespace
	operatorNamespace, _ := os.LookupEnv("OPERATOR_NAMESPACE")
	if len(operatorNamespace) == 0 {
		operatorNamespaceTemp, err := k8sutil.GetOperatorNamespace()
		if err != nil {
			if err == k8sutil.ErrNoNamespace {
				log.Info("Skipping leader election; not running in a cluster.")
			}
			log.Error(err, "Unable to get operator namespace")
		}
		operatorNamespace = operatorNamespaceTemp
	}

	apiToken, err := util.LoadSecretData(apiReader, JiraServiceDeskSecretName, operatorNamespace, JiraServiceDeskAPITokenSecretKey)
	if err != nil {
		log.Error(err, "Unable to fetch apiToken from secret")
	}

	apiBaseUrl, err := util.LoadSecretData(apiReader, JiraServiceDeskSecretName, operatorNamespace, JiraServiceDeskAPIBaseURLSecretKey)
	if err != nil {
		log.Error(err, "Unable to fetch apiBaseUrl from secret")
	}

	return apiToken, apiBaseUrl, err
}
