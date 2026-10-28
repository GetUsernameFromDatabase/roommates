# Graph Report - roommates  (2026-10-29)

## Corpus Check
- 152 files · ~64,910 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 35 file(s) not represented in the graph (top: .templ 18, (none) 6, .css 5)

## Summary
- 1276 nodes · 3701 edges · 80 communities (64 shown, 16 thin omitted)
- Extraction: 94% EXTRACTED · 6% INFERRED · 0% AMBIGUOUS · INFERRED: 206 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1ba2dfed`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- migrations.go
- en
- et
- db
- vendor-assets.go
- vars
- Session
- LKMessage
- github.com/a-h/templ.Component
- github.com/a-h/templ.Component
- components.json
- block-secret-files.sh
- go
- context.Context
- noRowsAs
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
- initRoutes
- validation-errors
- Queries
- notes
- Membership
- verify-email
- verify-email
- reset-password-page
- reset-password-page
- mail
- mail
- github.com/jackc/pgx/v5/pgtype.UUID
- forgot-password-page
- forgot-password-page
- mailer.go
- main
- appearance
- appearance
- EnvironmentVariableKey
- GetAppearance
- nav
- nav
- ListNotesRow
- password
- NormalizeInviteCode
- queries/README.md
- dev-signal.sh
- ErrorCode
- error-page

## God Nodes (most connected - your core abstractions)
1. `T()` - 42 edges
2. `CN()` - 32 edges
3. `LKMessage` - 29 edges
4. `T()` - 28 edges
5. `Membership` - 25 edges
6. `en` - 23 edges
7. `et` - 23 edges
8. `InternalError()` - 22 edges
9. `RenderTemplStatus()` - 22 edges
10. `shopping` - 22 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `InitGinEngine()`  [INFERRED]
  app/cmd/server/main.go → app/cmd/server/routes.go
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go
- `TestExpiredTokenIsInvalid()` --calls--> `hashToken()`  [INFERRED]
  app/internal/services/services-email_test.go → app/internal/services/services-email.go
- `newTestService()` --calls--> `NewDefaultService()`  [INFERRED]
  app/internal/services/services-households_test.go → app/internal/services/services.go
- `TestCN()` --calls--> `CN()`  [INFERRED]
  app/internal/utils/shadcn-templ_test.go → app/internal/utils/shadcn-templ.go

## Import Cycles
- None detected.

## Communities (80 total, 16 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.27
Nodes (22): TestSafeNextPath(), englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword(), TestVerifyEmail(), tokenOf(), DefaultService (+14 more)

### Community 1 - "handlers.go"
Cohesion: 0.15
Nodes (5): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), appearanceKey, directionKey

### Community 2 - "migrations.go"
Cohesion: 0.18
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.12
Nodes (15): title, en, app, email, error-page, logout, not-found-page, username (+7 more)

### Community 4 - "et"
Cohesion: 0.12
Nodes (15): title, et, app, email, logout, not-found-page, password, username (+7 more)

### Community 5 - "db"
Cohesion: 0.05
Nodes (43): build, depends_on, env_file, environment, ports, volumes, context, dockerfile (+35 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.05
Nodes (38): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+30 more)

### Community 7 - "vars"
Cohesion: 0.05
Nodes (38): cmds, desc, task, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV (+30 more)

### Community 8 - "Session"
Cohesion: 0.23
Nodes (4): LogValidationErrors(), gorillaSessionGetter, Session, SessionInfo

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.06
Nodes (110): ListHouseholdsOfUserRow, LK, HouseholdOverview, Membership, ShoppingList, CN(), DirectionAttributes(), flattenClasses() (+102 more)

### Community 11 - "github.com/a-h/templ.Component"
Cohesion: 0.07
Nodes (34): DeleteCookie(), GetCookie(), notice(), GetHeader(), householdPath(), membershipFromContext(), peoplePath(), sessionUserID() (+26 more)

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
Cohesion: 0.05
Nodes (18): Direction(), WithDirection(), GetLocaleString(), ConsumeEmailTokenParams, CountActiveMembersRow, DeleteInviteParams, DeleteNoteParams, DeleteShoppingItemParams (+10 more)

### Community 16 - "noRowsAs"
Cohesion: 0.27
Nodes (3): GetInviteByCodeRow, hashToken(), noRowsAs()

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
Cohesion: 0.16
Nodes (13): ListShoppingItemsRow, EmailToken, GetNoteParams, GetNoteRow, Household, InsertEmailTokenParams, Invite, ListShoppingItemsParams (+5 more)

### Community 38 - "shopping"
Cohesion: 0.09
Nodes (22): shopping, add, back, bought-by, buying, buying-by, check, clear (+14 more)

### Community 39 - "shopping"
Cohesion: 0.09
Nodes (22): shopping, add, back, bought-by, buying, buying-by, check, clear (+14 more)

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

### Community 51 - "initRoutes"
Cohesion: 0.16
Nodes (18): InitGinEngine(), initRoutes(), registerAssets(), GetAppBaseURL(), IsDevelopment(), DefaultHandler, NewDefaultHandler(), Mailer (+10 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "Queries"
Cohesion: 0.14
Nodes (4): Queries, New(), ownerActsOnMember(), DBTX

### Community 54 - "notes"
Cohesion: 0.10
Nodes (20): notes, back, body, body-hint, cancel, delete, edit, edited-by (+12 more)

### Community 55 - "Membership"
Cohesion: 0.15
Nodes (5): Note, NoteInput, rowsAffected(), trimToLength(), NoteForm

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
Cohesion: 0.09
Nodes (15): GetHouseholdRow, ListActiveInvitesRow, ListActiveMembersRow, GetShoppingListRow, ListShoppingListsRow, DeleteCheckedShoppingItemsParams, DeleteEmailTokensOfUserParams, GetActiveMembershipParams (+7 more)

### Community 63 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 64 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 65 - "mailer.go"
Cohesion: 0.14
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 66 - "main"
Cohesion: 0.25
Nodes (8): ensureDatabasePool(), main(), GetDatabaseURL(), TestGetDatabaseURLEscapes(), LoadLocales(), InitLogger(), TestTUnescapesTextButNotArgs(), unescapeJSON()

### Community 67 - "appearance"
Cohesion: 0.20
Nodes (10): mode, mode-dark, mode-light, mode-system, save, theme, theme-home, theme-plain (+2 more)

### Community 68 - "appearance"
Cohesion: 0.20
Nodes (10): mode, mode-dark, mode-light, mode-system, save, theme, theme-home, theme-plain (+2 more)

### Community 69 - "EnvironmentVariableKey"
Cohesion: 0.18
Nodes (8): GetEnv(), GetSessionKeyPairs(), MustGetEnv(), NewPostgreSessionStoreHandler(), startExpiredSessionCleanup(), ValidateCookieHashBlockPair(), EnvironmentVariableKey, sessionCleaner

### Community 70 - "GetAppearance"
Cohesion: 0.43
Nodes (6): NewAppearance(), GetAppearance(), NewAppearance(), TestNewAppearance(), WithAppearance(), Appearance

### Community 71 - "nav"
Cohesion: 0.33
Nodes (6): nav, label, notes, people, shopping, switch

### Community 72 - "nav"
Cohesion: 0.33
Nodes (6): nav, label, notes, people, shopping, switch

### Community 74 - "password"
Cohesion: 0.50
Nodes (4): password, confirm, enter, title

### Community 75 - "NormalizeInviteCode"
Cohesion: 0.50
Nodes (4): newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

### Community 79 - "error-page"
Cohesion: 0.50
Nodes (4): back, description, title, error-page

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 477 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `appearance`, `shopping`, `nav`, `home-page`, `password`, `login-page`, `notes`, `invite-page`, `registration-page`, `household`, `validation-errors`, `verify-email`, `reset-password-page`, `mail`, `forgot-password-page`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14702263779527558 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `forgot-password-page`, `appearance`, `shopping`, `nav`, `validation-errors`, `home-page`, `login-page`, `error-page`, `invite-page`, `household`, `registration-page`, `notes`, `verify-email`, `reset-password-page`, `mail`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._
- **Why does `LKMessage` connect `LKMessage` to `handlers.go`, `github.com/a-h/templ.Component`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `et` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._