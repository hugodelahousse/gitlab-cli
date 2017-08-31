package main

import (
	"github.com/xanzy/go-gitlab"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	now := time.Now()
	issues := []*gitlab.Issue{
		{
			ID:             1,
			IID:            1,
			ProjectID:      1,
			Title:          "Test Issue",
			Description:    "Some Description",
			Labels:         []string{"TODO", "FEATURE"},
			Milestone:      nil,
			State:          "opened",
			UpdatedAt:      &now,
			CreatedAt:      &now,
			Subscribed:     false,
			UserNotesCount: 0,
			Confidential:   false,
			DueDate:        "",
			WebURL:         "",
		},
	}
	t.Run("Saving one issue", func(t *testing.T) {
		Save(issues)
	})
}
