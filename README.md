To create a production-ready Go application that fetches data from the GitHub API and LeetCode, compares commit numbers to submission counts, and serves the result over its own REST API, you need to consider several important concepts and best practices. Here’s a comprehensive list:

### 1. **HTTP Client Configuration**

- **Timeouts**: Set reasonable timeouts for your HTTP clients to prevent hanging requests.
- **Retries and Backoff**: Implement retry logic with exponential backoff for transient errors.
- **Connection Pooling**: Optimize HTTP client performance with connection pooling.

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
```

### 2. **Rate Limiting**

- **GitHub API Rate Limiting**: Handle GitHub's rate limits gracefully by checking the headers in responses.
- **Application Rate Limiting**: Use a rate limiter for your application to prevent abuse.

```go
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(1, 5) // 1 request per second, burst of 5

if !limiter.Allow() {
    log.Println("Rate limit exceeded")
    // Handle rate limit exceeded
}
```

### 3. **Authentication and Authorization**

- **API Tokens**: Use environment variables or a secure vault to manage API tokens.
- **Secure Endpoints**: Ensure that your REST API endpoints are secure, using authentication methods like OAuth, JWT, etc.

### 4. **Error Handling**

- **Consistent Error Responses**: Ensure your API provides consistent and informative error responses.
- **Logging**: Log errors appropriately without exposing sensitive information.

### 5. **Concurrency and Goroutines**

- **Goroutines**: Use goroutines to handle concurrent operations.
- **Synchronization**: Use channels and sync primitives to manage concurrent access to shared resources.

```go
var wg sync.WaitGroup
wg.Add(2)

go func() {
    defer wg.Done()
    // Fetch GitHub data
}()

go func() {
    defer wg.Done()
    // Fetch LeetCode data
}()

wg.Wait()
```

### 6. **Configuration Management**

- **Environment Variables**: Use environment variables for configuration.
- **Configuration Files**: Use configuration files (YAML, JSON, etc.) if necessary.
- **Config Libraries**: Use libraries like `viper` for managing configurations.

### 7. **Dependency Management**

- **Go Modules**: Use Go modules for dependency management.
- **Versioning**: Pin dependencies to specific versions to ensure consistency.

```sh
go mod init your-module
go mod tidy
```

### 8. **Testing**

- **Unit Tests**: Write unit tests for your functions.
- **Integration Tests**: Write integration tests to test API interactions.
- **Mocking**: Use mocking frameworks like `gomock` for testing external dependencies.

```go
import "testing"

func TestCompareProfiles(t *testing.T) {
    // Test logic here
}
```

### 9. **Documentation**

- **API Documentation**: Use tools like Swagger or OpenAPI to document your REST API.
- **Code Documentation**: Comment your code thoroughly for maintainability.

### 10. **Monitoring and Logging**

- **Monitoring**: Use tools like Prometheus and Grafana for monitoring your application.
- **Structured Logging**: Use structured logging libraries like `logrus` or `zap`.

```go
import "github.com/sirupsen/logrus"

var log = logrus.New()

log.WithFields(logrus.Fields{
    "username": "john_doe",
}).Info("User logged in")
```

### 11. **Deployment and Scaling**

- **Docker**: Containerize your application using Docker.
- **Kubernetes**: Use Kubernetes for orchestration and scaling.
- **CI/CD**: Implement continuous integration and deployment pipelines.

```dockerfile
FROM golang:1.16

WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]
```

### 12. **Security**

- **Input Validation**: Validate all inputs to prevent injection attacks.
- **TLS/SSL**: Use HTTPS to secure communication.
- **Secrets Management**: Use secure methods for managing secrets (e.g., HashiCorp Vault).

### 13. **API Design**

- **RESTful Principles**: Follow RESTful design principles.
- **Pagination**: Implement pagination for endpoints returning large datasets.
- **Versioning**: Use versioning for your API to handle changes without breaking clients.

### Example Structure:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/shurcooL/graphql"
)

func main() {
    router := gin.Default()
    router.GET("/compare", compareHandler)
    router.Run(":8080")
}

func compareHandler(c *gin.Context) {
    githubUsername := c.Query("github")
    leetcodeUsername := c.Query("leetcode")

    if githubUsername == "" || leetcodeUsername == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username"})
        return
    }

    // Fetch data from GitHub and LeetCode
    githubData := fetchGitHubData(githubUsername)
    leetcodeData := fetchLeetCodeData(leetcodeUsername)

    // Compare data and return result
    comparison := compareData(githubData, leetcodeData)
    c.JSON(http.StatusOK, comparison)
}

func fetchGitHubData(username string) map[string]interface{} {
    // Implement GitHub API call
    return map[string]interface{}{"commits": 100}
}

func fetchLeetCodeData(username string) map[string]interface{} {
    // Implement LeetCode API call
    return map[string]interface{}{"submissions": 200}
}

func compareData(githubData, leetcodeData map[string]interface{}) map[string]interface{} {
    commits := githubData["commits"].(int)
    submissions := leetcodeData["submissions"].(int)
    return map[string]interface{}{
        "commits":     commits,
        "submissions": submissions,
        "ratio":       float64(commits) / float64(submissions),
    }
}
```

This example covers basic implementation and structure. Each section can be expanded with more robust implementations and additional features as needed.