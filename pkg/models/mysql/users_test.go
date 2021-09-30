package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/DataDavD/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	// t.Parallel()
	// Skip the test if the '-short' flag is provided when running the test.
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:      1,
				Name:    "Alice Jones2",
				Email:   "alice2@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			wantError: nil,
		},
		{
			name:      "Non-existent ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize a connection pool to our test database, and defer a call to the teardown
			// function, so it is always run immediately before this sub-test returns.
			db, teardown := newTestDB(t)
			defer teardown()

			// Create a new instance of the UserModel.
			m := UserModel{db}

			// Call the UserModel.Get() method and check that the return value and error match
			// the expected values for the sub-test.
			user, err := m.Get(tt.userID)

			t.Logf("testing %q for want-user %v and want-error %v", tt.name, tt.wantUser,
				tt.wantError)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v; got %v", tt.wantUser, user)
			}
		})
	}
}
