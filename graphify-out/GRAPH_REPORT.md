# Graph Report - roommates  (2026-10-28)

## Corpus Check
- 78 files · ~25,469 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 21 file(s) not represented in the graph (top: .templ 9, (none) 3, .css 3)

## Summary
- 655 nodes · 1547 edges · 37 communities (23 shown, 14 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 44 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `72c92338`
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
- github.com/gin-gonic/gin.Context
- LKMessage
- github.com/a-h/templ.Component
- field_templ.go
- components.json
- context.Context
- go
- shadcn-templ.go
- db.go
- card_templ.go
- button_templ.go
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
- CLAUDE.md
- .claude/CLAUDE.md
- roommates

## God Nodes (most connected - your core abstractions)
1. `Fe()` - 28 edges
2. `ee()` - 26 edges
3. `oe()` - 26 edges
4. `CN()` - 25 edges
5. `le()` - 24 edges
6. `re()` - 23 edges
7. `e()` - 23 edges
8. `c()` - 21 edges
9. `t()` - 19 edges
10. `vn()` - 19 edges

## Surprising Connections (you probably didn't know these)
- `Input()` --calls--> `RandomID()`  [EXTRACTED]
  app/web/components/input/input_templ.go → app/internal/utils/shadcn-templ.go
- `main()` --calls--> `InitGinEngine()`  [INFERRED]
  app/cmd/server/main.go → app/cmd/server/routes.go
- `main()` --calls--> `LoadLocales()`  [EXTRACTED]
  app/cmd/server/main.go → app/internal/locales/locales.go
- `main()` --calls--> `InitLogger()`  [EXTRACTED]
  app/cmd/server/main.go → app/internal/logger/logger.go
- `initRoutes()` --calls--> `NewInternationalization()`  [EXTRACTED]
  app/cmd/server/routes.go → app/internal/middleware/middleware-i18n.go

## Import Cycles
- None detected.

## Communities (37 total, 14 thin omitted)

### Community 0 - "htmx.org@2.0.11.min.js"
Cohesion: 0.08
Nodes (103): A(), Ae(), an(), at(), B(), Be(), bn(), bt() (+95 more)

### Community 1 - "handlers.go"
Cohesion: 0.10
Nodes (3): ErrorCode, BaseResponse, ErrorInfo

### Community 2 - "migrations.go"
Cohesion: 0.08
Nodes (34): ensureDatabasePool(), main(), ginTemplIntegration(), InitGinEngine(), initRoutes(), registerAssets(), GetDatabaseURL(), GetEnv() (+26 more)

### Community 3 - "en"
Cohesion: 0.06
Nodes (35): title, generic, en, app, email, login-page, password, registration-page (+27 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (35): title, generic, et, app, email, login-page, password, registration-page (+27 more)

### Community 5 - "db"
Cohesion: 0.06
Nodes (34): build, depends_on, env_file, environment, ports, context, dockerfile, build (+26 more)

### Community 6 - "generate-locales-keys.go"
Cohesion: 0.10
Nodes (14): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+6 more)

### Community 7 - "vars"
Cohesion: 0.06
Nodes (31): cmds, desc, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV, sh (+23 more)

### Community 8 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.11
Nodes (20): DeleteCookie(), GetCookie(), SetCookie(), Fail(), GetHeader(), InternalError(), IsRequestHTMX(), OK() (+12 more)

### Community 9 - "LKMessage"
Cohesion: 0.11
Nodes (12): GetUserCredentialsRow, InsertUserParams, LKMessage, IsModelValid(), Login, LoginRequest, Register, ValidateModel() (+4 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.26
Nodes (23): CN(), T(), AspectRatio(), Card(), Content(), Description(), Header(), Title() (+15 more)

### Community 11 - "field_templ.go"
Cohesion: 0.15
Nodes (20): childrenHTML(), Description(), Error(), fieldClasses(), Separator(), Set(), Props, ContentProps (+12 more)

### Community 12 - "components.json"
Cohesion: 0.12
Nodes (16): aliases, components, utils, iconLibrary, menuAccent, menuColor, rtl, $schema (+8 more)

### Community 13 - "context.Context"
Cohesion: 0.17
Nodes (9): LK, Direction(), DirectionAttributes(), WithDirection(), GetLocaleString(), N(), N(), New() (+1 more)

### Community 14 - "go"
Cohesion: 0.12
Nodes (15): go, emit_enum_valid_method, emit_exact_table_names, emit_json_tags, emit_pointers_for_null_types, out, package, sql_driver (+7 more)

### Community 15 - "shadcn-templ.go"
Cohesion: 0.15
Nodes (9): flattenClasses(), lastTopLevelColon(), MergeAttributes(), RandomID(), twImportantToV3(), twImportantToV4(), twMapTokens(), twReorder() (+1 more)

### Community 16 - "db.go"
Cohesion: 0.18
Nodes (3): Queries, New(), DBTX

### Community 17 - "card_templ.go"
Cohesion: 0.20
Nodes (10): Action(), Footer(), ActionProps, ContentProps, DescriptionProps, FooterProps, HeaderProps, Props (+2 more)

### Community 18 - "button_templ.go"
Cohesion: 0.44
Nodes (6): Button(), Variants(), Props, Size, Type, Variant

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

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 198 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **14 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `LKMessage` connect `LKMessage` to `handlers.go`, `github.com/a-h/templ.Component`, `context.Context`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `htmx.org@2.0.11.min.js` be split into smaller, more focused modules?**
  _Cohesion score 0.08159073935772965 - nodes in this community are weakly interconnected._
- **Why does `Renderer` connect `context.Context` to `handlers.go`, `github.com/a-h/templ.Component`?**
  _High betweenness centrality (0.006) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10186335403726708 - nodes in this community are weakly interconnected._
- **Should `migrations.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07738095238095238 - nodes in this community are weakly interconnected._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.05555555555555555 - nodes in this community are weakly interconnected._