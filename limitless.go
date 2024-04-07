package limitless

// features to be implemented:
// Middleware Integration: Develop a middleware component that can be easily integrated into popular Golang web frameworks (e.g., Gin, Echo, Fiber) to apply the rate limiting functionality
// Metrics and Reporting: Include options to expose metrics related to the rate limiting, such as the number of requests processed, the number of requests rejected due to rate limiting, and the current rate limiting state
// Burst Capacity: Allow users to configure a burst capacity, which is the maximum number of requests that can be processed immediately, even if it exceeds the rate limit. This can be useful for handling sudden spikes in traffic.
// Adaptive Rate Limiting: Implement an adaptive rate limiting algorithm that can dynamically adjust the rate limit based on factors such as the current load, resource utilization, or any other custom metrics
// Multi-Dimensional Rate Limiting: Provide the ability to apply rate limiting based on multiple dimensions, such as client IP, user identity, API endpoint, or any combination of these. This can help with more granular control over the rate limiting policies
// Contextual Information: Pass contextual information (e.g., request metadata, user identity) through the rate limiting system, allowing users to make more informed decisions or apply custom rate limiting rules
// Fallback Behavior: Define fallback behaviors for when the rate limiting storage (e.g., Redis, Memcached) becomes unavailable, such as allowing a default number of requests or completely blocking further requests
// Webhooks and Notifications: Allow users to configure webhooks or other notification mechanisms to be triggered when certain rate limiting thresholds are reached or exceeded, enabling them to take appropriate actions
// Grouping and Inheritance: Implement a mechanism for grouping rate limiting policies and allowing inheritance, so that users can define common rate limiting rules and apply them to multiple endpoints or clients
// Logging and Debugging: Provide robust logging and debugging capabilities, including the ability to log detailed information about rate limiting events, such as the request details, the applied rate limit, and the reason for any rejections
// HTTP/gRPC Integrations: Provide seamless integration with both HTTP-based and gRPC-based services, allowing users to apply rate limiting to their entire application stack
// Observability and Monitoring: Integrate with popular observability and monitoring tools (e.g., Prometheus, Grafana) to provide detailed metrics and dashboards for monitoring the rate limiting system's performance and health
// Backpressure Management: Implement mechanisms to handle backpressure, such as queuing and throttling, to ensure that the rate limiting system can gracefully handle sudden traffic spikes without causing cascading failures

// implements token bucket algorithm

type RateLimiter interface {
	allow() (*bool, error)
	load() error
	refill()
	available() *bool
	exhaust() error
}

func Allow(rateLimiter RateLimiter) (*bool, error) {
	return rateLimiter.allow()
}
