# llama-go Benchmark Report

## Test Environment
- Date: {date}
- Hardware: {cpu}, {memory}
- Model: {model}
- Backend: {backend}

## Test Scenarios

### Scenario 1: Non-Stream Short Text
- Concurrent Users: 10-50
- Request Rate: {rps}
- Average Latency: {latency_ms}ms
- P95 Latency: {p95_ms}ms
- P99 Latency: {p99_ms}ms
- TTFT: {ttft_ms}ms
- Error Rate: {error_rate}%

### Scenario 2: Stream Short Text
- Concurrent Users: 10-50
- Average Latency: {latency_ms}ms
- Stream Duration: {duration_ms}ms
- Tokens/sec: {tokens_per_sec}
- Error Rate: {error_rate}%

### Scenario 3: Multi-Session Concurrent
- Concurrent Sessions: 50-100
- Average Latency: {latency_ms}ms
- Cache Hit Rate: {cache_hit_rate}%
- Error Rate: {error_rate}%

## Resource Usage
- CPU: {cpu_usage}%
- Memory: {memory_usage}MB
- Goroutines: {goroutine_count}

## Recommendations
- Max Concurrent Requests: {max_concurrent}
- Recommended Rate Limit: {rate_limit} RPS
- Recommended Timeout: {timeout}s
- Instance Count: {instance_count}
