package artifactory

import "net/http"

type Client struct {
	baseUrl string
	http    *http.Client
}

type Artifactory interface {
	GetArchivesForRepo(repo string) (string, error)
}

type Results struct {
	Results []RepoResponse `json:"results"`
}

type RepoResponse struct {
	Repo       string       `json:"repo"`
	Path       string       `json:"path"`
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Size       int64        `json:"size"`
	Created    string       `json:"created"`
	CreatedBy  string       `json:"created_by"`
	Modified   string       `json:"modified"`
	ModifiedBy string       `json:"modified_by"`
	Updated    string       `json:"updated"`
	Stats      []Statistics `json:"stats"`
}

type Statistics struct {
	Downloads int64 `json:"downloads"`
}

type JARResponse struct {
	Name      string
	Downloads int64
}
