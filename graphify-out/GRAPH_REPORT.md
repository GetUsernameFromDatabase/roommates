# Graph Report - roommates  (2026-10-29)

## Corpus Check
- 163 files · ~75,344 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 37 file(s) not represented in the graph (top: .templ 19, (none) 7, .css 5)

## Summary
- 1448 nodes · 4315 edges · 84 communities (68 shown, 16 thin omitted)
- Extraction: 94% EXTRACTED · 6% INFERRED · 0% AMBIGUOUS · INFERRED: 264 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `01098e32`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- migrations.go
- en
- et
- environment
- vendor-assets.go
- vars
- generate-locales-keys.go
- LKMessage
- github.com/a-h/templ.Component
- github.com/a-h/templ.Component
- components.json
- block-secret-files.sh
- go
- context.Context
- github.com/a-h/templ.Attributes
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
- models.go
- shopping
- shopping
- newEmailToken
- validation-errors
- home-page
- home-page
- login-page
- login-page
- notes
- invite-page
- registration-page
- invite-page
- registration-page
- files
- validation-errors
- Queries
- notes
- noRowsAs
- verify-email
- verify-email
- reset-password-page
- reset-password-page
- mail
- mail
- github.com/jackc/pgx/v5/pgtype.UUID
- forgot-password-page
- forgot-password-page
- files
- S3Store
- appearance
- appearance
- .UploadFile
- environment
- nav
- nav
- ListNotesRow
- services
- db
- queries/README.md
- dev-signal.sh
- ErrorCode
- db
- context
- healthcheck
- app
- Renderer

## God Nodes (most connected - your core abstractions)
1. `T()` - 48 edges
2. `CN()` - 34 edges
3. `T()` - 33 edges
4. `Membership` - 31 edges
5. `membershipFromContext()` - 30 edges
6. `RenderTemplStatus()` - 30 edges
7. `LKMessage` - 29 edges
8. `InternalError()` - 26 edges
9. `Button()` - 25 edges
10. `en` - 24 edges

## Surprising Connections (you probably didn't know these)
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go
- `TestExpiredTokenIsInvalid()` --calls--> `hashToken()`  [INFERRED]
  app/internal/services/services-email_test.go → app/internal/services/services-email.go
- `TestCN()` --calls--> `CN()`  [INFERRED]
  app/internal/utils/shadcn-templ_test.go → app/internal/utils/shadcn-templ.go
- `TestTUnescapesTextButNotArgs()` --calls--> `T()`  [INFERRED]
  app/internal/utils/utils-i18n_test.go → app/internal/utils/utils-i18n.go
- `TestNoteHTMLIsSafe()` --calls--> `noteHTML()`  [INFERRED]
  app/web/pages/pages-notes_test.go → app/web/pages/pages-notes.go

## Import Cycles
- None detected.

## Communities (84 total, 16 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.06
Nodes (52): TestSafeNextPath(), exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage() (+44 more)

### Community 1 - "handlers.go"
Cohesion: 0.12
Nodes (5): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), appearanceKey, directionKey

### Community 2 - "migrations.go"
Cohesion: 0.05
Nodes (46): ensureDatabasePool(), main(), InitGinEngine(), initRoutes(), registerAssets(), GetAppBaseURL(), GetDatabaseURL(), GetEnv() (+38 more)

### Community 3 - "en"
Cohesion: 0.10
Nodes (19): title, en, app, email, error-page, logout, not-found-page, password (+11 more)

### Community 4 - "et"
Cohesion: 0.10
Nodes (19): title, back, description, title, et, app, email, error-page (+11 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.12
Nodes (24): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+16 more)

### Community 7 - "vars"
Cohesion: 0.05
Nodes (38): cmds, desc, task, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV (+30 more)

### Community 8 - "generate-locales-keys.go"
Cohesion: 0.07
Nodes (18): LogValidationErrors(), flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput() (+10 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.09
Nodes (86): ListHouseholdsOfUserRow, LK, AttachedFile, HouseholdOverview, Membership, ShoppingList, CN(), flattenClasses() (+78 more)

### Community 11 - "github.com/a-h/templ.Component"
Cohesion: 0.06
Nodes (37): DeleteCookie(), GetCookie(), notice(), filesPath(), shoppingItemPath(), GetHeader(), householdPath(), membershipFromContext() (+29 more)

### Community 12 - "components.json"
Cohesion: 0.12
Nodes (16): aliases, components, utils, iconLibrary, menuAccent, menuColor, rtl, $schema (+8 more)

### Community 13 - "block-secret-files.sh"
Cohesion: 0.60
Nodes (3): block(), deny(), block-secret-files.sh script

### Community 14 - "go"
Cohesion: 0.12
Nodes (15): go, emit_enum_valid_method, emit_exact_table_names, emit_json_tags, emit_pointers_for_null_types, out, package, sql_driver (+7 more)

### Community 15 - "context.Context"
Cohesion: 0.04
Nodes (23): Direction(), WithDirection(), GetLocaleString(), childrenHTML(), ConsumeEmailTokenParams, CountActiveMembersRow, DeleteCheckedShoppingItemsParams, DeleteEmailTokensOfUserParams (+15 more)

### Community 16 - "github.com/a-h/templ.Attributes"
Cohesion: 0.06
Nodes (31): DirectionAttributes(), MergeAttributes(), Variants(), fieldClasses(), Props, Props, Size, Type (+23 more)

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

### Community 27 - "dev-run.sh"
Cohesion: 0.67
Nodes (5): dlv_alive(), dlv_port_busy(), dev-run.sh script, start(), stop()

### Community 34 - "Roommates"
Cohesion: 0.40
Nodes (4): Commands, Conventions, graphify, Roommates

### Community 37 - "models.go"
Cohesion: 0.12
Nodes (17): ListShoppingItemsRow, EmailToken, File, FilePurge, GetFileParams, GetFileRow, Household, InsertEmailTokenParams (+9 more)

### Community 38 - "shopping"
Cohesion: 0.09
Nodes (23): shopping, add, back, bought-by, buying, buying-by, check, clear (+15 more)

### Community 39 - "shopping"
Cohesion: 0.09
Nodes (23): shopping, add, back, bought-by, buying, buying-by, check, clear (+15 more)

### Community 40 - "newEmailToken"
Cohesion: 0.12
Nodes (6): GetUserCredentialsRow, InsertUserParams, LoginRequest, NormalizeEmail(), newEmailToken(), timestamptz()

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

### Community 46 - "notes"
Cohesion: 0.10
Nodes (20): notes, back, body, body-hint, cancel, delete, edit, edited-by (+12 more)

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

### Community 51 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, delete, download, file (+7 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "Queries"
Cohesion: 0.14
Nodes (4): Queries, New(), ownerActsOnMember(), DBTX

### Community 54 - "notes"
Cohesion: 0.10
Nodes (20): notes, back, body, body-hint, cancel, delete, edit, edited-by (+12 more)

### Community 55 - "noRowsAs"
Cohesion: 0.11
Nodes (7): GetInviteByCodeRow, hashToken(), cleanFileName(), TestCleanFileName(), noRowsAs(), rowsAffected(), trimToLength()

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

### Community 62 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.07
Nodes (18): GetHouseholdRow, ListActiveInvitesRow, ListActiveMembersRow, GetShoppingListRow, ListShoppingListsRow, GetUserRow, GetActiveMembershipParams, GetActiveMembershipRow (+10 more)

### Community 63 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 64 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 65 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, delete, download, file (+7 more)

### Community 66 - "S3Store"
Cohesion: 0.16
Nodes (3): Store, New(), S3Store

### Community 67 - "appearance"
Cohesion: 0.20
Nodes (10): mode, mode-dark, mode-light, mode-system, save, theme, theme-home, theme-plain (+2 more)

### Community 68 - "appearance"
Cohesion: 0.20
Nodes (10): mode, mode-dark, mode-light, mode-system, save, theme, theme-home, theme-plain (+2 more)

### Community 69 - ".UploadFile"
Cohesion: 0.24
Nodes (4): ListFilesRow, FileMeta, FileTarget, ListFilesParams

### Community 70 - "environment"
Cohesion: 0.18
Nodes (11): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment, ports (+3 more)

### Community 71 - "nav"
Cohesion: 0.29
Nodes (7): nav, files, label, notes, people, shopping, switch

### Community 72 - "nav"
Cohesion: 0.29
Nodes (7): nav, files, label, notes, people, shopping, switch

### Community 73 - "ListNotesRow"
Cohesion: 0.19
Nodes (4): ListNotesRow, NoteInput, ListNotesParams, NoteForm

### Community 74 - "services"
Cohesion: 0.22
Nodes (8): ports, restart, services, mailpit, version, volumes, db_data, rustfs_data

### Community 75 - "db"
Cohesion: 0.25
Nodes (8): depends_on, condition, required, restart, db, rustfs, condition, required

### Community 79 - "db"
Cohesion: 0.25
Nodes (8): command, environment, ports, volumes, POSTGRES_DB, POSTGRES_PASSWORD, POSTGRES_USER, db

### Community 80 - "context"
Cohesion: 0.53
Nodes (6): build, context, dockerfile, build, build, build

### Community 81 - "healthcheck"
Cohesion: 0.53
Nodes (6): healthcheck, interval, retries, test, timeout, healthcheck

### Community 82 - "app"
Cohesion: 0.50
Nodes (4): env_file, ports, volumes, app

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 530 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `appearance`, `shopping`, `nav`, `home-page`, `login-page`, `notes`, `invite-page`, `registration-page`, `household`, `files`, `validation-errors`, `verify-email`, `reset-password-page`, `mail`, `forgot-password-page`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.056140350877192984 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `forgot-password-page`, `files`, `appearance`, `shopping`, `nav`, `validation-errors`, `home-page`, `login-page`, `invite-page`, `household`, `registration-page`, `notes`, `verify-email`, `reset-password-page`, `mail`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12274821952241308 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `github.com/a-h/templ.Component` to `testing.T`, `handlers.go`, `newEmailToken`, `github.com/a-h/templ.Component`, `context.Context`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `migrations.go` be split into smaller, more focused modules?**
  _Cohesion score 0.051560379918588875 - nodes in this community are weakly interconnected._