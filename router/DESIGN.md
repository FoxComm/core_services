# Router

The Router will be a critical piece of our infrastructure as **every single request and response** from outside clients will go thru its pipeline. Additionally, it will be responsible for such things as (non-exhaustive list and may not represent end result):

- session handling/transformation/expiration
- authorization (better to short-circuit immediately than offload to downstream services)
- rate limiting
- downstream request retries
- forwarding/teeing: forward request/response to another service for post-processing which does not affect response to client

## Current Design

The current Router is based off https://github.com/mailgun/vulcan which is not only deprecated but is quite buggy and not feature complete.

Known and *presumed* issues with current:

- AFAIK, 502s from downstreams are not retried with a different downstream in a given pool
- Insufficient instrumentation
- Panics which should be avoided or recovered

## New Design

@narmaru will research possible alternatives, weigh the pros/cons, and choose a better technology to serve as the foundation of our custom Router. Thus far, it looks like mailgun's newer https://github.com/mailgun/oxy represents a good and stable alternative.

Features which will exist in the new Router (non-exhaustive):

- Structured logging of request/response lifecycle
- Load-balancing algorithms (potentially more than Roundrobin such as leastconns)
- Rate-limiting
- Request/response forwarding (for post-processing via downstream services)
- Circuit breakers to handle 100% downstream timeouts, queues reaching capacity, etc.
- Full instrumentation
- Authorization
- Session handling
- Improved caching of configuration (e.g. store features or something queries PSQL for every single request. let's cache some data here)
- Load service endpoints and their configuration from a TOML config file
- Zero-downtime restarts/reload
- Binary upgrades
- No Panics :)

Furthermore, a consideration for the new design should be to reduce
allocations in our custom code, *as much as possible*, to mitigate
heap growth and GC pausing.

## Instrumentation

The new Router should implement both a push/pull metrics system whereby
pull would be an HTTP endpoint and takes developer priority over push.
Push shall be implemented later once we have metrics/graphing/monitoring
systems online.

The initial pull method could use go's https://golang.org/pkg/expvar/
for a quick implementation.

#### Logging

- Flag(s) to toggle full DEBUG logging of all HTTP requests/responses including middleware transformations.
- An HTTP endpoint to toggle the above flag for production environment crisis monitoring (AKA when we don't know what the hell is going on)
- Structured logging (logrus or something similar) which can ouptut to both stdout/rsyslog.

#### Metrics

- Metrics (my preference is https://github.com/rcrowley/go-metrics, but research your own) which record status codes & response latencies of request/response lifecycles keyed on downstream services.
- Global # requests/responses and statuses
- Internal metrics such as CPU, MEM profile, num goroutines, GC info (https://golang.org/pkg/runtime/)

## Testing

The new Router should be written with a TDD approach or, at a minimum,
adequate testing coverage from the start.

Additionally, we should build a test environment that allows us to load, stress, and capacity test the Router.

For stress test we can use next tools

- https://github.com/machinezone/tcpkali
- https://github.com/wg/wrk

## Resources

Here are some resources for reading if you want to learn about some of
these design considerations:

Go:

http://peter.bourgon.org/go-in-production/
https://blog.cloudflare.com/recycling-memory-buffers-in-go/

Circuit breakers:

http://techblog.netflix.com/2011/12/making-netflix-api-more-resilient.html

Metrics:

http://www.mikeperham.com/2014/12/17/expvar-metrics-for-golang/
https://www.youtube.com/watch?v=czes-oa0yik
http://www.slideshare.net/IzzetMustafaiev/metrics-by-coda-hale

Structured Logging:

http://blog.sematext.com/2013/05/28/structured-logging-with-rsyslog-and-elasticsearch/
http://gregoryszorc.com/blog/category/logging/
