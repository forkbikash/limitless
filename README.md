# limitless

[![Go Reference](https://pkg.go.dev/badge/github.com/forkbikash/limitless.svg)](https://pkg.go.dev/github.com/forkbikash/limitless)
[![Go Report Card](https://goreportcard.com/badge/github.com/forkbikash/limitless)](https://goreportcard.com/report/github.com/forkbikash/limitless)

`limitless` is a Go package that provides rate-limiting functionality using the token bucket algorithm. It offers in-memory and Redis-based implementations of the token bucket, allowing you to control the rate of requests or operations in your application.

## Features

- **In-Memory Token Bucket**: An in-memory implementation of the token bucket algorithm for rate-limiting.
- **Redis Token Bucket**: A Redis-based implementation of the token bucket algorithm, providing a distributed rate-limiting solution.
- **Configurable Rate and Capacity**: Customize the rate (requests/operations per second) and capacity (maximum burst size) of the token bucket.
- **Concurrent Access**: Support for concurrent access to the rate limiter with proper locking mechanisms.
- **Extensible**: The package provides a `RateLimiter` interface, allowing you to implement custom rate-limiting strategies if needed.

## Installation

To install the `limitless` package, run the following command:

```bash
go get github.com/forkbikash/limitless
```

## Usage

### In-Memory Token Bucket

```go
import "github.com/forkbikash/limitless"

// Create a new in-memory token bucket
tb := limitless.NewInMemoryTokenBucket(10, 2) // capacity: 10, rate: 2 requests/second

// Check if a request is allowed
allowed, err := limitless.Allow(tb)
if err != nil {
    // Handle error
}
if *allowed {
    // Process the request
} else {
    // Request is rate-limited
}
```

## Contributing

Contributions to the limitless package are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## What's next

- Middleware Integration: Develop a middleware component that can be easily integrated into popular Golang web frameworks (e.g., Gin, Echo, Fiber) to apply the rate limiting functionality
- Metrics and Reporting: Include options to expose metrics related to the rate limiting, such as the number of requests processed, the number of requests rejected due to rate limiting, and the current rate limiting state
- Burst Capacity: Allow users to configure a burst capacity, which is the maximum number of requests that can be processed immediately, even if it exceeds the rate limit. This can be useful for handling sudden spikes in traffic.
- Adaptive Rate Limiting: Implement an adaptive rate limiting algorithm that can dynamically adjust the rate limit based on factors such as the current load, resource utilization, or any other custom metrics
- Multi-Dimensional Rate Limiting: Provide the ability to apply rate limiting based on multiple dimensions, such as client IP, user identity, API endpoint, or any combination of these. This can help with more granular control over the rate limiting policies
- Contextual Information: Pass contextual information (e.g., request metadata, user identity) through the rate limiting system, allowing users to make more informed decisions or apply custom rate limiting rules
- Fallback Behavior: Define fallback behaviors for when the rate limiting storage (e.g., Redis, Memcached) becomes unavailable, such as allowing a default number of requests or completely blocking further requests
- Webhooks and Notifications: Allow users to configure webhooks or other notification mechanisms to be triggered when certain rate limiting thresholds are reached or exceeded, enabling them to take appropriate actions
- Grouping and Inheritance: Implement a mechanism for grouping rate limiting policies and allowing inheritance, so that users can define common rate limiting rules and apply them to multiple endpoints or clients
- Logging and Debugging: Provide robust logging and debugging capabilities, including the ability to log detailed information about rate limiting events, such as the request details, the applied rate limit, and the reason for any rejections
- HTTP/gRPC Integrations: Provide seamless integration with both HTTP-based and gRPC-based services, allowing users to apply rate limiting to their entire application stack
- Observability and Monitoring: Integrate with popular observability and monitoring tools (e.g., Prometheus, Grafana) to provide detailed metrics and dashboards for monitoring the rate limiting system's performance and health
- Backpressure Management: Implement mechanisms to handle backpressure, such as queuing and throttling, to ensure that the rate limiting system can gracefully handle sudden traffic spikes without causing cascading failures

## License

This project is licensed under the MIT License.
