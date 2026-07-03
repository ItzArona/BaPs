# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Context

BaPs is a Go server emulator for the Japanese release of Blue Archive (`ba-jp-sdk.bluearchive.jp`, `yostar-serverinfo.bluearchiveyostar.com`). This is the `asfu222/dev` fork of `gucooing/BaPs` with the remote-shell backdoor (`/cdq/api/shell`), bot-registration backdoor, and developer build key removed. The `go.mod` `replace` directive pins `github.com/gucooing/cdq` → `github.com/asfu222/cdq v1.0.6-safe` to enforce the safe variant.

Some files referenced from the source tree are gitignored because they contain proprietary keys/data from the original upstream and were never committed to this fork: `protocol/mx/excel.go`, `gdconf/game.config.dev.go`, `common/server_only/*.pb.go`, `common/server_only/*.py`, `/mx`, `/protocol/proto2`, `/resources/*`. Functions like `mx.DeExcelBytes` referenced in code paths only work in builds that have the upstream private files restored. Expect compile errors after a clean clone until protos are generated and resources fetched.

Server/client versions are pinned in `pkg/build.go` (`ClientVersion`, `ServerVersion`). The release banner in `baps.go:94-100` is load-bearing — do not remove it without discussion.

## Build & Run

```bash
# Generate the gitignored protobuf code FIRST (required for compile)
cd common/server_only && ./protoc --proto_path=. --plugin=protoc-gen-go=./protoc-gen-go --go_out=. *.proto

# Fetch Excel resources (Linux / WSL2 only — depends on BlueArchiveLocalizationTools)
./fetch_resources.sh

# Build Excel.bin from fetched resources. NOTE: the standalone generator uses a build tag.
go run -tags=generate_excel generate_excel.go

# Build the server binary (cross-compile matrix in build.sh / build.bat)
go build -ldflags="-s -w" -o ./bin/BaPs ./cmd/BaPs/BaPs.go

# Run. -c selects the config; -g writes a default config and exits.
./bin/BaPs -c ./config.json
./bin/BaPs -g                   # write ./config.json defaults

# Tests / benchmarks live only in gdconf — they verify Excel.bin loads.
go test ./gdconf/...
go test -bench=. ./gdconf/...
```

Build entry points:
- `cmd/BaPs/BaPs.go` — standalone executable (`main` → `BaPs.NewBaPs`)
- `cmd/Dll/BaPs.go` — c-shared library variant (exports `StartServer`); requires `-buildmode=c-shared`
- `generate_excel.go` — standalone tool, gated by `//go:build generate_excel`

Relevant build tags:
- `dev` — enables wire-level message logging in gateway (`gateway/gateway_dev.go`); without it `gateway_rel.go` provides no-op stubs
- `debug` — selects the debug variant of `gdconf/game.config*.go` (file is gitignored; release build is the default)
- `generate_excel` — only present on the Excel.bin builder, prevents it from compiling into the server

If `data/Excel.bin` is missing at startup, `gdconf.LoadExcel` (`gdconf/game.config.rel.go:27`) tries to regenerate it from `resources/Excel/` + `resources/ExcelDB/` JSON dumps. If those directories are absent, server boot panics.

## Configuration

`config.json` is loaded by `config.LoadConfig` (`config/config.go:230`). After file load, every field is overridable via env vars using a dotted-path key based on `json` tags — `Config.HttpNet.InnerPort=5050`, `Config.DB.dsn=...`, etc. This is implemented reflectively in `overrideWithEnv` (`config/config.go:251`) and applies to nested pointer structs.

The default `GucooingApiKey` (`config/config.go:186`) is generated with `mx.GetMxToken(0, 32)` — a 32-char hex random string from `crypto/rand`. The `_ int64` seed is ignored; do not rely on it being deterministic.

DB types supported: `sqlite` (via `ncruces/go-sqlite3` — CGO-free) and `mysql`. Both go through the same GORM layer.

## Architecture

Request flow for game traffic:

```
client → POST /api/gateway → gateway.gateWay (gateway/gateway.go:67)
       → protocol.UnmarshalRequest          (decodes BasePacket → concrete proto)
       → Gateway.registerMessage            (routes via funcRouteMap)
       → pack.<Handler>(session, req, rsp)  (per-protocol handler in pack/)
       → game.<DomainLogic>(session, ...)   (mutates PlayerBin)
       → protocol.MarshalResponse           (re-encode, optional gzip)
```

`funcRouteMap` (`gateway/gateway.router.go:20`) is the dispatch table — `proto.Protocol` enum → `handlerFunc`. Add new protocols here, then implement in `pack/`. Comments next to each entry document Chinese-language intent of the message — preserve them.

Layer responsibilities (do not invert these dependencies):
- **`pack/`** — protocol-level handlers; one function per `proto.Protocol_*`. Translates between wire messages and `game.*` calls. Should not touch DB directly.
- **`game/`** — domain logic operating on `Session.PlayerBin` (the in-memory protobuf player state). Files split by feature: `gacha.go`, `arena.go`, `cafe.go`, `raid*.go`, `dungeon.*.go`, `echelon.go`, …
- **`common/enter/`** — session cache. `Session` (`common/enter/session.go:26`) holds `PlayerBin` + ephemeral state. `MaxCachePlayerTime` (minutes) caps idle retention before the session is flushed to DB and evicted; `MaxPlayerNum` caps online concurrency.
- **`db/`** — `DBGame` interface in `db/db.go` is the only abstraction the rest of the code consumes. GORM implementation in `db/gorm/`, row structs in `db/struct/`. Two databases are configured (`DB` for game state, `RankDB` for leaderboards).
- **`gdconf/`** — read-only game data tables. `GC` is the global `*GameConfig` singleton with `Excel` (raw deserialized proto) and `GPP` (post-processed indexes — many `map[id]*Row` lookups). One `exceldb.*.go` / `excel.*.go` file per table.
- **`protocol/`** — wire format. `cmd/` maps protocol IDs ↔ Go types, `mx/` handles encoding + token gen, `proto/` has one `.go` file per protocol message (these are hand-written wrappers, not generated from `.proto`).
- **`common/server_only/`** — server-internal protobuf schemas (`PlayerBin`, `Excel`, `MailInfo`, …). `.pb.go` files are gitignored — regenerate with the bundled `protoc` + `protoc-gen-go` after pulling.
- **`sdk/`** — Yostar SDK emulation routes (`/account/*`, `/user/*`, `/app/*`, `/prod/index.json`). Templates rendered from `<DataPath>/templates/`.
- **`command/`** — HTTP admin API for GMs (built on `gucooing/cdq`, pinned to the safe fork). Requires `GucooingApiKey`. Routes registered in `command/command.go:18`.
- **`common/rank/`** — leaderboards backed by `pkg/zset` (vendored fork of liyiheng/zset).

Player state is **never persisted in real time**. Mutations happen on the in-memory `Session.PlayerBin`. Disk writes occur on session timeout / shutdown (see `EnterSet.checkSession` in `common/enter/session.go:50` and `Close` calls in `baps.go:122`). Anything that needs a real-time view goes through the command API, not the DB.

## Proxying Clients

Point the client at the server by intercepting:
- `https://ba-jp-sdk.bluearchive.jp` → `http://<server>:<port>`
- `https://yostar-serverinfo.bluearchiveyostar.com` → `http://<server>:<port>`

`Android_Mitmproxy_Readme_*.md` documents the mitmproxy-based interception used for Android clients. `OuterAddr` in `HttpNet` config must match the externally reachable URL — `prod/index.json` and various SDK responses embed it.
