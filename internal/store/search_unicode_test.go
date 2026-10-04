package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/openclaw/graincrawl/internal/model"
)

func TestSearchNotesUsesUnicodeCaseFold(t *testing.T) {
	ctx := context.Background()
	st, err := Open(ctx, filepath.Join(t.TempDir(), "graincrawl.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	now := time.Now().UTC().Round(0)
	notes := []struct {
		id, title string
	}{
		{"emile", "Émile sync"},
		{"oil", "Ölgeschäft review"},
		{"plain", "Plain Title"},
	}
	for _, item := range notes {
		title := item.title
		note := model.Note{
			ID: item.id, Title: &title, Type: "meeting",
			CreatedAt: now, UpdatedAt: now, Source: model.SourcePrivateAPI, LastSeenAt: now,
		}
		if err := st.UpsertNote(ctx, note); err != nil {
			t.Fatal(err)
		}
	}
	cases := []struct {
		query string
		want  string
	}{
		{"Émile", "emile"},
		{"émile", "emile"},
		{"ÉMILE", "emile"},
		{"Ölgeschäft", "oil"},
		{"ölgeschäft", "oil"},
		{"Plain", "plain"},
		{"plain", "plain"},
	}
	for _, tc := range cases {
		results, err := st.SearchNotes(ctx, tc.query, 10)
		if err != nil {
			t.Fatalf("search %q: %v", tc.query, err)
		}
		if len(results) != 1 || results[0].ID != tc.want {
			t.Fatalf("search %q = %#v", tc.query, results)
		}
	}
}
