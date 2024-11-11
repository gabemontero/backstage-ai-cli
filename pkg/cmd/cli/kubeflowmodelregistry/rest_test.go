package kubeflowmodelregistry

import (
	"bufio"
	"github.com/gabemontero/backstage-ai-cli/pkg/config"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const (
	MethodGet = "GET"

	TestJSONStringRegisteredModelOneLine    = `{"items":[{"createTimeSinceEpoch":"1731103949567","customProperties":{"foo":{"metadataType":"MetadataStringValue","string_value":"bar"}},"description":"dummy model 1","id":"1","lastUpdateTimeSinceEpoch":"1731103975700","name":"model-1","owner":"kube:admin","state":"LIVE"}],"nextPageToken":"","pageSize":0,"size":1}`
	TestJSONStringRegisteredModelOneLineGet = `{"createTimeSinceEpoch":"1731103949567","customProperties":{"foo":{"metadataType":"MetadataStringValue","string_value":"bar"}},"description":"dummy model 1","id":"1","lastUpdateTimeSinceEpoch":"1731103975700","name":"model-1","owner":"kube:admin","state":"LIVE"}`
	TestJSONStringModelVersionOneLine       = `{"items":[{"author":"kube:admin","createTimeSinceEpoch":"1731103949724","customProperties":{},"description":"version 1","id":"2","lastUpdateTimeSinceEpoch":"1731103949724","name":"v1","registeredModelId":"1","state":"LIVE"}],"nextPageToken":"","pageSize":0,"size":1}`
	TestJSONStringModelArtifactOneLine      = `{"items":[{"artifactType":"model-artifact","createTimeSinceEpoch":"1731103949909","customProperties":{},"description":"version 1","id":"1","lastUpdateTimeSinceEpoch":"1731103949909","modelFormatName":"tensorflow","modelFormatVersion":"v1","name":"model-1-v1-artifact","state":"LIVE","uri":"https://foo.com"}],"nextPageToken":"","pageSize":0,"size":1}`
)

var (
	hdrContentTypeKey = http.CanonicalHeaderKey("Content-Type")
)

func SetupKubeflowTestRESTClient(ts *httptest.Server, cfg *config.Config) {
	cfg.StoreURL = ts.URL
	cfg.KubeflowRESTClient = DC()
}

func CreateGetServer(t *testing.T) *httptest.Server {
	ts := CreateTestServer(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Method: %v", r.Method)
		t.Logf("Path: %v", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case MethodGet:
			switch {
			case strings.HasSuffix(r.URL.Path, LIST_REG_MODEL_URI):
				_, _ = w.Write([]byte(TestJSONStringRegisteredModelOneLine))
			case strings.HasSuffix(r.URL.Path, "/versions"):
				_, _ = w.Write([]byte(TestJSONStringModelVersionOneLine))
			case strings.HasSuffix(r.URL.Path, "/artifacts"):
				_, _ = w.Write([]byte(TestJSONStringModelArtifactOneLine))
			case strings.Contains(r.URL.Path, LIST_REG_MODEL_URI):
				_, _ = w.Write([]byte(TestJSONStringRegisteredModelOneLineGet))
			}
		}
	})

	return ts
}

func CreateTestServer(fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fn))
}

func AssertEqual(t *testing.T, e, g interface{}) (r bool) {
	t.Helper()
	if !Equal(e, g) {
		t.Errorf("Expected [%v], got [%v]", e, g)
	}

	return
}

func AssertContains(t *testing.T, s, sub string) {
	t.Helper()
	if !strings.Contains(s, sub) {
		t.Errorf("Expected '%s' to contain '%s'", s, sub)
	}
	return
}

func AssertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occurred [%v]", err)
	}
}

func AssertLineCompare(t *testing.T, str, expect string, minLine int) {
	scanner1 := bufio.NewScanner(strings.NewReader(str))
	scanner2 := bufio.NewScanner(strings.NewReader(expect))
	lineNum := 1
	for scanner1.Scan() && scanner2.Scan() {
		line1 := scanner1.Text()
		line2 := scanner2.Text()

		if line1 != line2 {
			if lineNum < minLine {
				return
			}

			t.Errorf("line diff at line %d between %s and %s", lineNum, line1, line2)
		}
		lineNum++
	}
}

func Equal(expected, got interface{}) bool {
	return reflect.DeepEqual(expected, got)
}

func DC() *resty.Client {
	c := resty.New()
	c.SetLogger(&logger{})
	return c
}

type logger struct{}

func (l *logger) Errorf(format string, v ...interface{}) {
}

func (l *logger) Warnf(format string, v ...interface{}) {
}

func (l *logger) Debugf(format string, v ...interface{}) {
}

func LogResponse(t *testing.T, resp *resty.Response) {
	t.Logf("Response Status: %v", resp.Status())
	t.Logf("Response Time: %v", resp.Time())
	t.Logf("Response Headers: %v", resp.Header())
	t.Logf("Response Cookies: %v", resp.Cookies())
	t.Logf("Response Body: %v", resp)
}
