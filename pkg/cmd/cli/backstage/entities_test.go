package backstage

import (
	"testing"
)

func TestListEntities(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	str, err := SetupBackstageTestRESTClient(ts).ListEntities()
	AssertError(t, err)
	AssertEqual(t, TestJSONStringIndented, str)
}
