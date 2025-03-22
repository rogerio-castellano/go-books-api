# To-Do List

- [ ] Add CORS Support

  - Uncomment and configure the CORS headers to allow cross-origin requests if needed.
  - Alternatively, use the github.com/rs/cors package for better CORS handling.

- [ ] Use a Persistent Storage

  - Replace the in-memory books slice with a database or file-based storage to persist data across server restarts.

- [ ] Use HTTP Status Codes Properly

  - Ensure all responses use appropriate HTTP status codes (e.g., 201 Created for successful POST).

- [ ] Add Unit Tests

  - Write unit tests for each handler function to ensure correctness and prevent regressions.

- [ ] Improve Code Comments

  - Add meaningful comments to explain the purpose of each function and complex logic.

- [ ] Use Context for Request Scoping

  - Use context.Context to handle request-scoped values like timeouts or user authentication.

- [ ] Handle Edge Cases

  - Handle cases like duplicate book IDs, invalid JSON payloads, and unsupported HTTP methods.

- [ ] Format and Lint Code

  - Use gofmt and a linter like golangci-lint to ensure consistent formatting and catch potential issues.

- [ ] Improve architecture

  - Create services.
