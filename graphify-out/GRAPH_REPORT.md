# Graph Report - roommates  (2026-10-28)

## Corpus Check
- 90 files · ~28,786 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 25 file(s) not represented in the graph (top: .templ 11, (none) 5, .css 3)

## Summary
- 717 nodes · 1730 edges · 34 communities (20 shown, 14 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 55 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a4dabf64`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- htmx.org@2.0.11.min.js
- handlers.go
- migrations.go
- en
- et
- db
- generate-locales-keys.go
- vars
- Session
- LKMessage
- github.com/a-h/templ.Component
- github.com/gin-gonic/gin.Context
- components.json
- block-secret-files.sh
- go
- shadcn-templ.go
- permissions
- linters
- Certificates
- HTMLTemplRenderer
- Roommates
- db/README.md
- handlers/README.md
- services/README.md
- dev-run.sh
- format.sh
- lint.sh
- tools/README.md
- assets/README.md
- components/README.md
- renderer/README.md
- Roommates
- .claude/CLAUDE.md
- roommates

## God Nodes (most connected - your core abstractions)
1. `Fe()` - 28 edges
2. `CN()` - 26 edges
3. `ee()` - 26 edges
4. `oe()` - 26 edges
5. `le()` - 24 edges
6. `re()` - 23 edges
7. `e()` - 23 edges
8. `c()` - 21 edges
9. `t()` - 19 edges
10. `vn()` - 19 edges

## Surprising Connections (you probably didn't know these)
- `TestRateLimit()` --calls--> `NewRateLimit()`  [INFERRED]
  app/internal/middleware/middleware-ratelimit_test.go → app/internal/middleware/middleware-ratelimit.go
- `TestCN()` --calls--> `CN()`  [INFERRED]
  app/internal/utils/shadcn-templ_test.go → app/internal/utils/shadcn-templ.go
- `Input()` --calls--> `RandomID()`  [EXTRACTED]
  app/web/components/input/input_templ.go → app/internal/utils/shadcn-templ.go
- `main()` --calls--> `InitGinEngine()`  [INFERRED]
  app/cmd/server/main.go → app/cmd/server/routes.go
- `main()` --calls--> `LoadLocales()`  [EXTRACTED]
  app/cmd/server/main.go → app/internal/locales/locales.go

## Import Cycles
- None detected.

## Communities (34 total, 14 thin omitted)

### Community 0 - "htmx.org@2.0.11.min.js"
Cohesion: 0.08
Nodes (103): A(), Ae(), an(), at(), B(), Be(), bn(), bt() (+95 more)

### Community 2 - "migrations.go"
Cohesion: 0.06
Nodes (41): ensureDatabasePool(), main(), InitGinEngine(), initRoutes(), registerAssets(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs() (+33 more)

### Community 3 - "en"
Cohesion: 0.04
Nodes (45): title, generic, en, app, email, error-page, home-page, login-page (+37 more)

### Community 4 - "et"
Cohesion: 0.04
Nodes (45): title, generic, back, description, title, et, app, email (+37 more)

### Community 5 - "db"
Cohesion: 0.06
Nodes (35): build, depends_on, env_file, environment, ports, volumes, context, dockerfile (+27 more)

### Community 6 - "generate-locales-keys.go"
Cohesion: 0.08
Nodes (14): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+6 more)

### Community 7 - "vars"
Cohesion: 0.06
Nodes (31): cmds, desc, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV, sh (+23 more)

### Community 8 - "Session"
Cohesion: 0.23
Nodes (4): LogValidationErrors(), gorillaSessionGetter, Session, SessionInfo

### Community 9 - "LKMessage"
Cohesion: 0.13
Nodes (12): LKMessage, TestRateLimit(), IsModelValid(), Login, Register, TestLoginIsValid(), TestRegisterIsValid(), ValidateModel() (+4 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.08
Nodes (63): LK, CN(), N(), T(), AspectRatio(), Button(), Variants(), Action() (+55 more)

### Community 11 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.10
Nodes (20): DeleteCookie(), GetCookie(), GetUserCredentialsRow, InsertUserParams, GetHeader(), InternalError(), IsRequestHTMX(), Redirect() (+12 more)

### Community 12 - "components.json"
Cohesion: 0.12
Nodes (16): aliases, components, utils, iconLibrary, menuAccent, menuColor, rtl, $schema (+8 more)

### Community 13 - "block-secret-files.sh"
Cohesion: 0.60
Nodes (3): block(), deny(), block-secret-files.sh script

### Community 14 - "go"
Cohesion: 0.12
Nodes (15): go, emit_enum_valid_method, emit_exact_table_names, emit_json_tags, emit_pointers_for_null_types, out, package, sql_driver (+7 more)

### Community 15 - "shadcn-templ.go"
Cohesion: 0.07
Nodes (17): Direction(), DirectionAttributes(), flattenClasses(), lastTopLevelColon(), MergeAttributes(), RandomID(), twImportantToV3(), twImportantToV4() (+9 more)

### Community 19 - "permissions"
Cohesion: 0.25
Nodes (7): hooks, PreToolUse, permissions, allow, ask, deny, $schema

### Community 20 - "linters"
Cohesion: 0.29
Nodes (6): allow-empty, linters, enable, settings, exhaustruct_v5, version

### Community 21 - "Certificates"
Cohesion: 0.40
Nodes (4): Certificates, Generate digital certificate, Generate RSA private key, Install Openssl

### Community 23 - "Roommates"
Cohesion: 0.50
Nodes (3): Project structure, Roommates, Setup

### Community 34 - "Roommates"
Cohesion: 0.40
Nodes (4): Commands, Conventions, graphify, Roommates

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 224 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **14 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `LKMessage` connect `LKMessage` to `handlers.go`, `github.com/a-h/templ.Component`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `htmx.org@2.0.11.min.js` be split into smaller, more focused modules?**
  _Cohesion score 0.08159073935772965 - nodes in this community are weakly interconnected._
- **Why does `Session` connect `Session` to `handlers.go`, `github.com/gin-gonic/gin.Context`?**
  _High betweenness centrality (0.006) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10613437195715676 - nodes in this community are weakly interconnected._
- **Why does `CN()` connect `github.com/a-h/templ.Component` to `LKMessage`, `shadcn-templ.go`?**
  _High betweenness centrality (0.005) - this node is a cross-community bridge._
- **Should `migrations.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05683060109289618 - nodes in this community are weakly interconnected._