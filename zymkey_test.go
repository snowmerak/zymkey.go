package zymkey

import (
	"bytes"
	"testing"
)

func TestZymkey(t *testing.T) {
	zk, err := NewZymkey()
	if err != nil {
		t.Fatalf("Failed to initialize Zymkey: %v", err)
	}
	defer zk.Close()

	t.Run("GenerateRandomBytes", func(t *testing.T) {
		length := 32
		data, err := zk.GenerateRandomBytes(length)
		if err != nil {
			t.Fatalf("GenerateRandomBytes failed: %v", err)
		}
		if len(data) != length {
			t.Errorf("Expected length %d, got %d", length, len(data))
		}
	})

	t.Run("Read", func(t *testing.T) {
		length := 16
		buf := make([]byte, length)
		n, err := zk.Read(buf)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if n != length {
			t.Errorf("Expected to read %d bytes, got %d", length, n)
		}
	})

	t.Run("LockUnlock", func(t *testing.T) {
		original := []byte("secret message")
		locked, err := zk.LockBuffer(original)
		if err != nil {
			t.Fatalf("LockBuffer failed: %v", err)
		}
		if len(locked) == 0 {
			t.Error("Locked buffer is empty")
		}

		// Ensure locked data is different from original (basic check)
		if bytes.Equal(original, locked) {
			t.Log("Warning: Locked data is identical to original data")
		}

		unlocked, err := zk.UnlockBuffer(locked)
		if err != nil {
			t.Fatalf("UnlockBuffer failed: %v", err)
		}

		if !bytes.Equal(original, unlocked) {
			t.Errorf("Original and unlocked data do not match.\nOriginal: %s\nUnlocked: %s", original, unlocked)
		}
	})

	t.Run("Now", func(t *testing.T) {
		ts, err := zk.Now()
		if err != nil {
			t.Fatalf("Now failed: %v", err)
		}
		if ts.IsZero() {
			t.Error("Returned zero time")
		}
		t.Logf("Current Zymkey time: %v", ts)
	})
}
