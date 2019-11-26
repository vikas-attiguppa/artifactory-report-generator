package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
)

func DefaultClient() *Client {
	return &Client{
		baseUrl: os.Getenv("ARTIFACTORY_BASE_URL"),
		http:    &http.Client{},
	}
}

func (client *Client) GetTopArchivesForRepo(repo string) ([]byte, error) {
	response := []byte("Please try with valid input")
	token := os.Getenv("ARTIFACTORY_API_TOKEN")
	req, _ := http.NewRequest("POST", client.baseUrl+"api/search/aql", bytes.NewBufferString(generateRequestBody(repo)))
	req.Header.Set("X-JFrog-Art-Api", token)

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while getting the archieves for repo %s - http status code: %d", repo, resp.StatusCode)
	}

	results := &Results{}
	err = json.NewDecoder(resp.Body).Decode(results)
	
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 0 {
		response, err = json.Marshal(populateServiceResponse(sortResults(results.Results)))
	}
	return response, nil
}

func generateRequestBody(repo string) string {
	return fmt.Sprintf("items.find({\"repo\":\"%s\",\"name\" : {\"$match\":\"*.jar\"}}).include(\"stat.downloads\").sort({\"$desc\" : [\"stat.downloads\",\"name\"]})", repo)
}

func sortResults(repos []Response) []Response {
	sort.Slice(repos, func(i, j int) bool {
		if repos[i].Stats[0].Downloads < repos[j].Stats[0].Downloads {
			return false
		}
		if repos[i].Stats[0].Downloads > repos[j].Stats[0].Downloads {
			return true
		}
		return repos[i].Name < repos[j].Name
	})
	return repos
}

func populateServiceResponse(repos []Response) *[]JARResponse {
	var responses []JARResponse
	if len(repos) > 0 {
		for i := 0; i < 2; i++ {
			responses = append(responses, JARResponse{repos[i].Name, repos[i].Stats[0].Downloads})
		}
	}
	return &responses
}
