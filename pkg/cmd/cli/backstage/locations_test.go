package backstage

import (
	"testing"
)

func TestListLocations(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	str, err := SetupBackstageTestRESTClient(ts).ListLocations()
	AssertError(t, err)
	AssertEqual(t, TestJSONStringIndented, str)
}

func TestGetLocations(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	key := "key1"
	str, err := SetupBackstageTestRESTClient(ts).GetLocation(key)
	AssertError(t, err)
	AssertContains(t, str, key)
}

func TestGetLocationsError(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	nsName := "404:404"
	_, err := SetupBackstageTestRESTClient(ts).GetLocation(nsName)
	if err == nil {
		t.Error("expected error")
	}
}

func TestImportLocation(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	arg := "https://my-repo/my.yaml"
	str, err := SetupBackstageTestRESTClient(ts).ImportLocation(arg)
	AssertError(t, err)
	AssertContains(t, str, arg)
}

func TestImportLocationError(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	arg := ":"
	_, err := SetupBackstageTestRESTClient(ts).ImportLocation(arg)
	if err == nil {
		t.Error("expected error")
	}
}

func TestDeleteLocation(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	arg := "my-location-id"
	str, err := SetupBackstageTestRESTClient(ts).DeleteLocation(arg)
	AssertError(t, err)
	AssertContains(t, str, arg)
}

func TestDeleteLocationsError(t *testing.T) {
	ts := CreateServer(t)
	defer ts.Close()

	nsName := "404:404"
	_, err := SetupBackstageTestRESTClient(ts).DeleteLocation(nsName)
	if err == nil {
		t.Error("expected error")
	}
}
