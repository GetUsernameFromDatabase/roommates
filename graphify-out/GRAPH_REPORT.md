# Graph Report - roommates  (2026-11-22)

## Corpus Check
- 362 files · ~428,448 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 73 file(s) not represented in the graph (top: .templ 53, (none) 9, .css 5)

## Summary
- 4410 nodes · 18978 edges · 211 communities (164 shown, 47 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1325 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b4d262a8`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- github.com/a-h/templ.Component
- en
- et
- ListActiveMembersRow
- github.com/jackc/pgx/v5/pgtype.Date
- vars
- RecurringRow
- LKMessage
- github.com/a-h/templ.Component
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
- membershipFromContext
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
- select_templ.go
- recurring
- Membership
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- Hub
- select.js
- floating_ui_dom.js
- environment
- files
- Debt
- github.com/a-h/templ.Attributes
- notes
- migrations.go
- verify-email
- verify-email
- DefaultService
- reset-password-page
- use_hover.js
- inputgroup_templ.go
- choreFromContext
- T
- isRoutesSel
- files
- Membership
- noRowsAs
- login-page
- calendar.js
- environment
- nav
- nav
- Membership
- services
- db
- queries/README.md
- data-export
- ErrorCode
- environment
- context
- healthcheck
- newRouteTest
- vendor-assets.go
- images.go
- export
- floating_ui_core.js
- data-export
- logNote
- tabbable.js
- useDismiss
- avatar
- sessions.go
- github.com/gin-gonic/gin.Context
- Subscription
- options
- RenderTemplStatus
- chat
- chat
- use_transition_status.js
- reset-password-page
- cloudflared
- composite.js
- appearance
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
- initRoutes
- confirm-email
- switch.js
- mounted
- scroll_lock.js
- account
- .renderEvent
- Event
- account
- TestEventScopes
- invite-page
- common
- popover.js
- mail
- forgot-password-page
- avatar.js
- .GetPaymentDetails
- avatar
- tailwind
- FormatDate
- floating_focus_manager.js
- SplitFileName
- common
- useListNavigation
- server
- restart-signal
- dev
- Browser check
- backup.sh
- HouseholdSettingsSection
- Message
- home-page
- expense-split.js
- checkbox.js
- popover_templ.go
- time.Time
- login-page
- useClick
- registration-page
- registerAssets
- mark_others.js
- mail
- useTypeahead
- count
- internal_backdrop.js
- date-picker.js
- payment
- count
- scripts-watch
- DefaultService
- DBTX
- confirm-email
- tailwind-version
- ListExportExpensesRow
- seed
- appearance
- playwright
- payment
- forgot-password-page
- deleted-with-keep
- deleted-with
- deleted-with-keep
- deleted-with
- restart-watch
- HTMLTemplRenderer
- testUUID
- LockChoreParams
- FieldError
- pushPreview
- jump-nav.js
- lastTopLevelColon
- home-page
- AccountDeleteBlocked
- app-key.sh

## God Nodes (most connected - your core abstractions)
1. `T()` - 177 edges
2. `Membership` - 114 edges
3. `CN()` - 108 edges
4. `T()` - 101 edges
5. `membershipFromContext()` - 91 edges
6. `newTestService()` - 86 edges
7. `expenses` - 76 edges
8. `expenses` - 76 edges
9. `InternalError()` - 70 edges
10. `newTestHousehold()` - 69 edges

## Surprising Connections (you probably didn't know these)
- `HTMX actions` --references--> `Redirect()`  [INFERRED]
  .claude/skills/browser-check/SKILL.md → app/internal/handlers/handlers.go
- `TestLimitBody()` --calls--> `limitBody()`  [INFERRED]
  app/cmd/server/routes-bodylimit_test.go → app/cmd/server/routes.go
- `newRouteTest()` --calls--> `initRoutes()`  [INFERRED]
  app/cmd/server/routes-household_test.go → app/cmd/server/routes.go
- `eventFromContext()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-calendar.go → app/internal/handlers/handlers.go
- `eventFromContext()` --calls--> `swapMainOnError()`  [INFERRED]
  app/internal/handlers/handlers-calendar.go → app/internal/handlers/handlers.go

## Import Cycles
- None detected.

## Communities (211 total, 47 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.06
Nodes (90): TestLimitBody(), TestEmptyEnvIsUnset(), TestSetAppearance(), TestSetLanguage(), TestSafeNextPath(), TestRateLimit(), TestSecurityHeaderHSTSOnlyInRelease(), FormatIBAN() (+82 more)

### Community 1 - "handlers.go"
Cohesion: 0.17
Nodes (9): membershipKey, ctxKey, ctxKey, Props, appearanceKey, directionKey, householdsKey, pendingSettlementsKey (+1 more)

### Community 2 - "github.com/a-h/templ.Component"
Cohesion: 0.04
Nodes (83): CN(), flattenClasses(), MergeAttributes(), twMapTokens(), twReorder(), Media(), AspectRatio(), Avatar() (+75 more)

### Community 3 - "en"
Cohesion: 0.06
Nodes (30): title, en, app, email, email-placeholder, error-page, language, logout (+22 more)

### Community 4 - "et"
Cohesion: 0.05
Nodes (42): title, back, description, title, et, app, email, email-placeholder (+34 more)

### Community 5 - "ListActiveMembersRow"
Cohesion: 0.07
Nodes (54): ListActiveMembersRow, EventDetail, Occurrence, balanceRows(), csvAmount(), csvText(), expenseRows(), FormatSize() (+46 more)

### Community 6 - "github.com/jackc/pgx/v5/pgtype.Date"
Cohesion: 0.03
Nodes (38): ClaimChoreNotifiedParams, ClaimChoreNudgedParams, ClaimEventReminderParams, CountUpcomingEventsParams, DeleteDetachedEventsParams, DeleteEventExceptionsParams, EventException, EventReminder (+30 more)

### Community 7 - "vars"
Cohesion: 0.25
Nodes (7): use_gitignore, vars, PORT, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_VERSION, version

### Community 8 - "RecurringRow"
Cohesion: 0.08
Nodes (21): InsertExpenseParams, UpdateExpenseParams, ListRecurringParams, ListRecurringRow, ListRecurringSharesRow, logExpenseEdit(), FormatCentsInput(), ExpenseInput (+13 more)

### Community 9 - "LKMessage"
Cohesion: 0.07
Nodes (23): LKMessage, TestChangePasswordIsValid(), TestChangeUsernameIsValid(), validateCurrentPassword(), EmailForm, ResetPassword, Household, IsModelValid() (+15 more)

### Community 10 - "github.com/a-h/templ.Component"
Cohesion: 0.06
Nodes (50): ActionProps, DescriptionProps, Props, TitleProps, Variant, readOnly(), PurgeDate(), Action() (+42 more)

### Community 11 - "InternalError"
Cohesion: 0.14
Nodes (8): settingsPath(), sessionUserID(), hxLocation(), InternalError(), IsRequestHTMX(), stay(), swapMainOnError(), ValidReminderHour()

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
Nodes (102): ListHouseholdsOfUserRow, ListBuyAgainRow, timestamptz(), csvComma(), writeCSV(), Store, New(), Direction() (+94 more)

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

### Community 23 - "membershipFromContext"
Cohesion: 0.12
Nodes (7): expenseFromForm(), expenseInputError(), membershipFromContext(), pageError(), RenderTempl(), ShoppingItem, FileTargetLoader

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
Cohesion: 0.08
Nodes (11): GetAppBaseURL(), GetUserCredentialsRow, GetUserRow, InsertUserParams, LoginRequest, NormalizeEmail(), zipJSON(), zipName() (+3 more)

### Community 41 - "recurring"
Cohesion: 0.17
Nodes (12): recurring, month, month-short, never-short, next-due, none, repeat, title (+4 more)

### Community 42 - "select_templ.go"
Cohesion: 0.10
Nodes (25): DirectionAttributes(), Content(), Group(), initialValue(), Item(), itemLabel(), Label(), Select() (+17 more)

### Community 43 - "recurring"
Cohesion: 0.17
Nodes (12): recurring, month, month-short, never-short, next-due, none, repeat, title (+4 more)

### Community 44 - "Membership"
Cohesion: 0.06
Nodes (72): ListShoppingItemsRow, ActivityRow, Expense, ExpenseFilter, AttachedFile, FileFilter, Membership, ShoppingList (+64 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.02
Nodes (98): ListExpenseSharesRow, LockExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, ListMessagesRow, RecurringExpense, GetShoppingListRow, ListShoppingListsRow (+90 more)

### Community 46 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 47 - "Hub"
Cohesion: 0.17
Nodes (7): Hub, DefaultHandler, NewDefaultHandler(), newFakeStore(), newFakePusher(), fakePusher, fakeStore

### Community 48 - "select.js"
Cohesion: 0.10
Nodes (53): allContents(), armCancelOpen(), arrowScrollStep(), clamp(), close(), closeAll(), commitValue(), contentFor() (+45 more)

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "environment"
Cohesion: 0.11
Nodes (19): environment, APP_BASE_URL, BACKUP_ALERT_EMAIL, BACKUP_KEEP_DAYS, PGDATABASE, PGHOST, PGPASSWORD, PGUSER (+11 more)

### Community 51 - "files"
Cohesion: 0.05
Nodes (38): files, added-by, all, any-place, any-type, attachments, delete-confirm, delete-keep (+30 more)

### Community 52 - "Debt"
Cohesion: 0.15
Nodes (11): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), iPaidAmount(), ListDebtItemsParams (+3 more)

### Community 53 - "github.com/a-h/templ.Attributes"
Cohesion: 0.05
Nodes (78): ActionProps, CancelProps, ContentProps, ctxKey, ctxState, DescriptionProps, FooterProps, HeaderProps (+70 more)

### Community 54 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 55 - "migrations.go"
Cohesion: 0.19
Nodes (12): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateToLatestConcurrently(), TestMigrateUpAndDown() (+4 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "DefaultService"
Cohesion: 0.10
Nodes (35): loggedKinds(), TestActivityMoney(), wantKinds(), choreTestHousehold(), everyDayChore(), historyKinds(), mustChore(), mustCreateChore() (+27 more)

### Community 59 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 60 - "use_hover.js"
Cohesion: 0.09
Nodes (41): applySafePolygonPointerEventsMutation(), clearSafePolygonPointerEventsMutation(), createHoverInteraction(), emitClosed(), getDelay(), getRestMs(), hasIntersectingEdge(), isClickLikeOpenEvent() (+33 more)

### Community 61 - "inputgroup_templ.go"
Cohesion: 0.17
Nodes (15): Props, Addon(), Button(), Input(), InputGroup(), inputGroupButtonClasses(), Text(), passwordInput() (+7 more)

### Community 62 - "choreFromContext"
Cohesion: 0.25
Nodes (4): choreFormError(), choreFromContext(), choreInput(), choreInputFromForm()

### Community 63 - "T"
Cohesion: 0.08
Nodes (75): LK, DataExportStatus, HouseholdOverview, PaymentDetails, T(), Button(), T(), ValidationErrors() (+67 more)

### Community 65 - "files"
Cohesion: 0.05
Nodes (38): files, added-by, all, any-place, any-type, attachments, delete-confirm, delete-keep (+30 more)

### Community 66 - "Membership"
Cohesion: 0.14
Nodes (14): Chore, choreNextDue(), cleanChore(), firstAssignee(), ChoreInput, ChoreTurn, lockChore(), nextWeekday() (+6 more)

### Community 67 - "noRowsAs"
Cohesion: 0.05
Nodes (29): Queries, GetExpenseRow, GetInviteByCodeRow, expenseOf(), listName(), logActivity(), logShare(), isUniqueViolation() (+21 more)

### Community 68 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 69 - "calendar.js"
Cohesion: 0.09
Nodes (44): addToRange(), build(), calculateFocusTarget(), CalendarDay, calendarFns(), addMonths(), createDateLib(), getISOWeek() (+36 more)

### Community 70 - "environment"
Cohesion: 0.25
Nodes (8): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, S3_BUCKET, environment

### Community 71 - "nav"
Cohesion: 0.12
Nodes (16): nav, calendar, chat, chores, expenses, files, households-add, label (+8 more)

### Community 72 - "nav"
Cohesion: 0.12
Nodes (16): nav, calendar, chat, chores, expenses, files, households-add, label (+8 more)

### Community 73 - "Membership"
Cohesion: 0.11
Nodes (9): ListFilesRow, Chat, ReplyPreview, FileMeta, FileTarget, insertStoredFile(), ListFilesParams, storedFile (+1 more)

### Community 74 - "services"
Cohesion: 0.17
Nodes (11): command, logging, ports, volumes, services, db, volumes, db_data (+3 more)

### Community 75 - "db"
Cohesion: 0.15
Nodes (14): depends_on, ports, restart, stop_grace_period, volumes, depends_on, condition, required (+6 more)

### Community 77 - "data-export"
Cohesion: 0.22
Nodes (9): description, download, new, pending, ready-until, request, requested, title (+1 more)

### Community 79 - "environment"
Cohesion: 0.10
Nodes (22): environment, environment, APP_S3_ACCESS_KEY, APP_S3_SECRET_KEY, DATABASE_HOST, DATABASE_MAX_CONNS, DATABASE_PORT, GIN_MODE (+14 more)

### Community 80 - "context"
Cohesion: 0.24
Nodes (12): build, build, context, dockerfile, build, build, build, ports (+4 more)

### Community 81 - "healthcheck"
Cohesion: 0.57
Nodes (7): healthcheck, healthcheck, interval, retries, test, timeout, healthcheck

### Community 82 - "newRouteTest"
Cohesion: 0.06
Nodes (42): ensureDatabasePool(), main(), newRouteTest(), randomUUID(), send(), TestDeletedHouseholdIsReadOnly(), TestHouseholdRoutesHideFromStrangers(), InitGinEngine() (+34 more)

### Community 83 - "vendor-assets.go"
Cohesion: 0.06
Nodes (38): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+30 more)

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (12): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+4 more)

### Community 85 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "data-export"
Cohesion: 0.22
Nodes (9): description, download, new, pending, ready-until, request, requested, title (+1 more)

### Community 88 - "logNote"
Cohesion: 0.23
Nodes (5): ListNotesRow, NoteInput, logNote(), ListNotesParams, NoteForm

### Community 89 - "tabbable.js"
Cohesion: 0.14
Nodes (28): activeElement(), appendCandidates(), appendMatchingElements(), contains(), disableFocusInside(), enableFocusInside(), focusable(), getComposedChildren() (+20 more)

### Community 90 - "useDismiss"
Cohesion: 0.19
Nodes (25): contains(), getTarget(), isLastTraversableNode(), normalizeBubbles(), useDismiss(), addTargetEventListenerOnce(), close(), closeOnEscapeKeyDown() (+17 more)

### Community 91 - "avatar"
Cohesion: 0.18
Nodes (11): cancel, choose, crop-hint, description, invalid, remove, remove-confirm, title (+3 more)

### Community 92 - "sessions.go"
Cohesion: 0.29
Nodes (3): gorillaSessionGetter, Session, SessionInfo

### Community 93 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.08
Nodes (3): notice(), sendCSV(), stayInPlace()

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (8): Message, Pusher, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

### Community 95 - "options"
Cohesion: 0.16
Nodes (16): logging, logging, restart, volumes, driver, options, logging, max-file (+8 more)

### Community 96 - "RenderTemplStatus"
Cohesion: 0.11
Nodes (11): backTo(), optionalUUID(), safeNextPath(), noteFromForm(), noteLengthMessage(), Redirect(), RenderTemplStatus(), Note (+3 more)

### Community 97 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 98 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 100 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 101 - "cloudflared"
Cohesion: 0.20
Nodes (10): condition, command, depends_on, environment, logging, profiles, restart, app (+2 more)

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

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
Cohesion: 0.13
Nodes (9): buildMessage(), Mailer, Message, isLocalHost(), TestBuildMessage(), TestIsLocalHost(), LogMailer, SMTPMailer (+1 more)

### Community 112 - "tasks"
Cohesion: 0.18
Nodes (11): cmds, desc, task, tasks, build, templ, vendor, cmds (+3 more)

### Community 113 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 114 - "New"
Cohesion: 0.47
Nodes (3): notFound(), New(), Renderer

### Community 115 - "initRoutes"
Cohesion: 0.07
Nodes (30): accountRoutes(), initRoutes(), limitBody(), DeleteCookie(), GetCookie(), GetHeader(), setLanguageCookie(), reload() (+22 more)

### Community 116 - "confirm-email"
Cohesion: 0.33
Nodes (6): done-description, done-title, invalid-description, taken-description, taken-title, confirm-email

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
Cohesion: 0.08
Nodes (25): current-password, delete-blocked-last-owner, delete-blocked-unsettled, delete-confirm, delete-description, delete-submit, delete-title, deleted-description (+17 more)

### Community 125 - ".renderEvent"
Cohesion: 0.19
Nodes (6): calendarDayPath(), eventFormError(), eventFromContext(), eventInputFromForm(), eventOn(), newEventInput()

### Community 126 - "Event"
Cohesion: 0.08
Nodes (40): Event, canChangeEvent(), checkOccurrence(), clockOf(), copyEvent(), DateOnly(), daysBetween(), dropEdits() (+32 more)

### Community 127 - "account"
Cohesion: 0.08
Nodes (25): current-password, delete-blocked-last-owner, delete-blocked-unsettled, delete-confirm, delete-description, delete-submit, delete-title, deleted-description (+17 more)

### Community 128 - "TestEventScopes"
Cohesion: 0.29
Nodes (14): calendarEvents(), dailyEvent(), eventInput(), eventOn(), eventTarget(), mustCreateEvent(), mustOccurrence(), TestEventCap() (+6 more)

### Community 129 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 131 - "common"
Cohesion: 0.20
Nodes (10): cancel, close, delete, former-member, link-invalid, request-failed, save, saved (+2 more)

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "mail"
Cohesion: 0.12
Nodes (17): mail, account-deleted-body, account-deleted-subject, account-erased-body, account-erased-subject, data-export-body, data-export-subject, email-change-body (+9 more)

### Community 137 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 141 - "avatar"
Cohesion: 0.18
Nodes (11): cancel, choose, crop-hint, description, invalid, remove, remove-confirm, title (+3 more)

### Community 142 - "tailwind"
Cohesion: 0.33
Nodes (7): preconditions, msg, sh, cmds, desc, preconditions, tailwind

### Community 146 - "FormatDate"
Cohesion: 0.13
Nodes (22): choreNextNotice(), Chore, ChoreDetail, ChoreLogEntry, ChoreMember, FormatDate(), Badge(), badgeClasses() (+14 more)

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "SplitFileName"
Cohesion: 0.50
Nodes (4): TestZipName(), uniqueZipName(), SplitFileName(), TestSplitFileName()

### Community 149 - "common"
Cohesion: 0.20
Nodes (10): cancel, close, delete, former-member, link-invalid, request-failed, save, saved (+2 more)

### Community 150 - "useListNavigation"
Cohesion: 0.24
Nodes (12): doSwitch(), isCrossOrientationCloseKey(), isCrossOrientationOpenKey(), isMainOrientationKey(), isMainOrientationToEndKey(), isVirtualClick(), isVirtualPointerEvent(), useListNavigation() (+4 more)

### Community 151 - "server"
Cohesion: 0.29
Nodes (7): SERVER_PORT, SMTP_HOST, cmds, desc, dotenv, env, server

### Community 152 - "restart-signal"
Cohesion: 0.33
Nodes (6): cmds, desc, sources, use_gitignore, exclude, restart-signal

### Community 153 - "dev"
Cohesion: 0.33
Nodes (6): cmds, desc, env, RUSTFS_CONSOLE_ENABLE, SHADCN_TEMPL_DEV, dev

### Community 154 - "Browser check"
Cohesion: 0.18
Nodes (10): Browser check, Check each page, Clean up, Email, HTMX actions, Limits, Report, Selectors (+2 more)

### Community 156 - "HouseholdSettingsSection"
Cohesion: 0.24
Nodes (13): HouseholdSettings, NotificationSetting, reminderHourLabel(), currencySelect(), HouseholdSettingsSection(), reminderHourSelect(), timeZoneSelect(), timeZoneLabel() (+5 more)

### Community 157 - "Message"
Cohesion: 0.48
Nodes (7): Message, MessageFile, endsRun(), sameRun(), splitFiles(), startsRun(), TestChatRuns()

### Community 158 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 162 - "popover_templ.go"
Cohesion: 0.11
Nodes (23): RandomID(), Content(), Description(), Header(), Popover(), state(), Title(), Trigger() (+15 more)

### Community 163 - "time.Time"
Cohesion: 0.09
Nodes (52): CalendarItem, Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON(), toJSON() (+44 more)

### Community 164 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 167 - "registerAssets"
Cohesion: 0.50
Nodes (5): cacheHeaders(), registerAssets(), TestAssetCaching(), Hash(), URL()

### Community 168 - "mark_others.js"
Cohesion: 0.52
Nodes (6): applyAttributeToOthers(), buildKeepSet(), collectOutsideElements(), correctElements(), markOthers(), unwrapHost()

### Community 170 - "mail"
Cohesion: 0.12
Nodes (17): mail, account-deleted-body, account-deleted-subject, account-erased-body, account-erased-subject, data-export-body, data-export-subject, email-change-body (+9 more)

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

### Community 181 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 183 - "DefaultService"
Cohesion: 0.26
Nodes (22): addTestSession(), countRows(), isDeleted(), newTestEmail(), setTestPassword(), TestChangeEmail(), TestChangePassword(), TestDeleteAccountBlocked() (+14 more)

### Community 185 - "confirm-email"
Cohesion: 0.33
Nodes (6): done-description, done-title, invalid-description, taken-description, taken-title, confirm-email

### Community 186 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 187 - "ListExportExpensesRow"
Cohesion: 0.29
Nodes (5): ListExportExpensesRow, ListExportSharesRow, exportFixture(), ListExportExpensesParams, ListExportSharesParams

### Community 188 - "seed"
Cohesion: 0.50
Nodes (4): cmds, desc, dotenv, seed

### Community 189 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 190 - "playwright"
Cohesion: 0.50
Nodes (3): npx, playwright, @playwright/mcp

### Community 192 - "payment"
Cohesion: 0.20
Nodes (10): payment, copy-iban, description, iban, iban-invalid, name, name-invalid, show (+2 more)

### Community 193 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

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

### Community 201 - "testUUID"
Cohesion: 0.40
Nodes (5): TestCopyInputFormOrder(), TestExpenseCleanField(), TestNetDebts(), TestSplitEqual(), testUUID()

### Community 205 - "jump-nav.js"
Cohesion: 0.70
Nodes (4): line(), queue(), section(), update()

### Community 206 - "lastTopLevelColon"
Cohesion: 0.67
Nodes (3): lastTopLevelColon(), twImportantToV3(), twImportantToV4()

### Community 208 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 221 - "app-key.sh"
Cohesion: 0.90
Nodes (4): admin(), must(), new_secret(), app-key.sh script

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1310 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **47 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `common`, `avatar`, `household`, `login-page`, `expenses`, `shopping`, `mail`, `recurring`, `notes`, `confirm-email`, `verify-email`, `reset-password-page`, `payment`, `files`, `forgot-password-page`, `nav`, `home-page`, `data-export`, `chat`, `appearance`, `push`, `export`, `account`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.05777491408934708 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `handlers.go`, `github.com/a-h/templ.Component`, `ListActiveMembersRow`, `RecurringRow`, `github.com/a-h/templ.Component`, `InternalError`, `context.Context`, `.handleUpload`, `FormatDate`, `membershipFromContext`, `time.Time`, `GetAppBaseURL`, `Membership`, `github.com/a-h/templ.Attributes`, `inputgroup_templ.go`, `choreFromContext`, `Membership`, `noRowsAs`, `Membership`, `newRouteTest`, `github.com/gin-gonic/gin.Context`, `RenderTemplStatus`, `initRoutes`, `.renderEvent`, `Event`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Should `github.com/a-h/templ.Component` be split into smaller, more focused modules?**
  _Cohesion score 0.03742314889067094 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `invite-page`, `mail`, `forgot-password-page`, `expenses`, `household`, `common`, `home-page`, `registration-page`, `shopping`, `recurring`, `notes`, `files`, `payment`, `verify-email`, `appearance`, `login-page`, `nav`, `data-export`, `export`, `avatar`, `chat`, `reset-password-page`, `push`, `confirm-email`, `account`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.06451612903225806 - nodes in this community are weakly interconnected._