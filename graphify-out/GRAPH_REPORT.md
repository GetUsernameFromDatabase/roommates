# Graph Report - roommates  (2026-11-02)

## Corpus Check
- 257 files · ~299,722 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 55 file(s) not represented in the graph (top: .templ 38, (none) 6, .css 5)

## Summary
- 3102 nodes · 9953 edges · 204 communities (147 shown, 57 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 727 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c7e9eb27`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- migrations.go
- en
- et
- environment
- Props
- vars
- householdURL
- LKMessage
- DefaultService
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
- dialog.js
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
- dialog_templ.go
- select.js
- floating_ui_dom.js
- field_templ.go
- files
- validation-errors
- select_templ.go
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- vendor-assets.go
- reset-password-page
- use_hover.js
- login-page
- ListNotesRow
- LK
- NewInternationalization
- files
- ExpensePage
- Props
- initRoutes
- calendar.js
- environment
- nav
- nav
- .UploadFile
- services
- db
- queries/README.md
- NotePage
- ErrorCode
- db
- context
- healthcheck
- app
- github.com/a-h/templ.Component
- images.go
- Props
- floating_ui_core.js
- chatMessages
- Debt
- tabbable.js
- useDismiss
- avatar
- Session
- Hub
- Subscription
- mail
- invite-page
- chat
- chat
- use_transition_status.js
- .clean
- forgot-password-page
- composite.js
- .handleUpload
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- mailer.go
- tasks
- New
- .settlementAction
- inputgroup_templ.go
- switch.js
- mounted
- scroll_lock.js
- InitGinEngine
- io.ReadSeeker
- calendar_templ.go
- T
- github.com/gin-gonic/gin.Context
- error-page
- appearance
- popover.js
- mail
- Assets
- avatar.js
- reset-password-page
- avatar
- tailwind
- noRowsAs
- registration-page
- floating_focus_manager.js
- NormalizeInviteCode
- .notifyChat
- useListNavigation
- server
- restart-signal
- dev
- Browser check
- .GetInvite
- tailwind-version
- scripts-watch
- password
- expense-split.js
- checkbox.js
- HTMLTemplRenderer
- popover_templ.go
- password
- useClick
- offline
- registration-page
- mark_others.js
- error-page
- useTypeahead
- ListBuyAgainRow
- internal_backdrop.js
- membershipFromContext
- invite-page
- icon.go
- TestNetDebts
- ConfirmSettlementParams
- appearance
- Redirect
- GetReplyPreviewParams
- build
- seed
- alertdialog_templ.go
- playwright
- pushPreview
- CountListOnlyFilesParams
- DeleteInviteParams
- DeleteSettlementParams
- DeleteSettlementSharesParams
- DeleteShoppingListParams
- EndMembershipParams
- InsertFileRefParams
- InsertSettlementParams
- SetMembershipRoleParams
- SetShoppingListNameParams
- UpsertPushSubscriptionParams

## God Nodes (most connected - your core abstractions)
1. `CN()` - 103 edges
2. `T()` - 93 edges
3. `Membership` - 61 edges
4. `T()` - 61 edges
5. `membershipFromContext()` - 57 edges
6. `expenses` - 53 edges
7. `expenses` - 53 edges
8. `RenderTemplStatus()` - 47 edges
9. `Button()` - 43 edges
10. `InternalError()` - 41 edges

## Surprising Connections (you probably didn't know these)
- `Check each page` --references--> `Redirect()`  [INFERRED]
  .claude/skills/browser-check/SKILL.md → app/internal/handlers/handlers.go
- `GetSessionKeyPairs()` --calls--> `ValidateCookieHashBlockPair()`  [EXTRACTED]
  app/internal/config/globals-env.go → app/internal/validators/validators.go
- `chatPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-chat.go → app/internal/handlers/handlers-households.go
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go
- `expensesPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-expenses.go → app/internal/handlers/handlers-households.go

## Import Cycles
- None detected.

## Communities (204 total, 57 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.12
Nodes (49): TestSetAppearance(), TestSafeNextPath(), TestAvatars(), TestChat(), TestChatFileRefs(), TestChatReply(), TestHouseholdEvents(), englishContext() (+41 more)

### Community 1 - "handlers.go"
Cohesion: 0.13
Nodes (11): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), DayOfWeek, ctxKey, ctxKey, Props, appearanceKey (+3 more)

### Community 2 - "migrations.go"
Cohesion: 0.18
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.08
Nodes (25): title, cancel, delete, link-invalid, save, en, app, common (+17 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (30): title, cancel, delete, link-invalid, save, et, app, common (+22 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "Props"
Cohesion: 0.31
Nodes (10): Content(), ContentClass(), variantClasses(), Align, ContentProps, GroupProps, Props, ReactionsProps (+2 more)

### Community 7 - "vars"
Cohesion: 0.20
Nodes (9): use_gitignore, vars, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_VERSION (+1 more)

### Community 8 - "householdURL"
Cohesion: 0.17
Nodes (19): GetAppearance(), GetLocaleString(), HeaderComponent(), HtmlWrap(), Scripts(), appShell(), ChatBadge(), HouseholdNav (+11 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "DefaultService"
Cohesion: 0.16
Nodes (10): Queries, closePaidDebt(), payShare(), DefaultService, check(), must(), seed(), seedExpenses() (+2 more)

### Community 11 - "membershipFromContext"
Cohesion: 0.19
Nodes (7): optionalUUID(), pageError(), RenderTemplStatus(), Note, ParseUUID(), NotFound(), FileTargetLoader

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
Nodes (49): ListHouseholdsOfUserRow, Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetHouseholds(), WithHouseholds(), state() (+41 more)

### Community 16 - "expenses"
Cohesion: 0.04
Nodes (53): expenses, add, amount, back, balances, cancel, category, category-groceries (+45 more)

### Community 17 - "household"
Cohesion: 0.06
Nodes (33): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+25 more)

### Community 18 - "household"
Cohesion: 0.06
Nodes (33): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+25 more)

### Community 19 - "permissions"
Cohesion: 0.25
Nodes (7): hooks, PreToolUse, permissions, allow, ask, deny, $schema

### Community 20 - "linters"
Cohesion: 0.29
Nodes (6): allow-empty, linters, enable, settings, exhaustruct_v5, version

### Community 21 - "Certificates"
Cohesion: 0.40
Nodes (4): Certificates, Generate digital certificate, Generate RSA private key, Install Openssl

### Community 22 - "dialog.js"
Cohesion: 0.20
Nodes (23): closeDialog(), createInternalBackdrop(), destroyDialog(), dialogFor(), ensureDialog(), getDialog(), init(), isDialogOpen() (+15 more)

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
Cohesion: 0.12
Nodes (35): GetExpenseRow, ListShoppingItemsRow, ExpenseFilter, Membership, householdPath(), ShoppingList, IfElse(), FormatMoney() (+27 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.04
Nodes (53): ListExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, ListMessagesRow, GetShoppingListRow, ListShoppingListsRow, avatarURL(), EmailToken (+45 more)

### Community 46 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 47 - "dialog_templ.go"
Cohesion: 0.16
Nodes (16): Close(), CloseFor(), initialOpen(), state(), Trigger(), TriggerFor(), ContentProps, ctxKey (+8 more)

### Community 48 - "select.js"
Cohesion: 0.10
Nodes (53): allContents(), armCancelOpen(), arrowScrollStep(), clamp(), close(), closeAll(), commitValue(), contentFor() (+45 more)

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "field_templ.go"
Cohesion: 0.16
Nodes (15): childrenHTML(), fieldClasses(), Separator(), ContentProps, DescriptionProps, ErrorProps, GroupProps, LabelProps (+7 more)

### Community 51 - "files"
Cohesion: 0.05
Nodes (42): one, other, one, other, files, added-by, all, any-place (+34 more)

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "select_templ.go"
Cohesion: 0.17
Nodes (14): initialValue(), Align, ContentProps, ctxKey, ctxState, GroupProps, ItemData, ItemProps (+6 more)

### Community 54 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.06
Nodes (42): ActionProps, DescriptionProps, Props, TitleProps, Variant, variantClasses(), Size, initialChecked() (+34 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "vendor-assets.go"
Cohesion: 0.06
Nodes (38): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+30 more)

### Community 59 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 60 - "use_hover.js"
Cohesion: 0.09
Nodes (41): applySafePolygonPointerEventsMutation(), clearSafePolygonPointerEventsMutation(), createHoverInteraction(), emitClosed(), getDelay(), getRestMs(), hasIntersectingEdge(), isClickLikeOpenEvent() (+33 more)

### Community 61 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 63 - "LK"
Cohesion: 0.16
Nodes (11): LK, N(), unescapeJSON(), N(), chatControlLabel(), chatControlTip(), optionSelect(), deleteConfirm() (+3 more)

### Community 64 - "NewInternationalization"
Cohesion: 0.17
Nodes (9): DeleteCookie(), GetCookie(), GetHeader(), localeMissing(), NewInternationalization(), setLocaleContext(), CookieKey, HTTPHeaderKey (+1 more)

### Community 65 - "files"
Cohesion: 0.05
Nodes (42): one, other, one, other, files, added-by, all, any-place (+34 more)

### Community 66 - "ExpensePage"
Cohesion: 0.13
Nodes (26): ListFilesRow, Expense, AttachedFile, FileFilter, ChatFilePicks(), ExpensePage(), attachedFiles(), attachedFilesContent() (+18 more)

### Community 67 - "Props"
Cohesion: 0.21
Nodes (11): ActionProps, CancelProps, Size, Variant, Variants(), Props, CaptionLayout, DateLib (+3 more)

### Community 68 - "initRoutes"
Cohesion: 0.14
Nodes (16): initRoutes(), registerAssets(), TestAssetCaching(), NewAppearance(), NewGuestOnly(), NewOnlyAuthenticated(), NewZerolog(), NewRateLimit() (+8 more)

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
Cohesion: 0.21
Nodes (6): cleanFileName(), FileMeta, FileTarget, SplitFileName(), TestCleanFileName(), TestSplitFileName()

### Community 74 - "services"
Cohesion: 0.22
Nodes (8): ports, restart, services, mailpit, version, volumes, db_data, rustfs_data

### Community 75 - "db"
Cohesion: 0.25
Nodes (8): depends_on, condition, required, restart, db, rustfs, condition, required

### Community 77 - "NotePage"
Cohesion: 0.31
Nodes (9): localTime(), liveRegion(), HouseholdNotes(), NoteBody(), NotePage(), notePin(), notePrivate(), noteTime() (+1 more)

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
Cohesion: 0.05
Nodes (96): CN(), DirectionAttributes(), flattenClasses(), Ptr(), RandomID(), twMapTokens(), twReorder(), Action() (+88 more)

### Community 84 - "images.go"
Cohesion: 0.11
Nodes (13): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+5 more)

### Community 85 - "Props"
Cohesion: 0.60
Nodes (4): Badge(), badgeClasses(), Props, Variant

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "chatMessages"
Cohesion: 0.27
Nodes (14): Chat, Message, MessageFile, MergeAttributes(), chatFilePicker(), chatMessages(), ChatPage(), replySnippet() (+6 more)

### Community 88 - "Debt"
Cohesion: 0.10
Nodes (14): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, checkSettled(), Debt, ExpensesOverview, netDebts(), ownerActsOnMember() (+6 more)

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
Cohesion: 0.17
Nodes (5): LogValidationErrors(), ValidateCookieHashBlockPair(), gorillaSessionGetter, Session, SessionInfo

### Community 93 - "Hub"
Cohesion: 0.18
Nodes (7): Hub, NewHub(), TestHub(), DefaultHandler, NewDefaultHandler(), newFakeStore(), fakeStore

### Community 94 - "Subscription"
Cohesion: 0.10
Nodes (15): Message, Pusher, Subscription, newFakePusher(), subscribeTestDevice(), takePushes(), TestPushNotifications(), TestValidPushSubscription() (+7 more)

### Community 95 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

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

### Community 100 - ".clean"
Cohesion: 0.18
Nodes (7): InsertExpenseParams, ListActiveMembersRow, FormatCentsInput(), ParseCents(), splitEqual(), TestParseCents(), expenseShare

### Community 101 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - ".handleUpload"
Cohesion: 0.18
Nodes (6): filesPath(), householdPath(), noteFromForm(), noteLengthMessage(), notePath(), notesPath()

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
Nodes (10): cmds, desc, tasks, restart-watch, templ, vendor, cmds, desc (+2 more)

### Community 114 - "New"
Cohesion: 0.38
Nodes (4): NewHouseholdMember(), notFound(), New(), Renderer

### Community 116 - "inputgroup_templ.go"
Cohesion: 0.24
Nodes (10): Type, inputGroupButtonClasses(), AddonProps, Align, ButtonProps, ButtonSize, InputProps, Props (+2 more)

### Community 117 - "switch.js"
Cohesion: 0.31
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "scroll_lock.js"
Cohesion: 0.11
Nodes (15): anchoredPopupScrollLock(), AnimationFrame, getViewportScroller(), hasInsetScrollbars(), isOverflowElement(), isPageScrollLocked(), preventScrollInsetScrollbars(), cleanup() (+7 more)

### Community 124 - "InitGinEngine"
Cohesion: 0.10
Nodes (28): ensureDatabasePool(), main(), InitGinEngine(), GetAppBaseURL(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), IsDevelopment() (+20 more)

### Community 125 - "io.ReadSeeker"
Cohesion: 0.12
Nodes (4): Store, New(), nullStore, S3Store

### Community 126 - "calendar_templ.go"
Cohesion: 0.26
Nodes (11): Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON(), toJSON(), DateAfter (+3 more)

### Community 127 - "T"
Cohesion: 0.16
Nodes (42): GetUserRow, HouseholdOverview, T(), Alert(), Button(), Checkbox(), T(), ValidationErrors() (+34 more)

### Community 128 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.12
Nodes (6): notice(), settingsPath(), sessionUserID(), InternalError(), RenderTempl(), stayInPlace()

### Community 129 - "error-page"
Cohesion: 0.50
Nodes (4): error-page, back, description, title

### Community 131 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

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
Cohesion: 0.10
Nodes (6): hashToken(), noRowsAs(), NoteInput, rowsAffected(), trimToLength(), NoteForm

### Community 146 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "NormalizeInviteCode"
Cohesion: 0.50
Nodes (4): newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

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
Cohesion: 0.40
Nodes (5): cmds, desc, env, SHADCN_TEMPL_DEV, dev

### Community 154 - "Browser check"
Cohesion: 0.29
Nodes (6): Browser check, Check each page, Email, Limits, Sign in, Start

### Community 156 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 157 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 158 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 163 - "popover_templ.go"
Cohesion: 0.22
Nodes (12): state(), Trigger(), triggerID(), Align, ContentProps, ctxKey, ctxState, DescriptionProps (+4 more)

### Community 164 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "offline"
Cohesion: 0.50
Nodes (4): offline, description, retry, title

### Community 167 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 168 - "mark_others.js"
Cohesion: 0.52
Nodes (6): applyAttributeToOthers(), buildKeepSet(), collectOutsideElements(), correctElements(), markOthers(), unwrapHost()

### Community 170 - "error-page"
Cohesion: 0.50
Nodes (4): back, description, title, error-page

### Community 171 - "useTypeahead"
Cohesion: 0.53
Nodes (4): useTypeahead(), getMatchingIndex(), isItemAvailable(), onKeyDown()

### Community 179 - "membershipFromContext"
Cohesion: 0.16
Nodes (7): expenseFromForm(), expensePath(), shoppingItemPath(), membershipFromContext(), shoppingListPath(), ExpenseInput, ShoppingItem

### Community 180 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 181 - "icon.go"
Cohesion: 0.73
Nodes (5): attrString(), generateSVG(), getIconContent(), Icon(), Props

### Community 182 - "TestNetDebts"
Cohesion: 0.67
Nodes (3): TestNetDebts(), TestSplitEqual(), testUUID()

### Community 184 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 185 - "Redirect"
Cohesion: 0.19
Nodes (7): chatPath(), peoplePath(), IsRequestHTMX(), safeNextPath(), Redirect(), swapMainOnError(), FromGinContext()

### Community 187 - "build"
Cohesion: 0.50
Nodes (4): cmds, desc, task, build

### Community 188 - "seed"
Cohesion: 0.50
Nodes (4): cmds, desc, dotenv, seed

### Community 189 - "alertdialog_templ.go"
Cohesion: 0.15
Nodes (18): ContentProps, ctxKey, ctxState, DescriptionProps, FooterProps, HeaderProps, MediaProps, OverlayProps (+10 more)

### Community 190 - "playwright"
Cohesion: 0.50
Nodes (3): npx, playwright, @playwright/mcp

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 968 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **57 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `error-page`, `appearance`, `reset-password-page`, `expenses`, `household`, `password`, `offline`, `registration-page`, `shopping`, `home-page`, `notes`, `files`, `invite-page`, `validation-errors`, `verify-email`, `login-page`, `nav`, `avatar`, `mail`, `chat`, `push`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.11683599419448476 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `mail`, `avatar`, `household`, `registration-page`, `password`, `expenses`, `shopping`, `validation-errors`, `error-page`, `home-page`, `notes`, `appearance`, `verify-email`, `reset-password-page`, `files`, `nav`, `invite-page`, `chat`, `forgot-password-page`, `push`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12962399478374267 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `github.com/gin-gonic/gin.Context`, `handlers.go`, `householdURL`, `membershipFromContext`, `context.Context`, `.notifyChat`, `newEmailToken`, `Membership`, `membershipFromContext`, `Redirect`, `LK`, `ExpensePage`, `NotePage`, `github.com/a-h/templ.Component`, `chatMessages`, `.handleUpload`, `.shoppingAction`, `.settlementAction`, `InitGinEngine`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._