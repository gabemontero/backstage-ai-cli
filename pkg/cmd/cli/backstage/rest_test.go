package backstage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodDelete = "DELETE"

	TestJSONStringOneLine  = `{"TestGet": "JSON response"}`
	TestJSONStringIndented = `{
    "TestGet": "JSON response"
}`
	TestJSONStringOneLinePlusPathParam  = `{"TestGet": "JSON response path %s"}`
	TestJSONStringOneLinePlusQueryParam = `{"TestGet": "JSON response query %s"}`

	TestPostJSONStringOneLinePlusBody = `{"TestPost": "JSON response body %s"}`

	TestDeleteJSONStringOneLinePlusPathParam = `{"TestDelete": "JSON response path %s"}`
)

var (
	hdrContentTypeKey = http.CanonicalHeaderKey("Content-Type")
)

func SetupBackstageTestRESTClient(ts *httptest.Server) *BackstageRESTClientWrapper {
	backstageTestRESTClient := &BackstageRESTClientWrapper{}
	backstageTestRESTClient.RESTClient = DC()
	backstageTestRESTClient.RootURL = ts.URL
	return backstageTestRESTClient
}

func CreateServer(t *testing.T) *httptest.Server {
	ts := CreateTestServer(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Method: %v", r.Method)
		t.Logf("Path: %v", r.URL.Path)

		switch r.Method {
		case MethodGet:
			switch r.URL.Path {
			case "/":
				_, _ = w.Write([]byte("TestGet: text response"))
				return
			case ENTITIES_URI:
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(TestJSONStringOneLine))
				return
			case LOCATION_URI:
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(TestJSONStringOneLine))
				return
			}

			switch {
			case strings.HasPrefix(r.URL.Path, QUERY_URI):
				w.Header().Set("Content-Type", "application/json")
				values := r.URL.Query()
				filter := values.Get("filter")
				switch {
				case strings.Contains(filter, "api"):
					if strings.Contains(filter, "metadata") {
						_, _ = w.Write([]byte(apisJsonFromTags))
					} else {
						_, _ = w.Write([]byte(apisJson))
					}
				case strings.Contains(filter, "component"):
					_, _ = w.Write([]byte(componentsJson))
				case strings.Contains(filter, "resource"):
					if strings.Contains(filter, "metadata") {
						_, _ = w.Write([]byte(resourcesFromTagsJson))
					} else {
						_, _ = w.Write([]byte(resourcesJson))
					}
				default:
					_, _ = w.Write([]byte(fmt.Sprintf(TestJSONStringOneLinePlusQueryParam, r.URL.Query().Encode())))
				}

			case strings.HasPrefix(r.URL.Path, LOCATION_URI):
				path := strings.TrimPrefix(r.URL.Path, LOCATION_URI)
				w.Header().Set("Content-Type", "application/json")
				if strings.Contains(path, "404") {
					w.WriteHeader(404)
					return
				}
				_, _ = w.Write([]byte(fmt.Sprintf(TestJSONStringOneLinePlusPathParam, path)))
			case strings.HasPrefix(r.URL.Path, API_URI):
				path := strings.TrimPrefix(r.URL.Path, LOCATION_URI)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(fmt.Sprintf(TestJSONStringOneLinePlusPathParam, path)))
			case strings.HasPrefix(r.URL.Path, ENTITIES_URI):
				w.Header().Set("Content-Type", "application/json")
				segs := strings.Split(r.URL.Path, "/")
				ns := segs[len(segs)-2]
				if ns == "404" {
					w.WriteHeader(404)
					return
				}
				_, _ = w.Write([]byte(fmt.Sprintf(TestJSONStringOneLinePlusPathParam, fmt.Sprintf("%s:%s", ns, segs[len(segs)-1]))))
			}
		case MethodPost:
			switch r.URL.Path {
			case LOCATION_URI:
				w.Header().Set("Content-Type", "application/json")
				bodyBuf, err := io.ReadAll(r.Body)
				if err != nil {
					_, _ = w.Write([]byte(fmt.Sprintf(TestPostJSONStringOneLinePlusBody, err.Error())))
					w.WriteHeader(500)
					return
				}
				if len(bodyBuf) == 0 {
					w.WriteHeader(500)
					return
				}
				data := Post{}
				err = json.Unmarshal(bodyBuf, &data)
				if err != nil {
					_, _ = w.Write([]byte(fmt.Sprintf(TestPostJSONStringOneLinePlusBody, err.Error())))
					w.WriteHeader(500)
					return
				}
				_, err = url.Parse(data.Target)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				_, _ = w.Write([]byte(fmt.Sprintf(TestPostJSONStringOneLinePlusBody, data.Target)))
			}
		case MethodDelete:
			switch {
			case strings.HasPrefix(r.URL.Path, LOCATION_URI):
				path := strings.TrimPrefix(r.URL.Path, LOCATION_URI)
				if strings.Contains(path, "404") {
					w.WriteHeader(404)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(fmt.Sprintf(TestDeleteJSONStringOneLinePlusPathParam, path)))
			}
		}
	})

	return ts
}

type Post struct {
	Target string `json:"target"`
	Type   string `json:"type"`
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

func TestSetContext(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	resp, err := DC().R().
		SetContext(context.Background()).
		Get(ts.URL + "/")

	AssertError(t, err)
	AssertEqual(t, http.StatusOK, resp.StatusCode())
	AssertEqual(t, "200 OK", resp.Status())
	AssertEqual(t, true, resp.Body() != nil)
	AssertEqual(t, "TestGet: text response", resp.String())

	LogResponse(t, resp)
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
