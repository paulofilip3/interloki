# Phase 4 Design: Fluent Bit Forward Source + S3 Persistence

## Context

Phase 4 was originally planned as Loki integration. After evaluation, Loki is unnecessary — interloki can receive logs directly from Fluent Bit sidecars and handle persistence itself.

## Architecture

### Multi-Tenant SaaS Topology

```
Platform namespace (1)
  Platform backend + frontend
  Frontend proxies WebSocket to user interloki via backend
  Backend resolves interloki.{namespace}.svc.cluster.local

User namespace (thousands)
  5-15 pods with Fluent Bit sidecars (stdout + file logs)
  1 interloki pod: receives Forward protocol, persists to S3, serves WebSocket UI
```

### Data Flow

```
Fluent Bit sidecars --Forward (msgpack/TCP)--> interloki:24224
                                                |
                                          pipeline (parse+enrich)
                                                |
                                          +-----+-----+
                                          |           |
                                    ring buffer    S3 writer
                                    (real-time)    (batch flush)
                                          |
                                    WebSocket --> browser
                                          |
                                    /api/history --> S3 reader (on scroll-back)
```

### Decisions Made

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Fluent Bit relation | Runs separately (DaemonSet/sidecar) | Battle-tested k8s pattern, rich input/filter ecosystem |
| Log transport | Forward protocol (msgpack/TCP) | Native Fluent Bit output, efficient binary encoding |
| S3 writer | Interloki writes to S3 | Frontend connects to interloki directly; interloki is long-lived, Fluent Bit sidecars are ephemeral |
| S3 purpose | Archival + historical browsing | Users scroll back through time; no full-text search needed |
| S3 format | Time-partitioned gzipped JSON | Natural fit for "scroll back in time", cheap sequential reads |
| Platform access | Backend proxies WebSocket | Keeps interloki unexposed, auth/authz centralized in platform |
| S3 prefix | Per-namespace | Each interloki manages its own prefix: s3://bucket/{namespace}/... |

### Fluent Bit Forward Protocol

The Forward protocol uses msgpack over TCP. Three message modes:

1. **Message mode**: `[tag, time, record]` — single event
2. **Forward mode**: `[tag, [[time, record], [time, record], ...]]` — batch of events
3. **PackedForward mode**: `[tag, msgpack_bin]` — binary-packed batch (most common in practice)

The `record` is a map with the log line plus Fluent Bit metadata (kubernetes pod name, namespace, container, labels).

### S3 Storage Layout

```
s3://bucket/{prefix}/
  2026/03/05/14/
    chunk-1709647200000.json.gz
    chunk-1709647210000.json.gz
  2026/03/05/15/
    chunk-1709650800000.json.gz
```

Each chunk is a gzipped JSON array of LogMessage objects. Chunks are created every 10 seconds or 1000 messages (whichever comes first).

### History API

When the user scrolls past the ring buffer's oldest message:

1. Frontend calls `GET /api/history?before={ISO_timestamp}&count=500`
2. Server lists S3 objects in reverse chronological order
3. Downloads and decompresses relevant chunks
4. Returns up to 500 messages older than the given timestamp

### New Go Dependencies

- `github.com/vmihailenco/msgpack/v5` — msgpack decoding for Forward protocol
- `github.com/aws/aws-sdk-go-v2` — S3 client (+ config, credentials, s3 service packages)

### New Files

```
internal/source/forward.go          # Forward protocol TCP server
internal/source/forward_test.go     # Tests with mock msgpack client
internal/storage/storage.go         # Storage interface
internal/storage/s3.go              # S3 implementation
internal/storage/writer.go          # Background flush goroutine
internal/storage/s3_test.go         # Tests (localstack or mock)
```

### Config Additions

```go
type S3Config struct {
    Bucket        string        // S3 bucket name
    Prefix        string        // Key prefix (typically namespace name)
    Region        string        // AWS region
    Endpoint      string        // Custom endpoint (MinIO, localstack)
    FlushInterval time.Duration // Flush interval (default 10s)
    FlushCount    int           // Flush count threshold (default 1000)
}
```

S3 storage is optional — if no bucket is configured, interloki works without persistence (ring buffer only, same as today).
