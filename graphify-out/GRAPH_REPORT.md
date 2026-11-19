# Graph Report - roommates  (2026-11-19)

## Corpus Check
- 267 files · ~310,129 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 59 file(s) not represented in the graph (top: .templ 39, (none) 9, .css 5)

## Summary
- 3245 nodes · 10462 edges · 190 communities (147 shown, 43 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 753 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9205bc84`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- migrations.go
- en
- et
- bubble_templ.go
- vars
- iPaidAmount
- LKMessage
- InitGinEngine
- models.go
- components.json
- block-secret-files.sh
- go
- context.Context
- expenses
- household
- household
- permissions
- linters
- Roommates
- dialog.js
- RenderTemplStatus
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
- GetAppBaseURL
- validation-errors
- home-page
- home-page
- Membership
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- vendor-assets.go
- select.js
- floating_ui_dom.js
- environment
- files
- validation-errors
- select_templ.go
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- generate-locales-keys.go
- reset-password-page
- use_hover.js
- login-page
- debtBreakdown
- T
- navIcon
- files
- io.ReadSeeker
- DefaultService
- initRoutes
- calendar.js
- environment
- nav
- nav
- .UploadFile
- services
- db
- queries/README.md
- login-page
- ErrorCode
- db
- context
- healthcheck
- newTestUser
- github.com/a-h/templ.Component
- images.go
- NormalizeInviteCode
- floating_ui_core.js
- chatMessages
- Debt
- tabbable.js
- useDismiss
- avatar
- Session
- Hub
- Subscription
- options
- invite-page
- chat
- chat
- use_transition_status.js
- .notifyChat
- cloudflared
- composite.js
- Queries
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- mailer.go
- tasks
- forgot-password-page
- Renderer
- password
- .loadShoppingItem
- switch.js
- mounted
- scroll_lock.js
- invite-page
- .handleUpload
- calendar_templ.go
- go_pkg_fmt
- github.com/gin-gonic/gin.Context
- app
- vendor-assets_test.go
- popover.js
- mail
- registerAssets
- avatar.js
- reset-password-page
- avatar
- tailwind
- noRowsAs
- registration-page
- floating_focus_manager.js
- Message
- common
- useListNavigation
- server
- restart-signal
- dev
- Browser check
- backup.sh
- common
- scripts-watch
- ListNotesRow
- expense-split.js
- checkbox.js
- HTMLTemplRenderer
- popover_templ.go
- appearance
- useClick
- error-page
- registration-page
- mark_others.js
- mail
- useTypeahead
- internal_backdrop.js
- .attachFilePicks
- icon.go
- tailwind-sha256-arm64
- appearance
- Redirect
- tailwind-sha256-x64
- vendor
- seed
- alertdialog_templ.go
- playwright

## God Nodes (most connected - your core abstractions)
1. `CN()` - 103 edges
2. `T()` - 93 edges
3. `T()` - 62 edges
4. `Membership` - 61 edges
5. `membershipFromContext()` - 57 edges
6. `expenses` - 53 edges
7. `expenses` - 53 edges
8. `RenderTemplStatus()` - 47 edges
9. `Button()` - 44 edges
10. `InternalError()` - 42 edges

## Surprising Connections (you probably didn't know these)
- `HTMX actions` --references--> `Redirect()`  [INFERRED]
  .claude/skills/browser-check/SKILL.md → app/internal/handlers/handlers.go
- `TestHouseholdRoutesHideFromStrangers()` --calls--> `initRoutes()`  [INFERRED]
  app/cmd/server/routes-household_test.go → app/cmd/server/routes.go
- `GetSessionKeyPairs()` --calls--> `ValidateCookieHashBlockPair()`  [EXTRACTED]
  app/internal/config/globals-env.go → app/internal/validators/validators.go
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go
- `expensesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-expenses.go → app/internal/handlers/handlers-households.go

## Import Cycles
- None detected.

## Communities (190 total, 43 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.12
Nodes (42): TestSetAppearance(), TestSetLanguage(), TestSafeNextPath(), TestAvatars(), TestChat(), TestChatFileRefs(), TestChatReply(), TestHouseholdEvents() (+34 more)

### Community 1 - "handlers.go"
Cohesion: 0.13
Nodes (10): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ctxKey, ctxKey, Props, appearanceKey, directionKey (+2 more)

### Community 2 - "migrations.go"
Cohesion: 0.18
Nodes (12): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateToLatestConcurrently(), TestMigrateUpAndDown() (+4 more)

### Community 3 - "en"
Cohesion: 0.06
Nodes (31): title, en, app, email, forgot-password-page, language, logout, not-found-page (+23 more)

### Community 4 - "et"
Cohesion: 0.08
Nodes (24): title, back, description, title, et, app, email, error-page (+16 more)

### Community 6 - "bubble_templ.go"
Cohesion: 0.23
Nodes (13): Bubble(), Content(), ContentClass(), Group(), Reactions(), variantClasses(), Align, ContentProps (+5 more)

### Community 7 - "vars"
Cohesion: 0.17
Nodes (11): use_gitignore, vars, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_SHA256_ARM64 (+3 more)

### Community 8 - "iPaidAmount"
Cohesion: 0.20
Nodes (9): FormatCentsInput(), splitEqual(), TestCopyInputFormOrder(), TestNetDebts(), TestParseCents(), TestSplitEqual(), testUUID(), iPaidAmount() (+1 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "InitGinEngine"
Cohesion: 0.08
Nodes (33): ensureDatabasePool(), main(), randomUUID(), send(), TestHouseholdRoutesHideFromStrangers(), InitGinEngine(), GetAppBaseURL(), GetDatabaseURL() (+25 more)

### Community 11 - "models.go"
Cohesion: 0.08
Nodes (30): ListMessagesRow, EmailToken, Expense, ExpenseCategory, ExpenseShare, File, FilePurge, FileRef (+22 more)

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
Cohesion: 0.02
Nodes (49): ListHouseholdsOfUserRow, ListBuyAgainRow, Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetHouseholds(), WithHouseholds() (+41 more)

### Community 16 - "expenses"
Cohesion: 0.04
Nodes (53): expenses, add, amount, back, balances, cancel, category, category-groceries (+45 more)

### Community 17 - "household"
Cohesion: 0.06
Nodes (34): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+26 more)

### Community 18 - "household"
Cohesion: 0.06
Nodes (34): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+26 more)

### Community 19 - "permissions"
Cohesion: 0.25
Nodes (7): hooks, PreToolUse, permissions, allow, ask, deny, $schema

### Community 20 - "linters"
Cohesion: 0.29
Nodes (6): allow-empty, linters, enable, settings, exhaustruct_v5, version

### Community 21 - "Roommates"
Cohesion: 0.07
Nodes (26): Certificates, Generate the certificate, Install mkcert, Assets, Third-party files, Alert, An older dump, Backup (+18 more)

### Community 22 - "dialog.js"
Cohesion: 0.20
Nodes (23): closeDialog(), createInternalBackdrop(), destroyDialog(), dialogFor(), ensureDialog(), getDialog(), init(), isDialogOpen() (+15 more)

### Community 23 - "RenderTemplStatus"
Cohesion: 0.21
Nodes (9): expensePath(), expensesPath(), membershipFromContext(), pageError(), RenderTemplStatus(), stay(), Note, ParseUUID() (+1 more)

### Community 27 - "dev-run.sh"
Cohesion: 0.67
Nodes (5): dlv_alive(), dlv_port_busy(), dev-run.sh script, start(), stop()

### Community 31 - "tooltip.js"
Cohesion: 0.26
Nodes (22): arrowOf(), close(), contentFor(), destroy(), emitOpenChange(), init(), joinGroup(), leaveGroup() (+14 more)

### Community 34 - "Roommates"
Cohesion: 0.33
Nodes (5): Commands, Conventions, graphify, Roommates, UI changes

### Community 37 - "expenses"
Cohesion: 0.04
Nodes (53): expenses, add, amount, back, balances, cancel, category, category-groceries (+45 more)

### Community 38 - "shopping"
Cohesion: 0.06
Nodes (33): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+25 more)

### Community 39 - "shopping"
Cohesion: 0.06
Nodes (33): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+25 more)

### Community 40 - "GetAppBaseURL"
Cohesion: 0.12
Nodes (7): GetUserCredentialsRow, InsertUserParams, LoginRequest, NormalizeEmail(), TestUserLimit(), newEmailToken(), userLimit

### Community 41 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 42 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 43 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 44 - "Membership"
Cohesion: 0.20
Nodes (23): ListShoppingItemsRow, ExpenseFilter, Membership, ShoppingList, IfElse(), ChatFilePicks(), expenseFilters(), AttachFilePicks() (+15 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.03
Nodes (45): InsertExpenseParams, ListExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, ListActiveMembersRow, GetShoppingListRow, ListShoppingListsRow, ParseCents() (+37 more)

### Community 46 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 47 - "vendor-assets.go"
Cohesion: 0.18
Nodes (19): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+11 more)

### Community 48 - "select.js"
Cohesion: 0.10
Nodes (53): allContents(), armCancelOpen(), arrowScrollStep(), clamp(), close(), closeAll(), commitValue(), contentFor() (+45 more)

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "environment"
Cohesion: 0.08
Nodes (28): environment, environment, APP_BASE_URL, BACKUP_ALERT_EMAIL, BACKUP_KEEP_DAYS, DATABASE_HOST, DATABASE_PORT, GIN_MODE (+20 more)

### Community 51 - "files"
Cohesion: 0.05
Nodes (42): one, other, one, other, files, added-by, all, any-place (+34 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "select_templ.go"
Cohesion: 0.09
Nodes (28): DirectionAttributes(), initialValue(), RadioGroup(), Content(), Group(), initialValue(), Item(), itemLabel() (+20 more)

### Community 54 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.04
Nodes (84): ActionProps, DescriptionProps, Props, TitleProps, Variant, ActionProps, CancelProps, variantClasses() (+76 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "generate-locales-keys.go"
Cohesion: 0.13
Nodes (14): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+6 more)

### Community 59 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 60 - "use_hover.js"
Cohesion: 0.09
Nodes (41): applySafePolygonPointerEventsMutation(), clearSafePolygonPointerEventsMutation(), createHoverInteraction(), emitClosed(), getDelay(), getRestMs(), hasIntersectingEdge(), isClickLikeOpenEvent() (+33 more)

### Community 61 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 62 - "debtBreakdown"
Cohesion: 0.25
Nodes (8): GetExpenseRow, householdPath(), FormatMoney(), debtBreakdown(), formatDay(), formatSignedMoney(), itemClass(), paidByLine()

### Community 63 - "T"
Cohesion: 0.08
Nodes (86): GetUserRow, LK, Expense, AttachedFile, FileFilter, HouseholdOverview, N(), T() (+78 more)

### Community 64 - "navIcon"
Cohesion: 0.14
Nodes (20): Badge(), badgeClasses(), ChatBadge(), HouseholdNav, householdMenu(), householdMenuItem(), householdNav(), householdNavItem() (+12 more)

### Community 65 - "files"
Cohesion: 0.05
Nodes (42): one, other, one, other, files, added-by, all, any-place (+34 more)

### Community 66 - "io.ReadSeeker"
Cohesion: 0.12
Nodes (5): Store, New(), nullStore, Upload, S3Store

### Community 67 - "DefaultService"
Cohesion: 0.30
Nodes (11): DefaultService, subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes(), check(), must(), seed() (+3 more)

### Community 68 - "initRoutes"
Cohesion: 0.08
Nodes (30): TestLimitBody(), initRoutes(), limitBody(), DeleteCookie(), GetCookie(), GetHeader(), setLanguageCookie(), safeNextPath() (+22 more)

### Community 69 - "calendar.js"
Cohesion: 0.09
Nodes (44): addToRange(), build(), calculateFocusTarget(), CalendarDay, calendarFns(), addMonths(), createDateLib(), getISOWeek() (+36 more)

### Community 70 - "environment"
Cohesion: 0.29
Nodes (7): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment

### Community 71 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 72 - "nav"
Cohesion: 0.17
Nodes (12): nav, chat, expenses, files, households-add, label, more, notes (+4 more)

### Community 73 - ".UploadFile"
Cohesion: 0.17
Nodes (8): ListFilesRow, cleanFileName(), FileMeta, FileTarget, SplitFileName(), TestCleanFileName(), TestSplitFileName(), ListFilesParams

### Community 74 - "services"
Cohesion: 0.29
Nodes (6): services, volumes, db_data, rustfs_data, x-logging, driver

### Community 75 - "db"
Cohesion: 0.25
Nodes (9): depends_on, depends_on, condition, required, restart, db, rustfs, condition (+1 more)

### Community 77 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 79 - "db"
Cohesion: 0.25
Nodes (8): command, environment, ports, volumes, POSTGRES_DB, POSTGRES_PASSWORD, POSTGRES_USER, db

### Community 80 - "context"
Cohesion: 0.17
Nodes (16): build, build, restart, volumes, context, dockerfile, build, build (+8 more)

### Community 81 - "healthcheck"
Cohesion: 0.57
Nodes (7): healthcheck, healthcheck, interval, retries, test, timeout, healthcheck

### Community 82 - "newTestUser"
Cohesion: 0.25
Nodes (15): englishContext(), sessionsOf(), TestDeleteExpiredEmailTokens(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword(), TestVerifyEmail(), tokenOf() (+7 more)

### Community 83 - "github.com/a-h/templ.Component"
Cohesion: 0.07
Nodes (57): CN(), flattenClasses(), twMapTokens(), twReorder(), Action(), Description(), Title(), AspectRatio() (+49 more)

### Community 84 - "images.go"
Cohesion: 0.12
Nodes (12): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+4 more)

### Community 85 - "NormalizeInviteCode"
Cohesion: 0.18
Nodes (6): GetInviteByCodeRow, timestamptz(), newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "chatMessages"
Cohesion: 0.13
Nodes (21): Chat, MergeAttributes(), Avatar(), Content(), Footer(), Group(), Header(), Message() (+13 more)

### Community 88 - "Debt"
Cohesion: 0.12
Nodes (10): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), ListDebtItemsParams, ListExpensesParams (+2 more)

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
Cohesion: 0.31
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 93 - "Hub"
Cohesion: 0.14
Nodes (7): Hub, DefaultHandler, NewDefaultHandler(), newFakeStore(), newFakePusher(), fakePusher, fakeStore

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (7): Message, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

### Community 95 - "options"
Cohesion: 0.20
Nodes (14): logging, logging, logging, logging, driver, options, max-file, max-size (+6 more)

### Community 96 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 97 - "chat"
Cohesion: 0.10
Nodes (21): attach, cancel-reply, copy, delete, empty, file-gone, files-too-large, length (+13 more)

### Community 98 - "chat"
Cohesion: 0.10
Nodes (21): attach, cancel-reply, copy, delete, empty, file-gone, files-too-large, length (+13 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 101 - "cloudflared"
Cohesion: 0.22
Nodes (9): condition, command, depends_on, environment, profiles, restart, app, TUNNEL_TOKEN (+1 more)

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "Queries"
Cohesion: 0.15
Nodes (4): Queries, checkSettled(), ownerActsOnMember(), DBTX

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
Cohesion: 0.20
Nodes (8): buildMessage(), Message, isLocalHost(), TestBuildMessage(), TestIsLocalHost(), LogMailer, SMTPMailer, fakeMailer

### Community 112 - "tasks"
Cohesion: 0.18
Nodes (11): cmds, desc, cmds, desc, silent, tasks, restart-watch, tailwind-version (+3 more)

### Community 113 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 114 - "Renderer"
Cohesion: 0.47
Nodes (3): notFound(), New(), Renderer

### Community 115 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 116 - ".loadShoppingItem"
Cohesion: 0.27
Nodes (3): shoppingItemPath(), shoppingListPath(), ShoppingItem

### Community 117 - "switch.js"
Cohesion: 0.31
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "scroll_lock.js"
Cohesion: 0.11
Nodes (15): anchoredPopupScrollLock(), AnimationFrame, getViewportScroller(), hasInsetScrollbars(), isOverflowElement(), isPageScrollLocked(), preventScrollInsetScrollbars(), cleanup() (+7 more)

### Community 124 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 125 - ".handleUpload"
Cohesion: 0.15
Nodes (5): filesPath(), noteFromForm(), noteLengthMessage(), notePath(), notesPath()

### Community 126 - "calendar_templ.go"
Cohesion: 0.21
Nodes (17): Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON(), toJSON(), CaptionLayout (+9 more)

### Community 128 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.08
Nodes (8): notice(), expenseFromForm(), settingsPath(), sessionUserID(), InternalError(), RenderTempl(), stayInPlace(), ExpenseInput

### Community 129 - "app"
Cohesion: 0.33
Nodes (6): env_file, ports, restart, stop_grace_period, volumes, app

### Community 131 - "vendor-assets_test.go"
Cohesion: 0.46
Nodes (5): fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 137 - "registerAssets"
Cohesion: 0.50
Nodes (5): cacheHeaders(), registerAssets(), TestAssetCaching(), Hash(), URL()

### Community 139 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 141 - "avatar"
Cohesion: 0.22
Nodes (9): cancel, crop-hint, description, invalid, remove, title, upload, zoom (+1 more)

### Community 142 - "tailwind"
Cohesion: 0.20
Nodes (11): cmds, desc, preconditions, task, msg, sh, cmds, desc (+3 more)

### Community 143 - "noRowsAs"
Cohesion: 0.10
Nodes (8): hashToken(), closePaidDebt(), payShare(), noRowsAs(), NoteInput, rowsAffected(), trimToLength(), NoteForm

### Community 146 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "Message"
Cohesion: 0.48
Nodes (7): Message, MessageFile, endsRun(), sameRun(), splitFiles(), startsRun(), TestChatRuns()

### Community 149 - "common"
Cohesion: 0.40
Nodes (5): cancel, delete, link-invalid, save, common

### Community 150 - "useListNavigation"
Cohesion: 0.24
Nodes (12): doSwitch(), isCrossOrientationCloseKey(), isCrossOrientationOpenKey(), isMainOrientationKey(), isMainOrientationToEndKey(), isVirtualClick(), isVirtualPointerEvent(), useListNavigation() (+4 more)

### Community 151 - "server"
Cohesion: 0.33
Nodes (6): SERVER_PORT, cmds, desc, dotenv, env, server

### Community 152 - "restart-signal"
Cohesion: 0.33
Nodes (6): cmds, desc, sources, use_gitignore, exclude, restart-signal

### Community 153 - "dev"
Cohesion: 0.33
Nodes (6): cmds, desc, env, RUSTFS_CONSOLE_ENABLE, SHADCN_TEMPL_DEV, dev

### Community 154 - "Browser check"
Cohesion: 0.18
Nodes (10): Browser check, Check each page, Clean up, Email, HTMX actions, Limits, Report, Selectors (+2 more)

### Community 156 - "common"
Cohesion: 0.40
Nodes (5): cancel, delete, link-invalid, save, common

### Community 157 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 163 - "popover_templ.go"
Cohesion: 0.09
Nodes (26): Ptr(), RandomID(), GetLocaleString(), Content(), Description(), Header(), Popover(), state() (+18 more)

### Community 164 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "error-page"
Cohesion: 0.50
Nodes (4): error-page, back, description, title

### Community 167 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 168 - "mark_others.js"
Cohesion: 0.52
Nodes (6): applyAttributeToOthers(), buildKeepSet(), collectOutsideElements(), correctElements(), markOthers(), unwrapHost()

### Community 170 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 171 - "useTypeahead"
Cohesion: 0.53
Nodes (4): useTypeahead(), getMatchingIndex(), isItemAvailable(), onKeyDown()

### Community 181 - "icon.go"
Cohesion: 0.73
Nodes (5): attrString(), generateSVG(), getIconContent(), Icon(), Props

### Community 183 - "tailwind-sha256-arm64"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-sha256-arm64

### Community 184 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 185 - "Redirect"
Cohesion: 0.19
Nodes (7): chatPath(), optionalUUID(), householdPath(), peoplePath(), IsRequestHTMX(), Redirect(), swapMainOnError()

### Community 186 - "tailwind-sha256-x64"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-sha256-x64

### Community 187 - "vendor"
Cohesion: 0.67
Nodes (3): vendor, cmds, desc

### Community 188 - "seed"
Cohesion: 0.50
Nodes (4): cmds, desc, dotenv, seed

### Community 189 - "alertdialog_templ.go"
Cohesion: 0.10
Nodes (30): ContentProps, ctxKey, ctxState, DescriptionProps, FooterProps, HeaderProps, MediaProps, OverlayProps (+22 more)

### Community 190 - "playwright"
Cohesion: 0.50
Nodes (3): npx, playwright, @playwright/mcp

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1025 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **43 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `mail`, `reset-password-page`, `expenses`, `household`, `common`, `appearance`, `error-page`, `registration-page`, `shopping`, `home-page`, `notes`, `files`, `validation-errors`, `verify-email`, `login-page`, `nav`, `avatar`, `chat`, `push`, `invite-page`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.1178743961352657 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `avatar`, `household`, `registration-page`, `common`, `expenses`, `shopping`, `validation-errors`, `mail`, `home-page`, `notes`, `appearance`, `verify-email`, `reset-password-page`, `files`, `nav`, `login-page`, `invite-page`, `chat`, `push`, `forgot-password-page`, `password`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13122171945701358 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `github.com/gin-gonic/gin.Context`, `handlers.go`, `navIcon`, `popover_templ.go`, `initRoutes`, `.notifyChat`, `GetAppBaseURL`, `Membership`, `context.Context`, `github.com/a-h/templ.Component`, `RenderTemplStatus`, `alertdialog_templ.go`, `github.com/a-h/templ.Attributes`, `chatMessages`, `.loadShoppingItem`, `Redirect`, `.handleUpload`, `debtBreakdown`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.0625 - nodes in this community are weakly interconnected._