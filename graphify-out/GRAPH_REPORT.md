# Graph Report - roommates  (2026-10-31)

## Corpus Check
- 279 files · ~281,120 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 53 file(s) not represented in the graph (top: .templ 36, (none) 6, .css 5)

## Summary
- 2915 nodes · 9255 edges · 184 communities (139 shown, 45 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 680 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9073fc95`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- migrations.go
- en
- et
- environment
- generate-locales-languages.go
- vars
- T
- LKMessage
- RenderTemplStatus
- pageError
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
- Membership
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- Redirect
- select.js
- floating_ui_dom.js
- DefaultService
- files
- validation-errors
- Queries
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- vendor-assets.go
- reset-password-page
- use_hover.js
- login-page
- login-page
- S3Store
- .GetInvite
- files
- membershipFromContext
- appearance
- appearance
- calendar.js
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
- Props
- floating_ui_core.js
- ParseUUID
- Debt
- tabbable.js
- useDismiss
- avatar
- Session
- NewInternationalization
- Subscription
- mail
- generate-locales-keys.go
- chat
- chat
- use_transition_status.js
- icon.go
- forgot-password-page
- composite.js
- go_pkg_fmt
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- mailer.go
- tasks
- password
- Renderer
- NewPostgreSessionStoreHandler
- ListActiveMembersRow
- switch.js
- mounted
- scroll_lock.js
- .UpdateNote
- .clean
- calendar_templ.go
- Button
- InternalError
- forgot-password-page
- Hub
- popover.js
- TestPushNotifications
- Assets
- avatar.js
- reset-password-page
- avatar
- tailwind
- noRowsAs
- registration-page
- floating_focus_manager.js
- NormalizeInviteCode
- select_templ.go
- useListNavigation
- server
- assets-signal
- dev
- error-page
- build
- tailwind-version
- scripts-watch
- github.com/gin-gonic/gin.Context
- expense-split.js
- checkbox.js
- registration-page
- popover_templ.go
- newTestUser
- useClick
- mail
- pushPreview
- mark_others.js
- vendor-assets_test.go
- useTypeahead
- ListBuyAgainRow
- GetAppearance
- internal_backdrop.js
- .handleUpload
- invite-page
- common
- DeleteShoppingItemParams
- SetMembershipRoleParams

## God Nodes (most connected - your core abstractions)
1. `T()` - 92 edges
2. `CN()` - 86 edges
3. `T()` - 59 edges
4. `Membership` - 54 edges
5. `expenses` - 53 edges
6. `expenses` - 53 edges
7. `membershipFromContext()` - 50 edges
8. `RenderTemplStatus()` - 45 edges
9. `InternalError()` - 41 edges
10. `Button()` - 41 edges

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

## Communities (184 total, 45 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.14
Nodes (37): TestSetAppearance(), TestSafeNextPath(), TestAvatars(), TestChat(), TestChatFileRefs(), TestChatReply(), TestHouseholdEvents(), settlementHelpers() (+29 more)

### Community 1 - "handlers.go"
Cohesion: 0.13
Nodes (10): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ctxKey, ctxKey, Props, appearanceKey, directionKey (+2 more)

### Community 2 - "migrations.go"
Cohesion: 0.15
Nodes (14): ensureDatabasePool(), GetDatabaseURL(), TestGetDatabaseURLEscapes(), ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations() (+6 more)

### Community 3 - "en"
Cohesion: 0.09
Nodes (22): title, en, app, email, error-page, logout, not-found-page, offline (+14 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (33): title, delete, link-invalid, save, et, app, common, email (+25 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "generate-locales-languages.go"
Cohesion: 0.24
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 7 - "vars"
Cohesion: 0.17
Nodes (11): sh, use_gitignore, vars, FREE_PORT, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT (+3 more)

### Community 8 - "T"
Cohesion: 0.09
Nodes (66): GetUserRow, LK, AttachedFile, FileFilter, N(), T(), unescapeJSON(), N() (+58 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "RenderTemplStatus"
Cohesion: 0.22
Nodes (7): chatPath(), optionalUUID(), RenderTemplStatus(), swapMainOnError(), notFound(), NotFound(), New()

### Community 11 - "pageError"
Cohesion: 0.33
Nodes (4): shoppingItemPath(), pageError(), shoppingListPath(), ShoppingItem

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
Nodes (50): ListHouseholdsOfUserRow, Direction(), WithDirection(), GetUnreadChat(), GetHouseholds(), WithHouseholds(), state(), state() (+42 more)

### Community 16 - "expenses"
Cohesion: 0.04
Nodes (53): expenses, add, amount, back, balances, cancel, category, category-groceries (+45 more)

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
Cohesion: 0.04
Nodes (53): expenses, add, amount, back, balances, cancel, category, category-groceries (+45 more)

### Community 38 - "shopping"
Cohesion: 0.06
Nodes (33): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+25 more)

### Community 39 - "shopping"
Cohesion: 0.06
Nodes (33): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+25 more)

### Community 40 - "newEmailToken"
Cohesion: 0.11
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

### Community 44 - "Membership"
Cohesion: 0.16
Nodes (32): ListShoppingItemsRow, ExpenseFilter, Membership, ShoppingList, IfElse(), FormatMoney(), debtBreakdown(), expenseFilters() (+24 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.04
Nodes (56): ListExpenseSharesRow, ListActiveInvitesRow, ListMessagesRow, ListNotesRow, GetShoppingListRow, avatarURL(), ConfirmSettlementParams, ConfirmSettlementRow (+48 more)

### Community 46 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 47 - "Redirect"
Cohesion: 0.23
Nodes (5): peoplePath(), IsRequestHTMX(), safeNextPath(), Redirect(), FromGinContext()

### Community 48 - "select.js"
Cohesion: 0.10
Nodes (53): allContents(), armCancelOpen(), arrowScrollStep(), clamp(), close(), closeAll(), commitValue(), contentFor() (+45 more)

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "DefaultService"
Cohesion: 0.12
Nodes (25): main(), InitGinEngine(), initRoutes(), registerAssets(), TestAssetCaching(), GetAppBaseURL(), IsDevelopment(), DefaultHandler (+17 more)

### Community 51 - "files"
Cohesion: 0.06
Nodes (31): files, added-by, all, any-place, any-type, attachments, count, download (+23 more)

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
Cohesion: 0.03
Nodes (90): ActionProps, DescriptionProps, TitleProps, Chat, Message, MessageFile, MergeAttributes(), Size (+82 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "vendor-assets.go"
Cohesion: 0.17
Nodes (18): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+10 more)

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

### Community 63 - "S3Store"
Cohesion: 0.14
Nodes (4): Store, New(), Upload, S3Store

### Community 65 - "files"
Cohesion: 0.06
Nodes (31): files, added-by, all, any-place, any-type, attachments, count, download (+23 more)

### Community 66 - "membershipFromContext"
Cohesion: 0.19
Nodes (3): membershipFromContext(), RenderTempl(), Note

### Community 67 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 68 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 69 - "calendar.js"
Cohesion: 0.09
Nodes (44): addToRange(), build(), calculateFocusTarget(), CalendarDay, calendarFns(), addMonths(), createDateLib(), getISOWeek() (+36 more)

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
Cohesion: 0.16
Nodes (8): ListFilesRow, cleanFileName(), FileMeta, FileTarget, SplitFileName(), TestCleanFileName(), TestSplitFileName(), ListFilesParams

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
Cohesion: 0.08
Nodes (54): CN(), DirectionAttributes(), flattenClasses(), twMapTokens(), twReorder(), Action(), Description(), Title() (+46 more)

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (11): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+3 more)

### Community 85 - "Props"
Cohesion: 0.39
Nodes (6): Badge(), badgeClasses(), ChatBadge(), unreadLabel(), Props, Variant

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "ParseUUID"
Cohesion: 0.19
Nodes (4): expenseFromForm(), expensePath(), expensesPath(), ParseUUID()

### Community 88 - "Debt"
Cohesion: 0.11
Nodes (13): GetExpenseRow, InsertExpenseParams, ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts() (+5 more)

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

### Community 93 - "NewInternationalization"
Cohesion: 0.17
Nodes (9): DeleteCookie(), GetCookie(), GetHeader(), localeMissing(), NewInternationalization(), setLocaleContext(), CookieKey, HTTPHeaderKey (+1 more)

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (8): Message, Pusher, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

### Community 95 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 96 - "generate-locales-keys.go"
Cohesion: 0.24
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 97 - "chat"
Cohesion: 0.10
Nodes (21): attach, cancel-reply, copy, delete, empty, file-gone, files-too-large, length (+13 more)

### Community 98 - "chat"
Cohesion: 0.10
Nodes (21): attach, cancel-reply, copy, delete, empty, file-gone, files-too-large, length (+13 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 100 - "icon.go"
Cohesion: 0.57
Nodes (5): attrString(), generateSVG(), getIconContent(), Icon(), Props

### Community 101 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

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
Cohesion: 0.24
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 112 - "tasks"
Cohesion: 0.20
Nodes (10): cmds, desc, tasks, assets-watch, templ, vendor, cmds, desc (+2 more)

### Community 113 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 115 - "NewPostgreSessionStoreHandler"
Cohesion: 0.19
Nodes (8): GetEnv(), GetSessionKeyPairs(), MustGetEnv(), NewPostgreSessionStoreHandler(), startExpiredSessionCleanup(), ValidateCookieHashBlockPair(), EnvironmentVariableKey, sessionCleaner

### Community 116 - "ListActiveMembersRow"
Cohesion: 0.22
Nodes (4): GetHouseholdRow, ListActiveMembersRow, ListShoppingListsRow, HouseholdOverview

### Community 117 - "switch.js"
Cohesion: 0.31
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "scroll_lock.js"
Cohesion: 0.11
Nodes (15): anchoredPopupScrollLock(), AnimationFrame, getViewportScroller(), hasInsetScrollbars(), isOverflowElement(), isPageScrollLocked(), preventScrollInsetScrollbars(), cleanup() (+7 more)

### Community 124 - ".UpdateNote"
Cohesion: 0.24
Nodes (4): noteFromForm(), noteLengthMessage(), notePath(), notesPath()

### Community 125 - ".clean"
Cohesion: 0.12
Nodes (11): ReplyPreview, FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestNetDebts(), TestParseCents(), TestSplitEqual() (+3 more)

### Community 126 - "calendar_templ.go"
Cohesion: 0.12
Nodes (25): Type, Variant, Variants(), Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON() (+17 more)

### Community 127 - "Button"
Cohesion: 0.16
Nodes (26): Props, Variant, Ptr(), GetLocaleString(), Alert(), variantClasses(), Button(), ValidationErrors() (+18 more)

### Community 128 - "InternalError"
Cohesion: 0.20
Nodes (5): notice(), settingsPath(), householdPath(), sessionUserID(), InternalError()

### Community 129 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 131 - "Hub"
Cohesion: 0.16
Nodes (7): Hub, NewHub(), TestHub(), newFakeStore(), newFakePusher(), fakePusher, fakeStore

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "TestPushNotifications"
Cohesion: 0.25
Nodes (9): LoadLocales(), subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes(), TestTUnescapesTextButNotArgs(), TestExpenseEditorCarriesTheItemLink(), TestHouseholdExpensesSettledStamp() (+1 more)

### Community 139 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 141 - "avatar"
Cohesion: 0.22
Nodes (9): cancel, crop-hint, description, invalid, remove, title, upload, zoom (+1 more)

### Community 142 - "tailwind"
Cohesion: 0.33
Nodes (7): preconditions, msg, sh, cmds, desc, preconditions, tailwind

### Community 143 - "noRowsAs"
Cohesion: 0.09
Nodes (9): hashToken(), closePaidDebt(), Expense, payShare(), noRowsAs(), NoteInput, rowsAffected(), trimToLength() (+1 more)

### Community 146 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "NormalizeInviteCode"
Cohesion: 0.50
Nodes (4): newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

### Community 149 - "select_templ.go"
Cohesion: 0.11
Nodes (23): Group(), initialValue(), Item(), itemLabel(), Label(), Select(), Separator(), state() (+15 more)

### Community 150 - "useListNavigation"
Cohesion: 0.24
Nodes (12): doSwitch(), isCrossOrientationCloseKey(), isCrossOrientationOpenKey(), isMainOrientationKey(), isMainOrientationToEndKey(), isVirtualClick(), isVirtualPointerEvent(), useListNavigation() (+4 more)

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
Nodes (4): back, description, title, error-page

### Community 155 - "build"
Cohesion: 0.50
Nodes (4): cmds, desc, task, build

### Community 156 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 157 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 162 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 163 - "popover_templ.go"
Cohesion: 0.11
Nodes (23): RandomID(), Content(), Description(), Header(), Popover(), state(), Title(), Trigger() (+15 more)

### Community 164 - "newTestUser"
Cohesion: 0.50
Nodes (9): englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword(), TestVerifyEmail(), tokenOf(), newTestUser(), nextEmail() (+1 more)

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 168 - "mark_others.js"
Cohesion: 0.52
Nodes (6): applyAttributeToOthers(), buildKeepSet(), collectOutsideElements(), correctElements(), markOthers(), unwrapHost()

### Community 170 - "vendor-assets_test.go"
Cohesion: 0.42
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 171 - "useTypeahead"
Cohesion: 0.53
Nodes (4): useTypeahead(), getMatchingIndex(), isItemAvailable(), onKeyDown()

### Community 173 - "GetAppearance"
Cohesion: 0.43
Nodes (6): NewAppearance(), GetAppearance(), NewAppearance(), TestNewAppearance(), WithAppearance(), Appearance

### Community 180 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 181 - "common"
Cohesion: 0.50
Nodes (4): delete, link-invalid, save, common

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 927 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **45 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `avatar`, `household`, `registration-page`, `error-page`, `expenses`, `mail`, `shopping`, `validation-errors`, `home-page`, `notes`, `verify-email`, `reset-password-page`, `login-page`, `files`, `appearance`, `nav`, `chat`, `forgot-password-page`, `push`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.13658536585365855 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `InternalError`, `handlers.go`, `TestPushNotifications`, `newEmailToken`, `RenderTemplStatus`, `pageError`, `Membership`, `context.Context`, `Redirect`, `github.com/a-h/templ.Component`, `.handleUpload`, `Props`, `ParseUUID`, `github.com/a-h/templ.Attributes`, `.UpdateNote`, `.clean`, `github.com/gin-gonic/gin.Context`, `Button`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13348342010224662 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `forgot-password-page`, `reset-password-page`, `expenses`, `household`, `registration-page`, `shopping`, `home-page`, `notes`, `files`, `invite-page`, `common`, `validation-errors`, `verify-email`, `login-page`, `appearance`, `nav`, `avatar`, `mail`, `chat`, `push`, `password`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.08695652173913043 - nodes in this community are weakly interconnected._