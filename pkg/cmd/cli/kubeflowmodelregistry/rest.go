package kubeflowmodelregistry

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gabemontero/backstage-ai-cli/pkg/config"
	"github.com/go-resty/resty/v2"
	"k8s.io/klog/v2"
	"os"
)

const (
	BASE_URI                         = "/api/model_registry/v1alpha3"
	GET_REG_MODEL_URI                = "/registered_models/%s"
	LIST_VERSIONS_OFF_REG_MODELS_URI = "/registered_models/%s/versions"
	LIST_ARTFIACTS_OFF_VERSIONS_URI  = "/model_versions/%s/artifacts"
	LIST_REG_MODEL_URI               = "/registered_models"
	LIST_ALL_MODEL_VERSIONS_URI      = "/model_versions"
	LIST_ALL_MODEL_ARTIFACTS_URI     = "/model_artifacts"
)

type KubeFlowRESTClientWrapper struct {
	RESTClient *resty.Client
	RootURL    string
	Token      string
}

func SetupKubeflowRESTClient(cfg *config.Config) *KubeFlowRESTClientWrapper {
	if cfg == nil {
		klog.Error("Command config is nil")
		klog.Flush()
		os.Exit(1)
	}
	kubeFlowRESTClient := &KubeFlowRESTClientWrapper{
		Token:      cfg.StoreToken,
		RootURL:    cfg.StoreURL + BASE_URI,
		RESTClient: cfg.KubeflowRESTClient,
	}
	if cfg.KubeflowRESTClient != nil {
		return kubeFlowRESTClient
	}
	cfg.KubeflowRESTClient = resty.New()
	kubeFlowRESTClient.RESTClient = cfg.KubeflowRESTClient
	if cfg.KubeflowRESTClient == nil {
		klog.Errorf("Unable to get Kubeflow REST client wrapper")
		klog.Flush()
		os.Exit(1)
	}
	tlsCfg := &tls.Config{}
	if cfg.StoreSkipTLS {
		tlsCfg.InsecureSkipVerify = true
	}
	kubeFlowRESTClient.RESTClient.SetTLSClientConfig(tlsCfg)

	return kubeFlowRESTClient
}

func (k *KubeFlowRESTClientWrapper) processUpdate(resp *resty.Response, action, url, body string) (string, error) {
	postResp := resp.String()
	rc := resp.StatusCode()
	if rc != 200 && rc != 201 {
		return "", fmt.Errorf("%s %s with body %s status code %d resp: %s\n", url, action, body, rc, postResp)
	} else {
		klog.V(4).Infof("%s %s with body %s status code %d resp: %s\n", url, action, body, rc, postResp)
	}
	return k.processBody(resp)
}

func (k *KubeFlowRESTClientWrapper) processBody(resp *resty.Response) (string, error) {
	retJSON := make(map[string]any)
	err := json.Unmarshal(resp.Body(), &retJSON)
	if err != nil {
		return "", fmt.Errorf("json unmarshall error for %s: %s\n", resp.Body(), err.Error())
	}
	id, ok := retJSON["id"]
	if !ok {
		return "", fmt.Errorf("id fetch did not work for %#v\n", retJSON)
	} else {
		klog.V(4).Infof("id %s\n", id)
	}
	return fmt.Sprintf("%s", id), nil
}

func (k *KubeFlowRESTClientWrapper) postToModelRegistry(url, body string) (string, error) {
	resp, err := k.RESTClient.R().SetAuthToken(k.Token).SetBody(body).Post(url)
	if err != nil {
		return "", err
	}

	return k.processUpdate(resp, "post", url, body)
}

func (k *KubeFlowRESTClientWrapper) patchToModelRegistry(url, body string) (string, error) {
	resp, err := k.RESTClient.R().SetAuthToken(k.Token).SetBody(body).Patch(url)
	if err != nil {
		return "", err
	}

	return k.processUpdate(resp, "patch", url, body)
}

func (k *KubeFlowRESTClientWrapper) processFetch(resp *resty.Response, url, action string) (string, error) {
	rc := resp.StatusCode()
	getResp := resp.String()
	if rc != 200 {
		return "", fmt.Errorf("%s for %s rc %d body %s\n", action, url, rc, getResp)
	} else {
		klog.V(4).Infof("%s for %s returned ok\n", action, url)
	}
	jb, err := json.MarshalIndent(getResp, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshall indent error for %s: %s", getResp, err.Error())
	}
	return string(jb), nil
}

func (k *KubeFlowRESTClientWrapper) getFromModelRegistry(url string) ([]byte, error) {
	resp, err := k.RESTClient.R().SetAuthToken(k.Token).Get(url)
	if err != nil {
		return nil, err
	}
	rc := resp.StatusCode()
	getResp := resp.String()
	if rc != 200 {
		return nil, fmt.Errorf("get for %s rc %d body %s\n", url, rc, getResp)
	} else {
		klog.V(4).Infof("get for %s returned ok\n", url)
	}
	return resp.Body(), err

}
