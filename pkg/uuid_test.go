package pkg

import "testing"

func TestDuplicatedUUID(t *testing.T) {
	t.Run("duplicated uuid", func(t *testing.T) {
		gotUid, err := NewUUID()
		if err != nil {
			t.Errorf("NewUUID() error = %v", err)
			return
		}
		if !isCollected(gotUid) {
			t.Errorf("duplicated uuid not collected")
			return
		}
	})
}
