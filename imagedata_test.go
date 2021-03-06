package collector

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"
)

var metadataSlice []ImageMetadataInfo
var tagSlice []TagInfo

func TestMain(m *testing.M) {
	fmt.Println("TestMain: Run First")
	// make sure environment vars have been setup
	_, _, _, e := dockerAuth()
	if e != nil {
		fmt.Println(e)
		os.Exit(55)
	}
	os.Exit(m.Run())
}

func TestPullImage(t *testing.T) {
	fmt.Println("TestPullImage")
	var e error
	DockerTransport, e = NewDockerTransport(DOCKERPROTO, DOCKERADDR)
	if e != nil {
		t.Fatal(e)
	}
	RegistrySpec = "index.docker.io"
	RegistryAPIURL, HubAPI, BasicAuth, XRegistryAuth = GetRegistryURL()
	metadata := ImageMetadataInfo{
		Repo: "busybox",
		Tag:  "latest",
	}
	fmt.Println("TestPullImage %v", metadata)
	PullImage(metadata)
	return
}

func TestRemoveImage(t *testing.T) {
	fmt.Println("TestRemoveImage")
	TestPullImage(t)
	metadata1 := ImageMetadataInfo{
		Repo: "busybox",
		Tag:  "latest",
	}
	/*
		metadata2 := ImageMetadataInfo{
			Repo: "busybox",
			Tag:  "buildroot-2014.02",
		}
	*/
	// fmt.Println("TestRemoveImage %v %v", metadata1, metadata2)
	fmt.Println("TestRemoveImage %v", metadata1)
	RemoveImages([]ImageMetadataInfo{metadata1}, GetImageToMDMap([]ImageMetadataInfo{metadata1 /*, metadata2*/}))
	return
}

func dockerAuth() (user, password, registry string, e error) {
	user = os.Getenv("DOCKER_USER")
	password = os.Getenv("DOCKER_PASSWORD")
	registry = os.Getenv("DOCKER_REGISTRY")
	if registry == "" {
		registry = "index.docker.io"
	}
	RegistryAPIURL = "https://" + registry
	s := user + ":" + password
	BasicAuth = base64.StdEncoding.EncodeToString([]byte(s))

	if user == "" || password == "" {
		e = fmt.Errorf("Please put valid credentials for registry " + registry + " in envvars DOCKER_USER and DOCKER_PASSWORD.")
		return
	}
	return
}

func TestGetReposHub(t *testing.T) {
	fmt.Println("TestGetReposHub")
	_, _, registry, e := dockerAuth()
	if e != nil {
		t.Fatal(e)
	}
	if registry != "index.docker.io" {
		t.Fatal("TestRegReposHub only works with DOCKER_REGISTRY=index.docker.io")
	}
	ReposToProcess["library/mysql"] = true
	//reposToProcess["ncarlier/redis"] = true
	repo := RepoType("library/mysql")
	client := &http.Client{}
	indexInfo, e := getReposTokenAuthV1(repo, client)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Print(indexInfo, e)
	return
}

func TestGetTagsMetadataHub(t *testing.T) {
	fmt.Println("TestGetTagsMetadataHub")
	_, _, registry, e := dockerAuth()
	if e != nil {
		t.Fatal(e)
	}
	if registry != "index.docker.io" {
		t.Fatal("TestGetTagsMetadataHub only works with DOCKER_REGISTRY=index.docker.io")
	}
	ReposToProcess["library/iojs"] = true
	repo := RepoType("library/iojs")
	client := &http.Client{}
	indexInfo, e := getReposTokenAuthV1(repo, client)
	if e != nil {
		t.Fatal(e)
	}
	tagSlice, e := getTagsTokenAuthV1(repo, client, indexInfo)
	if e != nil {
		t.Fatal(e)
	}
	oldMetadataSet := NewMetadataSet()
	metadataMap := NewImageToMetadataMap(oldMetadataSet)
	metadataSlice, e := getMetadataTokenAuthV1(tagSlice[0], metadataMap, client, indexInfo)
	//tagSlice, metadataSlice, e := getTagsMetadataTokenAuthV1(indexInfo, oldMetadataSet)
	if e != nil {
		t.Fatal(e)
	}
	if tagSlice == nil || len(tagSlice) == 0 {
		t.Fatal("tagSlice", tagSlice)
	}
	if metadataSlice == nil || len(metadataSlice) == 0 {
		t.Fatal("metadataSlice", metadataSlice)
	}
	fmt.Print(tagSlice)
	return
}

func TestParseDistro(t *testing.T) {
	fmt.Println("TestParseDistro")
	var tests = []struct {
		pretty   string
		codename string
	}{
		{"Ubuntu 14.04.1 LTS", "UBUNTU-trusty"},
		{"CentOS Linux 7 (Core)", "REDHAT-7Server"},
	}
	for _, trial := range tests {
		distro := getDistroID(trial.pretty)
		if distro != trial.codename {
			t.Fatal("input:", trial.pretty, "output", distro, "expected:", trial.codename)
		}
		fmt.Println("Found distro: ", distro)
	}
	return
}
