# AGENTS.md

Compact cheat-sheet for OpenCode sessions. `CLAUDE.md` has the full architecture,
request flow, and layer responsibilities — read it before non-trivial work. This
file only lists things you will get wrong without being told.

## Clean clone does NOT compile

Files required for compilation are gitignored (proprietary upstream data):
`protocol/mx/excel.go`, `gdconf/game.config.dev.go`, `common/server_only/*.pb.go`,
`common/server_only/*.py`, `/mx`, `/protocol/proto2`, `/resources/*`. A fresh
`go build` will fail until protos are generated and resources fetched. Code paths
referencing `mx.DeExcelBytes` etc. only work once those private files are restored.

## Required build order

```bash
# 1. Generate gitignored protobuf code (bundled protoc binaries in repo)
cd common/server_only && ./protoc --proto_path=. --plugin=protoc-gen-go=./protoc-gen-go --go_out=. *.proto

# 2. Fetch Excel resources — Linux/WSL2 only (needs BlueArchiveLocalizationTools)
./fetch_resources.sh

# 3. Build data/Excel.bin. MUST use the generate_excel build tag or it won't compile.
go run -tags=generate_excel generate_excel.go

# 4. Build server. CGO_ENABLED=0 (cross-compile matrix in build.sh / build.bat).
go build -ldflags="-s -w" -o ./bin/BaPs ./cmd/BaPs/BaPs.go
```

If `data/Excel.bin` is missing at boot, `gdconf.LoadExcel` (`gdconf/game.config.rel.go:27`)
tries to regenerate from `resources/Excel/` + `resources/ExcelDB/`. If those dirs
are absent, boot panics.

## Run / config

- `./bin/BaPs -g` writes a default `config.json` and exits; `-c <path>` selects it.
- Every config field is overridable via env vars using a dotted-path key based on
  `json` tags (e.g. `Config.HttpNet.InnerPort=5050`, `Config.DB.dsn=...`). See
  `overrideWithEnv` (`config/config.go:251`). Applies to nested pointer structs.
- Default `GucooingApiKey` is a 32-char hex random string from `crypto/rand` — the
  `_ int64` seed arg is ignored, do not treat it as deterministic.
- `OuterAddr` in `HttpNet` MUST match the externally reachable URL; `prod/index.json`
  and SDK responses embed it.
- DB: `sqlite` (ncruces, CGO-free) or `mysql`, both via GORM.

## Build tags

- `dev` — wire-level message logging in gateway (`gateway_dev.go`); without it
  `gateway_rel.go` is no-op stubs.
- `debug` — selects gitignored `gdconf/game.config.dev.go`; release is default.
- `generate_excel` — only on the Excel.bin builder; keeps it out of the server binary.

## Tests

Only `gdconf` has tests/benchmarks — they verify `Excel.bin` loads:
```bash
go test ./gdconf/...
go test -bench=. ./gdconf/...
```
No other packages have tests. Don't assume a repo-wide `go test` is meaningful.

## Do not change without asking

- The release banner at `baps.go:94-100` is load-bearing — do not remove.
- `go.mod` `replace` pins `github.com/gucooing/cdq` → `github.com/asfu222/cdq v1.0.6-safe`
  to enforce the safe fork (original had a remote-shell backdoor). Keep the pin.
- `pkg/build.go` `ClientVersion` / `ServerVersion` are version pins.

## Architecture rules you'd otherwise invert

- `funcRouteMap` (`gateway/gateway.router.go:20`) is the protocol dispatch table
  (`proto.Protocol` → handler). Add new protocols here, implement in `pack/`.
  Chinese comments next to entries document message intent — preserve them.
- Layer deps: `pack/` (protocol handlers) → `game/` (domain logic on
  `Session.PlayerBin`) → `db/` (only via `DBGame` interface). `pack/` must not
  touch the DB directly. Do not invert.
- Player state is NEVER persisted in real time. Mutations hit in-memory
  `Session.PlayerBin`; disk writes happen on session timeout / shutdown
  (`common/enter/session.go:50`, `baps.go:122`). Real-time player views go through
  the command API, not the DB.
- `gdconf.GC` is the global read-only game-data singleton (`Excel` raw +
  `GPP` post-processed indexes). One `exceldb.*.go` / `excel.*.go` per table.
- `protocol/proto/*.go` are hand-written wrappers, NOT generated from `.proto`.

## Build entry points

- `cmd/BaPs/BaPs.go` — standalone executable (`main` → `BaPs.NewBaPs`).
- `cmd/Dll/BaPs.go` — c-shared library (exports `StartServer`); needs
  `-buildmode=c-shared`.
- `generate_excel.go` — standalone tool, gated by `//go:build generate_excel`.

## Other

- No 32-bit support (per README); don't try to fix arm/x86 32-bit builds.
- Go 1.23.2 (pinned in `go.mod` and CI).
