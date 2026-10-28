# Graph Report - roommates  (2026-10-28)

## Corpus Check
- 104 files · ~39,170 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 27 file(s) not represented in the graph (top: .templ 13, (none) 5, .css 3)

## Summary
- 892 nodes · 2359 edges · 52 communities (39 shown, 13 thin omitted)
- Extraction: 96% EXTRACTED · 4% INFERRED · 0% AMBIGUOUS · INFERRED: 90 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e5aca3bf`
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
- github.com/a-h/templ.Component
- components.json
- block-secret-files.sh
- go
- password
- context.Context
- household
- household
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
- not-found-page
- not-found-page
- password
- validation-errors
- home-page
- home-page
- login-page
- login-page
- password
- invite-page
- registration-page
- invite-page
- registration-page
- validation-errors
- password

## God Nodes (most connected - your core abstractions)
1. `Fe()` - 28 edges
2. `CN()` - 26 edges
3. `ee()` - 26 edges
4. `oe()` - 26 edges
5. `le()` - 24 edges
6. `re()` - 23 edges
7. `e()` - 23 edges
8. `c()` - 21 edges
9. `LKMessage` - 19 edges
10. `t()` - 19 edges

## Surprising Connections (you probably didn't know these)
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go
- `Input()` --calls--> `RandomID()`  [EXTRACTED]
  app/web/components/input/input_templ.go → app/internal/utils/shadcn-templ.go
- `Household()` --calls--> `formatTime()`  [INFERRED]
  app/web/pages/household_templ.go → app/web/pages/pages-household.go
- `main()` --calls--> `InitGinEngine()`  [INFERRED]
  app/cmd/server/main.go → app/cmd/server/routes.go
- `main()` --calls--> `LoadLocales()`  [EXTRACTED]
  app/cmd/server/main.go → app/internal/locales/locales.go

## Import Cycles
- None detected.

## Communities (52 total, 13 thin omitted)

### Community 0 - "htmx.org@2.0.11.min.js"
Cohesion: 0.08
Nodes (103): A(), Ae(), an(), at(), B(), Be(), bn(), bt() (+95 more)

### Community 1 - "handlers.go"
Cohesion: 0.10
Nodes (10): Direction(), DirectionAttributes(), lastTopLevelColon(), MergeAttributes(), RandomID(), twImportantToV3(), twImportantToV4(), WithDirection() (+2 more)

### Community 2 - "migrations.go"
Cohesion: 0.06
Nodes (47): ensureDatabasePool(), main(), InitGinEngine(), initRoutes(), registerAssets(), GetDatabaseURL(), IsDevelopment(), TestGetDatabaseURLEscapes() (+39 more)

### Community 3 - "en"
Cohesion: 0.15
Nodes (12): title, en, app, email, error-page, logout, username, welcome-back (+4 more)

### Community 4 - "et"
Cohesion: 0.15
Nodes (12): title, back, description, title, et, app, email, error-page (+4 more)

### Community 5 - "db"
Cohesion: 0.06
Nodes (35): build, depends_on, env_file, environment, ports, volumes, context, dockerfile (+27 more)

### Community 6 - "generate-locales-keys.go"
Cohesion: 0.07
Nodes (19): GetEnv(), GetSessionKeyPairs(), MustGetEnv(), ValidateCookieHashBlockPair(), flattenYAML(), generateOutputContent(), main(), parseLocaleYml() (+11 more)

### Community 7 - "vars"
Cohesion: 0.06
Nodes (31): cmds, desc, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV, sh (+23 more)

### Community 8 - "Session"
Cohesion: 0.23
Nodes (4): LogValidationErrors(), gorillaSessionGetter, Session, SessionInfo

### Community 9 - "LKMessage"
Cohesion: 0.13
Nodes (11): LKMessage, Household, IsModelValid(), Login, Register, TestLoginIsValid(), TestRegisterIsValid(), ValidateModel() (+3 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.07
Nodes (77): LK, CN(), flattenClasses(), TestCN(), twMapTokens(), twReorder(), N(), T() (+69 more)

### Community 11 - "github.com/a-h/templ.Component"
Cohesion: 0.09
Nodes (27): DeleteCookie(), GetCookie(), GetHeader(), householdPath(), membershipFromContext(), sessionUserID(), InternalError(), IsRequestHTMX() (+19 more)

### Community 12 - "components.json"
Cohesion: 0.12
Nodes (16): aliases, components, utils, iconLibrary, menuAccent, menuColor, rtl, $schema (+8 more)

### Community 13 - "block-secret-files.sh"
Cohesion: 0.60
Nodes (3): block(), deny(), block-secret-files.sh script

### Community 14 - "go"
Cohesion: 0.12
Nodes (15): go, emit_enum_valid_method, emit_exact_table_names, emit_json_tags, emit_pointers_for_null_types, out, package, sql_driver (+7 more)

### Community 15 - "password"
Cohesion: 0.50
Nodes (4): password, confirm, enter, title

### Community 16 - "context.Context"
Cohesion: 0.06
Nodes (29): Queries, New(), GetHouseholdRow, GetInviteByCodeRow, GetUserCredentialsRow, InsertUserParams, ListActiveInvitesRow, ListActiveMembersRow (+21 more)

### Community 17 - "household"
Cohesion: 0.12
Nodes (17): household, back, create-invite, expires, invites, invites-description, last-owner, leave (+9 more)

### Community 18 - "household"
Cohesion: 0.12
Nodes (17): household, back, create-invite, expires, invites, invites-description, last-owner, leave (+9 more)

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

### Community 37 - "not-found-page"
Cohesion: 0.67
Nodes (3): not-found-page, description, title

### Community 38 - "not-found-page"
Cohesion: 0.67
Nodes (3): not-found-page, description, title

### Community 40 - "password"
Cohesion: 0.29
Nodes (7): case, length, match, max-length, required, symbol, password

### Community 41 - "validation-errors"
Cohesion: 0.22
Nodes (9): generic, validation-errors, name-length, length, max-length, spaces, email, household (+1 more)

### Community 42 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 43 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 44 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 45 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 46 - "password"
Cohesion: 0.29
Nodes (7): case, length, match, max-length, required, symbol, password

### Community 47 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 48 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 49 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 50 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 52 - "validation-errors"
Cohesion: 0.22
Nodes (9): generic, validation-errors, name-length, length, max-length, spaces, email, household (+1 more)

### Community 54 - "password"
Cohesion: 0.50
Nodes (4): password, confirm, enter, title

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 282 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `LKMessage` connect `LKMessage` to `handlers.go`, `github.com/a-h/templ.Component`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `htmx.org@2.0.11.min.js` be split into smaller, more focused modules?**
  _Cohesion score 0.08159073935772965 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `github.com/a-h/templ.Component` to `context.Context`, `handlers.go`, `github.com/a-h/templ.Component`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09802922188243289 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `not-found-page`, `home-page`, `login-page`, `invite-page`, `password`, `household`, `registration-page`, `validation-errors`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **Should `migrations.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05913461538461538 - nodes in this community are weakly interconnected._