# Graph Report - roommates  (2026-10-29)

## Corpus Check
- 239 files · ~149,323 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 47 file(s) not represented in the graph (top: .templ 30, (none) 6, .css 5)

## Summary
- 2307 nodes · 7599 edges · 147 communities (115 shown, 32 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 544 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `edc69d1a`
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
- Membership
- LKMessage
- github.com/a-h/templ.Attributes
- membershipFromContext
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
- models.go
- notes
- ParseUUID
- ListActiveMembersRow
- floating_ui_dom.js
- T
- files
- validation-errors
- Queries
- notes
- InternalError
- verify-email
- verify-email
- reset-password-page
- reset-password-page
- use_hover.js
- github.com/gin-gonic/gin.Context
- github.com/jackc/pgx/v5/pgtype.UUID
- CN
- github.com/a-h/templ.Component
- files
- .UploadFile
- appearance
- appearance
- NormalizeInviteCode
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
- Hub
- images.go
- generate-locales-keys.go
- floating_ui_core.js
- RenderTemplStatus
- Debt
- tabbable.js
- useDismiss
- avatar
- sessions.go
- InitGinEngine
- Subscription
- mail
- NewInternationalization
- chat
- chat
- use_transition_status.js
- Props
- Message
- composite.js
- .clean
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- Message
- invite-page
- ChatBadge
- Renderer
- GetAppBaseURL
- go_pkg_fmt
- mounted
- .notifyChat
- registration-page
- .UpdateNote
- mail
- generate-locales-languages.go
- Redirect
- .loadShoppingItem
- ListNotesRow
- login-page
- invite-page
- vendor-assets_test.go
- Assets
- avatar.js
- LK
- forgot-password-page
- avatar
- common
- error-page
- offline

## God Nodes (most connected - your core abstractions)
1. `T()` - 74 edges
2. `CN()` - 66 edges
3. `T()` - 50 edges
4. `membershipFromContext()` - 46 edges
5. `Membership` - 43 edges
6. `RenderTemplStatus()` - 42 edges
7. `expenses` - 42 edges
8. `expenses` - 42 edges
9. `InternalError()` - 39 edges
10. `Button()` - 33 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `InitGinEngine()`  [INFERRED]
  app/cmd/server/main.go → app/cmd/server/routes.go
- `chatPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-chat.go → app/internal/handlers/handlers-households.go
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go
- `expensesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-expenses.go → app/internal/handlers/handlers-households.go
- `filesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-files.go → app/internal/handlers/handlers-households.go

## Import Cycles
- None detected.

## Communities (147 total, 32 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.20
Nodes (30): TestSafeNextPath(), TestAvatars(), TestChat(), TestChatReply(), TestHouseholdEvents(), englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail() (+22 more)

### Community 1 - "handlers.go"
Cohesion: 0.13
Nodes (10): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ctxKey, ctxKey, Props, appearanceKey, directionKey (+2 more)

### Community 2 - "migrations.go"
Cohesion: 0.16
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.07
Nodes (28): title, en, app, email, logout, not-found-page, offline, password (+20 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (32): title, delete, link-invalid, save, back, description, title, et (+24 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.17
Nodes (18): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+10 more)

### Community 7 - "vars"
Cohesion: 0.04
Nodes (47): cmds, desc, preconditions, task, cmds, desc, env, SERVER_PORT (+39 more)

### Community 8 - "Membership"
Cohesion: 0.14
Nodes (38): GetExpenseRow, Expense, AttachedFile, Membership, ShoppingList, IfElse(), FormatMoney(), debtBreakdown() (+30 more)

### Community 9 - "LKMessage"
Cohesion: 0.09
Nodes (14): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), validateEmail(), validateNewPassword(), validatePasswordMatch() (+6 more)

### Community 10 - "github.com/a-h/templ.Attributes"
Cohesion: 0.04
Nodes (62): ActionProps, DescriptionProps, TitleProps, Size, variantClasses(), childrenHTML(), Content(), Description() (+54 more)

### Community 11 - "membershipFromContext"
Cohesion: 0.21
Nodes (3): membershipFromContext(), RenderTempl(), Note

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
Nodes (35): ListHouseholdsOfUserRow, Store, New(), Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetHouseholds() (+27 more)

### Community 16 - "expenses"
Cohesion: 0.05
Nodes (42): expenses, add, amount, back, balances, cancel, category, category-groceries (+34 more)

### Community 17 - "household"
Cohesion: 0.08
Nodes (25): household, back, balance-not-zero, copied, copy-code, copy-link, create-invite, expires (+17 more)

### Community 18 - "household"
Cohesion: 0.08
Nodes (25): household, back, balance-not-zero, copied, copy-code, copy-link, create-invite, expires (+17 more)

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

### Community 45 - "models.go"
Cohesion: 0.08
Nodes (26): ListMessagesRow, ListShoppingItemsRow, EmailToken, Expense, ExpenseShare, File, FilePurge, GetFileParams (+18 more)

### Community 46 - "notes"
Cohesion: 0.11
Nodes (18): notes, back, body, body-hint, cancel, delete, edit, edited-by (+10 more)

### Community 47 - "ParseUUID"
Cohesion: 0.21
Nodes (5): expenseFromForm(), expensePath(), expensesPath(), pageError(), ParseUUID()

### Community 48 - "ListActiveMembersRow"
Cohesion: 0.18
Nodes (5): GetHouseholdRow, ListActiveInvitesRow, ListActiveMembersRow, ListShoppingListsRow, HouseholdOverview

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "T"
Cohesion: 0.21
Nodes (34): GetUserRow, T(), Button(), T(), ValidationErrors(), Field(), Group(), Label() (+26 more)

### Community 51 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, download, file, invalid (+7 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "Queries"
Cohesion: 0.12
Nodes (10): Queries, New(), DefaultHandler, NewDefaultHandler(), Mailer, New(), Pusher, DefaultService (+2 more)

### Community 54 - "notes"
Cohesion: 0.11
Nodes (18): notes, back, body, body-hint, cancel, delete, edit, edited-by (+10 more)

### Community 55 - "InternalError"
Cohesion: 0.17
Nodes (6): notice(), settingsPath(), householdPath(), sessionUserID(), InternalError(), GetLocaleString()

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

### Community 62 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.05
Nodes (29): ListExpenseSharesRow, GetShoppingListRow, avatarURL(), ConsumeEmailTokenParams, CountActiveMembersRow, DeleteCheckedShoppingItemsParams, DeleteEmailTokensOfUserParams, DeleteMessageParams (+21 more)

### Community 63 - "CN"
Cohesion: 0.08
Nodes (49): CN(), DirectionAttributes(), flattenClasses(), MergeAttributes(), RandomID(), twMapTokens(), twReorder(), Action() (+41 more)

### Community 64 - "github.com/a-h/templ.Component"
Cohesion: 0.21
Nodes (21): Props, Variant, Login, Register, Alert(), variantClasses(), Card(), Content() (+13 more)

### Community 65 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, download, file, invalid (+7 more)

### Community 66 - ".UploadFile"
Cohesion: 0.22
Nodes (4): ListFilesRow, FileMeta, FileTarget, ListFilesParams

### Community 67 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 68 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 69 - "NormalizeInviteCode"
Cohesion: 0.50
Nodes (4): newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

### Community 70 - "environment"
Cohesion: 0.18
Nodes (11): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment, ports (+3 more)

### Community 71 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 72 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 73 - "ListNotesRow"
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

### Community 83 - "Hub"
Cohesion: 0.19
Nodes (5): Hub, NewHub(), TestHub(), newFakeStore(), fakeStore

### Community 84 - "images.go"
Cohesion: 0.12
Nodes (12): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+4 more)

### Community 85 - "generate-locales-keys.go"
Cohesion: 0.22
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "RenderTemplStatus"
Cohesion: 0.26
Nodes (6): chatPath(), optionalUUID(), IsRequestHTMX(), RenderTemplStatus(), swapMainOnError(), NotFound()

### Community 88 - "Debt"
Cohesion: 0.11
Nodes (13): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, checkSettled(), Debt, ExpensesOverview, netDebts(), TestNetDebts() (+5 more)

### Community 89 - "tabbable.js"
Cohesion: 0.14
Nodes (28): activeElement(), appendCandidates(), appendMatchingElements(), contains(), disableFocusInside(), enableFocusInside(), focusable(), getComposedChildren() (+20 more)

### Community 90 - "useDismiss"
Cohesion: 0.19
Nodes (25): contains(), getTarget(), isLastTraversableNode(), normalizeBubbles(), useDismiss(), addTargetEventListenerOnce(), close(), closeOnEscapeKeyDown() (+17 more)

### Community 91 - "avatar"
Cohesion: 0.25
Nodes (8): crop-hint, description, invalid, remove, title, upload, zoom, avatar

### Community 92 - "sessions.go"
Cohesion: 0.29
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 93 - "InitGinEngine"
Cohesion: 0.12
Nodes (20): InitGinEngine(), initRoutes(), registerAssets(), NewAppearance(), NewGuestOnly(), NewOnlyAuthenticated(), NewZerolog(), NewRateLimit() (+12 more)

### Community 94 - "Subscription"
Cohesion: 0.11
Nodes (14): Message, Subscription, newFakePusher(), subscribeTestDevice(), takePushes(), TestPushNotifications(), TestValidPushSubscription(), wantPushes() (+6 more)

### Community 95 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 96 - "NewInternationalization"
Cohesion: 0.17
Nodes (9): DeleteCookie(), GetCookie(), GetHeader(), localeMissing(), NewInternationalization(), setLocaleContext(), CookieKey, HTTPHeaderKey (+1 more)

### Community 97 - "chat"
Cohesion: 0.12
Nodes (17): attach, cancel-reply, copy, delete, empty, files-too-large, length, message (+9 more)

### Community 98 - "chat"
Cohesion: 0.12
Nodes (17): attach, cancel-reply, copy, delete, empty, files-too-large, length, message (+9 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 100 - "Props"
Cohesion: 0.38
Nodes (5): Variants(), Props, Size, Type, Variant

### Community 101 - "Message"
Cohesion: 0.48
Nodes (7): Message, MessageFile, endsRun(), sameRun(), splitFiles(), startsRun(), TestChatRuns()

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - ".clean"
Cohesion: 0.19
Nodes (9): InsertExpenseParams, FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestParseCents(), TestSplitEqual(), householdPath() (+1 more)

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

### Community 112 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 113 - "ChatBadge"
Cohesion: 0.29
Nodes (8): Badge(), badgeClasses(), sharePaidBadge(), ChatBadge(), sseTrigger(), unreadLabel(), Props, Variant

### Community 114 - "Renderer"
Cohesion: 0.38
Nodes (4): NewHouseholdMember(), notFound(), New(), Renderer

### Community 115 - "GetAppBaseURL"
Cohesion: 0.15
Nodes (15): ensureDatabasePool(), main(), GetAppBaseURL(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), IsDevelopment(), MustGetEnv() (+7 more)

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 124 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 125 - ".UpdateNote"
Cohesion: 0.38
Nodes (4): noteFromForm(), noteLengthMessage(), notePath(), notesPath()

### Community 126 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 127 - "generate-locales-languages.go"
Cohesion: 0.24
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 128 - "Redirect"
Cohesion: 0.22
Nodes (4): peoplePath(), safeNextPath(), Redirect(), FromGinContext()

### Community 129 - ".loadShoppingItem"
Cohesion: 0.32
Nodes (3): shoppingItemPath(), shoppingListPath(), ShoppingItem

### Community 134 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 135 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 136 - "vendor-assets_test.go"
Cohesion: 0.36
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 139 - "LK"
Cohesion: 0.20
Nodes (14): LK, N(), unescapeJSON(), N(), HouseholdNav, householdMenu(), householdMenuItem(), householdNav() (+6 more)

### Community 140 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 141 - "avatar"
Cohesion: 0.25
Nodes (8): crop-hint, description, invalid, remove, title, upload, zoom, avatar

### Community 142 - "common"
Cohesion: 0.50
Nodes (4): delete, link-invalid, save, common

### Community 143 - "error-page"
Cohesion: 0.50
Nodes (4): error-page, back, description, title

### Community 145 - "offline"
Cohesion: 0.50
Nodes (4): offline, description, retry, title

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 749 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **32 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `login-page`, `avatar`, `offline`, `household`, `expenses`, `shopping`, `validation-errors`, `home-page`, `notes`, `verify-email`, `reset-password-page`, `files`, `appearance`, `nav`, `chat`, `push`, `invite-page`, `registration-page`, `mail`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13391959798994976 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `github.com/a-h/templ.Component`, `handlers.go`, `.loadShoppingItem`, `Redirect`, `.clean`, `Membership`, `newEmailToken`, `LK`, `context.Context`, `ParseUUID`, `ChatBadge`, `GetAppBaseURL`, `github.com/gin-gonic/gin.Context`, `RenderTemplStatus`, `InternalError`, `.notifyChat`, `.UpdateNote`, `CN`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.06896551724137931 - nodes in this community are weakly interconnected._
- **Why does `LKMessage` connect `LKMessage` to `handlers.go`, `T`, `LK`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Should `et` be split into smaller, more focused modules?**
  _Cohesion score 0.06060606060606061 - nodes in this community are weakly interconnected._