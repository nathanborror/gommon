package tokens

import (
	"time"

	"github.com/jmoiron/modl"
)

// PreInsert sets the Created and Modified time before Token is saved.
func (t *Token) PreInsert(modl.SqlExecutor) error {
	if t.Created.IsZero() {
		t.Created = time.Now()
	}
	if t.Modified.IsZero() {
		t.Modified = time.Now()
	}
	return nil
}

// PreUpdate updates the Modified time before Token is updated.
func (t *Token) PreUpdate(modl.SqlExecutor) error {
	t.Modified = time.Now()
	return nil
}
