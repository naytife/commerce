## Service communication resilience plan (backend ↔ store-deployer, template-registry)

## Purpose

This document defines a concrete, prioritized implementation plan to improve resilience of the communication paths between the backend and the two deployment-oriented services: `store-deployer` and `template-registry`.

## Scope

- Focus: backend <-> `store-deployer` and backend <-> `template-registry`.
- Exclusions (per request): authentication and service identity/security hardening are intentionally excluded from this plan.

## Assumptions

- The repository contains: `backend/internal/api/handlers/proxy.handlers.go`, `backend/internal/api/handlers/template.handlers.go`, `backend/internal/api/handlers/shop.handlers.go`, and `services/store-deployer/main.go` which implement the calls and proxies described.
- Kubernetes overlay sets `TEMPLATE_REGISTRY_URL` and `STORE_DEPLOYER_URL`, defaulting to `http://template-registry:9001` and `http://store-deployer:9003`.

## Goal / Success criteria

- Improve availability and robustness of calls between backend and these services under transient failures and load.
- Reduce cascading failures (faster failures, soft degrade) and make critical flows retry-safe where appropriate.
- Provide observable signals (traces/metrics/logs) for cross-service debugging.

## Checklist (extracted requirements)

- Replace per-request HTTP clients and `http.Get`/`http.Post` with a shared, tuned `*http.Client` used across backend and store-deployer where appropriate.
- Ensure request contexts and deadlines are used (propagate incoming context or derive a bounded context).
- Add retry + exponential backoff for transient errors (with method and idempotency awareness).
- Add a circuit-breaker around calls to `store-deployer` and `template-registry`.
- Make `deploy` operations retry-safe: either asynchronous job model (preferred) or idempotency tokens for POSTs.
- Add basic distributed tracing (OpenTelemetry) and structured logs including trace/request ids.
- Add metrics (latency, error rate, circuit-breaker state) and define SLOs/alerts.
- Add integration tests that simulate transient downstream failures and validate retry/circuit behavior.

## Prioritized implementation plan

## Phase P0 — Immediate, high-impact (days)

1. Shared, tuned HTTP client

   - Implement `backend/internal/httpclient` exposing a `DefaultClient` (configured Transport: DialTimeouts, TLSHandshakeTimeout, IdleConnTimeout, MaxIdleConnsPerHost) and helper to build requests from a context.
   - Replace direct uses of `http.Get`/`http.Post` and per-request `&http.Client{}` instantiations in:
     - `backend/internal/api/handlers/template.handlers.go`
     - `backend/internal/api/handlers/proxy.handlers.go`
     - `backend/internal/api/handlers/shop.handlers.go`
     - `services/store-deployer/main.go` (calls to template-registry / GraphQL)
   - Acceptance: no `http.Get`/`http.Post` calls remain for internal service calls; all use context-aware requests and the shared client.

2. Context-aware requests

   - Change call sites to use `http.NewRequestWithContext(ctx, ...)` where `ctx` is the incoming request context or a derived context with a bounded deadline.
   - Ensure the proxy uses the Fiber request context (or a derived context) so cancelation propagates to downstream calls.
   - Acceptance: requests cancel when client disconnects or context times out; unit tests verify cancellation.

3. Small retry wrapper for transient failures
   - Add a minimal retry helper `DoWithRetry(ctx, client, req)` with exponential backoff and jitter (configurable max attempts and backoff cap). Use it for idempotent GET/HEAD requests and safe re-tries where applicable.
   - Choose a lightweight library (suggest `github.com/cenkalti/backoff/v4`) or implement a 20-line wrapper.
   - Acceptance: retriable errors (temporary network failures, 502/503/504) are retried, with logs indicating retries.

## Phase P1 — Near-term resilience (1–2 sprints)

4. Circuit-breaker

   - Wrap calls to `store-deployer` and `template-registry` with a circuit-breaker (e.g., `sony/gobreaker`). Tune thresholds for error rate and window.
   - Behavior: after repeated failures, breaker opens and calls fail fast with a clear error; breaker transitions logged and exposed as a metric.
   - Acceptance: simulated downstream flapping causes breaker to open and subsequent calls bail out quickly.

5. Make `deploy` retry-safe (async job model) — recommended

   - Convert `POST /deploy` to return 202 Accepted with a `deployment_id` and process the actual work asynchronously (in `store-deployer` worker or using a small durable queue). Provide `GET /deployments/{id}/status` for progress.
   - Short-term alternative: require a client-supplied `Idempotency-Key` header for deploy POSTs and track keys to avoid duplicate deployments.
   - Acceptance: repeated POSTs with same idempotency key or deployment_id do not cause duplicated work; system surfaces a job ID quickly.

6. OpenTelemetry tracing + structured logs
   - Add OpenTelemetry SDK initialization in backend and store-deployer (sampling configurable). Propagate trace headers through proxy paths and GraphQL calls.
   - Ensure logs include trace-id and request-id.
   - Acceptance: traces show end-to-end request spanning backend -> template-registry/store-deployer, and logs contain trace ids.

## Phase P2 — Operational hardening (later)

7. Metrics and SLOs

   - Export Prometheus-compatible metrics: request duration/histograms, error counts, retry counts, circuit-breaker state, queue depth for async deploys.
   - Define SLOs (e.g., 95% successful deployment API responses within X seconds) and alerting rules.
   - Acceptance: dashboards and alerts exist in staging and trigger on violations.

8. Integration and chaos tests

   - Add integration tests that intentionally inject transient errors (downstream 502/timeout) to verify retry/backoff and circuit-breaker behavior.
   - Add a small chaos experiment in staging to simulate template-registry slowdowns and validate system behavior.
   - Acceptance: tests pass in CI; alerts/logging show expected behavior.

9. DB & request timeouts tuning
   - Review repository-level DB `context.WithTimeout(..., 1*time.Second)` usages. Increase timeouts to realistic values for queries that may take longer, and make timeouts configurable.
   - Tune connection pooling (pgx) and consider a connection pooler if many concurrent workers are expected.
   - Acceptance: under load, DB timeouts don't cause spurious failures; pool sizes sustain peak load.

## Implementation details and file targets (suggested edits)

- New helper: `backend/internal/httpclient/httpclient.go`

  - Exposes `DefaultClient *http.Client` and a `DoWithRetry(ctx, req) (*http.Response, error)` helper (retries for idempotent methods)

- Replace direct calls in these files to use helpers above:
  - `backend/internal/api/handlers/proxy.handlers.go` (use `NewRequestWithContext`, use `DefaultClient.Do(req)`)
  - `backend/internal/api/handlers/template.handlers.go` (replace `http.Get`/`http.Post` with shared client + `DoWithRetry` for GETs)
  - `backend/internal/api/handlers/shop.handlers.go` (calls to cleanup endpoints)
  - `services/store-deployer/main.go` (use shared client patterns for template-registry and backend GraphQL calls; use `DoWithRetry` for GETs)

## Engineering contract (small)

- Inputs: context + HTTP request to downstream service.
- Outputs: JSON payload or well-formed error (with status code and reason). Retries limited and observable.
- Error modes: network timeout, HTTP 5xx, JSON decode errors. Retries apply only to network/timeouts/5xx for idempotent reads.
- Success criteria: correct response or a fast-failing error after circuit-breaker threshold.

## Edge cases and guidance

- Do not automatically retry non-idempotent POSTs unless an idempotency key is supplied and honored.
- Keep retry attempts and total retry window bounded (e.g., max attempts 3, max total retry time 30s).
- When converting to async deploys, ensure the job queue is durable (simple DB table or lightweight queue) so workers can recover after restarts.

## Quality gates and tests

- Add unit tests for `DoWithRetry` that simulate transient failures.
- Add integration test that runs a backend -> proxy -> store-deployer deploy flow while the deployer first returns 503 then succeeds; verify retries and final success.
- Run full build/lint/test as part of CI; ensure no compile errors after refactor.

## Timeline & rough estimates

- P0 (shared client, context, small retry helper): 1–2 days.
- P1 (circuit-breaker, async deploy changes or idempotency, tracing): 1–2 sprints (2–4 weeks) depending on team bandwidth.
- P2 (metrics, SLOs, chaos testing, DB tuning): next 2–4 weeks.

## Quick wins (to implement in the next PR)

1. Add `backend/internal/httpclient/httpclient.go` with `DefaultClient` and a `DoWithRetry` helper.
2. Replace `http.Get`/`http.Post` usage in `backend/internal/api/handlers/template.handlers.go` with the helper.
3. Update `proxy.handlers.go` to create requests with `NewRequestWithContext` and use the shared client.

## How to validate locally (smoke tests)

- Start the backend and store-deployer services locally with Kubernetes/Skaffold or run them locally.
- Trigger a template list and a deploy via backend endpoints and observe logs and traces.
- Simulate a flaky `template-registry` (sleep then 200) and verify retries and circuit-breaker logs.

## Acceptance checklist (before merge)

- [ ] Shared client added and used by targeted files.
- [ ] Context-aware requests implemented for proxy and direct service calls.
- [ ] Retry helper added and used for idempotent GETs.
- [ ] Circuit-breaker implemented and tested in integration.
- [ ] Deploy endpoint made async or idempotency supported and tested.
- [ ] Basic OpenTelemetry traces appear end-to-end in dev/staging.
- [ ] Unit and integration tests added and passing.

## Appendix: minimal `DoWithRetry` behavior (reference)

- Input: `ctx`, `client`, `req`, `maxAttempts`, `initialBackoff`.
- Loop up to `maxAttempts`:
  - Call `client.Do(req)`
  - If success (2xx) return
  - If network error or 5xx, sleep backoff\*jitter and retry
  - If 4xx return immediately

## Notes

- This plan intentionally omits authentication and header-forgery mitigations as requested; those items should be planned separately because they significantly affect design choices (e.g., header whitelisting, mTLS).

## Contact

If you want I can implement the Quick wins PR (shared client + replace `http.Get`/`http.Post` in `template.handlers.go` and update `proxy.handlers.go`) and run unit/integration tests. Reply with "Implement quick wins" and I'll start the changes.
