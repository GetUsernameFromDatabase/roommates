# Graph Report - roommates  (2026-10-31)

## Corpus Check
- 286 files · ~298,473 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 55 file(s) not represented in the graph (top: .templ 38, (none) 6, .css 5)

## Summary
- 3064 nodes · 9728 edges · 191 communities (151 shown, 40 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 723 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0f1c0542`
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
- Redirect
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
- categorySelect
- select.js
- floating_ui_dom.js
- TestPushNotifications
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
- models.go
- generate-locales-keys.go
- newTestUser
- files
- householdURL
- generate-locales-languages.go
- initRoutes
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
- card_templ.go
- New
- vendor-assets_test.go
- Props
- switch.js
- mounted
- scroll_lock.js
- NewPostgreSessionStoreHandler
- io.ReadSeeker
- calendar_templ.go
- T
- github.com/gin-gonic/gin.Context
- ListNotesRow
- go_pkg_fmt
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
- select_templ.go
- useListNavigation
- server
- restart-signal
- dev
- error-page
- vendor
- tailwind-version
- scripts-watch
- forgot-password-page
- expense-split.js
- checkbox.js
- HTMLTemplRenderer
- popover_templ.go
- appearance
- useClick
- dialog_templ.go
- registration-page
- mark_others.js
- Props
- useTypeahead
- DefaultService
- internal_backdrop.js
- .loadShoppingItem
- invite-page
- icon.go
- .notifyChat
- password
- appearance
- field_templ.go
- Props
- password
- alertdialog_templ.go

## God Nodes (most connected - your core abstractions)
1. `CN()` - 103 edges
2. `T()` - 93 edges
3. `T()` - 61 edges
4. `Membership` - 59 edges
5. `membershipFromContext()` - 57 edges
6. `expenses` - 53 edges
7. `expenses` - 53 edges
8. `RenderTemplStatus()` - 47 edges
9. `Button()` - 43 edges
10. `InternalError()` - 41 edges

## Surprising Connections (you probably didn't know these)
- `notice()` --calls--> `RenderTemplStatus()`  [INFERRED]
  app/internal/handlers/handlers-email.go → app/internal/handlers/handlers.go
- `settingsPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-home.go → app/internal/handlers/handlers-households.go
- `shoppingListPath()` --calls--> `householdPath()`  [INFERRED]
  app/internal/handlers/handlers-shopping.go → app/internal/handlers/handlers-households.go
- `TestSafeNextPath()` --calls--> `safeNextPath()`  [INFERRED]
  app/internal/handlers/handlers-next_test.go → app/internal/handlers/handlers-next.go
- `pageError()` --calls--> `InternalError()`  [INFERRED]
  app/internal/handlers/handlers-notes.go → app/internal/handlers/handlers.go

## Import Cycles
- None detected.

## Communities (191 total, 40 thin omitted)

### Community 0 - "testing.T"
Cohesion: 0.13
Nodes (39): TestSetAppearance(), TestSafeNextPath(), TestAvatars(), TestChat(), TestChatFileRefs(), TestChatReply(), TestHouseholdEvents(), settlementHelpers() (+31 more)

### Community 1 - "handlers.go"
Cohesion: 0.13
Nodes (11): lastTopLevelColon(), twImportantToV3(), twImportantToV4(), DayOfWeek, ctxKey, ctxKey, Props, appearanceKey (+3 more)

### Community 2 - "migrations.go"
Cohesion: 0.15
Nodes (14): ensureDatabasePool(), GetDatabaseURL(), TestGetDatabaseURLEscapes(), ensureMigrationRow(), ensureMigrationTable(), getMigrationVersion(), MigrateToLatest(), NewMigrations() (+6 more)

### Community 3 - "en"
Cohesion: 0.07
Nodes (27): title, cancel, delete, link-invalid, save, en, app, common (+19 more)

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
Cohesion: 0.17
Nodes (11): sh, use_gitignore, vars, FREE_PORT, PORT, PROXY_TLS_CRT, PROXY_TLS_KEY, TAILWIND_INPUT (+3 more)

### Community 8 - "householdURL"
Cohesion: 0.13
Nodes (29): FileFilter, HouseholdOverview, localTime(), fileFilters(), HouseholdFiles(), HouseholdShopping(), appShell(), ChatBadge() (+21 more)

### Community 9 - "LKMessage"
Cohesion: 0.10
Nodes (16): LKMessage, EmailForm, ResetPassword, Household, IsModelValid(), Login, validateEmail(), Register (+8 more)

### Community 10 - "Redirect"
Cohesion: 0.22
Nodes (5): peoplePath(), IsRequestHTMX(), safeNextPath(), Redirect(), swapMainOnError()

### Community 11 - "membershipFromContext"
Cohesion: 0.20
Nodes (7): optionalUUID(), membershipFromContext(), pageError(), RenderTemplStatus(), Note, ParseUUID(), NotFound()

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
Nodes (47): ListHouseholdsOfUserRow, Direction(), WithDirection(), GetUnreadChat(), WithUnreadChat(), GetHouseholds(), WithHouseholds(), state() (+39 more)

### Community 16 - "expenses"
Cohesion: 0.04
Nodes (53): expenses, add, amount, back, balances, cancel, category, category-groceries (+45 more)

### Community 17 - "household"
Cohesion: 0.06
Nodes (31): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+23 more)

### Community 18 - "household"
Cohesion: 0.06
Nodes (31): one, other, household, back, balance-not-zero, copied, copy-code, copy-link (+23 more)

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
Cohesion: 0.13
Nodes (36): GetExpenseRow, ListShoppingItemsRow, ExpenseFilter, Membership, ShoppingList, IfElse(), FormatMoney(), ChatFilePicks() (+28 more)

### Community 45 - "github.com/jackc/pgx/v5/pgtype.UUID"
Cohesion: 0.03
Nodes (41): ListExpenseSharesRow, GetHouseholdRow, ListActiveInvitesRow, GetShoppingListRow, ListBuyAgainRow, ListShoppingListsRow, avatarURL(), ConfirmSettlementParams (+33 more)

### Community 46 - "notes"
Cohesion: 0.06
Nodes (35): notes, back, body, body-hint, bold, cancel, delete, edit (+27 more)

### Community 47 - "categorySelect"
Cohesion: 0.40
Nodes (4): categorySelect(), optionSelect(), fileFilterSelect(), CategoryName()

### Community 48 - "select.js"
Cohesion: 0.10
Nodes (53): allContents(), armCancelOpen(), arrowScrollStep(), clamp(), close(), closeAll(), commitValue(), contentFor() (+45 more)

### Community 49 - "floating_ui_dom.js"
Cohesion: 0.13
Nodes (51): autoUpdate(), frameLoop(), convertOffsetParentRelativeRectToViewportRelativeRect(), getBoundingClientRect(), getClientRectFromClippingAncestor(), getClippingElementAncestors(), getClippingRect(), getComputedStyle$1() (+43 more)

### Community 50 - "TestPushNotifications"
Cohesion: 0.25
Nodes (9): LoadLocales(), subscribeTestDevice(), takePushes(), TestPushNotifications(), wantPushes(), TestTUnescapesTextButNotArgs(), TestExpenseEditorCarriesTheItemLink(), TestHouseholdExpensesSettledStamp() (+1 more)

### Community 51 - "files"
Cohesion: 0.05
Nodes (42): one, other, one, other, files, added-by, all, any-place (+34 more)

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
Cohesion: 0.07
Nodes (34): ActionProps, DescriptionProps, TitleProps, Size, initialChecked(), Props, initialValue(), DisableTrigger() (+26 more)

### Community 56 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 57 - "verify-email"
Cohesion: 0.22
Nodes (9): verify-email, banner, continue, done-description, done-title, invalid-description, resend, resent-description (+1 more)

### Community 58 - "vendor-assets.go"
Cohesion: 0.18
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

### Community 62 - "models.go"
Cohesion: 0.08
Nodes (29): ListMessagesRow, EmailToken, Expense, ExpenseCategory, ExpenseShare, File, FilePurge, FileRef (+21 more)

### Community 63 - "generate-locales-keys.go"
Cohesion: 0.20
Nodes (7): flattenYAML(), generateOutputContent(), main(), parseLocaleYml(), pathToConst(), toConstName(), writeOutput()

### Community 64 - "newTestUser"
Cohesion: 0.50
Nodes (9): englishContext(), TestExpiredTokenIsInvalid(), TestInviteNeedsConfirmedEmail(), TestResetPassword(), TestVerifyEmail(), tokenOf(), newTestUser(), nextEmail() (+1 more)

### Community 65 - "files"
Cohesion: 0.05
Nodes (42): one, other, one, other, files, added-by, all, any-place (+34 more)

### Community 66 - "householdURL"
Cohesion: 0.13
Nodes (27): LK, Expense, AttachedFile, N(), unescapeJSON(), N(), ExpensePage(), attachedFiles() (+19 more)

### Community 67 - "generate-locales-languages.go"
Cohesion: 0.27
Nodes (7): discoverLanguages(), generateLocales(), loadYAML(), main(), mapping(), missingKeys(), unwrapLocale()

### Community 68 - "initRoutes"
Cohesion: 0.08
Nodes (24): initRoutes(), DeleteCookie(), GetCookie(), GetHeader(), NewAppearance(), NewGuestOnly(), NewOnlyAuthenticated(), NewHouseholdMember() (+16 more)

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
Cohesion: 0.17
Nodes (5): ListFilesRow, FileMeta, SplitFileName(), TestSplitFileName(), ListFilesParams

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
Cohesion: 0.05
Nodes (85): CN(), DirectionAttributes(), flattenClasses(), Ptr(), RandomID(), twMapTokens(), twReorder(), Action() (+77 more)

### Community 84 - "images.go"
Cohesion: 0.12
Nodes (10): exifOrientation(), fit(), jpegOrientation(), orient(), encodePNG(), jpegWithOrientation(), testImage(), TestOrient() (+2 more)

### Community 85 - "Props"
Cohesion: 0.60
Nodes (4): Badge(), badgeClasses(), Props, Variant

### Community 86 - "floating_ui_core.js"
Cohesion: 0.13
Nodes (37): arrowMiddleware(), getLogicalSide(), hideMiddleware(), isRtlAt(), roundByDPR(), useAnchorPositioning(), getOffsetData(), middleware() (+29 more)

### Community 87 - "chatMessages"
Cohesion: 0.23
Nodes (16): Chat, Message, MessageFile, MergeAttributes(), chatControlLabel(), chatControlTip(), chatFilePicker(), chatMessages() (+8 more)

### Community 88 - "Debt"
Cohesion: 0.13
Nodes (12): ListDebtItemsRow, ListExpensesRow, ListPendingSettlementsRow, Debt, ExpensesOverview, netDebts(), TestNetDebts(), testUUID() (+4 more)

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

### Community 93 - "Hub"
Cohesion: 0.23
Nodes (5): Hub, NewHub(), TestHub(), newFakePusher(), fakePusher

### Community 94 - "Subscription"
Cohesion: 0.15
Nodes (7): Message, Subscription, TestValidPushSubscription(), validPushSubscription(), Keys, LogPusher, WebPusher

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
Cohesion: 0.16
Nodes (9): InsertExpenseParams, ListActiveMembersRow, FormatCentsInput(), ExpenseInput, ParseCents(), splitEqual(), TestParseCents(), TestSplitEqual() (+1 more)

### Community 101 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 102 - "composite.js"
Cohesion: 0.18
Nodes (18): findNonDisabledListIndex(), getMaxListIndex(), getMinListIndex(), getOffset(), getStyles(), isElementDisabled(), isListIndexDisabled(), isModifierKeySet() (+10 more)

### Community 103 - ".handleUpload"
Cohesion: 0.14
Nodes (10): chatPath(), expensePath(), expensesPath(), filesPath(), householdPath(), noteFromForm(), noteLengthMessage(), notePath() (+2 more)

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
Nodes (6): buildMessage(), Message, TestBuildMessage(), LogMailer, SMTPMailer, fakeMailer

### Community 112 - "tasks"
Cohesion: 0.18
Nodes (11): cmds, desc, task, cmds, desc, tasks, build, restart-watch (+3 more)

### Community 113 - "card_templ.go"
Cohesion: 0.11
Nodes (24): GetAppearance(), GetLocaleString(), Card(), Content(), Description(), Header(), Title(), HeaderComponent() (+16 more)

### Community 114 - "New"
Cohesion: 0.47
Nodes (3): notFound(), New(), Renderer

### Community 115 - "vendor-assets_test.go"
Cohesion: 0.42
Nodes (6): runSync(), fakeNPM(), readFile(), TestCheckAndUpdate(), TestSync(), writeManifest()

### Community 116 - "Props"
Cohesion: 0.24
Nodes (10): ActionProps, CancelProps, Size, Type, Variant, Variants(), inputGroupButtonClasses(), Props (+2 more)

### Community 117 - "switch.js"
Cohesion: 0.31
Nodes (5): forwardClick(), inputOf(), requestCheckedChange(), setup(), sync()

### Community 120 - "mounted"
Cohesion: 0.83
Nodes (3): mounted(), register(), sync()

### Community 123 - "scroll_lock.js"
Cohesion: 0.11
Nodes (15): anchoredPopupScrollLock(), AnimationFrame, getViewportScroller(), hasInsetScrollbars(), isOverflowElement(), isPageScrollLocked(), preventScrollInsetScrollbars(), cleanup() (+7 more)

### Community 124 - "NewPostgreSessionStoreHandler"
Cohesion: 0.19
Nodes (8): GetEnv(), GetSessionKeyPairs(), MustGetEnv(), NewPostgreSessionStoreHandler(), startExpiredSessionCleanup(), ValidateCookieHashBlockPair(), EnvironmentVariableKey, sessionCleaner

### Community 125 - "io.ReadSeeker"
Cohesion: 0.10
Nodes (8): Thumbnail(), cleanFileName(), newFakeStore(), TestCleanFileName(), Store, New(), fakeStore, S3Store

### Community 126 - "calendar_templ.go"
Cohesion: 0.16
Nodes (16): Calendar(), dateAttr(), isoDate(), matcherJSON(), modifiersJSON(), selectedJSON(), toJSON(), CaptionLayout (+8 more)

### Community 127 - "T"
Cohesion: 0.18
Nodes (37): GetUserRow, T(), Alert(), Button(), T(), ValidationErrors(), Field(), Group() (+29 more)

### Community 128 - "github.com/gin-gonic/gin.Context"
Cohesion: 0.11
Nodes (7): notice(), expenseFromForm(), settingsPath(), sessionUserID(), InternalError(), RenderTempl(), stayInPlace()

### Community 129 - "ListNotesRow"
Cohesion: 0.19
Nodes (4): ListNotesRow, NoteInput, ListNotesParams, NoteForm

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
Nodes (7): GetInviteByCodeRow, hashToken(), closePaidDebt(), payShare(), noRowsAs(), rowsAffected(), trimToLength()

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
Cohesion: 0.17
Nodes (14): initialValue(), Align, ContentProps, ctxKey, ctxState, GroupProps, ItemData, ItemProps (+6 more)

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

### Community 154 - "error-page"
Cohesion: 0.50
Nodes (4): back, description, title, error-page

### Community 155 - "vendor"
Cohesion: 0.67
Nodes (3): vendor, cmds, desc

### Community 156 - "tailwind-version"
Cohesion: 0.50
Nodes (4): cmds, desc, silent, tailwind-version

### Community 157 - "scripts-watch"
Cohesion: 0.67
Nodes (3): cmds, desc, scripts-watch

### Community 158 - "forgot-password-page"
Cohesion: 0.33
Nodes (6): forgot-password-page, description, sent-description, sent-title, submit, title

### Community 159 - "expense-split.js"
Cohesion: 0.83
Nodes (3): parseCents(), splitEqual(), update()

### Community 161 - "checkbox.js"
Cohesion: 0.31
Nodes (6): forwardClick(), inputOf(), requestCheckedChange(), setChecked(), setup(), sync()

### Community 163 - "popover_templ.go"
Cohesion: 0.22
Nodes (12): state(), Trigger(), triggerID(), Align, ContentProps, ctxKey, ctxState, DescriptionProps (+4 more)

### Community 164 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 165 - "useClick"
Cohesion: 0.39
Nodes (5): useClick(), getNextOpen(), onClick(), onMouseDown(), setOpenWithTouchDelay()

### Community 166 - "dialog_templ.go"
Cohesion: 0.16
Nodes (16): Close(), CloseFor(), initialOpen(), state(), Trigger(), TriggerFor(), ContentProps, ctxKey (+8 more)

### Community 167 - "registration-page"
Cohesion: 0.33
Nodes (6): registration-page, account-exists, already-have-account, description, submit, title

### Community 168 - "mark_others.js"
Cohesion: 0.52
Nodes (6): applyAttributeToOthers(), buildKeepSet(), collectOutsideElements(), correctElements(), markOthers(), unwrapHost()

### Community 170 - "Props"
Cohesion: 0.67
Nodes (3): Props, Variant, variantClasses()

### Community 171 - "useTypeahead"
Cohesion: 0.53
Nodes (4): useTypeahead(), getMatchingIndex(), isItemAvailable(), onKeyDown()

### Community 173 - "DefaultService"
Cohesion: 0.15
Nodes (17): main(), InitGinEngine(), registerAssets(), TestAssetCaching(), GetAppBaseURL(), IsDevelopment(), DefaultHandler, NewDefaultHandler() (+9 more)

### Community 179 - ".loadShoppingItem"
Cohesion: 0.27
Nodes (3): shoppingItemPath(), shoppingListPath(), ShoppingItem

### Community 180 - "invite-page"
Cohesion: 0.33
Nodes (6): invite-page, code, description, invalid, join, title

### Community 181 - "icon.go"
Cohesion: 0.73
Nodes (5): attrString(), generateSVG(), getIconContent(), Icon(), Props

### Community 182 - ".notifyChat"
Cohesion: 0.22
Nodes (5): ReplyPreview, householdPath(), pushPreview(), TestPushPreview(), Upload

### Community 183 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 184 - "appearance"
Cohesion: 0.25
Nodes (8): mode, mode-dark, mode-light, mode-system, theme, theme-home, theme-plain, appearance

### Community 186 - "field_templ.go"
Cohesion: 0.16
Nodes (15): childrenHTML(), fieldClasses(), Separator(), ContentProps, DescriptionProps, ErrorProps, GroupProps, LabelProps (+7 more)

### Community 187 - "Props"
Cohesion: 0.67
Nodes (3): initialChecked(), Props, Size

### Community 188 - "password"
Cohesion: 0.40
Nodes (5): password, confirm, enter, show, title

### Community 189 - "alertdialog_templ.go"
Cohesion: 0.15
Nodes (18): ContentProps, ctxKey, ctxState, DescriptionProps, FooterProps, HeaderProps, MediaProps, OverlayProps (+10 more)

## Knowledge Gaps
- **1 isolated node(s):** `roommates`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 955 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **40 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `en` connect `en` to `reset-password-page`, `expenses`, `household`, `forgot-password-page`, `appearance`, `shopping`, `registration-page`, `home-page`, `notes`, `files`, `invite-page`, `validation-errors`, `password`, `verify-email`, `login-page`, `nav`, `avatar`, `mail`, `chat`, `push`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `roommates` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.132890365448505 - nodes in this community are weakly interconnected._
- **Why does `et` connect `et` to `mail`, `avatar`, `household`, `registration-page`, `error-page`, `expenses`, `shopping`, `validation-errors`, `home-page`, `notes`, `appearance`, `verify-email`, `reset-password-page`, `password`, `files`, `nav`, `invite-page`, `chat`, `forgot-password-page`, `push`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `handlers.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12927104491085534 - nodes in this community are weakly interconnected._
- **Why does `T()` connect `T` to `github.com/gin-gonic/gin.Context`, `handlers.go`, `householdURL`, `Redirect`, `membershipFromContext`, `context.Context`, `newEmailToken`, `Membership`, `categorySelect`, `TestPushNotifications`, `.loadShoppingItem`, `.notifyChat`, `.shoppingAction`, `householdURL`, `github.com/a-h/templ.Component`, `chatMessages`, `.clean`, `.handleUpload`, `card_templ.go`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `en` be split into smaller, more focused modules?**
  _Cohesion score 0.07142857142857142 - nodes in this community are weakly interconnected._