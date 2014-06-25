package auth

import (
	"github.com/jmoiron/modl"
	"time"
)

// PreInsert sets the Created and Modified time before User is saved.
func (u *User) PreInsert(modl.SqlExecutor) error {
	if u.Created.IsZero() {
		u.Created = time.Now()
	}
	if u.Modified.IsZero() {
		u.Modified = time.Now()
	}
	return nil
}

// PreUpdate updates the Modified time before User is updated.
func (u *User) PreUpdate(modl.SqlExecutor) error {
	u.Modified = time.Now()
	return nil
}
