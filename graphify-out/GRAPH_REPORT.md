# Graph Report - roommates  (2026-10-28)

## Corpus Check
- 120 files · ~53,476 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 29 file(s) not represented in the graph (top: .templ 14, (none) 6, .css 3)

## Summary
- 1162 nodes · 3269 edges · 68 communities (52 shown, 16 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 163 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c620660a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- htmx.org@2.0.11.min.js
- handlers.go
- migrations.go
- en
- et
- db
- generate-locales-languages.go
- vars
- Session
- LKMessage
- github.com/a-h/templ.Component
- github.com/a-h/templ.Component
- components.json
- block-secret-files.sh
- go
- queries.sql.go
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
- github.com/jackc/pgx/v5/pgtype.UUID
- shopping
- shopping
- newEmailToken
- validation-errors
- home-page
- home-page
- login-page
- login-page
- generate-locales-keys.go
- invite-page
- registration-page
- invite-page
- registration-page
- sessions.go
- validation-errors
- Queries
- password
- Membership
- verify-email
- verify-email
- reset-password-page
- reset-password-page
- mail
- mail
- HouseholdOverview
- forgot-password-page
- forgot-password-page
- validators.go
- error-page

## God Nodes (most connected - your core abstractions)
1. `T()` - 34 edges
2. `LKMessage` - 29 edges
3. `Fe()` - 28 edges
4. `CN()` - 27 edges
5. `ee()` - 26 edges
6. `oe()` - 26 edges
7. `le()` - 24 edges
8. `re()` - 23 edges
9. `e()` - 23 edges
10. `household` - 21 edges

## Surprising Connections (you probably didn't know these)
- `TestExpiredTokenIsInvalid()` --calls--> `hashToken()`  [INFERRED]
  app/internal/services/services-email_test.go → app/internal/services/services-email.go
- `TestNewInviteCode()` --calls--> `newInviteCode()`  [INFERRED]
  app/internal/services/services-households_test.go → app/internal/services/services-households.go
- `Household()` --calls--> `formatDate()`  [INFERRED]
  app/web/pages/household_templ.go → app/web/pages/pages-household.go
- `main()` --calls--> `InitGinEngine()`  [INFERRED]
  app/cmd/server/main.go → app/cmd/server/routes.go
- `main()` --calls--> `InitLogger()`  [EXTRACTED]
  app/cmd/server/main.go → app/internal/logger/logger.go

## Import Cycles
- None detected.

## Communities (68 total, 16 thin omitted)

### Community 0 - "htmx.org@2.0.11.min.js"
Cohesion: 0.08
Nodes (103): A(), Ae(), an(), at(), B(), Be(), bn(), bt() (+95 more)

### Community 1 - "handlers.go"
Cohesion: 0.14
Nodes (5): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ErrorCode, directionKey

### Community 2 - "migrations.go"
Cohesion: 0.05
Nodes (66): ensureDatabasePool(), main(), InitGinEngine(), initRoutes(), registerAssets(), GetAppBaseURL(), GetDatabaseURL(), GetEnv() (+58 more)

### Community 3 - "en"
Cohesion: 0.10
Nodes (19): title, en, app, email, error-page, logout, not-found-page, password (+11 more)

### Community 4 - "et"
Cohesion: 0.17
Nodes (11): title, et, app, email, logout, not-found-page, username, welcome-back (+3 more)

### Community 5 - "db"
Cohesion: 0.05
Nodes (43): build, depends_on, env_file, environment, ports, volumes, context, dockerfile (+35 more)

### Community 6 - "generate-locales-languages.go"
Cohesion: 0.22
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 7 - "vars"
Cohesion: 0.06
Nodes (31): cmds, desc, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV, sh (+23 more)

### Community 8 - "Session"
Cohesion: 0.29
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.06
Nodes (92): LK, ShoppingList, CN(), DirectionAttributes(), flattenClasses(), IfElse(), MergeAttributes(), RandomID() (+84 more)

### Community 11 - "github.com/a-h/templ.Component"
Cohesion: 0.08
Nodes (28): DeleteCookie(), GetCookie(), notice(), GetHeader(), householdPath(), membershipFromContext(), sessionUserID(), InternalError() (+20 more)

### Community 12 - "components.json"
Cohesion: 0.12
Nodes (16): aliases, components, utils, iconLibrary, menuAccent, menuColor, rtl, $schema (+8 more)

### Community 13 - "block-secret-files.sh"
Cohesion: 0.60
Nodes (3): block(), deny(), block-secret-files.sh script

### Community 14 - "go"
Cohesion: 0.12
Nodes (15): go, emit_enum_valid_method, emit_exact_table_names, emit_json_tags, emit_pointers_for_null_types, out, package, sql_driver (+7 more)

### Community 15 - "queries.sql.go"
Cohesion: 0.06
Nodes (19): GetShoppingListRow, ConsumeEmailTokenParams, CountActiveMembersRow, DeleteCheckedShoppingItemsParams, DeleteEmailTokensOfUserParams, DeleteInviteParams, DeleteShoppingItemParams, DeleteShoppingListParams (+11 more)

### Community 16 - "context.Context"
Cohesion: 0.12
Nodes (9): GetInviteByCodeRow, hashToken(), Membership, newInviteCode(), noRowsAs(), ownerActsOnMember(), Direction(), WithDirection() (+1 more)

### Community 17 - "household"
Cohesion: 0.10
Nodes (21): household, back, create-invite, expires, invite-email-sent, invite-email-submit, invite-email-title, invites (+13 more)

### Community 18 - "household"
Cohesion: 0.10
Nodes (21): household, back, create-invite, expires, invite-email-sent, invite-email-submit, invite-email-title, invites (+13 more)

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

### Community 37 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.14
Nodes (17): ListActiveInvitesRow, ListActiveMembersRow, ListHouseholdsOfUserRow, ListShoppingItemsRow, formatDate(), EmailToken, Household, InsertEmailTokenParams (+9 more)

### Community 38 - "shopping"
Cohesion: 0.10
Nodes (21): shopping, add, bought-by, buying, buying-by, check, clear, create-list (+13 more)

### Community 39 - "shopping"
Cohesion: 0.10
Nodes (21): shopping, add, bought-by, buying, buying-by, check, clear, create-list (+13 more)

### Community 40 - "newEmailToken"
Cohesion: 0.12
Nodes (7): GetUserCredentialsRow, GetUserRow, InsertUserParams, LoginRequest, NormalizeEmail(), newEmailToken(), timestamptz()

### Community 41 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

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

### Community 46 - "generate-locales-keys.go"
Cohesion: 0.22
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

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

### Community 51 - "sessions.go"
Cohesion: 0.20
Nodes (3): NewPostgreSessionStoreHandler(), startExpiredSessionCleanup(), sessionCleaner

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "Queries"
Cohesion: 0.20
Nodes (3): Queries, New(), DBTX

### Community 54 - "password"
Cohesion: 0.50
Nodes (4): password, confirm, enter, title

### Community 56 - "verify-email"
Cohesion: 0.20
Nodes (10): verify-email, banner, continue, done-description, done-title, invalid-description, invalid-title, resend (+2 more)

### Community 57 - "verify-email"
Cohesion: 0.20
Nodes (10): verify-email, banner, continue, done-description, done-title, invalid-description, invalid-title, resend (+2 more)

### Community 58 - "reset-password-page"
Cohesion: 0.25
Nodes (8): reset-password-page, done-description, done-title, invalid-description, invalid-title, new-link, submit, title

### Community 59 - "reset-password-page"
Cohesion: 0.25
Nodes (8): reset-password-page, done-description, done-title, invalid-description, invalid-title, new-link, submit, title

### Community 60 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 61 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 62 - "HouseholdOverview"
Cohesion: 0.33
Nodes (3): GetHouseholdRow, ListShoppingListsRow, HouseholdOverview

### Community 63 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 64 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 67 - "error-page"
Cohesion: 0.50
Nodes (4): back, description, title, error-page

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 397 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `T()` connect `github.com/a-h/templ.Component` to `context.Context`, `handlers.go`, `newEmailToken`, `github.com/a-h/templ.Component`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `htmx.org@2.0.11.min.js` be split into smaller, more focused modules?**
  _Cohesion score 0.08159073935772965 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `forgot-password-page`, `error-page`, `shopping`, `validation-errors`, `home-page`, `login-page`, `invite-page`, `household`, `registration-page`, `password`, `verify-email`, `reset-password-page`, `mail`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13561776061776062 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `shopping`, `home-page`, `login-page`, `invite-page`, `registration-page`, `household`, `validation-errors`, `verify-email`, `reset-password-page`, `mail`, `forgot-password-page`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Should `migrations.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05067920585161965 - nodes in this community are weakly interconnected._