# Graph Report - roommates  (2026-10-29)

## Corpus Check
- 228 files · ~139,613 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 46 file(s) not represented in the graph (top: .templ 29, (none) 6, .css 5)

## Summary
- 2231 nodes · 7153 edges · 131 communities (105 shown, 26 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 515 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `81a7e03d`
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
- CN
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
- tooltip.js
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
- floating_ui_dom.js
- T
- files
- validation-errors
- Queries
- notes
- card_templ.go
- verify-email
- verify-email
- reset-password-page
- reset-password-page
- use_hover.js
- mail
- github.com/jackc/pgx/v5/pgtype.UUID
- forgot-password-page
- forgot-password-page
- files
- .UploadFile
- appearance
- appearance
- InitGinEngine
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
- generate-locales-keys.go
- images.go
- .clean
- floating_ui_core.js
- bubble_templ.go
- ListPendingSettlementsRow
- tabbable.js
- useDismiss
- HouseholdOverview
- Session
- initRoutes
- Subscription
- DefaultService
- common
- chat
- chat
- use_transition_status.js
- Props
- composite.js
- generate-locales-languages.go
- radiogroup.js
- portal.js
- useFocus
- push
- push
- Message
- LK
- ChatBadge
- vendor-assets_test.go
- mounted
- mail
- registration-page
- T
- noRowsAs
- Renderer
- go_pkg_fmt
- chatMessages
- fakeStore
- Assets
- avatar.js

## God Nodes (most connected - your core abstractions)
1. `T()` - 67 edges
2. `CN()` - 64 edges
3. `membershipFromContext()` - 46 edges
4. `T()` - 45 edges
5. `Membership` - 43 edges
6. `expenses` - 42 edges
7. `expenses` - 42 edges
8. `RenderTemplStatus()` - 41 edges
9. `InternalError()` - 35 edges
10. `Button()` - 31 edges

## Surprising Connections (you probably didn't know these)
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go
- `TestExpiredTokenIsInvalid()` --calls--> `hashToken()`  [INFERRED]
  app/internal/services/services-email_test.go → app/internal/services/services-email.go
- `newTestService()` --calls--> `newFakeStore()`  [INFERRED]
  app/internal/services/services-households_test.go → app/internal/services/services-files_test.go
- `TestNewInviteCode()` --calls--> `newInviteCode()`  [INFERRED]
  app/internal/services/services-households_test.go → app/internal/services/services-households.go
- `newTestService()` --calls--> `NewDefaultService()`  [INFERRED]
  app/internal/services/services-households_test.go → app/internal/services/services.go

## Import Cycles
- None detected.

## Communities (131 total, 26 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.16
Nodes (33): TestSafeNextPath(), TestChat(), TestChatReply(), TestHouseholdEvents(), englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword() (+25 more)

### Community 1 - "handlers.go"
Cohesion: 0.13
Nodes (9): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ctxKey, ctxKey, Props, appearanceKey, directionKey (+1 more)

### Community 2 - "migrations.go"
Cohesion: 0.16
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.08
Nodes (23): title, en, app, email, error-page, logout, not-found-page, offline (+15 more)

### Community 4 - "et"
Cohesion: 0.07
Nodes (27): title, delete, link-invalid, save, back, description, title, et (+19 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.17
Nodes (18): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+10 more)

### Community 7 - "vars"
Cohesion: 0.04
Nodes (47): cmds, desc, preconditions, task, cmds, desc, env, SERVER_PORT (+39 more)

### Community 8 - "CN"
Cohesion: 0.10
Nodes (25): CN(), flattenClasses(), twMapTokens(), twReorder(), Action(), Description(), Title(), AspectRatio() (+17 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "field_templ.go"
Cohesion: 0.06
Nodes (50): ActionProps, DescriptionProps, Props, TitleProps, Variant, variantClasses(), Action(), childrenHTML() (+42 more)

### Community 11 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.05
Nodes (42): DeleteCookie(), GetCookie(), chatPath(), optionalUUID(), notice(), expenseFromForm(), expensePath(), expensesPath() (+34 more)

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
Nodes (32): Direction(), WithDirection(), GetUnreadChat(), ConfirmSettlementParams, ConsumeEmailTokenParams, CountActiveMembersRow, DeleteCheckedShoppingItemsParams, DeleteEmailTokensOfUserParams (+24 more)

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

### Community 31 - "tooltip.js"
Cohesion: 0.26
Nodes (22): arrowOf(), close(), contentFor(), destroy(), emitOpenChange(), init(), joinGroup(), leaveGroup() (+14 more)

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
Nodes (8): GetUserCredentialsRow, GetUserRow, InsertUserParams, LoginRequest, NormalizeEmail(), newEmailToken(), timestamptz(), newInviteCode()

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

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "T"
Cohesion: 0.15
Nodes (36): AttachedFile, Membership, ShoppingList, IfElse(), FormatMoney(), debtBreakdown(), ExpensePage(), HouseholdExpenses() (+28 more)

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

### Community 55 - "card_templ.go"
Cohesion: 0.29
Nodes (20): GetAppearance(), Alert(), Card(), Content(), Description(), Footer(), Header(), Title() (+12 more)

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

### Community 60 - "use_hover.js"
Cohesion: 0.09
Nodes (41): applySafePolygonPointerEventsMutation(), clearSafePolygonPointerEventsMutation(), createHoverInteraction(), emitClosed(), getDelay(), getRestMs(), hasIntersectingEdge(), isClickLikeOpenEvent() (+33 more)

### Community 61 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 62 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.05
Nodes (50): ListExpenseSharesRow, ListFilesRow, ListActiveInvitesRow, ListMessagesRow, ListNotesRow, GetShoppingListRow, ListShoppingItemsRow, EmailToken (+42 more)

### Community 63 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 64 - "forgot-password-page"
Cohesion: 0.15
Nodes (12): forgot-password-page, invite-page, description, sent-description, sent-title, submit, title, code (+4 more)

### Community 65 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, download, file, invalid (+7 more)

### Community 66 - ".UploadFile"
Cohesion: 0.11
Nodes (5): ReplyPreview, Store, New(), Upload, S3Store

### Community 67 - "appearance"
Cohesion: 0.22
Nodes (9): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, title (+1 more)

### Community 68 - "appearance"
Cohesion: 0.22
Nodes (9): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, title (+1 more)

### Community 69 - "InitGinEngine"
Cohesion: 0.13
Nodes (18): ensureDatabasePool(), main(), InitGinEngine(), GetAppBaseURL(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), IsDevelopment() (+10 more)

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

### Community 83 - "generate-locales-keys.go"
Cohesion: 0.22
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (11): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+3 more)

### Community 85 - ".clean"
Cohesion: 0.13
Nodes (12): GetExpenseRow, InsertExpenseParams, FormatCentsInput(), Expense, ExpenseInput, ParseCents(), splitEqual(), TestParseCents() (+4 more)

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "bubble_templ.go"
Cohesion: 0.23
Nodes (13): Bubble(), Content(), ContentClass(), Group(), Reactions(), variantClasses(), Align, ContentProps (+5 more)

### Community 88 - "ListPendingSettlementsRow"
Cohesion: 0.14
Nodes (11): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), TestNetDebts(), testUUID() (+3 more)

### Community 89 - "tabbable.js"
Cohesion: 0.14
Nodes (28): activeElement(), appendCandidates(), appendMatchingElements(), contains(), disableFocusInside(), enableFocusInside(), focusable(), getComposedChildren() (+20 more)

### Community 90 - "useDismiss"
Cohesion: 0.19
Nodes (25): contains(), getTarget(), isLastTraversableNode(), normalizeBubbles(), useDismiss(), addTargetEventListenerOnce(), close(), closeOnEscapeKeyDown() (+17 more)

### Community 91 - "HouseholdOverview"
Cohesion: 0.17
Nodes (5): GetHouseholdRow, ListHouseholdsOfUserRow, ListActiveMembersRow, ListShoppingListsRow, HouseholdOverview

### Community 92 - "Session"
Cohesion: 0.29
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 93 - "initRoutes"
Cohesion: 0.14
Nodes (17): initRoutes(), registerAssets(), NewAppearance(), NewGuestOnly(), NewOnlyAuthenticated(), NewHouseholdMember(), NewZerolog(), NewRateLimit() (+9 more)

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (8): Message, Pusher, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

### Community 95 - "DefaultService"
Cohesion: 0.20
Nodes (8): Hub, NewHub(), TestHub(), DefaultHandler, NewDefaultHandler(), Mailer, DefaultService, NewDefaultService()

### Community 96 - "common"
Cohesion: 0.50
Nodes (4): delete, link-invalid, save, common

### Community 97 - "chat"
Cohesion: 0.12
Nodes (16): attach, cancel-reply, copy, delete, empty, length, message, newest (+8 more)

### Community 98 - "chat"
Cohesion: 0.12
Nodes (16): attach, cancel-reply, copy, delete, empty, length, message, newest (+8 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 100 - "Props"
Cohesion: 0.38
Nodes (5): Variants(), Props, Size, Type, Variant

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "generate-locales-languages.go"
Cohesion: 0.27
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 105 - "radiogroup.js"
Cohesion: 0.28
Nodes (10): groupOf(), initGroup(), inputOf(), itemOf(), itemsOf(), requestValueChange(), setupItem(), syncByInput() (+2 more)

### Community 106 - "portal.js"
Cohesion: 0.39
Nodes (6): containerFor(), createOutsideGuards(), onFocus(), remove(), render(), setFocusManagerState()

### Community 107 - "useFocus"
Cohesion: 0.43
Nodes (6): isInsideEnabledTrigger(), matchesFocusVisible(), useFocus(), onBlur(), onFocus(), resetBlockedFocus()

### Community 108 - "push"
Cohesion: 0.18
Nodes (11): push, chat-file, chat-title, denied, disable, enable, expense-body, expense-title (+3 more)

### Community 109 - "push"
Cohesion: 0.18
Nodes (11): push, chat-file, chat-title, denied, disable, enable, expense-body, expense-title (+3 more)

### Community 111 - "Message"
Cohesion: 0.24
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 112 - "LK"
Cohesion: 0.14
Nodes (14): LK, DirectionAttributes(), N(), unescapeJSON(), N(), initialValue(), RadioGroup(), Content() (+6 more)

### Community 113 - "ChatBadge"
Cohesion: 0.48
Nodes (5): Badge(), badgeClasses(), sharePaidBadge(), Props, Variant

### Community 114 - "vendor-assets_test.go"
Cohesion: 0.36
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 124 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 125 - "T"
Cohesion: 0.21
Nodes (34): T(), Button(), T(), ValidationErrors(), Field(), Group(), Label(), Input() (+26 more)

### Community 126 - "noRowsAs"
Cohesion: 0.13
Nodes (5): GetInviteByCodeRow, hashToken(), FileMeta, FileTarget, noRowsAs()

### Community 135 - "chatMessages"
Cohesion: 0.14
Nodes (21): Chat, Message, MergeAttributes(), RandomID(), Avatar(), Content(), Footer(), Header() (+13 more)

### Community 136 - "fakeStore"
Cohesion: 0.14
Nodes (11): LoadLocales(), newFakeStore(), newFakePusher(), subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes(), TestTUnescapesTextButNotArgs() (+3 more)

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 715 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **26 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `forgot-password-page`, `files`, `chat`, `appearance`, `expenses`, `shopping`, `nav`, `validation-errors`, `home-page`, `login-page`, `push`, `household`, `notes`, `verify-email`, `reset-password-page`, `registration-page`, `mail`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13193740853315322 - nodes in this community are weakly interconnected._
- **Why does `LKMessage` connect `LKMessage` to `LK`, `handlers.go`, `T`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.08333333333333333 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `handlers.go`, `.UploadFile`, `chatMessages`, `fakeStore`, `newEmailToken`, `github.com/gin-gonic/gin.Context`, `context.Context`, `LK`, `T`, `.clean`, `card_templ.go`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `et` be split into smaller, more focused modules?**
  _Cohesion score 0.07142857142857142 - nodes in this community are weakly interconnected._