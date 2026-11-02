# Graph Report - roommates  (2026-11-02)

## Corpus Check
- 258 files · ~300,049 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 55 file(s) not represented in the graph (top: .templ 38, (none) 6, .css 5)

## Summary
- 3103 nodes · 9971 edges · 189 communities (149 shown, 40 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 728 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d5108ebd`
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
- ExpensePage
- T
- files
- appearance
- Props
- InitGinEngine
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
- .UpdateNote
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- mailer.go
- tasks
- New
- models.go
- NewInternationalization
- switch.js
- mounted
- scroll_lock.js
- NewDefaultService
- io.ReadSeeker
- calendar_templ.go
- T
- github.com/gin-gonic/gin.Context
- error-page
- DBTX
- popover.js
- forgot-password-page
- Assets
- avatar.js
- reset-password-page
- avatar
- tailwind
- noRowsAs
- registration-page
- floating_focus_manager.js
- NormalizeInviteCode
- .PostMessage
- useListNavigation
- server
- restart-signal
- dev
- Browser check
- mail
- tailwind-version
- scripts-watch
- TestNetDebts
- expense-split.js
- checkbox.js
- HTMLTemplRenderer
- popover_templ.go
- pushPreview
- useClick
- common
- registration-page
- mark_others.js
- useTypeahead
- Props
- internal_backdrop.js
- membershipFromContext
- invite-page
- icon.go
- appearance
- Redirect
- build
- seed
- alertdialog_templ.go
- playwright
- Props

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

## Communities (189 total, 40 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.12
Nodes (48): TestSetAppearance(), TestSafeNextPath(), TestAvatars(), TestChat(), TestChatFileRefs(), TestChatReply(), TestHouseholdEvents(), englishContext() (+40 more)

### Community 1 - "handlers.go"
Cohesion: 0.13
Nodes (10): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ctxKey, ctxKey, Props, appearanceKey, directionKey (+2 more)

### Community 2 - "migrations.go"
Cohesion: 0.18
Nodes (11): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateUpAndDown(), updateMigrationVersion() (+3 more)

### Community 3 - "en"
Cohesion: 0.07
Nodes (28): title, cancel, delete, link-invalid, save, en, app, common (+20 more)

### Community 4 - "et"
Cohesion: 0.07
Nodes (27): title, back, description, title, et, app, email, error-page (+19 more)

### Community 5 - "environment"
Cohesion: 0.20
Nodes (10): environment, APP_BASE_URL, DATABASE_HOST, DATABASE_PORT, GIN_MODE, S3_ENDPOINT, SERVER_PORT, SMTP_FROM (+2 more)

### Community 6 - "Props"
Cohesion: 0.33
Nodes (6): variantClasses(), Align, Props, ReactionsProps, Side, Variant

### Community 7 - "vars"
Cohesion: 0.20
Nodes (9): use_gitignore, vars, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_VERSION (+1 more)

### Community 8 - "householdURL"
Cohesion: 0.13
Nodes (28): FileFilter, localTime(), fileFilters(), fileUploadForm(), HouseholdFiles(), appShell(), ChatBadge(), HouseholdNav (+20 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "DefaultService"
Cohesion: 0.23
Nodes (14): DefaultHandler, NewDefaultHandler(), sessionsOf(), DefaultService, subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes() (+6 more)

### Community 11 - "membershipFromContext"
Cohesion: 0.12
Nodes (9): expenseFromForm(), expensePath(), expensesPath(), filesPath(), membershipFromContext(), RenderTemplStatus(), stay(), ParseUUID() (+1 more)

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
Nodes (48): ListHouseholdsOfUserRow, GetInviteByCodeRow, ListBuyAgainRow, Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetHouseholds() (+40 more)

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
Cohesion: 0.50
Nodes (3): Certificates, Generate the certificate, Install mkcert

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
Nodes (36): GetExpenseRow, ListShoppingItemsRow, ExpenseFilter, Membership, ShoppingList, IfElse(), FormatMoney(), ChatFilePicks() (+28 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.03
Nodes (41): ListExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, GetShoppingListRow, ListShoppingListsRow, avatarURL(), ConfirmSettlementParams, ConfirmSettlementRow (+33 more)

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
Cohesion: 0.14
Nodes (14): initialValue(), Align, ContentProps, ctxKey, ctxState, GroupProps, ItemData, ItemProps (+6 more)

### Community 54 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.05
Nodes (44): ActionProps, DescriptionProps, TitleProps, Size, initialChecked(), Props, initialValue(), DisableTrigger() (+36 more)

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

### Community 62 - "ListNotesRow"
Cohesion: 0.19
Nodes (4): ListNotesRow, NoteInput, ListNotesParams, NoteForm

### Community 63 - "ExpensePage"
Cohesion: 0.11
Nodes (28): LK, Expense, AttachedFile, N(), N(), chatControlLabel(), ExpensePage(), optionSelect() (+20 more)

### Community 64 - "T"
Cohesion: 0.15
Nodes (20): GetUserRow, Ptr(), GetAppearance(), GetLocaleString(), T(), TestTUnescapesTextButNotArgs(), unescapeJSON(), HeaderComponent() (+12 more)

### Community 65 - "files"
Cohesion: 0.05
Nodes (42): one, other, one, other, files, added-by, all, any-place (+34 more)

### Community 66 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 67 - "Props"
Cohesion: 0.24
Nodes (10): ActionProps, CancelProps, Size, Type, Variant, Variants(), inputGroupButtonClasses(), Props (+2 more)

### Community 68 - "InitGinEngine"
Cohesion: 0.11
Nodes (19): initRoutes(), registerAssets(), TestAssetCaching(), NewAppearance(), NewGuestOnly(), NewOnlyAuthenticated(), NewRateLimit(), TestRateLimit() (+11 more)

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

### Community 77 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

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
Nodes (93): CN(), DirectionAttributes(), flattenClasses(), RandomID(), twMapTokens(), twReorder(), Action(), Description() (+85 more)

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (12): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+4 more)

### Community 85 - "Props"
Cohesion: 0.60
Nodes (4): Badge(), badgeClasses(), Props, Variant

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "chatMessages"
Cohesion: 0.23
Nodes (16): Chat, Message, MessageFile, MergeAttributes(), Content(), ContentClass(), chatFilePicker(), chatMessages() (+8 more)

### Community 88 - "Debt"
Cohesion: 0.10
Nodes (13): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, checkSettled(), Debt, ExpensesOverview, netDebts(), iPaidAmount() (+5 more)

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
Cohesion: 0.18
Nodes (4): LogValidationErrors(), gorillaSessionGetter, Session, SessionInfo

### Community 93 - "Hub"
Cohesion: 0.23
Nodes (5): Hub, NewHub(), TestHub(), newFakePusher(), fakePusher

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (8): Message, Pusher, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

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
Cohesion: 0.11
Nodes (10): InsertExpenseParams, ListActiveMembersRow, ReplyPreview, FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestParseCents() (+2 more)

### Community 101 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - ".UpdateNote"
Cohesion: 0.15
Nodes (5): noteFromForm(), noteLengthMessage(), notePath(), notesPath(), Note

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

### Community 115 - "models.go"
Cohesion: 0.08
Nodes (30): ListMessagesRow, EmailToken, Expense, ExpenseCategory, ExpenseShare, File, FilePurge, FileRef (+22 more)

### Community 116 - "NewInternationalization"
Cohesion: 0.17
Nodes (9): DeleteCookie(), GetCookie(), GetHeader(), localeMissing(), NewInternationalization(), setLocaleContext(), CookieKey, HTTPHeaderKey (+1 more)

### Community 117 - "switch.js"
Cohesion: 0.31
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "scroll_lock.js"
Cohesion: 0.11
Nodes (15): anchoredPopupScrollLock(), AnimationFrame, getViewportScroller(), hasInsetScrollbars(), isOverflowElement(), isPageScrollLocked(), preventScrollInsetScrollbars(), cleanup() (+7 more)

### Community 124 - "NewDefaultService"
Cohesion: 0.12
Nodes (25): ensureDatabasePool(), main(), InitGinEngine(), GetAppBaseURL(), GetDatabaseURL(), GetEnv(), GetSessionKeyPairs(), IsDevelopment() (+17 more)

### Community 125 - "io.ReadSeeker"
Cohesion: 0.09
Nodes (7): newFakeStore(), Store, New(), nullStore, fakeStore, Upload, S3Store

### Community 126 - "calendar_templ.go"
Cohesion: 0.23
Nodes (17): Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON(), toJSON(), CaptionLayout (+9 more)

### Community 127 - "T"
Cohesion: 0.19
Nodes (33): HouseholdOverview, Alert(), Button(), Checkbox(), T(), ValidationErrors(), Field(), Group() (+25 more)

### Community 128 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.15
Nodes (7): notice(), settingsPath(), householdPath(), sessionUserID(), InternalError(), RenderTempl(), stayInPlace()

### Community 129 - "error-page"
Cohesion: 0.50
Nodes (4): error-page, back, description, title

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

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
Cohesion: 0.11
Nodes (8): Queries, hashToken(), closePaidDebt(), payShare(), noRowsAs(), ownerActsOnMember(), rowsAffected(), trimToLength()

### Community 146 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "NormalizeInviteCode"
Cohesion: 0.50
Nodes (4): newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

### Community 149 - ".PostMessage"
Cohesion: 0.43
Nodes (3): chatPath(), optionalUUID(), swapMainOnError()

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

### Community 155 - "mail"
Cohesion: 0.29
Nodes (7): mail, invite-body, invite-subject, reset-body, reset-subject, verify-body, verify-subject

### Community 156 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 157 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 158 - "TestNetDebts"
Cohesion: 0.67
Nodes (3): TestNetDebts(), TestSplitEqual(), testUUID()

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 163 - "popover_templ.go"
Cohesion: 0.22
Nodes (12): state(), Trigger(), triggerID(), Align, ContentProps, ctxKey, ctxState, DescriptionProps (+4 more)

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "common"
Cohesion: 0.40
Nodes (5): cancel, delete, link-invalid, save, common

### Community 167 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 168 - "mark_others.js"
Cohesion: 0.52
Nodes (6): applyAttributeToOthers(), buildKeepSet(), collectOutsideElements(), correctElements(), markOthers(), unwrapHost()

### Community 171 - "useTypeahead"
Cohesion: 0.53
Nodes (4): useTypeahead(), getMatchingIndex(), isItemAvailable(), onKeyDown()

### Community 173 - "Props"
Cohesion: 0.67
Nodes (3): initialChecked(), Props, Size

### Community 179 - "membershipFromContext"
Cohesion: 0.20
Nodes (5): shoppingItemPath(), pageError(), shoppingListPath(), ShoppingItem, FileTargetLoader

### Community 180 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 181 - "icon.go"
Cohesion: 0.73
Nodes (5): attrString(), generateSVG(), getIconContent(), Icon(), Props

### Community 184 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 185 - "Redirect"
Cohesion: 0.23
Nodes (5): peoplePath(), IsRequestHTMX(), safeNextPath(), Redirect(), FromGinContext()

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

### Community 205 - "Props"
Cohesion: 0.67
Nodes (3): Props, Variant, variantClasses()

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 967 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **40 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `error-page`, `forgot-password-page`, `reset-password-page`, `expenses`, `household`, `shopping`, `registration-page`, `home-page`, `notes`, `files`, `invite-page`, `validation-errors`, `verify-email`, `login-page`, `appearance`, `nav`, `avatar`, `mail`, `chat`, `push`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.11990950226244344 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `avatar`, `household`, `registration-page`, `mail`, `expenses`, `common`, `shopping`, `validation-errors`, `home-page`, `notes`, `appearance`, `verify-email`, `reset-password-page`, `files`, `nav`, `login-page`, `invite-page`, `chat`, `forgot-password-page`, `push`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13062377743968703 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `github.com/gin-gonic/gin.Context`, `handlers.go`, `.clean`, `.UpdateNote`, `householdURL`, `newEmailToken`, `membershipFromContext`, `Membership`, `context.Context`, `.shoppingAction`, `github.com/a-h/templ.Component`, `membershipFromContext`, `.PostMessage`, `chatMessages`, `Redirect`, `T`, `ExpensePage`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.06896551724137931 - nodes in this community are weakly interconnected._