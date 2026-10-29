# Graph Report - roommates  (2026-10-29)

## Corpus Check
- 172 files · ~87,019 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 38 file(s) not represented in the graph (top: .templ 20, (none) 7, .css 5)

## Summary
- 1625 nodes · 4971 edges · 97 communities (79 shown, 18 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 334 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2867bdcc`
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
- DefaultService
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
- HouseholdOverview
- vendor-assets_test.go
- Session
- NewAppearance
- go_pkg_fmt
- fakeStore
- NormalizeInviteCode

## God Nodes (most connected - your core abstractions)
1. `T()` - 54 edges
2. `membershipFromContext()` - 41 edges
3. `RenderTemplStatus()` - 38 edges
4. `expenses` - 38 edges
5. `expenses` - 38 edges
6. `CN()` - 37 edges
7. `T()` - 37 edges
8. `Membership` - 35 edges
9. `InternalError()` - 30 edges
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
- `TestNetDebts()` --calls--> `netDebts()`  [INFERRED]
  app/internal/services/services-expenses_test.go → app/internal/services/services-expenses.go

## Import Cycles
- None detected.

## Communities (97 total, 18 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.20
Nodes (27): TestSafeNextPath(), englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword(), TestVerifyEmail(), tokenOf(), TestExpenses() (+19 more)

### Community 1 - "handlers.go"
Cohesion: 0.15
Nodes (5): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), appearanceKey, directionKey

### Community 2 - "migrations.go"
Cohesion: 0.16
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

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
Cohesion: 0.17
Nodes (18): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+10 more)

### Community 7 - "vars"
Cohesion: 0.05
Nodes (38): cmds, desc, task, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV (+30 more)

### Community 8 - "generate-locales-keys.go"
Cohesion: 0.14
Nodes (8): InitLogger(), flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.05
Nodes (127): LK, AttachedFile, Membership, ShoppingList, CN(), DirectionAttributes(), flattenClasses(), IfElse() (+119 more)

### Community 11 - "github.com/a-h/templ.Component"
Cohesion: 0.06
Nodes (40): DeleteCookie(), GetCookie(), notice(), expenseFromForm(), expensePath(), expensesPath(), filesPath(), shoppingItemPath() (+32 more)

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
Nodes (27): Direction(), WithDirection(), ConfirmSettlementParams, ConsumeEmailTokenParams, CountActiveMembersRow, DeleteCheckedShoppingItemsParams, DeleteEmailTokensOfUserParams, DeleteExpenseParams (+19 more)

### Community 16 - "expenses"
Cohesion: 0.05
Nodes (38): expenses, add-from-item, amount, back, balances, cancel, category, category-groceries (+30 more)

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
Nodes (38): expenses, add-from-item, amount, back, balances, cancel, category, category-groceries (+30 more)

### Community 38 - "shopping"
Cohesion: 0.09
Nodes (23): shopping, add, back, bought-by, buying, buying-by, check, clear (+15 more)

### Community 39 - "shopping"
Cohesion: 0.09
Nodes (23): shopping, add, back, bought-by, buying, buying-by, check, clear (+15 more)

### Community 40 - "newEmailToken"
Cohesion: 0.11
Nodes (8): GetAppBaseURL(), GetUserCredentialsRow, GetUserRow, InsertUserParams, LoginRequest, NormalizeEmail(), newEmailToken(), timestamptz()

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
Cohesion: 0.12
Nodes (16): files, added-by, all, attachments, count, delete, download, file (+8 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "Queries"
Cohesion: 0.14
Nodes (5): Queries, New(), checkSettled(), ownerActsOnMember(), DBTX

### Community 54 - "notes"
Cohesion: 0.10
Nodes (20): notes, back, body, body-hint, cancel, delete, edit, edited-by (+12 more)

### Community 55 - "noRowsAs"
Cohesion: 0.21
Nodes (3): GetInviteByCodeRow, hashToken(), noRowsAs()

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
Cohesion: 0.06
Nodes (39): ListExpenseSharesRow, ListFilesRow, ListActiveInvitesRow, ListNotesRow, GetShoppingListRow, ListShoppingItemsRow, EmailToken, Expense (+31 more)

### Community 63 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 64 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 65 - "files"
Cohesion: 0.12
Nodes (16): files, added-by, all, attachments, count, delete, download, file (+8 more)

### Community 66 - "S3Store"
Cohesion: 0.11
Nodes (5): FileMeta, FileTarget, Store, New(), S3Store

### Community 67 - "appearance"
Cohesion: 0.20
Nodes (10): mode, mode-dark, mode-light, mode-system, save, theme, theme-home, theme-plain (+2 more)

### Community 68 - "appearance"
Cohesion: 0.20
Nodes (10): mode, mode-dark, mode-light, mode-system, save, theme, theme-home, theme-plain (+2 more)

### Community 69 - "DefaultService"
Cohesion: 0.13
Nodes (22): main(), InitGinEngine(), initRoutes(), registerAssets(), IsDevelopment(), DefaultHandler, NewDefaultHandler(), LoadLocales() (+14 more)

### Community 70 - "environment"
Cohesion: 0.18
Nodes (11): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment, ports (+3 more)

### Community 71 - "nav"
Cohesion: 0.25
Nodes (8): nav, expenses, files, label, notes, people, shopping, switch

### Community 72 - "nav"
Cohesion: 0.25
Nodes (8): nav, expenses, files, label, notes, people, shopping, switch

### Community 73 - "rowsAffected"
Cohesion: 0.11
Nodes (6): cleanFileName(), TestCleanFileName(), NoteInput, rowsAffected(), trimToLength(), NoteForm

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
Cohesion: 0.14
Nodes (10): GetExpenseRow, InsertExpenseParams, FormatCentsInput(), Expense, ExpenseInput, ParseCents(), splitEqual(), TestParseCents() (+2 more)

### Community 86 - "mailer.go"
Cohesion: 0.15
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 87 - "GetDatabaseURL"
Cohesion: 0.15
Nodes (11): ensureDatabasePool(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), MustGetEnv(), TestGetDatabaseURLEscapes(), NewPostgreSessionStoreHandler(), startExpiredSessionCleanup() (+3 more)

### Community 88 - "ListPendingSettlementsRow"
Cohesion: 0.18
Nodes (8): ListDebtFlowsRow, ListExpensesRow, ListPendingSettlementsRow, ExpensesOverview, netDebts(), ListDebtFlowsParams, ListPendingSettlementsParams, Debt

### Community 89 - "generate-locales-languages.go"
Cohesion: 0.27
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 90 - "HouseholdOverview"
Cohesion: 0.17
Nodes (5): GetHouseholdRow, ListHouseholdsOfUserRow, ListActiveMembersRow, ListShoppingListsRow, HouseholdOverview

### Community 91 - "vendor-assets_test.go"
Cohesion: 0.36
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 92 - "Session"
Cohesion: 0.29
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 93 - "NewAppearance"
Cohesion: 0.40
Nodes (5): NewAppearance(), NewAppearance(), TestNewAppearance(), WithAppearance(), Appearance

### Community 96 - "NormalizeInviteCode"
Cohesion: 0.50
Nodes (4): newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 612 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `appearance`, `shopping`, `nav`, `home-page`, `login-page`, `notes`, `invite-page`, `expenses`, `household`, `registration-page`, `files`, `validation-errors`, `verify-email`, `reset-password-page`, `mail`, `forgot-password-page`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14638128060946853 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `github.com/a-h/templ.Component` to `handlers.go`, `DefaultService`, `newEmailToken`, `github.com/a-h/templ.Component`, `context.Context`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `forgot-password-page`, `files`, `appearance`, `expenses`, `shopping`, `nav`, `validation-errors`, `home-page`, `login-page`, `invite-page`, `household`, `registration-page`, `notes`, `verify-email`, `reset-password-page`, `mail`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `et` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._