# Graph Report - roommates  (2026-11-19)

## Corpus Check
- 319 files · ~381,817 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 68 file(s) not represented in the graph (top: .templ 48, (none) 9, .css 5)

## Summary
- 4059 nodes · 15852 edges · 231 communities (175 shown, 56 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1181 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8e5ba85d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- migrations.go
- en
- et
- T
- choreFields
- vars
- github.com/jackc/pgx/v5/pgtype.Date
- LKMessage
- NewDefaultService
- InternalError
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
- recurring
- generate-locales-languages.go
- recurring
- Membership
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- vendor-assets.go
- select.js
- floating_ui_dom.js
- environment
- files
- newTestUser
- balanceRows
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- generate-locales-keys.go
- reset-password-page
- use_hover.js
- home-page
- noRowsAs
- T
- navIcon
- files
- Membership
- LoadLocales
- go_pkg_fmt
- calendar.js
- environment
- nav
- nav
- Membership
- services
- db
- queries/README.md
- login-page
- ErrorCode
- db
- context
- healthcheck
- vendor-assets_test.go
- github.com/a-h/templ.Component
- images.go
- export
- floating_ui_core.js
- chores
- Debt
- tabbable.js
- useDismiss
- avatar
- Session
- Hub
- Subscription
- options
- Props
- chat
- chat
- use_transition_status.js
- .notifyChat
- cloudflared
- composite.js
- ExpenseInput
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- mailer.go
- tasks
- export
- New
- DefaultService
- membershipFromContext
- switch.js
- mounted
- scroll_lock.js
- account
- .changeFileRef
- calendar_templ.go
- account
- github.com/gin-gonic/gin.Context
- app
- forgot-password-page
- popover.js
- mail
- chores
- avatar.js
- reset-password-page
- avatar
- tailwind
- logNote
- confirm-email
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
- forgot-password-page
- NotePage
- expense-split.js
- checkbox.js
- .UpdateNote
- LK
- appearance
- useClick
- NewInternationalization
- registration-page
- mark_others.js
- mail
- useTypeahead
- count
- internal_backdrop.js
- date-picker.js
- payment
- count
- icon.go
- DefaultService
- appearance
- Redirect
- tailwind-version
- vendor
- seed
- ContentProps
- playwright
- payment
- timestamptz
- deleted-with-keep
- deleted-with
- deleted-with-keep
- deleted-with
- restart-watch
- HTMLTemplRenderer
- ContentProps
- Props
- FieldError
- .handleUpload
- confirm-email
- registration-page
- FormatIBAN
- home-page
- login-page
- invite-page
- pushPreview
- invite-page
- Props
- state
- Props
- testUUID
- Props
- fieldClasses
- state
- Size
- initialChecked
- initialValue
- Props
- LegendProps
- AddonProps
- Props
- ContentProps
- TriggerProps
- ContentProps

## God Nodes (most connected - your core abstractions)
1. `T()` - 146 edges
2. `CN()` - 106 edges
3. `Membership` - 100 edges
4. `T()` - 85 edges
5. `membershipFromContext()` - 84 edges
6. `expenses` - 76 edges
7. `expenses` - 76 edges
8. `newTestService()` - 75 edges
9. `InternalError()` - 63 edges
10. `RenderTemplStatus()` - 62 edges

## Surprising Connections (you probably didn't know these)
- `HTMX actions` --references--> `Redirect()`  [INFERRED]
  .claude/skills/browser-check/SKILL.md → app/internal/handlers/handlers.go
- `newRouteTest()` --calls--> `initRoutes()`  [INFERRED]
  app/cmd/server/routes-household_test.go → app/cmd/server/routes.go
- `chatPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-chat.go → app/internal/handlers/handlers-households.go
- `choresPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-chores.go → app/internal/handlers/handlers-households.go
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go

## Import Cycles
- None detected.

## Communities (231 total, 56 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.08
Nodes (69): TestSetAppearance(), TestSetLanguage(), TestSafeNextPath(), TestActivity(), TestActivityMembers(), TestActivityNotes(), TestActivityShopping(), TestAvatars() (+61 more)

### Community 1 - "handlers.go"
Cohesion: 0.14
Nodes (31): ctxKey, ctxState, flattenClasses(), lastTopLevelColon(), Ptr(), twImportantToV3(), twImportantToV4(), twMapTokens() (+23 more)

### Community 2 - "migrations.go"
Cohesion: 0.18
Nodes (12): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateToLatestConcurrently(), TestMigrateUpAndDown() (+4 more)

### Community 3 - "en"
Cohesion: 0.06
Nodes (30): title, en, app, email, email-placeholder, error-page, language, logout (+22 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (30): title, back, description, title, et, app, email, email-placeholder (+22 more)

### Community 5 - "T"
Cohesion: 0.10
Nodes (29): ExpenseFilter, PaymentDetails, householdPath(), printer(), GetLocaleString(), T(), FormatMoney(), AppearanceForm() (+21 more)

### Community 6 - "choreFields"
Cohesion: 0.12
Nodes (34): ListActiveMembersRow, choreNextNotice(), readOnly(), Chore, ChoreDetail, ChoreLogEntry, ChoreMember, PurgeDate() (+26 more)

### Community 7 - "vars"
Cohesion: 0.20
Nodes (9): use_gitignore, vars, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_VERSION (+1 more)

### Community 8 - "github.com/jackc/pgx/v5/pgtype.Date"
Cohesion: 0.04
Nodes (28): Chore, RecurringExpense, ListRecurringParams, ListRecurringRow, ListRecurringSharesRow, advanceRecurring(), RecurringRow, lockRecurringToChange() (+20 more)

### Community 9 - "LKMessage"
Cohesion: 0.06
Nodes (24): LKMessage, TestChangePasswordIsValid(), TestChangeUsernameIsValid(), validateCurrentPassword(), EmailForm, ResetPassword, Household, IsModelValid() (+16 more)

### Community 10 - "NewDefaultService"
Cohesion: 0.08
Nodes (31): ensureDatabasePool(), main(), newRouteTest(), randomUUID(), send(), TestDeletedHouseholdIsReadOnly(), TestHouseholdRoutesHideFromStrangers(), InitGinEngine() (+23 more)

### Community 11 - "InternalError"
Cohesion: 0.19
Nodes (5): settingsPath(), householdPath(), sessionUserID(), InternalError(), stay()

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
Cohesion: 0.01
Nodes (83): ListBuyAgainRow, GetUserCredentialsRow, InsertUserParams, NewHouseholdMember(), Store, New(), Direction(), WithDirection() (+75 more)

### Community 16 - "expenses"
Cohesion: 0.03
Nodes (76): expenses, add, amount, amount-locked-currency, amount-locked-left, back, balances, cancel (+68 more)

### Community 17 - "household"
Cohesion: 0.05
Nodes (42): household, back, balance-not-zero, copied, copy-code, copy-link, create-invite, expires (+34 more)

### Community 18 - "household"
Cohesion: 0.05
Nodes (42): household, back, balance-not-zero, copied, copy-code, copy-link, create-invite, expires (+34 more)

### Community 19 - "permissions"
Cohesion: 0.25
Nodes (7): hooks, PreToolUse, permissions, allow, ask, deny, $schema

### Community 20 - "linters"
Cohesion: 0.29
Nodes (6): allow-empty, linters, enable, settings, exhaustruct_v5, version

### Community 21 - "Roommates"
Cohesion: 0.05
Nodes (34): Certificates, Generate the certificate, Install mkcert, Assets, Third-party files, Alert, An older dump, Backup (+26 more)

### Community 22 - "dialog.js"
Cohesion: 0.20
Nodes (23): closeDialog(), createInternalBackdrop(), destroyDialog(), dialogFor(), ensureDialog(), getDialog(), init(), isDialogOpen() (+15 more)

### Community 23 - "RenderTemplStatus"
Cohesion: 0.11
Nodes (13): chatPath(), optionalUUID(), choreFormError(), choreFromContext(), choreInputFromForm(), expensesPath(), RenderTemplStatus(), swapMainOnError() (+5 more)

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
Cohesion: 0.03
Nodes (76): expenses, add, amount, amount-locked-currency, amount-locked-left, back, balances, cancel (+68 more)

### Community 38 - "shopping"
Cohesion: 0.06
Nodes (36): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+28 more)

### Community 39 - "shopping"
Cohesion: 0.06
Nodes (36): shopping, add, add-bag-expense, all-bought, back, bought-by, buy-again, buying (+28 more)

### Community 40 - "GetAppBaseURL"
Cohesion: 0.27
Nodes (3): GetAppBaseURL(), NormalizeEmail(), newEmailToken()

### Community 41 - "recurring"
Cohesion: 0.11
Nodes (19): recurring, add, delete-confirm, due, edit-title, month, never, next-due (+11 more)

### Community 42 - "generate-locales-languages.go"
Cohesion: 0.24
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 43 - "recurring"
Cohesion: 0.11
Nodes (19): recurring, add, delete-confirm, due, edit-title, month, never, next-due (+11 more)

### Community 44 - "Membership"
Cohesion: 0.14
Nodes (36): ListShoppingItemsRow, FileFilter, Membership, NotificationSetting, ShoppingList, IfElse(), ChatFilePicks(), expenseFilters() (+28 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.02
Nodes (86): ListExpenseSharesRow, LockExpenseSharesRow, GetHouseholdRow, ListHouseholdsOfUserRow, ListActiveInvitesRow, ListMessagesRow, GetShoppingListRow, ListShoppingListsRow (+78 more)

### Community 46 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

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
Cohesion: 0.06
Nodes (36): files, added-by, all, any-place, any-type, attachments, delete-confirm, delete-keep (+28 more)

### Community 52 - "newTestUser"
Cohesion: 0.29
Nodes (17): addTestSession(), newTestEmail(), setTestPassword(), TestChangeEmail(), TestChangePassword(), TestEmailChangeTakenAtConfirm(), englishContext(), sessionsOf() (+9 more)

### Community 53 - "balanceRows"
Cohesion: 0.14
Nodes (14): ListExportExpensesRow, ListExportSharesRow, TestChoreClean(), balanceRows(), csvAmount(), csvComma(), csvText(), expenseRows() (+6 more)

### Community 54 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.03
Nodes (62): ActionProps, DescriptionProps, TitleProps, DescriptionProps, FooterProps, HeaderProps, MediaProps, OverlayProps (+54 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "generate-locales-keys.go"
Cohesion: 0.24
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 59 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 60 - "use_hover.js"
Cohesion: 0.09
Nodes (41): applySafePolygonPointerEventsMutation(), clearSafePolygonPointerEventsMutation(), createHoverInteraction(), emitClosed(), getDelay(), getRestMs(), hasIntersectingEdge(), isClickLikeOpenEvent() (+33 more)

### Community 61 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 62 - "noRowsAs"
Cohesion: 0.06
Nodes (21): Queries, GetExpenseRow, expenseOf(), listName(), logActivity(), logShare(), isUniqueViolation(), hashToken() (+13 more)

### Community 63 - "T"
Cohesion: 0.09
Nodes (61): GetUserRow, HouseholdOverview, HouseholdSettings, Alert(), Action(), Cancel(), Button(), T() (+53 more)

### Community 64 - "navIcon"
Cohesion: 0.18
Nodes (23): Props, GetHouseholds(), AlertDialog(), initialOpen(), Badge(), appShell(), ChatBadge(), confirmDialog() (+15 more)

### Community 65 - "files"
Cohesion: 0.06
Nodes (36): files, added-by, all, any-place, any-type, attachments, delete-confirm, delete-keep (+28 more)

### Community 66 - "Membership"
Cohesion: 0.12
Nodes (17): choreInput(), choreNextDue(), cleanChore(), firstAssignee(), ChoreInput, ChoreTurn, lockChore(), nextWeekday() (+9 more)

### Community 67 - "LoadLocales"
Cohesion: 0.14
Nodes (19): LoadLocales(), notificationOn(), TestNotificationSettings(), subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes(), daysAgo() (+11 more)

### Community 69 - "calendar.js"
Cohesion: 0.09
Nodes (44): addToRange(), build(), calculateFocusTarget(), CalendarDay, calendarFns(), addMonths(), createDateLib(), getISOWeek() (+36 more)

### Community 70 - "environment"
Cohesion: 0.29
Nodes (7): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment

### Community 71 - "nav"
Cohesion: 0.14
Nodes (14): nav, chat, chores, expenses, files, households-add, label, more (+6 more)

### Community 72 - "nav"
Cohesion: 0.14
Nodes (14): nav, chat, chores, expenses, files, households-add, label, more (+6 more)

### Community 73 - "Membership"
Cohesion: 0.16
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

### Community 82 - "vendor-assets_test.go"
Cohesion: 0.39
Nodes (5): fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 83 - "github.com/a-h/templ.Component"
Cohesion: 0.04
Nodes (107): CN(), DirectionAttributes(), MergeAttributes(), RandomID(), Action(), Description(), Title(), Content() (+99 more)

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (12): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+4 more)

### Community 85 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "chores"
Cohesion: 0.08
Nodes (25): create, done, due-date, due-overdue, due-today, due-tomorrow, empty, new (+17 more)

### Community 88 - "Debt"
Cohesion: 0.13
Nodes (14): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, householdSettled(), netDebts(), TestNetDebts() (+6 more)

### Community 89 - "tabbable.js"
Cohesion: 0.14
Nodes (28): activeElement(), appendCandidates(), appendMatchingElements(), contains(), disableFocusInside(), enableFocusInside(), focusable(), getComposedChildren() (+20 more)

### Community 90 - "useDismiss"
Cohesion: 0.19
Nodes (25): contains(), getTarget(), isLastTraversableNode(), normalizeBubbles(), useDismiss(), addTargetEventListenerOnce(), close(), closeOnEscapeKeyDown() (+17 more)

### Community 91 - "avatar"
Cohesion: 0.18
Nodes (11): cancel, choose, crop-hint, description, invalid, remove, remove-confirm, title (+3 more)

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

### Community 96 - "Props"
Cohesion: 0.17
Nodes (13): ActionProps, CancelProps, Size, Type, Variant, inputGroupButtonClasses(), CaptionLayout, DateLib (+5 more)

### Community 97 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 98 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 101 - "cloudflared"
Cohesion: 0.22
Nodes (9): condition, command, depends_on, environment, profiles, restart, app, TUNNEL_TOKEN (+1 more)

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "ExpenseInput"
Cohesion: 0.13
Nodes (12): InsertExpenseParams, UpdateExpenseParams, logExpenseEdit(), FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestParseCents() (+4 more)

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
Cohesion: 0.17
Nodes (12): push, chat-file, chat-title, confirmed-body, confirmed-title, denied, expense-body, expense-title (+4 more)

### Community 109 - "push"
Cohesion: 0.17
Nodes (12): push, chat-file, chat-title, confirmed-body, confirmed-title, denied, expense-body, expense-title (+4 more)

### Community 111 - "mailer.go"
Cohesion: 0.20
Nodes (8): buildMessage(), Message, isLocalHost(), TestBuildMessage(), TestIsLocalHost(), LogMailer, SMTPMailer, fakeMailer

### Community 112 - "tasks"
Cohesion: 0.18
Nodes (11): cmds, desc, task, cmds, desc, tasks, build, scripts-watch (+3 more)

### Community 113 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 115 - "DefaultService"
Cohesion: 0.08
Nodes (25): accountRoutes(), TestLimitBody(), cacheHeaders(), initRoutes(), limitBody(), registerAssets(), TestAssetCaching(), NewAppearance() (+17 more)

### Community 116 - "membershipFromContext"
Cohesion: 0.15
Nodes (9): choresPath(), expenseFromForm(), expenseInputError(), expensePath(), shoppingItemPath(), membershipFromContext(), pageError(), recurringPath() (+1 more)

### Community 117 - "switch.js"
Cohesion: 0.31
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "scroll_lock.js"
Cohesion: 0.11
Nodes (15): anchoredPopupScrollLock(), AnimationFrame, getViewportScroller(), hasInsetScrollbars(), isOverflowElement(), isPageScrollLocked(), preventScrollInsetScrollbars(), cleanup() (+7 more)

### Community 124 - "account"
Cohesion: 0.13
Nodes (15): current-password, email-description, email-same, email-sent, email-submit, new-email, new-password, password-changed (+7 more)

### Community 125 - ".changeFileRef"
Cohesion: 0.16
Nodes (7): backTo(), setLanguageCookie(), safeNextPath(), reload(), Language, FromGinContext(), FileTargetLoader

### Community 126 - "calendar_templ.go"
Cohesion: 0.18
Nodes (15): addIntervals(), nextDue(), TestNextDue(), Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON() (+7 more)

### Community 127 - "account"
Cohesion: 0.13
Nodes (15): current-password, email-description, email-same, email-sent, email-submit, new-email, new-password, password-changed (+7 more)

### Community 128 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.09
Nodes (4): notice(), sendCSV(), RenderTempl(), stayInPlace()

### Community 129 - "app"
Cohesion: 0.33
Nodes (6): env_file, ports, restart, stop_grace_period, volumes, app

### Community 131 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "mail"
Cohesion: 0.18
Nodes (11): mail, email-change-body, email-change-subject, email-changed-body, email-changed-subject, invite-body, invite-subject, reset-body (+3 more)

### Community 137 - "chores"
Cohesion: 0.08
Nodes (25): create, done, due-date, due-overdue, due-today, due-tomorrow, empty, new (+17 more)

### Community 139 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 141 - "avatar"
Cohesion: 0.18
Nodes (11): cancel, choose, crop-hint, description, invalid, remove, remove-confirm, title (+3 more)

### Community 142 - "tailwind"
Cohesion: 0.33
Nodes (7): preconditions, msg, sh, cmds, desc, preconditions, tailwind

### Community 143 - "logNote"
Cohesion: 0.23
Nodes (5): ListNotesRow, NoteInput, logNote(), ListNotesParams, NoteForm

### Community 146 - "confirm-email"
Cohesion: 0.33
Nodes (6): done-description, done-title, invalid-description, taken-description, taken-title, confirm-email

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "Message"
Cohesion: 0.48
Nodes (7): Message, MessageFile, endsRun(), sameRun(), splitFiles(), startsRun(), TestChatRuns()

### Community 149 - "common"
Cohesion: 0.29
Nodes (7): cancel, delete, link-invalid, request-failed, save, saved, common

### Community 150 - "useListNavigation"
Cohesion: 0.24
Nodes (12): doSwitch(), isCrossOrientationCloseKey(), isCrossOrientationOpenKey(), isMainOrientationKey(), isMainOrientationToEndKey(), isVirtualClick(), isVirtualPointerEvent(), useListNavigation() (+4 more)

### Community 151 - "server"
Cohesion: 0.25
Nodes (8): SERVER_PORT, SMTP_HOST, SMTP_PORT, cmds, desc, dotenv, env, server

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
Cohesion: 0.29
Nodes (7): cancel, delete, link-invalid, request-failed, save, saved, common

### Community 157 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 158 - "NotePage"
Cohesion: 0.24
Nodes (13): ActivityRow, ActivityPage(), liveRegion(), HouseholdNotes(), NoteBody(), NotePage(), notePin(), notePrivate() (+5 more)

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 162 - ".UpdateNote"
Cohesion: 0.16
Nodes (5): noteFromForm(), noteLengthMessage(), notePath(), notesPath(), Note

### Community 163 - "LK"
Cohesion: 0.09
Nodes (38): LK, Expense, AttachedFile, FormatSize(), N(), unescapeJSON(), Variants(), Checkbox() (+30 more)

### Community 164 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "NewInternationalization"
Cohesion: 0.22
Nodes (8): DeleteCookie(), GetCookie(), GetHeader(), localeMissing(), NewInternationalization(), setLocaleContext(), CookieKey, HTTPHeaderKey

### Community 167 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 168 - "mark_others.js"
Cohesion: 0.52
Nodes (6): applyAttributeToOthers(), buildKeepSet(), collectOutsideElements(), correctElements(), markOthers(), unwrapHost()

### Community 170 - "mail"
Cohesion: 0.18
Nodes (11): mail, email-change-body, email-change-subject, email-changed-body, email-changed-subject, invite-body, invite-subject, reset-body (+3 more)

### Community 171 - "useTypeahead"
Cohesion: 0.53
Nodes (4): useTypeahead(), getMatchingIndex(), isItemAvailable(), onKeyDown()

### Community 173 - "count"
Cohesion: 0.67
Nodes (4): one, other, count, count

### Community 175 - "date-picker.js"
Cohesion: 0.67
Nodes (5): fill(), isoDay(), pad(), setDay(), shownDay()

### Community 179 - "payment"
Cohesion: 0.20
Nodes (10): payment, copy-iban, description, iban, iban-invalid, name, name-invalid, show (+2 more)

### Community 180 - "count"
Cohesion: 0.67
Nodes (4): one, other, count, count

### Community 181 - "icon.go"
Cohesion: 0.73
Nodes (5): attrString(), generateSVG(), getIconContent(), Icon(), Props

### Community 183 - "DefaultService"
Cohesion: 0.19
Nodes (21): loggedKinds(), TestActivityMoney(), wantKinds(), choreTestHousehold(), everyDayChore(), historyKinds(), mustChore(), mustCreateChore() (+13 more)

### Community 184 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 185 - "Redirect"
Cohesion: 0.27
Nodes (5): peoplePath(), hxLocation(), IsRequestHTMX(), Redirect(), shoppingListPath()

### Community 186 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 187 - "vendor"
Cohesion: 0.67
Nodes (3): vendor, cmds, desc

### Community 188 - "seed"
Cohesion: 0.50
Nodes (4): cmds, desc, dotenv, seed

### Community 190 - "playwright"
Cohesion: 0.50
Nodes (3): npx, playwright, @playwright/mcp

### Community 192 - "payment"
Cohesion: 0.20
Nodes (10): payment, copy-iban, description, iban, iban-invalid, name, name-invalid, show (+2 more)

### Community 193 - "timestamptz"
Cohesion: 0.18
Nodes (6): GetInviteByCodeRow, timestamptz(), newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode()

### Community 194 - "deleted-with-keep"
Cohesion: 0.67
Nodes (3): one, other, deleted-with-keep

### Community 195 - "deleted-with"
Cohesion: 0.67
Nodes (3): one, other, deleted-with

### Community 196 - "deleted-with-keep"
Cohesion: 0.67
Nodes (3): one, other, deleted-with-keep

### Community 197 - "deleted-with"
Cohesion: 0.67
Nodes (3): one, other, deleted-with

### Community 198 - "restart-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, restart-watch

### Community 201 - "ContentProps"
Cohesion: 0.67
Nodes (3): Align, ContentProps, Side

### Community 202 - "Props"
Cohesion: 0.67
Nodes (3): initialChecked(), Props, Size

### Community 205 - "confirm-email"
Cohesion: 0.33
Nodes (6): done-description, done-title, invalid-description, taken-description, taken-title, confirm-email

### Community 206 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 207 - "FormatIBAN"
Cohesion: 0.20
Nodes (5): FormatIBAN(), NormalizeIBAN(), TestIBAN(), ValidIBAN(), payeeDetails()

### Community 208 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 209 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 210 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 213 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 214 - "Props"
Cohesion: 0.33
Nodes (6): variantClasses(), Align, Props, ReactionsProps, Side, Variant

### Community 215 - "state"
Cohesion: 0.33
Nodes (6): Close(), CloseFor(), state(), Trigger(), TriggerFor(), ctxState

### Community 216 - "Props"
Cohesion: 0.70
Nodes (3): badgeClasses(), Props, Variant

### Community 217 - "testUUID"
Cohesion: 0.50
Nodes (4): TestCopyInputFormOrder(), TestExpenseCleanField(), TestSplitEqual(), testUUID()

### Community 218 - "Props"
Cohesion: 0.67
Nodes (3): Props, Variant, variantClasses()

### Community 219 - "fieldClasses"
Cohesion: 0.67
Nodes (3): fieldClasses(), Orientation, Props

### Community 220 - "state"
Cohesion: 0.67
Nodes (3): state(), Trigger(), ctxState

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1299 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **56 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `forgot-password-page`, `mail`, `reset-password-page`, `expenses`, `household`, `common`, `appearance`, `shopping`, `registration-page`, `recurring`, `notes`, `files`, `payment`, `verify-email`, `home-page`, `nav`, `confirm-email`, `login-page`, `invite-page`, `export`, `chores`, `avatar`, `chat`, `push`, `account`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.08432432432432432 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `chores`, `avatar`, `confirm-email`, `household`, `common`, `forgot-password-page`, `expenses`, `shopping`, `mail`, `recurring`, `notes`, `appearance`, `verify-email`, `reset-password-page`, `payment`, `files`, `nav`, `login-page`, `registration-page`, `home-page`, `invite-page`, `chat`, `push`, `export`, `account`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1403567447045708 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `github.com/gin-gonic/gin.Context`, `handlers.go`, `choreFields`, `InternalError`, `context.Context`, `RenderTemplStatus`, `NotePage`, `.UpdateNote`, `LK`, `GetAppBaseURL`, `Membership`, `balanceRows`, `Redirect`, `noRowsAs`, `T`, `navIcon`, `LoadLocales`, `.handleUpload`, `github.com/a-h/templ.Component`, `.notifyChat`, `ExpenseInput`, `membershipFromContext`, `.changeFileRef`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.06451612903225806 - nodes in this community are weakly interconnected._