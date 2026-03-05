 Plan to implement 
 
 Interloki - Implementation Plan 
 
 Context 
 
 Build an open-source logs viewer/explorer/analyzer. Go backend + Vue.js frontend monorepo. Inspired by logdy-core
 (UI/architecture) and built on the github.com/paulofilip3/pipeline library (backend processing).

 Designed for multi-tenant SaaS: each user namespace runs its own interloki pod receiving logs from
 Fluent Bit sidecars via the Forward protocol, persisting to S3, and serving a real-time WebSocket UI.
 A platform namespace hosts the SaaS frontend which proxies WebSocket connections to per-namespace interloki instances.
 Also supports local use via stdin/file/socket/demo modes.

 Key decisions

 - Web Component (<interloki-viewer>) for embedding in the SaaS platform frontend
 - Log sources: Fluent Bit Forward protocol + stdin + file follow + TCP socket + demo
 - S3 persistence: time-partitioned gzipped JSON files, read back for historical browsing
 - Monorepo: Go at root, Vue in /web, frontend embedded via go:embed
 - Vue 3 + Vite + TypeScript + Pinia
 - Dual light/dark themes via CSS custom properties, controllable by host page 
 
 --- 
 Project Structure 
 
 interloki/ 
 ├── cmd/interloki/main.go # Cobra CLI entry point 
 ├── internal/ 
 │ ├── config/config.go # Config struct, env/flag parsing 
 │ ├── models/message.go # LogMessage type (pipeline item

 │ ├── source/
 │ │ ├── source.go              # Source interface
 │ │ ├── forward.go             # Fluent Bit Forward protocol server (msgpack over TCP)
 │ │ ├── stdin.go               # Stdin line reader
 │ │ ├── file.go                # File follow (nxadm/tail)
 │ │ ├── socket.go              # TCP listener
 │ │ └── demo.go                # Fake data generator
 │ ├── processing/
 │ │ ├── pipeline.go            # Pipeline construction (parse + enrich stages)
 │ │ ├── parse.go               # JSON detection stage worker
 │ │ └── enrich.go              # Timestamp/ID/level extraction stage worker
 │ ├── buffer/ring.go           # Thread-safe generic ring buffer
 │ ├── storage/
 │ │ ├── storage.go             # Storage interface (Write chunks, Read ranges)
 │ │ ├── s3.go                  # S3 backend (time-partitioned gzipped JSON)
 │ │ └── writer.go              # Background flush goroutine (ring buffer -> S3)
 │ ├── server/
 │ │ ├── server.go              # HTTP server + routes
 │ │ ├── handlers.go            # REST API handlers (+ /api/history for S3 reads)
 │ │ ├── websocket.go           # WebSocket handler + client manager
 │ │ ├── client.go              # Per-client buffering + flush loop
 │ │ └── embed.go               # go:embed frontend assets
 │ └── app/app.go               # Orchestrator: source -> pipeline -> buffer -> storage -> server 
 ├── web/ 
 │ ├── vite.config.ts # Dual build: standalone app + web component lib 
 │ ├── src/ 
 │ │ ├── main.ts # Standalone app mount 
 │ │ ├── main-wc.ts # Web component registration (defineCustomElement
 │ │ ├── App.vue # Root standalone layout 
 │ │ ├── components/ 
 │ │ │ ├── LogViewer.vue # Core log viewer (shared by both outputs
 │ │ │ ├── LogRow.vue # Single row in virtual scroll 
 │ │ │ ├── LogDetail.vue # Expanded row: JSON tree + raw content 
 │ │ │ ├── HeaderBar.vue # Top bar (standalone only
 │ │ │ ├── Sidebar.vue # LogQL editor + filters (standalone only
 │ │ │ ├── StatusBar.vue # Connection status + message count 
 │ │ │ ├── SearchBar.vue # Text/regex filter 
 │ │ │ ├── ColumnConfig.vue # Column visibility/order 
 │ │ │ ├── LogQLEditor.vue # Monaco LogQL editor 
 │ │ │ └── ThemeToggle.vue # Light/dark switch 
 │ │ ├── composables/ 
 │ │ │ ├── useWebSocket.ts # WS connection + reconnect 
 │ │ │ └── useTheme.ts # Theme switching 
 │ │ ├── stores/ 
 │ │ │ ├── logs.ts # Log messages + filtering state 
 │ │ │ ├── connection.ts # WebSocket state 
 │ │ │ ├── settings.ts # Display preferences (persisted
 │ │ │ └── query.ts # LogQL query state 
 │ │ ├── types/ # TS types mirroring Go models 
 │ │ └── styles/ 
 │ │ ├── variables.css # CSS custom properties (theme tokens
 │ │ └── themes.css # Light/dark definitions 
 ├── Makefile 
 ├── Dockerfile 
 ├── docker-compose.yml # interloki + Loki for dev/testing 
 ├── .pre-commit-config.yaml 
 ├── .gitignore 
 ├── go.mod 
 ├── CLAUDE.md 
 ├── LICENSE 
 └── README.md 
 
 --- 
 Backend Architecture 
 
 Core Types 
 
 LogMessage (internal/models/message.go) — the pipeline item type T
 type LogMessage struct { 
 ID string            `json:"id"
 Content string            `json:"content"
 JsonContent json.RawMessage   `json:"json_content,omitempty"
 IsJson bool              `json:"is_json"
 Timestamp time.Time         `json:"ts"
 Source SourceType        `json:"source"`      // "loki"|"stdin"|"file"|"socket"|"demo
 Origin Origin            `json:"origin"
 Labels map[string]string `json:"labels,omitempty"
 Level string            `json:"level,omitempty"
 } 
 
 Source interface (internal/source/source.go)
 type Source interface { 
 Name() string 
 Start(ctx context.Context) (<-chan models.LogMessage, error
 } 
 
 Data Flow (using pipeline library)

 [Sources: forward, stdin, file, socket, demo]
     |              |             |
     v              v             v
   <-chan LogMessage (per source)
     \              |            /
      pipeline.Merge(ctx, ...)        <-- fan-in
               |
     pipeline.New(parseStage, enrichStage)
     pipe.Run(ctx, merged)
               |
        Stage "parse": JSON detect via fastjson.Validate()
               |
        Stage "enrich": assign ID, timestamp, extract level
               |
          out <-chan LogMessage     errs <-chan error
               |                         |
     ClientManager.consumeLoop()    log errors
        - ring.Push(msg)
        - storage.Writer.Enqueue(msg)   <-- S3 persistence
        - distribute to all clients
               |
     Per-client flush loop (100ms batch window)
               |
     WebSocket -> Browser

     Storage Writer (background goroutine):
        - Buffers messages, flushes to S3 every 10s or 1000 messages
        - Writes as: {prefix}/{YYYY}/{MM}/{DD}/{HH}/chunk-{timestamp}.json.gz
        - Each chunk is a gzipped JSON array of LogMessage

     S3 Reader (on-demand):
        - Frontend requests /api/history?before={timestamp}&count=500
        - Server lists S3 objects by time prefix, reads chunks in reverse chronological order
        - Returns messages older than the ring buffer's oldest entry

 Both pipeline stages use Concurrency: 1 to preserve message ordering (critical for log viewing).

 Fluent Bit Forward Source

 - Implements the Fluent Bit Forward protocol (msgpack over TCP)
 - Listens on a configurable port (default 24224)
 - Decodes msgpack events: [tag, time, record] and [tag, [[time, record], ...]]
 - Extracts Fluent Bit metadata (tag, kubernetes labels) into LogMessage.Labels
 - Each connected Fluent Bit sidecar is a separate TCP connection
 
 WebSocket Protocol 
 
 Bidirectional (improvement over logdy's REST+WS split)

 
 ┌───────────┬───────────────┬───────────────────────────────────────────┐ 
 │ Direction │ Type │ Purpose │ 
 ├───────────┼───────────────┼───────────────────────────────────────────┤ 
 │ S->C │ client_joined │ Connection confirmed + config │ 
 ├───────────┼───────────────┼───────────────────────────────────────────┤ 
 │ S->C │ log_bulk │ Batched messages + stats (every ~100ms)   │ 
 ├───────────┼───────────────┼───────────────────────────────────────────┤ 
 │ S->C │ status │ Periodic stats update │ 
 ├───────────┼───────────────┼───────────────────────────────────────────┤ 
 │ C->S │ set_status │ Pause/resume (following/stopped)          │ 
 ├───────────┼───────────────┼───────────────────────────────────────────┤ 
 │ C->S │ load_range │ Request historical range from ring buffer │ 
 ├───────────┼───────────────┼───────────────────────────────────────────┤ 
 │ C->S │ ping │ Keepalive │ 
 └───────────┴───────────────┴───────────────────────────────────────────┘ 
 
 REST API kept for programmatic access: /api/status, /api/client/load, /api/history.

 CLI (Cobra)

 interloki forward --listen=:24224 [--s3-bucket, --s3-prefix, --s3-region, --s3-flush-interval]
 interloki stdin
 interloki follow --file=/var/log/*.log
 interloki socket --port=9999
 interloki demo   [--rate=100]

 Global flags: --port, --ip, --ui-pass, --max-messages, --bulk-window-ms, --verbose
 S3 flags: --s3-bucket, --s3-prefix, --s3-region, --s3-endpoint, --s3-flush-interval, --s3-flush-count
 Env vars: INTERLOKI_ prefix (e.g., INTERLOKI_PORT, INTERLOKI_S3_BUCKET)
 
 --- 
 Frontend Architecture 
 
 Web Component API 
 
 <interloki-viewer 
 ws-url="ws://localhost:8080/ws
 theme="dark
 auto-follow="true
 show-search="true
 show-columns="timestamp,level,source,content
 height="600px
 ></interloki-viewer
 
 JS API: connect(), disconnect(), pause(), resume(), clear(), setTheme(), setFilter(), exportLogs(
 Events: interloki:connected, interloki:disconnected, interloki:message, interloki:error, interloki:row-click 
 
 CSS custom properties for theming: --interloki-bg, --interloki-fg, --interloki-accent, --interloki-level-*
-interloki-font-family, --interloki-row-height, etc. 
 
 Component Hierarchy 
 
 Standalone wraps in full layout (HeaderBar, Sidebar with LogQL editor, LogViewer, StatusBar). 
 Web Component wraps minimal shell (SearchBar, LogViewer, StatusBar) in shadow DOM. 
 Both share the same LogViewer core. 
 
 --- 
 Implementation Phases 
 
 Phase 1: Skeleton + Stdin + Basic Viewer (MVP
 
 Goal: echo "hello" | interloki stdin opens browser with real-time log display. 
 
 1. Init Go module, directory structure, .gitignore, Makefile, .pre-commit-config.yaml 
 2. internal/models/message.go — LogMessage type 
 3. internal/source/source.go + stdin.go — Source interface + stdin reader 
 4. internal/processing/ — pipeline with parse + enrich stages (import github.com/paulofilip3/pipeline
 
 5. internal/buffer/ring.go — thread-safe ring buffer 
 6. internal/server/ — HTTP server, WebSocket handler, client management, 100ms flush loop 
 7. internal/app/app.go — orchestrator wiring 
 8. cmd/interloki/main.go — Cobra root + stdin subcommand 
 9. Scaffold Vue project (Vite + TS + Pinia
 
 10. useWebSocket.ts composable + connection.ts store 
 11. logs.ts store + LogViewer.vue + LogRow.vue with virtual scrolling 
 12. StatusBar.vue + basic styling 
 13. go:embed integration + Makefile build target 
 14. End-to-end test: cat /var/log/syslog | ./bin/interloki stdin 
 
 Phase 2: All Non-Loki Sources + Control 
 
 Goal: file follow, socket, demo modes + pause/resume. 
 
 1. internal/source/file.go, socket.go, demo.go 
 2. Cobra subcommands: follow, socket, demo 
 3. internal/config/config.go — full config parsing 
 4. Bidirectional WS messages (pause/resume, load range
 
 5. REST API endpoints 
 6. Tests for all sources 
 
 Phase 3: Frontend Polish + Theming 
 
 Goal: Production-quality UI. 
 
 1. SearchBar.vue — text + regex filtering 
 2. LogDetail.vue — JSON tree viewer + raw content 
 3. CSS custom properties + themes.css (light/dark
 
 4. ThemeToggle.vue + useTheme.ts 
 5. HeaderBar.vue + Sidebar.vue (standalone layout
 
 6. ColumnConfig.vue — column visibility/reorder 
 7. Virtual scrolling (vue-virtual-scroller
 
 8. Settings persistence (localStorage
 
 Phase 4: Fluent Bit Forward Source + S3 Persistence

 Goal: interloki forward --listen=:24224 --s3-bucket=my-logs --s3-prefix=ns1/
 Fluent Bit sidecars send logs via Forward protocol, interloki persists to S3, frontend scrolls back through S3 history.

 k8s Architecture:
   Platform namespace (1):
     - Platform backend + frontend
     - Frontend proxies WebSocket to user interloki via backend
     - Backend resolves interloki.{namespace}.svc.cluster.local

   User namespace (thousands):
     - 5-15 pods with Fluent Bit sidecars (stdout + file logs)
     - 1 interloki pod: receives Forward, persists S3, serves UI
     - Fluent Bit config: output forward to interloki.{namespace}.svc.cluster.local:24224

 Backend tasks:
 1. internal/source/forward.go — Forward protocol server (msgpack TCP)
    - Accept multiple Fluent Bit connections concurrently
    - Decode msgpack: Message mode [tag, time, record], Forward mode [tag, [entries...]]
    - Extract tag, kubernetes labels/annotations from record into LogMessage.Labels
    - Cobra subcommand: `interloki forward --listen=:24224`
 2. internal/storage/storage.go — Storage interface
    - Writer: Enqueue(msg), Start(ctx), flush loop
    - Reader: ReadBefore(timestamp, count) -> []LogMessage
    - ListChunks(before timestamp) -> chunk metadata
 3. internal/storage/s3.go — S3 implementation
    - Write: gzipped JSON arrays to {prefix}/{YYYY}/{MM}/{DD}/{HH}/chunk-{unix_ms}.json.gz
    - Read: list objects by prefix (reverse chronological), download and decompress
    - Config: bucket, prefix, region, endpoint (for MinIO/localstack), flush interval (10s), flush count (1000)
    - Uses AWS SDK v2 (github.com/aws/aws-sdk-go-v2)
 4. internal/storage/writer.go — Background flush goroutine
    - Buffers messages from ConsumeLoop, flushes to S3 on interval or count threshold
    - Graceful shutdown: flush remaining on context cancellation
 5. Wire S3 storage into app.go orchestrator
    - ConsumeLoop enqueues to both ring buffer and storage writer
    - S3 flags added to config.go
 6. /api/history endpoint — serves historical logs from S3
    - Query params: before (ISO timestamp), count (int, default 500)
    - Reads from S3 in reverse chronological order
    - Returns messages older than what the ring buffer holds
 7. Frontend: update connection store to fall back to /api/history when ring buffer is exhausted
    - When scrolling up past ring buffer, fetch from /api/history?before={oldest_ts}&count=500
 8. Tests: Forward source tests (mock msgpack client), S3 storage tests (localstack or mocked)

 Phase 5: Web Component Build

 Goal: <interloki-viewer> works as drop-in custom element, embeddable in the SaaS platform frontend.
 The platform frontend opens WebSocket to its backend, which proxies to the user's interloki pod.

 1. main-wc.ts — defineCustomElement wrapper
 2. Vite dual build config (app mode + lib mode)
 3. Shadow DOM encapsulation + CSS property passthrough
 4. Attribute/property bindings (ws-url, theme, auto-follow), JS methods, Custom Events
 5. Example embedding page + multi-instance test
 6. NPM package publish config for @interloki/viewer

 Phase 6: k8s Deployment + DevOps + Release

 Goal: Helm chart, CI/CD, Docker, docs, v0.1.0.

 1. Dockerfile (multi-stage: Node build -> Go build -> Alpine)
 2. Helm chart for user-namespace interloki deployment
    - Deployment: 1 replica, interloki forward mode
    - Service: ClusterIP, ports 8080 (HTTP/WS) + 24224 (Forward)
    - ConfigMap: S3 bucket/prefix/region from namespace
    - ServiceAccount + IAM role for S3 access (IRSA or equivalent)
 3. Fluent Bit sidecar config template
    - Input: tail (stdout/stderr logs + application log files)
    - Output: forward to interloki.{namespace}.svc.cluster.local:24224
 4. .github/workflows/ci.yml
 5. README.md with installation, quick start, k8s guide, embedding guide, config reference
 6. CLAUDE.md
 7. LICENSE
 8. Tag v0.1.0

 ---
 Key Dependencies

 Go:
 - github.com/paulofilip3/pipeline — core processing
 - github.com/spf13/cobra — CLI
 - github.com/gorilla/websocket — WebSocket server
 - github.com/valyala/fastjson — fast JSON validation
 - github.com/nxadm/tail — file following
 - github.com/sirupsen/logrus — logging
 - github.com/vmihailenco/msgpack/v5 — msgpack decoding (Forward protocol)
 - github.com/aws/aws-sdk-go-v2 — S3 client

 Frontend:
 - vue@3, pinia
 - (no monaco-editor — LogQL editor removed)

 ---
 Verification

 1. Phase 1 test: echo "hello world" | ./bin/interloki stdin — browser at :8080 shows the message
 2. Phase 2 test: ./bin/interloki demo --rate=100 — 100 msgs/sec streaming in browser, pause/resume works
 3. Phase 3 test: Theme toggle, search filter, JSON expand, column reorder all work
 4. Phase 4 test: Fluent Bit container sends logs via Forward to interloki, logs appear in browser and in S3; scroll up past ring buffer loads from S3
 5. Phase 5 test: Host page with <interloki-viewer ws-url="...">, theme controlled from host, events fired
 6. Phase 6 test: Helm install in k8s, Fluent Bit sidecars deliver logs, interloki serves them via proxied WebSocket from platform namespace
