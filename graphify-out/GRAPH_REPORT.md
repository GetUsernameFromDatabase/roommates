# Graph Report - roommates  (2026-11-22)

## Corpus Check
- 342 files · ~423,467 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 73 file(s) not represented in the graph (top: .templ 53, (none) 9, .css 5)

## Summary
- 4387 nodes · 18484 edges · 217 communities (170 shown, 47 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1388 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ea065b1f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- CN
- en
- et
- github.com/jackc/pgx/v5/pgtype.Date
- ListActiveMembersRow
- vars
- RecurringRow
- LKMessage
- DefaultService
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
- newEmailToken
- recurring
- select_templ.go
- recurring
- householdURL
- github.com/jackc/pgx/v5/pgtype.UUID
- notes
- calendar/calendar_templ.go
- select.js
- floating_ui_dom.js
- environment
- files
- logNote
- balanceRows
- notes
- generate-locales-keys.go
- verify-email
- verify-email
- choreTestHousehold
- reset-password-page
- use_hover.js
- home-page
- NewInternationalization
- T
- migrations.go
- files
- Membership
- generate-locales-languages.go
- popover_templ.go
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
- .UpdateNote
- github.com/a-h/templ.Component
- images.go
- export
- floating_ui_core.js
- data-export
- noRowsAs
- tabbable.js
- useDismiss
- avatar
- Session
- github.com/gin-gonic/gin.Context
- Subscription
- options
- RenderTemplStatus
- chat
- chat
- use_transition_status.js
- Membership
- cloudflared
- composite.js
- .loadShoppingItem
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
- vendor-assets_test.go
- switch.js
- mounted
- scroll_lock.js
- account
- .renderEvent
- time.Time
- account
- go_pkg_fmt
- DBTX
- common
- popover.js
- mail
- DefaultService
- avatar.js
- reset-password-page
- avatar
- tailwind
- confirm-email
- Badge
- floating_focus_manager.js
- Message
- common
- useListNavigation
- server
- restart-signal
- dev
- Browser check
- backup.sh
- forgot-password-page
- testUUID
- invite-page
- expense-split.js
- checkbox.js
- schedulePicker
- time.Time
- login-page
- useClick
- registration-page
- newRouteTest
- mark_others.js
- mail
- useTypeahead
- count
- internal_backdrop.js
- date-picker.js
- payment
- count
- scripts-watch
- registration-page
- login-page
- Redirect
- tailwind-version
- .writeHouseholdExport
- seed
- Debt
- playwright
- payment
- NewHouseholdMember
- deleted-with-keep
- deleted-with
- deleted-with-keep
- deleted-with
- restart-watch
- HTMLTemplRenderer
- .GetPaymentDetails
- password
- FieldError
- seed
- jump-nav.js
- lastTopLevelColon
- FormatIBAN
- home-page
- FormatCentsInput
- bubble_templ.go
- appearance
- dialog_templ.go
- AccountDeleteBlocked
- pushPreview
- app-key.sh

## God Nodes (most connected - your core abstractions)
1. `T()` - 175 edges
2. `Membership` - 116 edges
3. `CN()` - 108 edges
4. `T()` - 101 edges
5. `membershipFromContext()` - 91 edges
6. `newTestService()` - 85 edges
7. `expenses` - 76 edges
8. `expenses` - 76 edges
9. `InternalError()` - 70 edges
10. `newTestHousehold()` - 68 edges

## Surprising Connections (you probably didn't know these)
- `HTMX actions` --references--> `Redirect()`  [INFERRED]
  .claude/skills/browser-check/SKILL.md → app/internal/handlers/handlers.go
- `newRouteTest()` --calls--> `initRoutes()`  [INFERRED]
  app/cmd/server/routes-household_test.go → app/cmd/server/routes.go
- `GetSessionKeyPairs()` --calls--> `ValidateCookieHashBlockPair()`  [EXTRACTED]
  app/internal/config/globals-env.go → app/internal/validators/validators.go
- `calendarPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-calendar.go → app/internal/handlers/handlers-households.go
- `eventFromContext()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-calendar.go → app/internal/handlers/handlers.go

## Import Cycles
- None detected.

## Communities (217 total, 47 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.07
Nodes (77): TestEmptyEnvIsUnset(), TestSetAppearance(), TestSetLanguage(), TestSafeNextPath(), loggedKinds(), TestActivity(), TestActivityMembers(), TestActivityMoney() (+69 more)

### Community 1 - "handlers.go"
Cohesion: 0.18
Nodes (7): membershipKey, ctxKey, appearanceKey, directionKey, householdsKey, pendingSettlementsKey, unreadChatKey

### Community 2 - "CN"
Cohesion: 0.03
Nodes (94): ActionProps, DescriptionProps, TitleProps, ActionProps, CancelProps, ContentProps, ctxKey, ctxState (+86 more)

### Community 3 - "en"
Cohesion: 0.05
Nodes (39): title, mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain (+31 more)

### Community 4 - "et"
Cohesion: 0.05
Nodes (36): title, back, description, title, et, app, email, email-placeholder (+28 more)

### Community 5 - "github.com/jackc/pgx/v5/pgtype.Date"
Cohesion: 0.03
Nodes (41): copyEvent(), insertEvent(), ClaimChoreNotifiedParams, ClaimChoreNudgedParams, ClaimEventReminderParams, CountUpcomingEventsParams, DeleteDetachedEventsParams, DeleteEventExceptionsParams (+33 more)

### Community 6 - "ListActiveMembersRow"
Cohesion: 0.14
Nodes (29): ListActiveMembersRow, Chore, ChoreDetail, ChoreLogEntry, ChoreMember, Item(), state(), choreAction() (+21 more)

### Community 7 - "vars"
Cohesion: 0.25
Nodes (7): use_gitignore, vars, PORT, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_VERSION, version

### Community 8 - "RecurringRow"
Cohesion: 0.08
Nodes (20): InsertExpenseParams, UpdateExpenseParams, ListRecurringParams, ListRecurringRow, ListRecurringSharesRow, logExpenseEdit(), ExpenseInput, ParseCents() (+12 more)

### Community 9 - "LKMessage"
Cohesion: 0.06
Nodes (24): LKMessage, TestChangePasswordIsValid(), TestChangeUsernameIsValid(), validateCurrentPassword(), EmailForm, ResetPassword, Household, IsModelValid() (+16 more)

### Community 10 - "DefaultService"
Cohesion: 0.10
Nodes (44): addTestSession(), countRows(), isDeleted(), newTestEmail(), setTestPassword(), TestChangeEmail(), TestChangePassword(), TestDeleteAccountBlocked() (+36 more)

### Community 11 - "InternalError"
Cohesion: 0.12
Nodes (11): settingsPath(), householdPath(), sessionUserID(), hxLocation(), InternalError(), IsRequestHTMX(), setLanguageCookie(), reload() (+3 more)

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
Nodes (98): ListHouseholdsOfUserRow, ListBuyAgainRow, timestamptz(), Store, New(), Direction(), WithDirection(), GetPendingSettlements() (+90 more)

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
Cohesion: 0.10
Nodes (15): backTo(), choreFromContext(), choreInput(), choresPath(), expenseFromForm(), expenseInputError(), expensePath(), expensesPath() (+7 more)

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
Cohesion: 0.13
Nodes (4): GetUserRow, TestUserLimit(), newEmailToken(), userLimit

### Community 41 - "recurring"
Cohesion: 0.14
Nodes (14): recurring, delete-confirm, edit-title, month, month-short, never-short, next-due, none (+6 more)

### Community 42 - "select_templ.go"
Cohesion: 0.09
Nodes (28): DirectionAttributes(), initialValue(), RadioGroup(), Content(), Group(), initialValue(), Item(), itemLabel() (+20 more)

### Community 43 - "recurring"
Cohesion: 0.14
Nodes (14): recurring, delete-confirm, edit-title, month, month-short, never-short, next-due, none (+6 more)

### Community 44 - "householdURL"
Cohesion: 0.05
Nodes (102): ListShoppingItemsRow, LK, ActivityRow, Expense, ExpenseFilter, AttachedFile, FileFilter, Membership (+94 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.02
Nodes (99): ListExpenseSharesRow, LockExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, ListMessagesRow, RecurringExpense, GetShoppingListRow, ListShoppingListsRow (+91 more)

### Community 46 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 47 - "calendar/calendar_templ.go"
Cohesion: 0.18
Nodes (18): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+10 more)

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

### Community 52 - "logNote"
Cohesion: 0.23
Nodes (5): ListNotesRow, NoteInput, logNote(), ListNotesParams, NoteForm

### Community 53 - "balanceRows"
Cohesion: 0.16
Nodes (12): ListExportExpensesRow, ListExportSharesRow, balanceRows(), csvAmount(), csvComma(), csvText(), expenseRows(), exportFixture() (+4 more)

### Community 54 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 55 - "generate-locales-keys.go"
Cohesion: 0.22
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

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
Cohesion: 0.07
Nodes (91): choreNextNotice(), EventDetail, Occurrence, HouseholdOverview, PaymentDetails, Ptr(), FormatDate(), FormatSize() (+83 more)

### Community 64 - "migrations.go"
Cohesion: 0.18
Nodes (12): ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations(), SelectMigrationVersion(), TestMigrateToLatestConcurrently(), TestMigrateUpAndDown() (+4 more)

### Community 65 - "files"
Cohesion: 0.05
Nodes (38): files, added-by, all, any-place, any-type, attachments, delete-confirm, delete-keep (+30 more)

### Community 66 - "Membership"
Cohesion: 0.16
Nodes (9): Chore, cleanChore(), firstAssignee(), ChoreInput, ChoreTurn, lockChore(), setChoreRota(), choreMemberPicked() (+1 more)

### Community 67 - "generate-locales-languages.go"
Cohesion: 0.27
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 68 - "popover_templ.go"
Cohesion: 0.16
Nodes (16): Content(), Description(), Header(), state(), Title(), Trigger(), triggerID(), Align (+8 more)

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
Cohesion: 0.20
Nodes (5): ListFilesRow, cleanFileName(), FileMeta, FileTarget, ListFilesParams

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

### Community 82 - ".UpdateNote"
Cohesion: 0.16
Nodes (5): noteFromForm(), noteLengthMessage(), notePath(), notesPath(), Note

### Community 83 - "github.com/a-h/templ.Component"
Cohesion: 0.04
Nodes (105): Props, Variant, readOnly(), HouseholdSettings, PurgeDate(), CN(), flattenClasses(), MergeAttributes() (+97 more)

### Community 84 - "images.go"
Cohesion: 0.12
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

### Community 88 - "noRowsAs"
Cohesion: 0.05
Nodes (27): Queries, GetExpenseRow, GetInviteByCodeRow, expenseOf(), listName(), logActivity(), logShare(), hashToken() (+19 more)

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

### Community 93 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.08
Nodes (9): InsertUserParams, choreFormError(), choreInputFromForm(), notice(), sendCSV(), RenderTempl(), stayInPlace(), isUniqueViolation() (+1 more)

### Community 94 - "Subscription"
Cohesion: 0.09
Nodes (12): Message, Pusher, Subscription, newFakeStore(), newFakePusher(), TestValidPushSubscription(), validPushSubscription(), Keys (+4 more)

### Community 95 - "options"
Cohesion: 0.16
Nodes (16): logging, logging, restart, volumes, driver, options, logging, max-file (+8 more)

### Community 96 - "RenderTemplStatus"
Cohesion: 0.12
Nodes (6): filesPath(), RenderTemplStatus(), notFound(), NotFound(), New(), listForm

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
Nodes (11): cmds, desc, task, tasks, build, templ, vendor, cmds (+3 more)

### Community 113 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 115 - "initRoutes"
Cohesion: 0.09
Nodes (26): accountRoutes(), TestLimitBody(), cacheHeaders(), initRoutes(), limitBody(), registerAssets(), TestAssetCaching(), Hub (+18 more)

### Community 116 - "vendor-assets_test.go"
Cohesion: 0.42
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

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
Cohesion: 0.17
Nodes (8): calendarDayPath(), calendarPath(), eventFormError(), eventFromContext(), eventInputFromForm(), eventOn(), eventPath(), newEventInput()

### Community 126 - "time.Time"
Cohesion: 0.08
Nodes (42): Event, canChangeEvent(), checkOccurrence(), clockOf(), DateOnly(), daysBetween(), dropEdits(), eventRepeats() (+34 more)

### Community 127 - "account"
Cohesion: 0.08
Nodes (25): current-password, delete-blocked-last-owner, delete-blocked-unsettled, delete-confirm, delete-description, delete-submit, delete-title, deleted-description (+17 more)

### Community 131 - "common"
Cohesion: 0.25
Nodes (8): cancel, delete, former-member, link-invalid, request-failed, save, saved, common

### Community 134 - "popover.js"
Cohesion: 0.22
Nodes (29): allContents(), close(), closeAll(), closeNearest(), contentFor(), destroy(), idOf(), init() (+21 more)

### Community 136 - "mail"
Cohesion: 0.12
Nodes (17): mail, account-deleted-body, account-deleted-subject, account-erased-body, account-erased-subject, data-export-body, data-export-subject, email-change-body (+9 more)

### Community 137 - "DefaultService"
Cohesion: 0.31
Nodes (10): calendarEvents(), dailyEvent(), eventInput(), eventOn(), eventTarget(), mustOccurrence(), TestEventCap(), TestEventRights() (+2 more)

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

### Community 146 - "Badge"
Cohesion: 0.48
Nodes (5): Badge(), badgeClasses(), fromEditBadge(), Props, Variant

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "Message"
Cohesion: 0.48
Nodes (7): Message, MessageFile, endsRun(), sameRun(), splitFiles(), startsRun(), TestChatRuns()

### Community 149 - "common"
Cohesion: 0.25
Nodes (8): cancel, delete, former-member, link-invalid, request-failed, save, saved, common

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

### Community 156 - "forgot-password-page"
Cohesion: 0.15
Nodes (12): forgot-password-page, invite-page, description, sent-description, sent-title, submit, title, code (+4 more)

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

### Community 162 - "schedulePicker"
Cohesion: 0.25
Nodes (7): optionSelect(), fileFilterSelect(), scheduleEvery(), schedulePicker(), segmented(), schedulePick, segmentOption

### Community 163 - "time.Time"
Cohesion: 0.09
Nodes (50): CalendarItem, DataExportStatus, Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON() (+42 more)

### Community 164 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 167 - "newRouteTest"
Cohesion: 0.07
Nodes (41): ensureDatabasePool(), main(), newRouteTest(), randomUUID(), send(), TestDeletedHouseholdIsReadOnly(), TestHouseholdRoutesHideFromStrangers(), InitGinEngine() (+33 more)

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

### Community 183 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 184 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 185 - "Redirect"
Cohesion: 0.17
Nodes (5): chatPath(), optionalUUID(), peoplePath(), Redirect(), shoppingListPath()

### Community 186 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 187 - ".writeHouseholdExport"
Cohesion: 0.31
Nodes (6): TestZipName(), uniqueZipName(), zipJSON(), zipName(), SplitFileName(), TestSplitFileName()

### Community 188 - "seed"
Cohesion: 0.50
Nodes (4): cmds, desc, dotenv, seed

### Community 189 - "Debt"
Cohesion: 0.15
Nodes (10): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), ListDebtItemsParams, ListExpensesParams (+2 more)

### Community 190 - "playwright"
Cohesion: 0.50
Nodes (3): npx, playwright, @playwright/mcp

### Community 192 - "payment"
Cohesion: 0.20
Nodes (10): payment, copy-iban, description, iban, iban-invalid, name, name-invalid, show (+2 more)

### Community 193 - "NewHouseholdMember"
Cohesion: 0.50
Nodes (4): NewHouseholdMember(), WithPendingSettlements(), WithUnreadChat(), WithMembership()

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

### Community 202 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 204 - "seed"
Cohesion: 0.90
Nodes (5): check(), must(), seed(), seedExpenses(), seedShopping()

### Community 205 - "jump-nav.js"
Cohesion: 0.70
Nodes (4): line(), queue(), section(), update()

### Community 206 - "lastTopLevelColon"
Cohesion: 0.67
Nodes (3): lastTopLevelColon(), twImportantToV3(), twImportantToV4()

### Community 207 - "FormatIBAN"
Cohesion: 0.67
Nodes (4): FormatIBAN(), NormalizeIBAN(), TestIBAN(), ValidIBAN()

### Community 208 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 209 - "FormatCentsInput"
Cohesion: 0.50
Nodes (4): FormatCentsInput(), TestParseCents(), iPaidAmount(), TestIPaidAmountLeavesOutPending()

### Community 211 - "bubble_templ.go"
Cohesion: 0.26
Nodes (11): ContentClass(), Group(), Reactions(), variantClasses(), Align, ContentProps, GroupProps, Props (+3 more)

### Community 213 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 214 - "dialog_templ.go"
Cohesion: 0.10
Nodes (29): Close(), CloseFor(), Content(), Description(), Dialog(), Footer(), Header(), initialOpen() (+21 more)

### Community 221 - "app-key.sh"
Cohesion: 0.90
Nodes (4): admin(), must(), new_secret(), app-key.sh script

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1303 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **47 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `common`, `avatar`, `confirm-email`, `household`, `invite-page`, `login-page`, `expenses`, `shopping`, `mail`, `recurring`, `notes`, `registration-page`, `verify-email`, `reset-password-page`, `payment`, `files`, `nav`, `home-page`, `appearance`, `data-export`, `chat`, `push`, `export`, `account`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0717154331612163 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `handlers.go`, `CN`, `ListActiveMembersRow`, `RecurringRow`, `InternalError`, `context.Context`, `membershipFromContext`, `schedulePicker`, `time.Time`, `newRouteTest`, `newEmailToken`, `householdURL`, `balanceRows`, `Redirect`, `Membership`, `.UpdateNote`, `github.com/a-h/templ.Component`, `noRowsAs`, `github.com/gin-gonic/gin.Context`, `RenderTemplStatus`, `Membership`, `.loadShoppingItem`, `.renderEvent`, `time.Time`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `CN` be split into smaller, more focused modules?**
  _Cohesion score 0.032226752086973405 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `mail`, `reset-password-page`, `expenses`, `household`, `common`, `forgot-password-page`, `registration-page`, `shopping`, `recurring`, `notes`, `files`, `payment`, `login-page`, `verify-email`, `home-page`, `nav`, `password`, `data-export`, `export`, `avatar`, `chat`, `push`, `account`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.05 - nodes in this community are weakly interconnected._