Creating a simple yet production-ready Go project involves incorporating best practices, tools, and deployment strategies. Here’s a structured guide to setting up such a project:

### Project Structure
A standard Go project structure:

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   └── myapp/
│       ├── server.go
│       ├── handler.go
│       └── service.go
├── pkg/
│   └── myapp/
│       └── middleware/
│           └── logging.go
├── go.mod
└── go.sum
```

### Dependencies
Use Go modules to manage dependencies. Initialize with:

```sh
go mod init myapp
```

### Main Application
**`cmd/myapp/main.go`**:
```go
package main

import (
    "log"
    "myapp/internal/myapp"
    "net/http"
)

func main() {
    server := myapp.NewServer()
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", server))
}
```

### Server and Handlers
**`internal/myapp/server.go`**:
```go
package myapp

import (
    "net/http"
    "myapp/pkg/myapp/middleware"
)

type Server struct {
    mux *http.ServeMux
}

func NewServer() *Server {
    mux := http.NewServeMux()
    server := &Server{mux: mux}
    server.routes()
    return server
}

func (s *Server) routes() {
    s.mux.Handle("/compare", middleware.Logging(http.HandlerFunc(s.compareHandler)))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mux.ServeHTTP(w, r)
}
```

**`internal/myapp/handler.go`**:
```go
package myapp

import (
    "encoding/json"
    "net/http"
)

type CompareResponse struct {
    GitHubCommits    int `json:"github_commits"`
    LeetCodeSubmissions int `json:"leetcode_submissions"`
}

func (s *Server) compareHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    githubUser := r.URL.Query().Get("github")
    leetcodeUser := r.URL.Query().Get("leetcode")

    githubCommits, err := fetchGitHubCommits(githubUser)
    if err != nil {
        http.Error(w, "Failed to fetch GitHub commits", http.StatusInternalServerError)
        return
    }

    leetcodeSubmissions, err := fetchLeetCodeSubmissions(leetcodeUser)
    if err != nil {
        http.Error(w, "Failed to fetch LeetCode submissions", http.StatusInternalServerError)
        return
    }

    response := CompareResponse{
        GitHubCommits:    githubCommits,
        LeetCodeSubmissions: leetcodeSubmissions,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### Middleware
**`pkg/myapp/middleware/logging.go`**:
```go
package middleware

import (
    "log"
    "net/http"
)

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received %s request for %s", r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}
```

### Fetch Functions
**`internal/myapp/service.go`**:
```go
package myapp

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

type Repository struct {
    Name string `json:"name"`
}

type Commit struct {
    Sha string `json:"sha"`
}

func fetchGitHubCommits(username string) (int, error) {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    reposURL := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
    req, _ := http.NewRequest("GET", reposURL, nil)
    req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

    resp, err := client.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("failed to fetch repositories for user %s", username)
    }

    var repos []Repository
    if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
        return 0, err
    }

    totalCommits := 0
    for _, repo := range repos {
        commitsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", username, repo.Name)
        for {
            req, _ := http.NewRequest("GET", commitsURL, nil)
            req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

            resp, err := client.Do(req)
            if err != nil {
                return 0, err
            }
            defer resp.Body.Close()

            if resp.StatusCode != http.StatusOK {
                return 0, fmt.Errorf("failed to fetch commits for repo %s", repo.Name)
            }

            var commits []Commit
            if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
                return 0, err
            }

            totalCommits += len(commits)

            linkHeader := resp.Header.Get("Link")
            if linkHeader == "" || !strings.Contains(linkHeader, "rel=\"next\"") {
                break
            }

            links := strings.Split(linkHeader, ",")
            for _, link := range links {
                if strings.Contains(link, "rel=\"next\"") {
                    nextLink := strings.Split(link, ";")[0]
                    nextLink = strings.Trim(nextLink, " <>")
                    commitsURL = nextLink
                    break
                }
            }
        }
    }

    return totalCommits, nil
}

type LeetCodeProfile struct {
    MatchedUser struct {
        SubmitStats struct {
            AcSubmissionNum []struct {
                Count int `json:"count"`
            } `json:"acSubmissionNum"`
        } `json:"submitStats"`
    } `json:"matchedUser"`
}

func fetchLeetCodeSubmissions(username string) (int, error) {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    query := fmt.Sprintf(`{"query":"query getUserProfile($username: String!) { matchedUser(username: $username) { submitStats { acSubmissionNum { count } } } }","variables":{"username":"%s"}}`, username)
    req, _ := http.NewRequest("POST", "https://leetcode.com/graphql", strings.NewReader(query))
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("failed to fetch LeetCode submissions for user %s", username)
    }

    var profile LeetCodeProfile
    if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
        return 0, err
    }

    totalSubmissions := 0
    for _, submission := range profile.MatchedUser.SubmitStats.AcSubmissionNum {
        totalSubmissions += submission.Count
    }

    return totalSubmissions, nil
}
```

### Rate Limiting
Integrate rate limiting middleware:

**`pkg/myapp/middleware/ratelimit.go`**:
```go
package middleware

import (
    "net/http"
    "sync"
    "time"
)

func RateLimiting(next http.Handler) http.Handler {
    var rateLimit = 10 // Limit to 10 requests per minute
    var interval = time.Minute
    var requestCounts = make(map[string]int)
    var lastReset = time.Now()
    var mu sync.Mutex

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        mu.Lock()
        defer mu.Unlock()

        ip := r.RemoteAddr
        if time.Since(lastReset) > interval {
            lastReset = time.Now()
            requestCounts = make(map[string]int)
        }

        if requestCounts[ip] >= rateLimit {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        requestCounts[ip]++
        next.ServeHTTP(w, r)
    })
}
```

### Applying Middlewares
Update the server to use logging and rate limiting middlewares:

**`internal/myapp/server.go`**:
```go
package myapp

import (
    "net/http"
    "myapp/pkg/myapp/middleware"
)

type Server struct {
    mux *http.ServeMux
}

func NewServer() *Server {
    mux := http.NewServeMux()
    server := &Server{mux: mux}
    server.routes()
    return server
}

func (s *Server) routes() {
    s.mux.Handle("/compare", middleware.RateLimiting(middleware.Logging(http.HandlerFunc(s.compareHandler))))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mux.ServeHTTP(w, r)
}
```

### Dockerfile
Create a Dockerfile for containerization:

```

Here's how you can complete your production-ready Go application by adding Docker support, setting up a Makefile, and integrating CI/CD with GitHub Actions.

### Dockerfile

Create a Dockerfile to containerize your application:

```Dockerfile
# syntax=docker/dockerfile:1

# Stage 1: Build the application
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /myapp cmd/myapp/main.go

# Stage 2: Create the runtime image
FROM scratch

COPY --from=builder /myapp /myapp

EXPOSE 8080

ENTRYPOINT ["/myapp"]
```

### Makefile

A Makefile simplifies common tasks like building, running, and testing your application.

**Makefile**:
```Makefile
.PHONY: build run clean

build:
    docker build -t myapp .

run:
    docker run -p 8080:8080 myapp

clean:
    docker rmi myapp
```

### GitHub Actions for CI/CD

Create a GitHub Actions workflow file to automate the build and test process.

**.github/workflows/ci.yml**:
```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:

    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.18

    - name: Cache Go modules
      uses: actions/cache@v2
      with:
        path: |
          ~/go/pkg/mod
          ~/.cache/go-build
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test ./...

    - name: Build Docker image
      run: docker build -t myapp .

    - name: Push Docker image
      run: |
        echo "${{ secrets.GITHUB_TOKEN }}" | docker login ghcr.io -u ${{ github.actor }} --password-stdin
        docker tag myapp ghcr.io/${{ github.repository_owner }}/myapp:latest
        docker push ghcr.io/${{ github.repository_owner }}/myapp:latest
```

### Running and Testing

1. **Run the Application**:
    ```sh
    make build
    make run
    ```

2. **Test the Application**:
    ```sh
    go test ./...
    ```

3. **Clean Up**:
    ```sh
    make clean
    ```

### Summary

By setting up a structured project layout, using middleware for logging and rate limiting, containerizing with Docker, automating with Makefile, and setting up CI/CD with GitHub Actions, you've created a production-ready Go application. These practices ensure your application is maintainable, scalable, and deployable.