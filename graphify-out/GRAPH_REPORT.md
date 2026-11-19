# Graph Report - roommates  (2026-11-20)

## Corpus Check
- 332 files · ~411,460 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 71 file(s) not represented in the graph (top: .templ 51, (none) 9, .css 5)

## Summary
- 4278 nodes · 17531 edges · 221 communities (178 shown, 43 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1324 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4fb2221b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- github.com/a-h/templ.Component
- en
- et
- github.com/jackc/pgx/v5/pgtype.Date
- Chore
- vars
- RecurringRow
- LKMessage
- newRouteTest
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
- newEmailToken
- recurring
- calendar
- recurring
- householdURL
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- calendar/calendar_templ.go
- select.js
- floating_ui_dom.js
- environment
- files
- Props
- balanceRows
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- choreTestHousehold
- reset-password-page
- use_hover.js
- home-page
- NewInternationalization
- T
- Membership
- files
- Membership
- Variants
- generate-locales-keys.go
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
- environment
- context
- healthcheck
- .UpdateNote
- tooltip_templ.go
- images.go
- export
- floating_ui_core.js
- calendar
- noRowsAs
- tabbable.js
- useDismiss
- avatar
- Session
- github.com/gin-gonic/gin.Context
- Subscription
- options
- newTestUser
- chat
- chat
- use_transition_status.js
- Membership
- cloudflared
- composite.js
- ListActiveMembersRow
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
- membershipFromContext
- switch.js
- mounted
- scroll_lock.js
- account
- .renderEvent
- time.Time
- account
- confirm-email
- registration-page
- common
- popover.js
- mail
- DefaultService
- avatar.js
- reset-password-page
- avatar
- tailwind
- confirm-email
- HouseholdExpenses
- floating_focus_manager.js
- Message
- common
- useListNavigation
- server
- restart-signal
- dev
- Browser check
- backup.sh
- popover_templ.go
- testUUID
- invite-page
- expense-split.js
- checkbox.js
- .loadShoppingItem
- time.Time
- appearance
- useClick
- alertdialog_templ.go
- migrations.go
- mark_others.js
- mail
- useTypeahead
- count
- internal_backdrop.js
- date-picker.js
- payment
- count
- scripts-watch
- LoadLocales
- login-page
- Redirect
- tailwind-version
- SplitFileName
- seed
- Debt
- playwright
- payment
- registration-page
- deleted-with-keep
- deleted-with
- deleted-with-keep
- deleted-with
- restart-watch
- HTMLTemplRenderer
- forgot-password-page
- calendar/calendar_templ.go
- FieldError
- Hub
- dialog_templ.go
- FormatIBAN
- home-page
- select_templ.go
- invite-page
- bubble_templ.go
- appearance
- icon.go
- FormatCentsInput
- Props
- fieldClasses
- Props
- pushPreview
- app-key.sh

## God Nodes (most connected - your core abstractions)
1. `T()` - 167 edges
2. `Membership` - 114 edges
3. `CN()` - 108 edges
4. `T()` - 99 edges
5. `membershipFromContext()` - 91 edges
6. `newTestService()` - 80 edges
7. `expenses` - 76 edges
8. `expenses` - 76 edges
9. `InternalError()` - 67 edges
10. `RenderTemplStatus()` - 65 edges

## Surprising Connections (you probably didn't know these)
- `HTMX actions` --references--> `Redirect()`  [INFERRED]
  .claude/skills/browser-check/SKILL.md → app/internal/handlers/handlers.go
- `newRouteTest()` --calls--> `initRoutes()`  [INFERRED]
  app/cmd/server/routes-household_test.go → app/cmd/server/routes.go
- `calendarPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-calendar.go → app/internal/handlers/handlers-households.go
- `eventFromContext()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-calendar.go → app/internal/handlers/handlers.go
- `eventFromContext()` --calls--> `swapMainOnError()`  [INFERRED]
  app/internal/handlers/handlers-calendar.go → app/internal/handlers/handlers.go

## Import Cycles
- None detected.

## Communities (221 total, 43 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.08
Nodes (74): TestEmptyEnvIsUnset(), TestSetAppearance(), TestSetLanguage(), TestSafeNextPath(), TestActivity(), TestActivityMembers(), TestActivityNotes(), TestActivityShopping() (+66 more)

### Community 1 - "handlers.go"
Cohesion: 0.20
Nodes (10): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), membershipKey, ctxKey, appearanceKey, directionKey, householdsKey (+2 more)

### Community 2 - "github.com/a-h/templ.Component"
Cohesion: 0.04
Nodes (111): CN(), DirectionAttributes(), flattenClasses(), MergeAttributes(), RandomID(), twMapTokens(), twReorder(), Action() (+103 more)

### Community 3 - "en"
Cohesion: 0.05
Nodes (36): title, en, app, email, email-placeholder, error-page, forgot-password-page, language (+28 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (30): title, back, description, title, et, app, email, email-placeholder (+22 more)

### Community 5 - "github.com/jackc/pgx/v5/pgtype.Date"
Cohesion: 0.03
Nodes (40): InsertExpenseParams, UpdateExpenseParams, logExpenseEdit(), ExpenseInput, ParseCents(), splitEqual(), TestSplitEqual(), advanceRecurring() (+32 more)

### Community 6 - "Chore"
Cohesion: 0.10
Nodes (41): readOnly(), EventDetail, Occurrence, Chore, ChoreDetail, ChoreLogEntry, ChoreMember, PurgeDate() (+33 more)

### Community 7 - "vars"
Cohesion: 0.25
Nodes (7): use_gitignore, vars, PORT, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_VERSION, version

### Community 8 - "RecurringRow"
Cohesion: 0.17
Nodes (7): ListRecurringParams, ListRecurringRow, ListRecurringSharesRow, RecurringInput, RecurringRow, lockRecurringToChange(), ListRecurringSharesParams

### Community 9 - "LKMessage"
Cohesion: 0.06
Nodes (24): LKMessage, TestChangePasswordIsValid(), TestChangeUsernameIsValid(), validateCurrentPassword(), EmailForm, ResetPassword, Household, IsModelValid() (+16 more)

### Community 10 - "newRouteTest"
Cohesion: 0.07
Nodes (41): ensureDatabasePool(), main(), newRouteTest(), randomUUID(), send(), TestDeletedHouseholdIsReadOnly(), TestHouseholdRoutesHideFromStrangers(), InitGinEngine() (+33 more)

### Community 11 - "InternalError"
Cohesion: 0.12
Nodes (8): settingsPath(), sessionUserID(), InternalError(), setLanguageCookie(), reload(), Language, NormalizeEmail(), FromGinContext()

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
Nodes (94): ListHouseholdsOfUserRow, ListBuyAgainRow, GetUserCredentialsRow, payeeDetails(), Store, New(), Direction(), WithDirection() (+86 more)

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
Cohesion: 0.19
Nodes (10): backTo(), choreFromContext(), choreInput(), expensesPath(), safeNextPath(), pageError(), RenderTemplStatus(), stay() (+2 more)

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

### Community 40 - "newEmailToken"
Cohesion: 0.10
Nodes (7): GetUserRow, TestUserLimit(), newEmailToken(), timestamptz(), newInviteCode(), TestNewInviteCode(), userLimit

### Community 41 - "recurring"
Cohesion: 0.09
Nodes (22): recurring, add, delete-confirm, due, edit-title, month, month-short, never-short (+14 more)

### Community 42 - "calendar"
Cohesion: 0.11
Nodes (19): create, day-empty, day-more, empty, month-1, month-2, month-3, month-4 (+11 more)

### Community 43 - "recurring"
Cohesion: 0.09
Nodes (22): recurring, add, delete-confirm, due, edit-title, month, month-short, never-short (+14 more)

### Community 44 - "householdURL"
Cohesion: 0.08
Nodes (60): choreNextNotice(), LK, ActivityRow, Expense, AttachedFile, FileFilter, N(), unescapeJSON() (+52 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.02
Nodes (87): ListExpenseSharesRow, LockExpenseSharesRow, ListActiveInvitesRow, ListMessagesRow, RecurringExpense, ListNotesRow, GetShoppingListRow, ListShoppingListsRow (+79 more)

### Community 46 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 47 - "calendar/calendar_templ.go"
Cohesion: 0.14
Nodes (24): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+16 more)

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

### Community 52 - "Props"
Cohesion: 0.70
Nodes (3): badgeClasses(), Props, Variant

### Community 53 - "balanceRows"
Cohesion: 0.17
Nodes (12): ListExportExpensesRow, ListExportSharesRow, balanceRows(), csvAmount(), csvComma(), csvText(), expenseRows(), exportFixture() (+4 more)

### Community 54 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.05
Nodes (49): ActionProps, DescriptionProps, TitleProps, Size, initialChecked(), Props, initialValue(), Props (+41 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "choreTestHousehold"
Cohesion: 0.34
Nodes (12): choreTestHousehold(), everyDayChore(), historyKinds(), mustChore(), mustCreateChore(), rotaOf(), TestChoreFixedRotation(), TestChoreLeastRecent() (+4 more)

### Community 59 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 60 - "use_hover.js"
Cohesion: 0.09
Nodes (41): applySafePolygonPointerEventsMutation(), clearSafePolygonPointerEventsMutation(), createHoverInteraction(), emitClosed(), getDelay(), getRestMs(), hasIntersectingEdge(), isClickLikeOpenEvent() (+33 more)

### Community 61 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 62 - "NewInternationalization"
Cohesion: 0.22
Nodes (8): DeleteCookie(), GetCookie(), GetHeader(), localeMissing(), NewInternationalization(), setLocaleContext(), CookieKey, HTTPHeaderKey

### Community 63 - "T"
Cohesion: 0.14
Nodes (43): Button(), T(), ValidationErrors(), Field(), Group(), Label(), Input(), Button() (+35 more)

### Community 64 - "Membership"
Cohesion: 0.13
Nodes (36): ListShoppingItemsRow, HouseholdSettings, Membership, NotificationSetting, ShoppingList, IfElse(), ChatFilePicks(), expenseFilters() (+28 more)

### Community 65 - "files"
Cohesion: 0.05
Nodes (38): files, added-by, all, any-place, any-type, attachments, delete-confirm, delete-keep (+30 more)

### Community 66 - "Membership"
Cohesion: 0.16
Nodes (9): Chore, cleanChore(), firstAssignee(), ChoreInput, ChoreTurn, lockChore(), setChoreRota(), choreMemberPicked() (+1 more)

### Community 67 - "Variants"
Cohesion: 0.31
Nodes (9): ActionProps, CancelProps, Size, Type, Variant, inputGroupButtonClasses(), Props, ButtonProps (+1 more)

### Community 68 - "generate-locales-keys.go"
Cohesion: 0.08
Nodes (14): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput(), discoverLanguages() (+6 more)

### Community 69 - "calendar.js"
Cohesion: 0.09
Nodes (44): addToRange(), build(), calculateFocusTarget(), CalendarDay, calendarFns(), addMonths(), createDateLib(), getISOWeek() (+36 more)

### Community 70 - "environment"
Cohesion: 0.25
Nodes (8): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, S3_BUCKET, environment

### Community 71 - "nav"
Cohesion: 0.13
Nodes (15): nav, calendar, chat, chores, expenses, files, households-add, label (+7 more)

### Community 72 - "nav"
Cohesion: 0.13
Nodes (15): nav, calendar, chat, chores, expenses, files, households-add, label (+7 more)

### Community 73 - "Membership"
Cohesion: 0.20
Nodes (5): ListFilesRow, cleanFileName(), FileMeta, FileTarget, ListFilesParams

### Community 74 - "services"
Cohesion: 0.17
Nodes (11): command, logging, ports, volumes, services, db, volumes, db_data (+3 more)

### Community 75 - "db"
Cohesion: 0.15
Nodes (14): depends_on, ports, restart, stop_grace_period, volumes, depends_on, condition, required (+6 more)

### Community 77 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 79 - "environment"
Cohesion: 0.10
Nodes (22): environment, environment, APP_S3_ACCESS_KEY, APP_S3_SECRET_KEY, DATABASE_HOST, DATABASE_MAX_CONNS, DATABASE_PORT, GIN_MODE (+14 more)

### Community 80 - "context"
Cohesion: 0.24
Nodes (12): build, build, context, dockerfile, build, build, build, ports (+4 more)

### Community 81 - "healthcheck"
Cohesion: 0.57
Nodes (7): healthcheck, healthcheck, interval, retries, test, timeout, healthcheck

### Community 82 - ".UpdateNote"
Cohesion: 0.19
Nodes (5): noteFromForm(), noteLengthMessage(), notePath(), notesPath(), Note

### Community 83 - "tooltip_templ.go"
Cohesion: 0.31
Nodes (8): DisableTrigger(), state(), Trigger(), ContentProps, ctxKey, ctxState, Props, Side

### Community 84 - "images.go"
Cohesion: 0.12
Nodes (11): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+3 more)

### Community 85 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "calendar"
Cohesion: 0.11
Nodes (19): create, day-empty, day-more, empty, month-1, month-2, month-3, month-4 (+11 more)

### Community 88 - "noRowsAs"
Cohesion: 0.05
Nodes (30): Queries, GetExpenseRow, GetInviteByCodeRow, InsertUserParams, expenseOf(), listName(), logActivity(), logShare() (+22 more)

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
Cohesion: 0.24
Nodes (4): LogValidationErrors(), gorillaSessionGetter, Session, SessionInfo

### Community 93 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.08
Nodes (5): choreFormError(), choreInputFromForm(), notice(), sendCSV(), stayInPlace()

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (7): Message, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

### Community 95 - "options"
Cohesion: 0.16
Nodes (16): logging, logging, restart, volumes, driver, options, logging, max-file (+8 more)

### Community 96 - "newTestUser"
Cohesion: 0.33
Nodes (16): newTestEmail(), setTestPassword(), TestChangeEmail(), TestChangePassword(), TestEmailChangeTakenAtConfirm(), englishContext(), sessionsOf(), TestDeleteExpiredEmailTokens() (+8 more)

### Community 97 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 98 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 100 - "Membership"
Cohesion: 0.29
Nodes (3): Chat, ReplyPreview, Upload

### Community 101 - "cloudflared"
Cohesion: 0.20
Nodes (10): condition, command, depends_on, environment, logging, profiles, restart, app (+2 more)

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "ListActiveMembersRow"
Cohesion: 0.08
Nodes (43): ListActiveMembersRow, HouseholdOverview, Ptr(), Variants(), Checkbox(), accountOutcome(), eventFieldError(), eventFields() (+35 more)

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
Cohesion: 0.08
Nodes (28): accountRoutes(), TestLimitBody(), cacheHeaders(), initRoutes(), limitBody(), registerAssets(), TestAssetCaching(), NewAppearance() (+20 more)

### Community 116 - "membershipFromContext"
Cohesion: 0.17
Nodes (5): expenseFromForm(), expenseInputError(), expensePath(), membershipFromContext(), RenderTempl()

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

### Community 125 - ".renderEvent"
Cohesion: 0.18
Nodes (7): calendarDayPath(), calendarPath(), eventFormError(), eventFromContext(), eventInputFromForm(), eventOn(), eventPath()

### Community 126 - "time.Time"
Cohesion: 0.08
Nodes (43): GetHouseholdRow, Event, canChangeEvent(), checkOccurrence(), clockOf(), copyEvent(), DateOnly(), daysBetween() (+35 more)

### Community 127 - "account"
Cohesion: 0.13
Nodes (15): current-password, email-description, email-same, email-sent, email-submit, new-email, new-password, password-changed (+7 more)

### Community 128 - "confirm-email"
Cohesion: 0.33
Nodes (6): done-description, done-title, invalid-description, taken-description, taken-title, confirm-email

### Community 129 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 131 - "common"
Cohesion: 0.29
Nodes (7): cancel, delete, link-invalid, request-failed, save, saved, common

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "mail"
Cohesion: 0.18
Nodes (11): mail, email-change-body, email-change-subject, email-changed-body, email-changed-subject, invite-body, invite-subject, reset-body (+3 more)

### Community 137 - "DefaultService"
Cohesion: 0.13
Nodes (26): newEventInput(), addTestSession(), loggedKinds(), TestActivityMoney(), wantKinds(), EventInput, calendarEvents(), dailyEvent() (+18 more)

### Community 139 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 141 - "avatar"
Cohesion: 0.18
Nodes (11): cancel, choose, crop-hint, description, invalid, remove, remove-confirm, title (+3 more)

### Community 142 - "tailwind"
Cohesion: 0.33
Nodes (7): preconditions, msg, sh, cmds, desc, preconditions, tailwind

### Community 143 - "confirm-email"
Cohesion: 0.33
Nodes (6): done-description, done-title, invalid-description, taken-description, taken-title, confirm-email

### Community 146 - "HouseholdExpenses"
Cohesion: 0.10
Nodes (23): ExpenseFilter, PaymentDetails, FormatSize(), printer(), GetLocaleString(), FormatMoney(), debtBreakdown(), HouseholdExpenses() (+15 more)

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

### Community 156 - "popover_templ.go"
Cohesion: 0.22
Nodes (12): state(), Trigger(), triggerID(), Align, ContentProps, ctxKey, ctxState, DescriptionProps (+4 more)

### Community 157 - "testUUID"
Cohesion: 0.50
Nodes (4): TestCopyInputFormOrder(), TestExpenseCleanField(), TestNetDebts(), testUUID()

### Community 158 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 162 - ".loadShoppingItem"
Cohesion: 0.21
Nodes (3): filesPath(), shoppingItemPath(), ShoppingItem

### Community 163 - "time.Time"
Cohesion: 0.14
Nodes (35): CalendarItem, calendarAgenda(), calendarAgendaDay(), calendarAgendaMore(), calendarCell(), calendarChip(), calendarDayList(), calendarItem() (+27 more)

### Community 164 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "alertdialog_templ.go"
Cohesion: 0.14
Nodes (16): ContentProps, ctxKey, ctxState, DescriptionProps, FooterProps, HeaderProps, MediaProps, OverlayProps (+8 more)

### Community 167 - "migrations.go"
Cohesion: 0.18
Nodes (12): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateToLatestConcurrently(), TestMigrateUpAndDown() (+4 more)

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

### Community 181 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 183 - "LoadLocales"
Cohesion: 0.18
Nodes (15): TestNotifyCalendarDue(), notificationOn(), TestNotificationSettings(), subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes(), daysAgo() (+7 more)

### Community 184 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 185 - "Redirect"
Cohesion: 0.11
Nodes (12): chatPath(), optionalUUID(), choresPath(), householdPath(), peoplePath(), hxLocation(), IsRequestHTMX(), recurringPath() (+4 more)

### Community 186 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 188 - "seed"
Cohesion: 0.50
Nodes (4): cmds, desc, dotenv, seed

### Community 189 - "Debt"
Cohesion: 0.10
Nodes (11): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), DBTX, ListDebtItemsParams (+3 more)

### Community 190 - "playwright"
Cohesion: 0.50
Nodes (3): npx, playwright, @playwright/mcp

### Community 192 - "payment"
Cohesion: 0.20
Nodes (10): payment, copy-iban, description, iban, iban-invalid, name, name-invalid, show (+2 more)

### Community 193 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

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

### Community 201 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 202 - "calendar/calendar_templ.go"
Cohesion: 0.18
Nodes (17): Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON(), toJSON(), CaptionLayout (+9 more)

### Community 204 - "Hub"
Cohesion: 0.15
Nodes (7): Hub, DefaultHandler, NewDefaultHandler(), newFakeStore(), newFakePusher(), fakePusher, fakeStore

### Community 206 - "dialog_templ.go"
Cohesion: 0.16
Nodes (16): Close(), CloseFor(), initialOpen(), state(), Trigger(), TriggerFor(), ContentProps, ctxKey (+8 more)

### Community 207 - "FormatIBAN"
Cohesion: 0.67
Nodes (4): FormatIBAN(), NormalizeIBAN(), TestIBAN(), ValidIBAN()

### Community 208 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 209 - "select_templ.go"
Cohesion: 0.15
Nodes (14): initialValue(), Align, ContentProps, ctxKey, ctxState, GroupProps, ItemData, ItemProps (+6 more)

### Community 210 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 211 - "bubble_templ.go"
Cohesion: 0.31
Nodes (10): Content(), ContentClass(), variantClasses(), Align, ContentProps, GroupProps, Props, ReactionsProps (+2 more)

### Community 213 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 214 - "icon.go"
Cohesion: 0.73
Nodes (5): attrString(), generateSVG(), getIconContent(), Icon(), Props

### Community 215 - "FormatCentsInput"
Cohesion: 0.50
Nodes (4): FormatCentsInput(), TestParseCents(), iPaidAmount(), TestIPaidAmountLeavesOutPending()

### Community 216 - "Props"
Cohesion: 0.67
Nodes (3): Props, Variant, variantClasses()

### Community 217 - "fieldClasses"
Cohesion: 0.67
Nodes (3): fieldClasses(), Orientation, Props

### Community 218 - "Props"
Cohesion: 0.67
Nodes (3): initialChecked(), Props, Size

### Community 221 - "app-key.sh"
Cohesion: 0.90
Nodes (4): admin(), must(), new_secret(), app-key.sh script

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1303 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **43 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `T()` connect `Chore` to `handlers.go`, `github.com/a-h/templ.Component`, `github.com/jackc/pgx/v5/pgtype.Date`, `newRouteTest`, `InternalError`, `context.Context`, `HouseholdExpenses`, `RenderTemplStatus`, `.loadShoppingItem`, `time.Time`, `newEmailToken`, `householdURL`, `balanceRows`, `Redirect`, `T`, `Membership`, `Membership`, `.UpdateNote`, `noRowsAs`, `github.com/gin-gonic/gin.Context`, `Membership`, `ListActiveMembersRow`, `membershipFromContext`, `.renderEvent`, `time.Time`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.07530864197530865 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `registration-page`, `common`, `avatar`, `confirm-email`, `household`, `invite-page`, `expenses`, `shopping`, `mail`, `recurring`, `notes`, `verify-email`, `reset-password-page`, `payment`, `files`, `nav`, `forgot-password-page`, `login-page`, `home-page`, `appearance`, `calendar`, `chat`, `push`, `export`, `account`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Should `github.com/a-h/templ.Component` be split into smaller, more focused modules?**
  _Cohesion score 0.043758043758043756 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `confirm-email`, `mail`, `reset-password-page`, `expenses`, `household`, `common`, `appearance`, `shopping`, `recurring`, `calendar`, `notes`, `files`, `payment`, `login-page`, `verify-email`, `home-page`, `registration-page`, `nav`, `invite-page`, `export`, `avatar`, `chat`, `push`, `account`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.05405405405405406 - nodes in this community are weakly interconnected._