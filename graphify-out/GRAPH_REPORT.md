# Graph Report - roommates  (2026-11-19)

## Corpus Check
- 292 files · ~336,672 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 63 file(s) not represented in the graph (top: .templ 43, (none) 9, .css 5)

## Summary
- 3665 nodes · 12645 edges · 213 communities (165 shown, 48 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 935 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `30c2eb1d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- testing.T
- handlers.go
- migrations.go
- en
- et
- householdURL
- Props
- vars
- ExpenseInput
- LKMessage
- InitGinEngine
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
- validation-errors
- generate-locales-languages.go
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
- balanceRows
- notes
- github.com/a-h/templ.Attributes
- verify-email
- verify-email
- generate-locales-keys.go
- reset-password-page
- use_hover.js
- login-page
- FormatMoney
- T
- navIcon
- files
- io.ReadSeeker
- DefaultService
- go_pkg_fmt
- calendar.js
- environment
- nav
- nav
- SplitFileName
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
- invite-page
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
- Queries
- radiogroup.js
- portal.js
- useFocus
- push
- push
- chat-keyboard.js
- mailer.go
- tasks
- export
- renderer.go
- initRoutes
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
- dialog_templ.go
- avatar.js
- reset-password-page
- avatar
- tailwind
- noRowsAs
- settlementHelpers
- floating_focus_manager.js
- chatMessages
- common
- useListNavigation
- server
- restart-signal
- dev
- Browser check
- backup.sh
- common
- Variants
- ListNotesRow
- expense-split.js
- checkbox.js
- .UpdateNote
- select_templ.go
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
- error-page
- appearance
- Redirect
- tailwind-version
- vendor
- seed
- alertdialog_templ.go
- playwright
- payment
- ListBuyAgainRow
- deleted-with-keep
- deleted-with
- deleted-with-keep
- deleted-with
- restart-watch
- field_templ.go
- New
- popover_templ.go
- FieldError
- .handleUpload
- confirm-email
- registration-page
- .GetPaymentDetails
- password
- password
- FormatIBAN
- pushPreview

## God Nodes (most connected - your core abstractions)
1. `T()` - 112 edges
2. `CN()` - 103 edges
3. `expenses` - 76 edges
4. `expenses` - 76 edges
5. `T()` - 71 edges
6. `Membership` - 67 edges
7. `membershipFromContext()` - 66 edges
8. `newTestService()` - 58 edges
9. `InternalError()` - 55 edges
10. `RenderTemplStatus()` - 51 edges

## Surprising Connections (you probably didn't know these)
- `HTMX actions` --references--> `Redirect()`  [INFERRED]
  .claude/skills/browser-check/SKILL.md → app/internal/handlers/handlers.go
- `main()` --calls--> `InitLogger()`  [EXTRACTED]
  app/cmd/server/main.go → app/internal/logger/logger.go
- `TestHouseholdRoutesHideFromStrangers()` --calls--> `initRoutes()`  [INFERRED]
  app/cmd/server/routes-household_test.go → app/cmd/server/routes.go
- `GetSessionKeyPairs()` --calls--> `ValidateCookieHashBlockPair()`  [EXTRACTED]
  app/internal/config/globals-env.go → app/internal/validators/validators.go
- `chatPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-chat.go → app/internal/handlers/handlers-households.go

## Import Cycles
- None detected.

## Communities (213 total, 48 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.10
Nodes (67): TestSetAppearance(), TestSetLanguage(), TestSafeNextPath(), newTestEmail(), setTestPassword(), TestChangeEmail(), TestChangePassword(), TestEmailChangeTakenAtConfirm() (+59 more)

### Community 1 - "handlers.go"
Cohesion: 0.15
Nodes (12): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), ExpenseCategory, ctxKey, ctxKey, Props, appearanceKey (+4 more)

### Community 2 - "migrations.go"
Cohesion: 0.11
Nodes (23): ensureDatabasePool(), randomUUID(), TestHouseholdRoutesHideFromStrangers(), GetDatabaseURL(), TestGetDatabaseURLEscapes(), New(), ensureMigrationRow(), ensureMigrationTable() (+15 more)

### Community 3 - "en"
Cohesion: 0.05
Nodes (39): title, en, app, email, email-placeholder, error-page, home-page, invite-page (+31 more)

### Community 4 - "et"
Cohesion: 0.06
Nodes (33): title, done-description, done-title, invalid-description, taken-description, taken-title, et, app (+25 more)

### Community 5 - "householdURL"
Cohesion: 0.16
Nodes (25): Expense, AttachedFile, localTime(), ExpensePage(), sharePaidStamp(), attachedFiles(), attachedFilesContent(), attachedFilesTill() (+17 more)

### Community 6 - "Props"
Cohesion: 0.31
Nodes (10): Content(), ContentClass(), variantClasses(), Align, ContentProps, GroupProps, Props, ReactionsProps (+2 more)

### Community 7 - "vars"
Cohesion: 0.20
Nodes (9): use_gitignore, vars, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT, TAILWIND_OUTPUT, TAILWIND_VERSION (+1 more)

### Community 8 - "ExpenseInput"
Cohesion: 0.11
Nodes (14): GetExpenseRow, InsertExpenseParams, ExpenseInput, ParseCents(), splitEqual(), TestCopyInputFormOrder(), TestExpenseCleanField(), TestNetDebts() (+6 more)

### Community 9 - "LKMessage"
Cohesion: 0.06
Nodes (24): LKMessage, TestChangePasswordIsValid(), TestChangeUsernameIsValid(), validateCurrentPassword(), EmailForm, ResetPassword, Household, IsModelValid() (+16 more)

### Community 10 - "InitGinEngine"
Cohesion: 0.10
Nodes (22): main(), send(), InitGinEngine(), GetAppBaseURL(), GetEnv(), GetSessionKeyPairs(), IsDevelopment(), MustGetEnv() (+14 more)

### Community 11 - "InternalError"
Cohesion: 0.21
Nodes (5): settingsPath(), householdPath(), sessionUserID(), InternalError(), FromGinContext()

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
Nodes (71): ListHouseholdsOfUserRow, GetUserCredentialsRow, InsertUserParams, Direction(), WithDirection(), GetPendingSettlements(), GetUnreadChat(), WithPendingSettlements() (+63 more)

### Community 16 - "expenses"
Cohesion: 0.03
Nodes (76): expenses, add, amount, amount-locked-currency, amount-locked-left, back, balances, cancel (+68 more)

### Community 17 - "household"
Cohesion: 0.05
Nodes (37): household, back, balance-not-zero, copied, copy-code, copy-link, create-invite, expires (+29 more)

### Community 18 - "household"
Cohesion: 0.05
Nodes (37): household, back, balance-not-zero, copied, copy-code, copy-link, create-invite, expires (+29 more)

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
Cohesion: 0.13
Nodes (7): expensesPath(), RenderTemplStatus(), NewHouseholdMember(), notFound(), ParseUUID(), NotFound(), listForm

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
Cohesion: 0.08
Nodes (9): GetInviteByCodeRow, TestUserLimit(), newEmailToken(), timestamptz(), newInviteCode(), NormalizeInviteCode(), TestNewInviteCode(), TestNormalizeInviteCode() (+1 more)

### Community 41 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 42 - "generate-locales-languages.go"
Cohesion: 0.24
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 43 - "home-page"
Cohesion: 0.25
Nodes (8): home-page, create-submit, create-title, greeting, households, join-submit, join-title, no-households

### Community 44 - "Membership"
Cohesion: 0.12
Nodes (47): ListShoppingItemsRow, ExpenseFilter, FileFilter, Membership, ShoppingList, IfElse(), T(), ChatFilePicks() (+39 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.03
Nodes (64): ListExpenseSharesRow, LockExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, ListMessagesRow, GetShoppingListRow, ListShoppingListsRow, avatarURL() (+56 more)

### Community 46 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 47 - "vendor-assets.go"
Cohesion: 0.18
Nodes (18): download(), fileSHA256(), latestVersion(), load(), main(), pinLatest(), readJSON(), runCheck() (+10 more)

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

### Community 52 - "validation-errors"
Cohesion: 0.12
Nodes (16): generic, validation-errors, name-length, case, length, match, max-length, required (+8 more)

### Community 53 - "balanceRows"
Cohesion: 0.12
Nodes (14): ListExportExpensesRow, ListExportSharesRow, ListActiveMembersRow, balanceRows(), csvAmount(), csvComma(), csvText(), expenseRows() (+6 more)

### Community 54 - "notes"
Cohesion: 0.05
Nodes (38): notes, back, body, body-hint, bold, cancel, delete, delete-confirm (+30 more)

### Community 55 - "github.com/a-h/templ.Attributes"
Cohesion: 0.05
Nodes (48): ActionProps, DescriptionProps, Props, TitleProps, Variant, variantClasses(), Size, initialChecked() (+40 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "generate-locales-keys.go"
Cohesion: 0.27
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 59 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 60 - "use_hover.js"
Cohesion: 0.09
Nodes (41): applySafePolygonPointerEventsMutation(), clearSafePolygonPointerEventsMutation(), createHoverInteraction(), emitClosed(), getDelay(), getRestMs(), hasIntersectingEdge(), isClickLikeOpenEvent() (+33 more)

### Community 61 - "login-page"
Cohesion: 0.29
Nodes (7): login-page, create-account, enter-credentials-description, forgot-password, invalid-credentials, no-account, sign-in

### Community 62 - "FormatMoney"
Cohesion: 0.17
Nodes (11): householdPath(), FormatDate(), FormatSize(), printer(), TestFormatPerLanguage(), GetLocaleString(), FormatMoney(), HtmlWrap() (+3 more)

### Community 63 - "T"
Cohesion: 0.10
Nodes (56): GetUserRow, HouseholdOverview, HouseholdSettings, PaymentDetails, Alert(), Button(), T(), ValidationErrors() (+48 more)

### Community 64 - "navIcon"
Cohesion: 0.10
Nodes (32): LK, GetHouseholds(), N(), unescapeJSON(), Badge(), badgeClasses(), N(), deleteKeepConfirm() (+24 more)

### Community 65 - "files"
Cohesion: 0.06
Nodes (36): files, added-by, all, any-place, any-type, attachments, delete-confirm, delete-keep (+28 more)

### Community 66 - "io.ReadSeeker"
Cohesion: 0.08
Nodes (12): ListFilesRow, fit(), Thumbnail(), cleanFileName(), FileMeta, FileTarget, newFakeStore(), TestCleanFileName() (+4 more)

### Community 67 - "DefaultService"
Cohesion: 0.16
Nodes (15): Pusher, addTestSession(), DefaultService, lockHouseholdInTx(), subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes() (+7 more)

### Community 69 - "calendar.js"
Cohesion: 0.09
Nodes (44): addToRange(), build(), calculateFocusTarget(), CalendarDay, calendarFns(), addMonths(), createDateLib(), getISOWeek() (+36 more)

### Community 70 - "environment"
Cohesion: 0.29
Nodes (7): RUSTFS_ACCESS_KEY, RUSTFS_ADDRESS, RUSTFS_CONSOLE_ADDRESS, RUSTFS_CONSOLE_ENABLE, RUSTFS_SECRET_KEY, RUSTFS_VOLUMES, environment

### Community 71 - "nav"
Cohesion: 0.15
Nodes (13): nav, chat, expenses, files, households-add, label, more, more-waiting (+5 more)

### Community 72 - "nav"
Cohesion: 0.15
Nodes (13): nav, chat, expenses, files, households-add, label, more, more-waiting (+5 more)

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
Cohesion: 0.42
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 83 - "github.com/a-h/templ.Component"
Cohesion: 0.04
Nodes (98): CN(), DirectionAttributes(), flattenClasses(), Ptr(), RandomID(), twMapTokens(), twReorder(), Action() (+90 more)

### Community 84 - "images.go"
Cohesion: 0.13
Nodes (10): exifOrientation(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient(), TestThumbnail() (+2 more)

### Community 85 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 88 - "Debt"
Cohesion: 0.15
Nodes (10): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), ListDebtItemsParams, ListExpensesParams (+2 more)

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

### Community 93 - "Hub"
Cohesion: 0.23
Nodes (5): Hub, NewHub(), TestHub(), newFakePusher(), fakePusher

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (7): Message, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

### Community 95 - "options"
Cohesion: 0.20
Nodes (14): logging, logging, logging, logging, driver, options, max-file, max-size (+6 more)

### Community 96 - "Props"
Cohesion: 0.29
Nodes (8): ActionProps, CancelProps, Size, Type, Variant, inputGroupButtonClasses(), ButtonProps, ButtonSize

### Community 97 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 98 - "chat"
Cohesion: 0.08
Nodes (24): attach, cancel-reply, copy, delete, delete-confirm, empty, file-gone, files-too-large (+16 more)

### Community 99 - "use_transition_status.js"
Cohesion: 0.40
Nodes (7): animationsFinished(), close(), open(), reset(), resolve(), setOpen(), whenAnimationsFinish()

### Community 100 - ".notifyChat"
Cohesion: 0.29
Nodes (3): Chat, ReplyPreview, Upload

### Community 101 - "cloudflared"
Cohesion: 0.22
Nodes (9): condition, command, depends_on, environment, profiles, restart, app, TUNNEL_TOKEN (+1 more)

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - "Queries"
Cohesion: 0.11
Nodes (10): Queries, checkSettled(), closeEvenDebts(), closePaidDebt(), editShares(), payIfCovered(), householdSettled(), payShare() (+2 more)

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
Cohesion: 0.16
Nodes (9): buildMessage(), Mailer, Message, isLocalHost(), TestBuildMessage(), TestIsLocalHost(), LogMailer, SMTPMailer (+1 more)

### Community 112 - "tasks"
Cohesion: 0.18
Nodes (11): cmds, desc, task, cmds, desc, tasks, build, scripts-watch (+3 more)

### Community 113 - "export"
Cohesion: 0.11
Nodes (19): export, balance, balances, category, counterpart, currency, date, expenses-long (+11 more)

### Community 114 - "renderer.go"
Cohesion: 0.24
Nodes (3): New(), HTMLTemplRenderer, Renderer

### Community 115 - "initRoutes"
Cohesion: 0.10
Nodes (25): accountRoutes(), TestLimitBody(), cacheHeaders(), initRoutes(), limitBody(), registerAssets(), TestAssetCaching(), DefaultHandler (+17 more)

### Community 116 - "membershipFromContext"
Cohesion: 0.18
Nodes (8): expenseFromForm(), expenseInputError(), expensePath(), shoppingItemPath(), membershipFromContext(), pageError(), stay(), ShoppingItem

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
Cohesion: 0.17
Nodes (6): backTo(), setLanguageCookie(), safeNextPath(), reload(), Language, FileTargetLoader

### Community 126 - "calendar_templ.go"
Cohesion: 0.23
Nodes (17): Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON(), toJSON(), CaptionLayout (+9 more)

### Community 127 - "account"
Cohesion: 0.13
Nodes (15): current-password, email-description, email-same, email-sent, email-submit, new-email, new-password, password-changed (+7 more)

### Community 128 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.12
Nodes (5): notice(), sendCSV(), RenderTempl(), stayInPlace(), NormalizeEmail()

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

### Community 137 - "dialog_templ.go"
Cohesion: 0.16
Nodes (16): Close(), CloseFor(), initialOpen(), state(), Trigger(), TriggerFor(), ContentProps, ctxKey (+8 more)

### Community 139 - "reset-password-page"
Cohesion: 0.29
Nodes (7): reset-password-page, done-description, done-title, invalid-description, new-link, submit, title

### Community 141 - "avatar"
Cohesion: 0.18
Nodes (11): cancel, choose, crop-hint, description, invalid, remove, remove-confirm, title (+3 more)

### Community 142 - "tailwind"
Cohesion: 0.33
Nodes (7): preconditions, msg, sh, cmds, desc, preconditions, tailwind

### Community 143 - "noRowsAs"
Cohesion: 0.09
Nodes (7): isUniqueViolation(), hashToken(), noRowsAs(), NoteInput, rowsAffected(), trimToLength(), NoteForm

### Community 146 - "settlementHelpers"
Cohesion: 0.21
Nodes (11): editHelpers(), TestEditExpenseMovesBalanceByTheChange(), TestEditExpenseOverpaidInApp(), TestEditExpenseOverpaidOutside(), TestEditExpenseParticipants(), TestEditExpenseUnderpaid(), FormatCentsInput(), settlementHelpers() (+3 more)

### Community 147 - "floating_focus_manager.js"
Cohesion: 0.23
Nodes (17): addPreviouslyFocusedElement(), clearDisconnectedPreviouslyFocusedElements(), createFocusGuard(), enqueueFocus(), getEventType(), getFirstTabbableElement(), getFloatingFocusElement(), getPreviouslyFocusedElement() (+9 more)

### Community 148 - "chatMessages"
Cohesion: 0.23
Nodes (15): Message, MessageFile, MergeAttributes(), chatControlLabel(), chatControlTip(), chatFilePicker(), chatMessages(), ChatPage() (+7 more)

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

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 162 - ".UpdateNote"
Cohesion: 0.21
Nodes (5): noteFromForm(), noteLengthMessage(), notePath(), notesPath(), Note

### Community 163 - "select_templ.go"
Cohesion: 0.14
Nodes (14): initialValue(), Align, ContentProps, ctxKey, ctxState, GroupProps, ItemData, ItemProps (+6 more)

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

### Community 183 - "error-page"
Cohesion: 0.50
Nodes (4): back, description, title, error-page

### Community 184 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 185 - "Redirect"
Cohesion: 0.17
Nodes (8): chatPath(), optionalUUID(), peoplePath(), hxLocation(), IsRequestHTMX(), Redirect(), shoppingListPath(), swapMainOnError()

### Community 186 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 187 - "vendor"
Cohesion: 0.67
Nodes (3): vendor, cmds, desc

### Community 188 - "seed"
Cohesion: 0.50
Nodes (4): cmds, desc, dotenv, seed

### Community 189 - "alertdialog_templ.go"
Cohesion: 0.15
Nodes (18): ContentProps, ctxKey, ctxState, DescriptionProps, FooterProps, HeaderProps, MediaProps, OverlayProps (+10 more)

### Community 190 - "playwright"
Cohesion: 0.50
Nodes (3): npx, playwright, @playwright/mcp

### Community 192 - "payment"
Cohesion: 0.20
Nodes (10): payment, copy-iban, description, iban, iban-invalid, name, name-invalid, show (+2 more)

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

### Community 199 - "field_templ.go"
Cohesion: 0.16
Nodes (15): childrenHTML(), fieldClasses(), Separator(), ContentProps, DescriptionProps, ErrorProps, GroupProps, LabelProps (+7 more)

### Community 202 - "popover_templ.go"
Cohesion: 0.22
Nodes (12): state(), Trigger(), triggerID(), Align, ContentProps, ctxKey, ctxState, DescriptionProps (+4 more)

### Community 205 - "confirm-email"
Cohesion: 0.33
Nodes (6): done-description, done-title, invalid-description, taken-description, taken-title, confirm-email

### Community 206 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 208 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 209 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 210 - "FormatIBAN"
Cohesion: 0.60
Nodes (4): FormatIBAN(), NormalizeIBAN(), TestIBAN(), ValidIBAN()

## Knowledge Gaps
- **3 isolated node(s):** `npx`, `@playwright/mcp`, `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1216 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **48 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `et` connect `et` to `avatar`, `household`, `common`, `expenses`, `shopping`, `validation-errors`, `mail`, `home-page`, `notes`, `error-page`, `appearance`, `verify-email`, `reset-password-page`, `payment`, `files`, `nav`, `login-page`, `registration-page`, `password`, `invite-page`, `chat`, `push`, `export`, `account`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `npx`, `@playwright/mcp`, `roommates` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.10017050298380221 - nodes in this community are weakly interconnected._
- **Why does `en` connect `en` to `forgot-password-page`, `mail`, `reset-password-page`, `expenses`, `household`, `common`, `appearance`, `shopping`, `registration-page`, `notes`, `files`, `payment`, `validation-errors`, `verify-email`, `login-page`, `nav`, `confirm-email`, `password`, `export`, `avatar`, `chat`, `push`, `account`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `migrations.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10952380952380952 - nodes in this community are weakly interconnected._
- **Why does `Redirect()` connect `Redirect` to `github.com/gin-gonic/gin.Context`, `handlers.go`, `.UpdateNote`, `InternalError`, `membershipFromContext`, `RenderTemplStatus`, `Browser check`, `.changeFileRef`?**
  _High betweenness centrality (0.006) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.05 - nodes in this community are weakly interconnected._