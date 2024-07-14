package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestPayload struct {
	Username string `json:"username"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type LeetCodeResponse struct {
	Data struct {
		MatchedUser struct {
			ContestBadge struct {
				Name      string `json:"name"`
				Expired   bool   `json:"expired"`
				HoverText string `json:"hoverText"`
				Icon      string `json:"icon"`
			} `json:"contestBadge"`
			Username    string `json:"username"`
			GithubUrl   string `json:"githubUrl"`
			TwitterUrl  string `json:"twitterUrl"`
			LinkedinUrl string `json:"linkedinUrl"`
			Profile     struct {
				Ranking                  int      `json:"ranking"`
				UserAvatar               string   `json:"userAvatar"`
				RealName                 string   `json:"realName"`
				AboutMe                  string   `json:"aboutMe"`
				School                   string   `json:"school"`
				Websites                 []string `json:"websites"`
				CountryName              string   `json:"countryName"`
				Company                  string   `json:"company"`
				JobTitle                 string   `json:"jobTitle"`
				SkillTags                []string `json:"skillTags"`
				PostViewCount            int      `json:"postViewCount"`
				PostViewCountDiff        int      `json:"postViewCountDiff"`
				Reputation               int      `json:"reputation"`
				ReputationDiff           int      `json:"reputationDiff"`
				SolutionCount            int      `json:"solutionCount"`
				SolutionCountDiff        int      `json:"solutionCountDiff"`
				CategoryDiscussCount     int      `json:"categoryDiscussCount"`
				CategoryDiscussCountDiff int      `json:"categoryDiscussCountDiff"`
			} `json:"profile"`
		} `json:"matchedUser"`
	} `json:"data"`
}

func leetCodeStatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	graphqlQuery := `query userPublicProfile($username: String!) {
        matchedUser(username: $username) {
            contestBadge {
                name
                expired
                hoverText
                icon
            }
            username
            githubUrl
            twitterUrl
            linkedinUrl
            profile {
                ranking
                userAvatar
                realName
                aboutMe
                school
                websites
                countryName
                company
                jobTitle
                skillTags
                postViewCount
                postViewCountDiff
                reputation
                reputationDiff
                solutionCount
                solutionCountDiff
                categoryDiscussCount
                categoryDiscussCountDiff
            }
        }
    }`

	graphqlRequest := GraphQLRequest{
		Query: graphqlQuery,
		Variables: map[string]interface{}{
			"username": payload.Username,
		},
	}

	requestBody, err := json.Marshal(graphqlRequest)
	if err != nil {
		http.Error(w, "Error creating GraphQL request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewBuffer(requestBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Error fetching data from LeetCode", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var leetcodeResponse LeetCodeResponse
	err = json.NewDecoder(resp.Body).Decode(&leetcodeResponse)
	if err != nil {
		http.Error(w, "Error decoding LeetCode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leetcodeResponse)
}

func main() {
	http.HandleFunc("/leetcodestat", leetCodeStatHandler)
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
