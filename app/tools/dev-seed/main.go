// adds two users sharing a household with example data for browser checks,
// the password is public so it only runs in development
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"time"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/mailer"
	"roommates/internal/push"
	"roommates/internal/services"
	"roommates/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	emailMari = "mari@roommates.test"
	emailJaan = "jaan@roommates.test"
	password  = "Roommates-dev-1"
)

func main() {
	if err := checkDevelopment(); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, g.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	locales.LoadLocales()
	s := services.NewDefaultService(pool, mailer.New(), storage.New(ctx), push.New())

	if seed(ctx, s, dbqueries.New(pool)) {
		fmt.Printf("seeded %s and %s, password %s\n", emailMari, emailJaan, password)
	} else {
		fmt.Println("already seeded, nothing changed")
	}
}

func checkDevelopment() error {
	if !g.IsDevelopment() {
		return errors.New("dev-seed only runs in development, its password is public")
	}
	return nil
}

// seed returns false when the data already exists
//
// shortcut: a failed run is not resumed since the first user already exists, delete that user to retry
func seed(ctx context.Context, s *services.DefaultService, q *dbqueries.Queries) bool {
	mari, err := s.RegisterUser(&gin.Context{}, dbqueries.InsertUserParams{Email: emailMari, Username: "Mari", Password: password})
	if errors.Is(err, g.ErrorAccountAlreadyExists) {
		return false
	}
	check(err)
	jaan := must(s.RegisterUser(&gin.Context{}, dbqueries.InsertUserParams{Email: emailJaan, Username: "Jaan", Password: password}))
	check(q.SetEmailValidated(ctx, *mari))
	check(q.SetEmailValidated(ctx, *jaan))

	householdID := must(s.CreateHousehold(ctx, *mari, "Tähe 3 korter"))
	mariM := must(s.GetMembership(ctx, householdID, *mari))
	check(s.CreateInvite(ctx, mariM))
	overview := must(s.GetHouseholdOverview(ctx, mariM))
	must(s.JoinWithInvite(ctx, overview.Invites[0].Code, *jaan))
	jaanM := must(s.GetMembership(ctx, householdID, *jaan))

	seedShopping(ctx, s, mariM)
	seedExpenses(ctx, s, mariM, jaanM)

	noteID := must(s.CreateNote(ctx, mariM, services.NoteInput{
		Title:   "Maja reeglid",
		BodyMD:  "## Prügi\n\n- **Esmaspäeviti** viib Jaan\n- **Neljapäeviti** viib Mari",
		Private: false,
	}))
	check(s.SetNotePinned(ctx, mariM, noteID, true))

	check(s.PostMessage(ctx, mariM, "Kas keegi saab homme kullerile ukse avada?", pgtype.UUID{}, nil, nil))
	check(s.PostMessage(ctx, jaanM, "Mina olen kodus, pole probleemi.", pgtype.UUID{}, nil, nil))

	var picture bytes.Buffer
	check(png.Encode(&picture, examplePNG()))
	must(s.UploadFile(ctx, mariM, services.FileTarget{}, "arve.png", bytes.NewReader(picture.Bytes()), int64(picture.Len())))
	return true
}

func seedShopping(ctx context.Context, s *services.DefaultService, m services.Membership) {
	listID := must(s.CreateShoppingList(ctx, m, "Nädala toit", "groceries"))
	check(s.AddShoppingItem(ctx, m, listID, "Piim", "2"))
	check(s.AddShoppingItem(ctx, m, listID, "Leib", "1"))
	check(s.AddShoppingItem(ctx, m, listID, "Kohv", ""))
	list := must(s.GetShoppingList(ctx, m, listID))
	check(s.SetShoppingItemChecked(ctx, m, listID, list.Items[0].ID, true))
}

// one share is paid, one unpaid, and Jaan's payment waits for Mari to confirm it
func seedExpenses(ctx context.Context, s *services.DefaultService, mariM, jaanM services.Membership) {
	members := []pgtype.UUID{mariM.UserID, jaanM.UserID}
	today := time.Now().Format(services.DateLayout)
	split := func(title, amount, category string) services.ExpenseInput {
		return services.ExpenseInput{
			Title: title, Amount: amount, Category: category, IncurredOn: today, Members: members,
			Custom: false, Shares: nil, ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
		}
	}
	paid := must(s.CreateExpense(ctx, mariM, split("Elekter", "64,20", "utilities")))
	check(s.SetExpenseSharePaid(ctx, mariM, paid, jaanM.UserID, true))
	must(s.CreateExpense(ctx, mariM, split("Toidupood", "38,50", "groceries")))
	must(s.CreateExpense(ctx, jaanM, split("Internet", "25", "household")))
	check(s.CreateSettlement(ctx, jaanM, mariM.UserID, "3"))
}

func examplePNG() image.Image {
	return image.NewGray(image.Rect(0, 0, 320, 200))
}

// an example-data failure has nothing to recover, stop with the stack
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](value T, err error) T {
	check(err)
	return value
}
