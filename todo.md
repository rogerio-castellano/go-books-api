# To-Do List

- [x] Add CORS Support

  - Uncomment and configure the CORS headers to allow cross-origin requests if needed.
  - Alternatively, use the github.com/rs/cors package for better CORS handling.

- [x] Use a Persistent Storage

  - Replace the in-memory books slice with a database or file-based storage to persist data across server restarts.

- [x] Use HTTP Status Codes Properly

  - Ensure all responses use appropriate HTTP status codes (e.g., 201 Created for successful POST).

- [x] Set connection string and data repository in environment variables

- [ ] Handle Edge Cases

  - Handle cases like duplicate book IDs, invalid JSON payloads, and unsupported HTTP methods.

- [ ] Improve Code Comments

  - Add meaningful comments to explain the purpose of each function and complex logic when it cannot be done through clean code.

- [ ] Improve architecture

  - Create services.
  - Add ORM

- [ ] Handle errors

  - Handle main errors (including currently disregarded with \_ ) analyzing who (the caller or the called) should raise it, according to the situation (the caller can workaround it or not?)

- [x] Remove all images that are no longer within the scope of the Docker Compose configuration, including Nginx

MOVE TO OTHER GITHUB REPOSITORIES TO REDUCE COGNITIVE LOAD ON A SINGLE REPOSITORY

- [ ] Add Unit Tests

  - Write unit tests for each handler function to ensure correctness and prevent regressions.

- [ ] Use Context for Request Scoping

  - Use context.Context to handle request-scoped values like timeouts or user authentication.

- [ ] Format and Lint Code

  - Use gofmt and a linter like golangci-lint to ensure consistent formatting and catch potential issues.

- [ ] Build containers health check

- [ ] Integrate MongoDB, RabbitMQ and Redis to the API, increasing the scope of the API

  - Relocate connection data to environment variables configured within the service in the Docker Compose file

- [x] Add MongoDB, RabbitMQ and Redis, with connectivity tested
