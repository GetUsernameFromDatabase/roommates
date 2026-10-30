# Graph Report - roommates  (2026-10-30)

## Corpus Check
- 245 files · ~159,411 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 49 file(s) not represented in the graph (top: .templ 32, (none) 6, .css 5)

## Summary
- 2396 nodes · 7773 edges · 150 communities (116 shown, 34 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 559 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `30e8d1b0`
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
- T
- RenderTemplStatus
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
- io.ReadSeeker
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- Redirect
- householdPath
- floating_ui_dom.js
- models.go
- files
- validation-errors
- Queries
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- reset-password-page
- reset-password-page
- use_hover.js
- login-page
- login-page
- github.com/a-h/templ.Component
- forgot-password-page
- files
- password
- appearance
- appearance
- householdURL
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
- chatMessages
- images.go
- offline
- floating_ui_core.js
- generate-locales-keys.go
- Debt
- tabbable.js
- useDismiss
- avatar
- Session
- DefaultService
- Subscription
- mail
- NewInternationalization
- chat
- chat
- use_transition_status.js
- Props
- registration-page
- composite.js
- bubble_templ.go
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- mailer.go
- Redirect
- DefaultService
- Renderer
- GetDatabaseURL
- InternalError
- switch.js
- mounted
- ChatBadge
- registration-page
- .clean
- mail
- generate-locales-languages.go
- github.com/gin-gonic/gin.Context
- vendor-assets_test.go
- Hub
- .notifyChat
- go_pkg_fmt
- Assets
- avatar.js
- .renderShoppingItem
- avatar
- password
- ListNotesRow
- error-page
- invite-page
- forgot-password-page
- noRowsAs

## God Nodes (most connected - your core abstractions)
1. `T()` - 77 edges
2. `CN()` - 72 edges
3. `T()` - 53 edges
4. `membershipFromContext()` - 46 edges
5. `Membership` - 43 edges
6. `RenderTemplStatus()` - 42 edges
7. `expenses` - 42 edges
8. `expenses` - 42 edges
9. `InternalError()` - 39 edges
10. `notes` - 35 edges

## Surprising Connections (you probably didn't know these)
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go
- `expensesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-expenses.go → app/internal/handlers/handlers-households.go
- `settingsPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-home.go → app/internal/handlers/handlers-households.go
- `shoppingListPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-shopping.go → app/internal/handlers/handlers-households.go
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go

## Import Cycles
- None detected.

## Communities (150 total, 34 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.15
Nodes (35): TestSafeNextPath(), TestAvatars(), TestChat(), TestChatReply(), TestHouseholdEvents(), englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail() (+27 more)

### Community 1 - "handlers.go"
Cohesion: 0.14
Nodes (10): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ctxKey, ctxKey, Props, appearanceKey, directionKey (+2 more)

### Community 2 - "migrations.go"
Cohesion: 0.16
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.07
Nodes (26): title, delete, link-invalid, save, en, app, common, email (+18 more)

### Community 4 - "et"
Cohesion: 0.08
Nodes (24): title, delete, link-invalid, save, et, app, common, email (+16 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.18
Nodes (17): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+9 more)

### Community 7 - "vars"
Cohesion: 0.04
Nodes (47): cmds, desc, preconditions, task, cmds, desc, env, SERVER_PORT (+39 more)

### Community 8 - "Membership"
Cohesion: 0.14
Nodes (28): GetExpenseRow, Expense, AttachedFile, Membership, ShoppingList, FormatMoney(), localTime(), debtBreakdown() (+20 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "T"
Cohesion: 0.15
Nodes (44): GetUserRow, HouseholdOverview, T(), Button(), T(), ValidationErrors(), Field(), Group() (+36 more)

### Community 11 - "RenderTemplStatus"
Cohesion: 0.21
Nodes (10): optionalUUID(), expenseFromForm(), expensePath(), expensesPath(), membershipFromContext(), pageError(), RenderTemplStatus(), Note (+2 more)

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
Nodes (37): ListHouseholdsOfUserRow, Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetHouseholds(), WithHouseholds(), state() (+29 more)

### Community 16 - "expenses"
Cohesion: 0.05
Nodes (42): expenses, add, amount, back, balances, cancel, category, category-groceries (+34 more)

### Community 17 - "household"
Cohesion: 0.07
Nodes (30): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+22 more)

### Community 18 - "household"
Cohesion: 0.07
Nodes (30): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+22 more)

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
Cohesion: 0.10
Nodes (7): GetUserCredentialsRow, InsertUserParams, LoginRequest, NormalizeEmail(), newEmailToken(), timestamptz(), newInviteCode()

### Community 41 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 42 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 43 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 44 - "io.ReadSeeker"
Cohesion: 0.13
Nodes (4): Store, New(), Upload, S3Store

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.04
Nodes (30): ListExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, GetShoppingListRow, ListShoppingListsRow, avatarURL(), ConsumeEmailTokenParams, DeleteCheckedShoppingItemsParams (+22 more)

### Community 46 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 47 - "Redirect"
Cohesion: 0.21
Nodes (3): peoplePath(), IsRequestHTMX(), swapMainOnError()

### Community 48 - "householdPath"
Cohesion: 0.16
Nodes (7): chatPath(), filesPath(), householdPath(), noteFromForm(), noteLengthMessage(), notePath(), notesPath()

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "models.go"
Cohesion: 0.08
Nodes (26): ListMessagesRow, ListShoppingItemsRow, EmailToken, Expense, ExpenseShare, File, FilePurge, GetFileParams (+18 more)

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
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.05
Nodes (52): ActionProps, DescriptionProps, Props, TitleProps, Variant, variantClasses(), Badge(), Size (+44 more)

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

### Community 61 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 62 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 63 - "github.com/a-h/templ.Component"
Cohesion: 0.07
Nodes (59): CN(), DirectionAttributes(), flattenClasses(), RandomID(), twMapTokens(), twReorder(), GetAppearance(), Action() (+51 more)

### Community 64 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 65 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, download, file, invalid (+7 more)

### Community 66 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 67 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 68 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 69 - "householdURL"
Cohesion: 0.11
Nodes (30): LK, IfElse(), N(), unescapeJSON(), Reactions(), N(), Description(), chatControlLabel() (+22 more)

### Community 70 - "environment"
Cohesion: 0.18
Nodes (11): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment, ports (+3 more)

### Community 71 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 72 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 73 - "rowsAffected"
Cohesion: 0.09
Nodes (7): GetInviteByCodeRow, hashToken(), noRowsAs(), NoteInput, rowsAffected(), trimToLength(), NoteForm

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

### Community 83 - "chatMessages"
Cohesion: 0.12
Nodes (26): Chat, Message, MessageFile, MergeAttributes(), Avatar(), Content(), Footer(), Group() (+18 more)

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (11): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+3 more)

### Community 85 - "offline"
Cohesion: 0.50
Nodes (4): offline, description, retry, title

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "generate-locales-keys.go"
Cohesion: 0.24
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 88 - "Debt"
Cohesion: 0.13
Nodes (13): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), TestNetDebts(), testUUID() (+5 more)

### Community 89 - "tabbable.js"
Cohesion: 0.14
Nodes (28): activeElement(), appendCandidates(), appendMatchingElements(), contains(), disableFocusInside(), enableFocusInside(), focusable(), getComposedChildren() (+20 more)

### Community 90 - "useDismiss"
Cohesion: 0.19
Nodes (25): contains(), getTarget(), isLastTraversableNode(), normalizeBubbles(), useDismiss(), addTargetEventListenerOnce(), close(), closeOnEscapeKeyDown() (+17 more)

### Community 91 - "avatar"
Cohesion: 0.25
Nodes (8): crop-hint, description, invalid, remove, title, upload, zoom, avatar

### Community 92 - "Session"
Cohesion: 0.29
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 93 - "DefaultService"
Cohesion: 0.13
Nodes (20): main(), InitGinEngine(), initRoutes(), registerAssets(), GetAppBaseURL(), IsDevelopment(), InitLogger(), New() (+12 more)

### Community 94 - "Subscription"
Cohesion: 0.12
Nodes (9): Message, Subscription, newFakePusher(), TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher (+1 more)

### Community 95 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 96 - "NewInternationalization"
Cohesion: 0.12
Nodes (14): DeleteCookie(), GetCookie(), GetHeader(), NewAppearance(), localeMissing(), NewInternationalization(), setLocaleContext(), NewAppearance() (+6 more)

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
Cohesion: 0.36
Nodes (6): Type, Variant, Variants(), Props, Size, ButtonProps

### Community 101 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "bubble_templ.go"
Cohesion: 0.26
Nodes (12): Bubble(), Content(), ContentClass(), Group(), variantClasses(), Align, ContentProps, GroupProps (+4 more)

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

### Community 111 - "mailer.go"
Cohesion: 0.24
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 112 - "Redirect"
Cohesion: 0.29
Nodes (6): safeNextPath(), Redirect(), NewHouseholdMember(), notFound(), FromGinContext(), New()

### Community 113 - "DefaultService"
Cohesion: 0.17
Nodes (13): DefaultHandler, NewDefaultHandler(), LoadLocales(), Mailer, Pusher, DefaultService, NewDefaultService(), subscribeTestDevice() (+5 more)

### Community 115 - "GetDatabaseURL"
Cohesion: 0.24
Nodes (8): ensureDatabasePool(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), MustGetEnv(), TestGetDatabaseURLEscapes(), ValidateCookieHashBlockPair(), EnvironmentVariableKey

### Community 116 - "InternalError"
Cohesion: 0.16
Nodes (5): notice(), settingsPath(), sessionUserID(), InternalError(), GetLocaleString()

### Community 117 - "switch.js"
Cohesion: 0.24
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "ChatBadge"
Cohesion: 0.60
Nodes (4): Badge(), badgeClasses(), Props, Variant

### Community 124 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 125 - ".clean"
Cohesion: 0.15
Nodes (10): InsertExpenseParams, ListActiveMembersRow, FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestParseCents(), TestSplitEqual() (+2 more)

### Community 126 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 127 - "generate-locales-languages.go"
Cohesion: 0.27
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 129 - "vendor-assets_test.go"
Cohesion: 0.26
Nodes (7): fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest(), writeAtomic(), writeJSON()

### Community 131 - "Hub"
Cohesion: 0.21
Nodes (5): Hub, NewHub(), TestHub(), newFakeStore(), fakeStore

### Community 139 - ".renderShoppingItem"
Cohesion: 0.29
Nodes (3): shoppingItemPath(), shoppingListPath(), ShoppingItem

### Community 141 - "avatar"
Cohesion: 0.25
Nodes (8): crop-hint, description, invalid, remove, title, upload, zoom, avatar

### Community 142 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 146 - "error-page"
Cohesion: 0.50
Nodes (4): back, description, title, error-page

### Community 147 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 148 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 149 - "noRowsAs"
Cohesion: 0.22
Nodes (4): ListFilesRow, FileMeta, FileTarget, ListFilesParams

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 799 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **34 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `avatar`, `error-page`, `household`, `forgot-password-page`, `expenses`, `shopping`, `validation-errors`, `home-page`, `notes`, `verify-email`, `reset-password-page`, `login-page`, `files`, `password`, `appearance`, `nav`, `offline`, `chat`, `push`, `registration-page`, `mail`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13505614520842948 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `handlers.go`, `householdURL`, `.notifyChat`, `Membership`, `newEmailToken`, `RenderTemplStatus`, `.renderShoppingItem`, `context.Context`, `householdPath`, `DefaultService`, `Redirect`, `chatMessages`, `InternalError`, `Redirect`, `.clean`, `github.com/a-h/templ.Component`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.07407407407407407 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `password`, `expenses`, `household`, `invite-page`, `shopping`, `home-page`, `notes`, `files`, `validation-errors`, `verify-email`, `reset-password-page`, `login-page`, `forgot-password-page`, `appearance`, `nav`, `avatar`, `mail`, `chat`, `registration-page`, `push`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Should `et` be split into smaller, more focused modules?**
  _Cohesion score 0.08 - nodes in this community are weakly interconnected._