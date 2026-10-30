# Graph Report - roommates  (2026-10-30)

## Corpus Check
- 253 files · ~168,111 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 49 file(s) not represented in the graph (top: .templ 32, (none) 6, .css 5)

## Summary
- 2479 nodes · 8092 edges · 161 communities (128 shown, 33 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 604 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8d887244`
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
- Expense
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- Redirect
- membershipFromContext
- floating_ui_dom.js
- chatMessages
- files
- validation-errors
- Queries
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- ListNotesRow
- reset-password-page
- use_hover.js
- login-page
- login-page
- io.ReadSeeker
- forgot-password-page
- files
- password
- appearance
- appearance
- models.go
- environment
- nav
- nav
- .UploadFile
- services
- db
- queries/README.md
- dev-signal.sh
- ErrorCode
- db
- context
- healthcheck
- app
- github.com/a-h/templ.Component
- images.go
- Badge
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
- ParseUUID
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
- tasks
- DefaultService
- Renderer
- GetDatabaseURL
- InternalError
- switch.js
- mounted
- T
- .UpdateNote
- .clean
- mail
- LK
- github.com/gin-gonic/gin.Context
- vendor-assets_test.go
- Hub
- .notifyChat
- go_pkg_fmt
- Assets
- avatar.js
- GetAppearance
- avatar
- tailwind
- noRowsAs
- registration-page
- generate-locales-languages.go
- .PostMessage
- password
- invite-page
- server
- assets-signal
- dev
- error-page
- build
- tailwind-version
- scripts-watch
- invite-page
- expense-split.js

## God Nodes (most connected - your core abstractions)
1. `T()` - 84 edges
2. `CN()` - 76 edges
3. `T()` - 56 edges
4. `membershipFromContext()` - 48 edges
5. `Membership` - 48 edges
6. `RenderTemplStatus()` - 43 edges
7. `expenses` - 43 edges
8. `expenses` - 43 edges
9. `InternalError()` - 40 edges
10. `Button()` - 36 edges

## Surprising Connections (you probably didn't know these)
- `chatPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-chat.go → app/internal/handlers/handlers-households.go
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go
- `expensesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-expenses.go → app/internal/handlers/handlers-households.go
- `filesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-files.go → app/internal/handlers/handlers-households.go
- `notesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-notes.go → app/internal/handlers/handlers-households.go

## Import Cycles
- None detected.

## Communities (161 total, 33 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.11
Nodes (46): TestSetAppearance(), TestSafeNextPath(), LoadLocales(), TestAvatars(), TestChat(), TestChatReply(), TestHouseholdEvents(), englishContext() (+38 more)

### Community 1 - "handlers.go"
Cohesion: 0.14
Nodes (10): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ctxKey, ctxKey, Props, appearanceKey, directionKey (+2 more)

### Community 2 - "migrations.go"
Cohesion: 0.16
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.07
Nodes (29): title, delete, link-invalid, save, en, app, common, email (+21 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (32): title, delete, link-invalid, save, back, description, title, et (+24 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "vendor-assets.go"
Cohesion: 0.18
Nodes (17): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+9 more)

### Community 7 - "vars"
Cohesion: 0.17
Nodes (11): sh, use_gitignore, vars, FREE_PORT, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT (+3 more)

### Community 8 - "Membership"
Cohesion: 0.18
Nodes (27): ListShoppingItemsRow, Membership, ShoppingList, IfElse(), FormatMoney(), categoryOptions(), debtBreakdown(), HouseholdExpenses() (+19 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "T"
Cohesion: 0.30
Nodes (19): Alert(), T(), ValidationErrors(), Field(), Group(), Label(), Input(), ExpenseEditor() (+11 more)

### Community 11 - "RenderTemplStatus"
Cohesion: 0.18
Nodes (6): shoppingItemPath(), pageError(), RenderTemplStatus(), shoppingListPath(), ShoppingItem, NotFound()

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
Nodes (38): ListHouseholdsOfUserRow, GetInviteByCodeRow, ListBuyAgainRow, Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetHouseholds() (+30 more)

### Community 16 - "expenses"
Cohesion: 0.05
Nodes (43): expenses, add, amount, back, balances, cancel, category, category-groceries (+35 more)

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
Nodes (43): expenses, add, amount, back, balances, cancel, category, category-groceries (+35 more)

### Community 38 - "shopping"
Cohesion: 0.06
Nodes (33): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+25 more)

### Community 39 - "shopping"
Cohesion: 0.06
Nodes (33): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+25 more)

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

### Community 44 - "Expense"
Cohesion: 0.29
Nodes (4): GetExpenseRow, ListExpenseSharesRow, Expense, GetExpenseParams

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.04
Nodes (33): GetHouseholdRow, ListActiveInvitesRow, GetShoppingListRow, ListShoppingListsRow, avatarURL(), DeleteExpenseParams, DeleteFileParams, DeleteNoteParams (+25 more)

### Community 46 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 47 - "Redirect"
Cohesion: 0.24
Nodes (5): peoplePath(), IsRequestHTMX(), safeNextPath(), Redirect(), FromGinContext()

### Community 48 - "membershipFromContext"
Cohesion: 0.19
Nodes (3): membershipFromContext(), RenderTempl(), Note

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "chatMessages"
Cohesion: 0.08
Nodes (33): Chat, Message, MessageFile, RandomID(), Avatar(), Content(), Footer(), Group() (+25 more)

### Community 51 - "files"
Cohesion: 0.13
Nodes (15): files, added-by, all, attachments, count, download, file, invalid (+7 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "Queries"
Cohesion: 0.20
Nodes (3): Queries, New(), DBTX

### Community 54 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.04
Nodes (66): ActionProps, DescriptionProps, Props, TitleProps, Variant, DirectionAttributes(), MergeAttributes(), variantClasses() (+58 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "ListNotesRow"
Cohesion: 0.19
Nodes (4): ListNotesRow, NoteInput, ListNotesParams, NoteForm

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

### Community 63 - "io.ReadSeeker"
Cohesion: 0.14
Nodes (4): Store, New(), Upload, S3Store

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

### Community 69 - "models.go"
Cohesion: 0.09
Nodes (25): ListMessagesRow, EmailToken, Expense, ExpenseCategory, ExpenseShare, File, FilePurge, GetFileParams (+17 more)

### Community 70 - "environment"
Cohesion: 0.18
Nodes (11): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment, ports (+3 more)

### Community 71 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 72 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 73 - ".UploadFile"
Cohesion: 0.22
Nodes (4): ListFilesRow, FileMeta, FileTarget, ListFilesParams

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

### Community 83 - "github.com/a-h/templ.Component"
Cohesion: 0.09
Nodes (43): CN(), flattenClasses(), twMapTokens(), twReorder(), Action(), Description(), Title(), AspectRatio() (+35 more)

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (11): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+3 more)

### Community 85 - "Badge"
Cohesion: 0.29
Nodes (8): Badge(), badgeClasses(), sharePaidBadge(), ChatBadge(), sseTrigger(), unreadLabel(), Props, Variant

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "generate-locales-keys.go"
Cohesion: 0.24
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 88 - "Debt"
Cohesion: 0.12
Nodes (11): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, checkSettled(), Debt, ExpensesOverview, netDebts(), ownerActsOnMember() (+3 more)

### Community 89 - "tabbable.js"
Cohesion: 0.14
Nodes (28): activeElement(), appendCandidates(), appendMatchingElements(), contains(), disableFocusInside(), enableFocusInside(), focusable(), getComposedChildren() (+20 more)

### Community 90 - "useDismiss"
Cohesion: 0.19
Nodes (25): contains(), getTarget(), isLastTraversableNode(), normalizeBubbles(), useDismiss(), addTargetEventListenerOnce(), close(), closeOnEscapeKeyDown() (+17 more)

### Community 91 - "avatar"
Cohesion: 0.22
Nodes (9): cancel, crop-hint, description, invalid, remove, title, upload, zoom (+1 more)

### Community 92 - "Session"
Cohesion: 0.29
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 93 - "DefaultService"
Cohesion: 0.13
Nodes (20): main(), InitGinEngine(), initRoutes(), registerAssets(), GetAppBaseURL(), IsDevelopment(), InitLogger(), New() (+12 more)

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (7): Message, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

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

### Community 100 - "ParseUUID"
Cohesion: 0.21
Nodes (4): expenseFromForm(), expensePath(), expensesPath(), ParseUUID()

### Community 101 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "bubble_templ.go"
Cohesion: 0.23
Nodes (13): Bubble(), Content(), ContentClass(), Group(), Reactions(), variantClasses(), Align, ContentProps (+5 more)

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
Cohesion: 0.20
Nodes (10): push, chat-file, chat-title, denied, expense-body, expense-title, install, label (+2 more)

### Community 109 - "push"
Cohesion: 0.20
Nodes (10): push, chat-file, chat-title, denied, expense-body, expense-title, install, label (+2 more)

### Community 111 - "mailer.go"
Cohesion: 0.22
Nodes (7): buildMessage(), Message, TestBuildMessage(), nextEmail(), LogMailer, SMTPMailer, fakeMailer

### Community 112 - "tasks"
Cohesion: 0.20
Nodes (10): cmds, desc, tasks, assets-watch, templ, vendor, cmds, desc (+2 more)

### Community 113 - "DefaultService"
Cohesion: 0.20
Nodes (11): DefaultHandler, NewDefaultHandler(), Mailer, Pusher, DefaultService, NewDefaultService(), subscribeTestDevice(), takePushes() (+3 more)

### Community 114 - "Renderer"
Cohesion: 0.38
Nodes (4): NewHouseholdMember(), notFound(), New(), Renderer

### Community 115 - "GetDatabaseURL"
Cohesion: 0.24
Nodes (8): ensureDatabasePool(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), MustGetEnv(), TestGetDatabaseURLEscapes(), ValidateCookieHashBlockPair(), EnvironmentVariableKey

### Community 116 - "InternalError"
Cohesion: 0.17
Nodes (6): notice(), settingsPath(), householdPath(), sessionUserID(), InternalError(), GetLocaleString()

### Community 117 - "switch.js"
Cohesion: 0.24
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "T"
Cohesion: 0.13
Nodes (45): GetUserRow, AttachedFile, HouseholdOverview, T(), Button(), Textarea(), AppearanceForm(), themeSwatch() (+37 more)

### Community 124 - ".UpdateNote"
Cohesion: 0.38
Nodes (4): noteFromForm(), noteLengthMessage(), notePath(), notesPath()

### Community 125 - ".clean"
Cohesion: 0.13
Nodes (12): InsertExpenseParams, ListActiveMembersRow, FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestNetDebts(), TestParseCents() (+4 more)

### Community 126 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 127 - "LK"
Cohesion: 0.16
Nodes (18): LK, N(), unescapeJSON(), N(), chatControlLabel(), HouseholdNav, householdMenu(), householdMenuItem() (+10 more)

### Community 129 - "vendor-assets_test.go"
Cohesion: 0.26
Nodes (7): fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest(), writeAtomic(), writeJSON()

### Community 131 - "Hub"
Cohesion: 0.15
Nodes (7): Hub, NewHub(), TestHub(), newFakeStore(), newFakePusher(), fakePusher, fakeStore

### Community 139 - "GetAppearance"
Cohesion: 0.43
Nodes (6): NewAppearance(), GetAppearance(), NewAppearance(), TestNewAppearance(), WithAppearance(), Appearance

### Community 141 - "avatar"
Cohesion: 0.22
Nodes (9): cancel, crop-hint, description, invalid, remove, title, upload, zoom (+1 more)

### Community 142 - "tailwind"
Cohesion: 0.33
Nodes (7): preconditions, msg, sh, cmds, desc, preconditions, tailwind

### Community 143 - "noRowsAs"
Cohesion: 0.12
Nodes (6): hashToken(), cleanFileName(), TestCleanFileName(), noRowsAs(), rowsAffected(), trimToLength()

### Community 146 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 147 - "generate-locales-languages.go"
Cohesion: 0.27
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 148 - ".PostMessage"
Cohesion: 0.36
Nodes (3): chatPath(), optionalUUID(), swapMainOnError()

### Community 149 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 150 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 151 - "server"
Cohesion: 0.33
Nodes (6): SERVER_PORT, cmds, desc, dotenv, env, server

### Community 152 - "assets-signal"
Cohesion: 0.33
Nodes (6): cmds, desc, sources, use_gitignore, exclude, assets-signal

### Community 153 - "dev"
Cohesion: 0.40
Nodes (5): cmds, desc, env, SHADCN_TEMPL_DEV, dev

### Community 154 - "error-page"
Cohesion: 0.50
Nodes (4): error-page, back, description, title

### Community 155 - "build"
Cohesion: 0.50
Nodes (4): cmds, desc, task, build

### Community 156 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 157 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 158 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 832 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `avatar`, `household`, `registration-page`, `invite-page`, `expenses`, `shopping`, `validation-errors`, `home-page`, `notes`, `verify-email`, `reset-password-page`, `login-page`, `files`, `password`, `appearance`, `nav`, `chat`, `push`, `mail`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.11394557823129252 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `testing.T`, `handlers.go`, `github.com/gin-gonic/gin.Context`, `.notifyChat`, `Membership`, `T`, `RenderTemplStatus`, `context.Context`, `.PostMessage`, `newEmailToken`, `Redirect`, `chatMessages`, `github.com/a-h/templ.Attributes`, `github.com/a-h/templ.Component`, `Badge`, `ParseUUID`, `InternalError`, `.UpdateNote`, `.clean`, `LK`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13623115577889447 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `expenses`, `household`, `password`, `invite-page`, `error-page`, `shopping`, `home-page`, `notes`, `files`, `validation-errors`, `verify-email`, `login-page`, `forgot-password-page`, `appearance`, `nav`, `avatar`, `mail`, `chat`, `registration-page`, `push`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.06666666666666667 - nodes in this community are weakly interconnected._