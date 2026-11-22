package pages

import (
	"context"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"

	"github.com/invopop/ctxi18n/i18n"
)

// names of the [services.FileKinds] and [services.FilePlaces] options
var (
	fileKindNames = map[string]locales.LK{
		"image": locales.LK_FilesTypeImage,
		"pdf":   locales.LK_FilesTypePdf,
		"other": locales.LK_FilesTypeOther,
	}
	filePlaceNames = map[string]locales.LK{
		"note":    locales.LK_FilesPlaceNote,
		"item":    locales.LK_FilesPlaceItem,
		"expense": locales.LK_FilesPlaceExpense,
		"message": locales.LK_FilesPlaceMessage,
		"event":   locales.LK_FilesPlaceEvent,
		"none":    locales.LK_FilesPlaceNone,
	}
)

// deleteConfirm adds to question how many files would go with the deletion
func deleteConfirm(ctx context.Context, question string, files int) string {
	if files == 0 {
		return question
	}
	return question + "\n" + utils.N(ctx, locales.LK_FilesDeletedWith, files, i18n.M{"count": files})
}

// deleteKeepConfirm is [deleteConfirm] for a form with data-confirm-keep, the dialog itself saves the files
func deleteKeepConfirm(ctx context.Context, question string, files int) string {
	if files == 0 {
		return question
	}
	return question + "\n" + utils.N(ctx, locales.LK_FilesDeletedWithKeep, files, i18n.M{"count": files})
}

// deleteFileConfirm says where else a file on refs places goes from, the Files list shows only its first
func deleteFileConfirm(ctx context.Context, name string, refs int64) string {
	question := utils.T(ctx, locales.LK_FilesDeleteConfirm, name)
	if refs < 2 {
		return question
	}
	return question + "\n" + utils.T(ctx, locales.LK_FilesDeleteShared, refs)
}

// lastRefs counts the files that go with what they are attached to
func lastRefs(files []services.AttachedFile) int {
	n := 0
	for _, f := range files {
		if f.LastRef {
			n++
		}
	}
	return n
}

const fileContextClass = "block text-xs text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"

const filePickClass = "flex min-h-11 w-full items-center gap-2 rounded-lg px-2 text-left text-sm outline-none hover:bg-muted focus-visible:ring-3 focus-visible:ring-ring/50"
