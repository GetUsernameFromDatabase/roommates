# Graph Report - roommates  (2026-10-29)

## Corpus Check
- 186 files · ~96,769 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 39 file(s) not represented in the graph (top: .templ 21, (none) 7, .css 5)

## Summary
- 1723 nodes · 5561 edges · 105 communities (86 shown, 19 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 396 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a7d3b3fa`
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
- field_templ.go
- github.com/a-h/templ.Component
- components.json
- block-secret-files.sh
- go
- context.Context
- expenses
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
- expenses
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
- T
- files
- validation-errors
- Queries
- notes
- Invite
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
- initRoutes
- environment
- nav
- nav
- rowsAffected
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
- images.go
- .clean
- mailer.go
- GetDatabaseURL
- ListPendingSettlementsRow
- generate-locales-languages.go
- github.com/a-h/templ.Component
- vendor-assets_test.go
- Session
- NewAppearance
- InitGinEngine
- Hub
- NormalizeInviteCode
- chat
- chat
- .UploadFile
- Props
- LK
- password
- ListMessageFilesParams

## God Nodes (most connected - your core abstractions)
1. `T()` - 61 edges
2. `membershipFromContext()` - 46 edges
3. `Membership` - 43 edges
4. `expenses` - 42 edges
5. `expenses` - 42 edges
6. `RenderTemplStatus()` - 41 edges
7. `CN()` - 40 edges
8. `T()` - 40 edges
9. `InternalError()` - 33 edges
10. `LKMessage` - 29 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `InitLogger()`  [EXTRACTED]
  app/cmd/server/main.go → app/internal/logger/logger.go
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go
- `TestExpiredTokenIsInvalid()` --calls--> `hashToken()`  [INFERRED]
  app/internal/services/services-email_test.go → app/internal/services/services-email.go
- `TestSplitEqual()` --calls--> `splitEqual()`  [INFERRED]
  app/internal/services/services-expenses_test.go → app/internal/services/services-expenses.go
- `newTestService()` --calls--> `newFakeStore()`  [INFERRED]
  app/internal/services/services-households_test.go → app/internal/services/services-files_test.go

## Import Cycles
- None detected.

## Communities (105 total, 19 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.22
Nodes (28): TestSafeNextPath(), TestChat(), TestHouseholdEvents(), englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword(), TestVerifyEmail() (+20 more)

### Community 1 - "handlers.go"
Cohesion: 0.15
Nodes (6): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), appearanceKey, directionKey, unreadChatKey

### Community 2 - "migrations.go"
Cohesion: 0.16
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.08
Nodes (23): title, delete, link-invalid, save, en, app, common, email (+15 more)

### Community 4 - "et"
Cohesion: 0.08
Nodes (25): title, delete, link-invalid, save, back, description, title, et (+17 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.17
Nodes (18): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+10 more)

### Community 7 - "vars"
Cohesion: 0.05
Nodes (38): cmds, desc, task, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV (+30 more)

### Community 8 - "generate-locales-keys.go"
Cohesion: 0.10
Nodes (8): InitLogger(), flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "field_templ.go"
Cohesion: 0.08
Nodes (41): CN(), DirectionAttributes(), flattenClasses(), MergeAttributes(), twMapTokens(), twReorder(), AspectRatio(), Action() (+33 more)

### Community 11 - "github.com/a-h/templ.Component"
Cohesion: 0.06
Nodes (41): DeleteCookie(), GetCookie(), chatPath(), notice(), expenseFromForm(), expensePath(), expensesPath(), filesPath() (+33 more)

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
Cohesion: 0.03
Nodes (33): Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetLocaleString(), ConfirmSettlementParams, ConsumeEmailTokenParams, CountActiveMembersRow (+25 more)

### Community 16 - "expenses"
Cohesion: 0.05
Nodes (42): expenses, add, amount, back, balances, cancel, category, category-groceries (+34 more)

### Community 17 - "household"
Cohesion: 0.09
Nodes (22): household, back, balance-not-zero, create-invite, expires, invite-email-sent, invite-email-submit, invite-email-title (+14 more)

### Community 18 - "household"
Cohesion: 0.09
Nodes (22): household, back, balance-not-zero, create-invite, expires, invite-email-sent, invite-email-submit, invite-email-title (+14 more)

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

### Community 37 - "expenses"
Cohesion: 0.05
Nodes (42): expenses, add, amount, back, balances, cancel, category, category-groceries (+34 more)

### Community 38 - "shopping"
Cohesion: 0.09
Nodes (22): shopping, add, back, bought-by, buying, buying-by, check, clear (+14 more)

### Community 39 - "shopping"
Cohesion: 0.09
Nodes (22): shopping, add, back, bought-by, buying, buying-by, check, clear (+14 more)

### Community 40 - "newEmailToken"
Cohesion: 0.11
Nodes (9): GetAppBaseURL(), GetUserCredentialsRow, GetUserRow, InsertUserParams, LoginRequest, NormalizeEmail(), newEmailToken(), timestamptz() (+1 more)

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
Cohesion: 0.11
Nodes (18): notes, back, body, body-hint, cancel, delete, edit, edited-by (+10 more)

### Community 47 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 48 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 49 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 50 - "T"
Cohesion: 0.11
Nodes (54): GetExpenseRow, Expense, AttachedFile, Membership, ShoppingList, IfElse(), RandomID(), T() (+46 more)

### Community 51 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, download, file, invalid (+7 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "Queries"
Cohesion: 0.14
Nodes (5): Queries, New(), checkSettled(), ownerActsOnMember(), DBTX

### Community 54 - "notes"
Cohesion: 0.11
Nodes (18): notes, back, body, body-hint, cancel, delete, edit, edited-by (+10 more)

### Community 55 - "Invite"
Cohesion: 0.48
Nodes (12): Card(), Content(), Description(), Header(), Title(), Notice(), Invite(), guestShell() (+4 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 59 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 60 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 61 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 62 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.05
Nodes (46): ListExpenseSharesRow, ListFilesRow, GetHouseholdRow, ListActiveInvitesRow, ListMessagesRow, ListNotesRow, GetShoppingListRow, ListShoppingItemsRow (+38 more)

### Community 63 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 64 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 65 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, download, file, invalid (+7 more)

### Community 66 - "S3Store"
Cohesion: 0.10
Nodes (7): Chat, Store, New(), Message, MessageFile, Upload, S3Store

### Community 67 - "appearance"
Cohesion: 0.22
Nodes (9): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, title (+1 more)

### Community 68 - "appearance"
Cohesion: 0.22
Nodes (9): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, title (+1 more)

### Community 69 - "initRoutes"
Cohesion: 0.13
Nodes (13): initRoutes(), registerAssets(), NewGuestOnly(), NewOnlyAuthenticated(), NewZerolog(), NewRateLimit(), TestRateLimit(), NewCrossOriginProtection() (+5 more)

### Community 70 - "environment"
Cohesion: 0.18
Nodes (11): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment, ports (+3 more)

### Community 71 - "nav"
Cohesion: 0.22
Nodes (9): nav, chat, expenses, files, label, notes, people, shopping (+1 more)

### Community 72 - "nav"
Cohesion: 0.22
Nodes (9): nav, chat, expenses, files, label, notes, people, shopping (+1 more)

### Community 73 - "rowsAffected"
Cohesion: 0.08
Nodes (9): GetInviteByCodeRow, hashToken(), cleanFileName(), TestCleanFileName(), noRowsAs(), NoteInput, rowsAffected(), trimToLength() (+1 more)

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

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (11): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+3 more)

### Community 85 - ".clean"
Cohesion: 0.16
Nodes (8): InsertExpenseParams, ListActiveMembersRow, FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestParseCents(), expenseShare

### Community 86 - "mailer.go"
Cohesion: 0.13
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 87 - "GetDatabaseURL"
Cohesion: 0.24
Nodes (8): ensureDatabasePool(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), MustGetEnv(), TestGetDatabaseURLEscapes(), ValidateCookieHashBlockPair(), EnvironmentVariableKey

### Community 88 - "ListPendingSettlementsRow"
Cohesion: 0.13
Nodes (12): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), TestNetDebts(), TestSplitEqual() (+4 more)

### Community 89 - "generate-locales-languages.go"
Cohesion: 0.40
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 90 - "github.com/a-h/templ.Component"
Cohesion: 0.26
Nodes (26): ListHouseholdsOfUserRow, HouseholdOverview, Button(), T(), ValidationErrors(), Field(), Group(), Label() (+18 more)

### Community 91 - "vendor-assets_test.go"
Cohesion: 0.36
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 92 - "Session"
Cohesion: 0.23
Nodes (4): LogValidationErrors(), gorillaSessionGetter, Session, SessionInfo

### Community 93 - "NewAppearance"
Cohesion: 0.43
Nodes (6): NewAppearance(), GetAppearance(), NewAppearance(), TestNewAppearance(), WithAppearance(), Appearance

### Community 94 - "InitGinEngine"
Cohesion: 0.24
Nodes (11): main(), InitGinEngine(), IsDevelopment(), NewHub(), TestHub(), LoadLocales(), Mailer, New() (+3 more)

### Community 95 - "Hub"
Cohesion: 0.22
Nodes (5): Hub, DefaultHandler, NewDefaultHandler(), newFakeStore(), fakeStore

### Community 96 - "NormalizeInviteCode"
Cohesion: 0.67
Nodes (3): NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

### Community 97 - "chat"
Cohesion: 0.18
Nodes (11): attach, delete, empty, length, message, newest, older, send (+3 more)

### Community 98 - "chat"
Cohesion: 0.18
Nodes (11): attach, delete, empty, length, message, newest, older, send (+3 more)

### Community 100 - "Props"
Cohesion: 0.38
Nodes (5): Variants(), Props, Size, Type, Variant

### Community 101 - "LK"
Cohesion: 0.47
Nodes (4): LK, N(), N(), householdNavItem()

### Community 102 - "password"
Cohesion: 0.50
Nodes (4): password, confirm, enter, title

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 635 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `forgot-password-page`, `files`, `chat`, `appearance`, `expenses`, `password`, `shopping`, `nav`, `validation-errors`, `home-page`, `login-page`, `invite-page`, `household`, `notes`, `verify-email`, `reset-password-page`, `mail`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14807146691204662 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `handlers.go`, `LK`, `newEmailToken`, `github.com/a-h/templ.Component`, `context.Context`, `Invite`, `github.com/a-h/templ.Component`, `InitGinEngine`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.08333333333333333 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `chat`, `appearance`, `shopping`, `nav`, `home-page`, `login-page`, `notes`, `invite-page`, `expenses`, `household`, `registration-page`, `files`, `validation-errors`, `verify-email`, `reset-password-page`, `mail`, `forgot-password-page`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `et` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._