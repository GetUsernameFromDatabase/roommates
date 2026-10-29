# Graph Report - roommates  (2026-10-29)

## Corpus Check
- 198 files · ~101,053 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 39 file(s) not represented in the graph (top: .templ 21, (none) 7, .css 5)

## Summary
- 1813 nodes · 5952 edges · 120 communities (97 shown, 23 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 407 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c1a50c54`
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
- RenderTemplStatus
- LKMessage
- field_templ.go
- github.com/gin-gonic/gin.Context
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
- github.com/a-h/templ.Component
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
- membershipFromContext
- Button
- HouseholdExpenses
- Session
- NewAppearance
- Subscription
- Hub
- NormalizeInviteCode
- chat
- chat
- .UploadFile
- Props
- LK
- password
- Redirect
- NewInternationalization
- TestPushNotifications
- .shoppingAction
- push
- push
- .UpdateNote
- Membership
- registration-page
- ExpenseInput
- common
- offline
- Chat

## God Nodes (most connected - your core abstractions)
1. `T()` - 67 edges
2. `membershipFromContext()` - 46 edges
3. `Membership` - 43 edges
4. `expenses` - 42 edges
5. `expenses` - 42 edges
6. `RenderTemplStatus()` - 41 edges
7. `T()` - 41 edges
8. `CN()` - 40 edges
9. `InternalError()` - 35 edges
10. `Button()` - 30 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `InitGinEngine()`  [INFERRED]
  app/cmd/server/main.go → app/cmd/server/routes.go
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go
- `expensesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-expenses.go → app/internal/handlers/handlers-households.go
- `shoppingListPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-shopping.go → app/internal/handlers/handlers-households.go
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go

## Import Cycles
- None detected.

## Communities (120 total, 23 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.21
Nodes (27): TestSafeNextPath(), TestChat(), TestHouseholdEvents(), englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword(), TestVerifyEmail() (+19 more)

### Community 1 - "handlers.go"
Cohesion: 0.15
Nodes (6): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), appearanceKey, directionKey, unreadChatKey

### Community 2 - "migrations.go"
Cohesion: 0.16
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.10
Nodes (19): title, en, app, email, error-page, logout, not-found-page, password (+11 more)

### Community 4 - "et"
Cohesion: 0.08
Nodes (23): title, delete, link-invalid, save, back, description, title, et (+15 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.06
Nodes (38): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+30 more)

### Community 7 - "vars"
Cohesion: 0.05
Nodes (38): cmds, desc, task, cmds, desc, env, SERVER_PORT, SHADCN_TEMPL_DEV (+30 more)

### Community 8 - "RenderTemplStatus"
Cohesion: 0.29
Nodes (6): expensePath(), expensesPath(), pageError(), RenderTemplStatus(), ParseUUID(), NotFound()

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "field_templ.go"
Cohesion: 0.08
Nodes (40): CN(), DirectionAttributes(), flattenClasses(), MergeAttributes(), twMapTokens(), twReorder(), AspectRatio(), Action() (+32 more)

### Community 11 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.12
Nodes (5): notice(), sessionUserID(), InternalError(), RenderTempl(), GetLocaleString()

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
Nodes (34): Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), ConfirmSettlementParams, ConsumeEmailTokenParams, CountActiveMembersRow, DeleteCheckedShoppingItemsParams (+26 more)

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
Cohesion: 0.09
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
Cohesion: 0.16
Nodes (47): ListHouseholdsOfUserRow, AttachedFile, HouseholdOverview, Membership, ShoppingList, IfElse(), T(), T() (+39 more)

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

### Community 55 - "github.com/a-h/templ.Component"
Cohesion: 0.22
Nodes (24): Card(), Content(), Description(), Footer(), Header(), Title(), HeaderComponent(), HtmlWrap() (+16 more)

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
Nodes (50): ListExpenseSharesRow, ListFilesRow, GetHouseholdRow, ListActiveInvitesRow, ListMessagesRow, ListNotesRow, GetShoppingListRow, ListShoppingItemsRow (+42 more)

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
Cohesion: 0.11
Nodes (5): pushPreview(), Store, New(), Upload, S3Store

### Community 67 - "appearance"
Cohesion: 0.22
Nodes (9): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, title (+1 more)

### Community 68 - "appearance"
Cohesion: 0.22
Nodes (9): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, title (+1 more)

### Community 69 - "DefaultService"
Cohesion: 0.13
Nodes (21): InitGinEngine(), initRoutes(), registerAssets(), GetAppBaseURL(), IsDevelopment(), DefaultHandler, NewDefaultHandler(), Mailer (+13 more)

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
Cohesion: 0.12
Nodes (4): GetInviteByCodeRow, hashToken(), noRowsAs(), rowsAffected()

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

### Community 83 - "Renderer"
Cohesion: 0.38
Nodes (4): NewHouseholdMember(), notFound(), New(), Renderer

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (11): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+3 more)

### Community 85 - ".clean"
Cohesion: 0.13
Nodes (10): InsertExpenseParams, ListActiveMembersRow, splitEqual(), TestNetDebts(), TestSplitEqual(), testUUID(), cleanFileName(), TestCleanFileName() (+2 more)

### Community 86 - "mailer.go"
Cohesion: 0.15
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 87 - "GetDatabaseURL"
Cohesion: 0.12
Nodes (15): ensureDatabasePool(), main(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), MustGetEnv(), TestGetDatabaseURLEscapes(), LoadLocales() (+7 more)

### Community 88 - "ListPendingSettlementsRow"
Cohesion: 0.17
Nodes (9): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), ListDebtItemsParams, ListPendingSettlementsParams (+1 more)

### Community 89 - "membershipFromContext"
Cohesion: 0.16
Nodes (5): shoppingItemPath(), membershipFromContext(), shoppingListPath(), Note, ShoppingItem

### Community 90 - "Button"
Cohesion: 0.36
Nodes (14): RandomID(), Button(), ValidationErrors(), Field(), Group(), Label(), Input(), Textarea() (+6 more)

### Community 91 - "HouseholdExpenses"
Cohesion: 0.17
Nodes (14): GetExpenseRow, FormatCentsInput(), Expense, ParseCents(), TestParseCents(), householdPath(), FormatMoney(), debtBreakdown() (+6 more)

### Community 92 - "Session"
Cohesion: 0.17
Nodes (4): LogValidationErrors(), gorillaSessionGetter, Session, SessionInfo

### Community 93 - "NewAppearance"
Cohesion: 0.43
Nodes (6): NewAppearance(), GetAppearance(), NewAppearance(), TestNewAppearance(), WithAppearance(), Appearance

### Community 94 - "Subscription"
Cohesion: 0.17
Nodes (7): Message, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

### Community 95 - "Hub"
Cohesion: 0.39
Nodes (3): Hub, NewHub(), TestHub()

### Community 96 - "NormalizeInviteCode"
Cohesion: 0.50
Nodes (4): newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

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
Cohesion: 0.38
Nodes (5): LK, N(), unescapeJSON(), N(), appearanceSelect()

### Community 102 - "password"
Cohesion: 0.50
Nodes (4): password, confirm, enter, title

### Community 103 - "Redirect"
Cohesion: 0.23
Nodes (4): peoplePath(), safeNextPath(), Redirect(), FromGinContext()

### Community 105 - "NewInternationalization"
Cohesion: 0.17
Nodes (9): DeleteCookie(), GetCookie(), GetHeader(), localeMissing(), NewInternationalization(), setLocaleContext(), CookieKey, HTTPHeaderKey (+1 more)

### Community 106 - "TestPushNotifications"
Cohesion: 0.17
Nodes (9): newFakeStore(), newFakePusher(), subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes(), fakePusher, fakeStore (+1 more)

### Community 108 - "push"
Cohesion: 0.18
Nodes (11): push, chat-file, chat-title, denied, disable, enable, expense-body, expense-title (+3 more)

### Community 109 - "push"
Cohesion: 0.18
Nodes (11): push, chat-file, chat-title, denied, disable, enable, expense-body, expense-title (+3 more)

### Community 110 - ".UpdateNote"
Cohesion: 0.24
Nodes (7): chatPath(), filesPath(), householdPath(), noteFromForm(), noteLengthMessage(), notePath(), notesPath()

### Community 112 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 114 - "common"
Cohesion: 0.50
Nodes (4): delete, link-invalid, save, common

### Community 115 - "offline"
Cohesion: 0.50
Nodes (4): offline, description, retry, title

### Community 116 - "Chat"
Cohesion: 0.67
Nodes (3): Chat, Message, MessageFile

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 673 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `T()` connect `T` to `handlers.go`, `S3Store`, `LK`, `Redirect`, `RenderTemplStatus`, `newEmailToken`, `github.com/gin-gonic/gin.Context`, `.shoppingAction`, `.UpdateNote`, `context.Context`, `ExpenseInput`, `github.com/a-h/templ.Component`, `GetDatabaseURL`, `membershipFromContext`, `Button`, `HouseholdExpenses`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14853801169590644 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `household`, `expenses`, `shopping`, `validation-errors`, `home-page`, `login-page`, `invite-page`, `notes`, `verify-email`, `reset-password-page`, `mail`, `forgot-password-page`, `files`, `appearance`, `nav`, `chat`, `password`, `push`, `registration-page`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `expenses`, `household`, `shopping`, `home-page`, `login-page`, `notes`, `invite-page`, `registration-page`, `files`, `validation-errors`, `verify-email`, `reset-password-page`, `mail`, `forgot-password-page`, `appearance`, `nav`, `chat`, `push`, `common`, `offline`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Should `et` be split into smaller, more focused modules?**
  _Cohesion score 0.08333333333333333 - nodes in this community are weakly interconnected._