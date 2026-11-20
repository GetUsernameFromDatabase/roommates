package pages

import (
	"context"
	"strings"
	"testing"
	"time"

	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"

	"github.com/invopop/ctxi18n"
)

// "for X" follows the assignee id, an erased assignee's empty name still gets it
func TestChoreLogLineForOther(t *testing.T) {
	locales.LoadLocales()
	ctx, err := ctxi18n.WithLocale(context.Background(), "en")
	if err != nil {
		t.Fatal(err)
	}
	former := utils.MemberName(ctx, "")
	entry := func(assignee string, forOther bool) services.ChoreLogEntry {
		return services.ChoreLogEntry{Kind: g.ChoreDone, UserName: "Mari", AssigneeName: assignee, ForOther: forOther, Due: time.Time{}, At: time.Now()}
	}
	if got := choreLogLine(ctx, entry("", true)); !strings.Contains(got, former) {
		t.Errorf("erased assignee: %q lacks %q", got, former)
	}
	if got := choreLogLine(ctx, entry("Jaan", true)); !strings.Contains(got, "Jaan") {
		t.Errorf("other assignee: %q lacks Jaan", got)
	}
	if got := choreLogLine(ctx, entry("Mari", false)); strings.Count(got, "Mari") != 1 {
		t.Errorf("own turn: %q names Mari twice", got)
	}
}
