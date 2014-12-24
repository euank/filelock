package lock

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewLock(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating tmpfile: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	l, err := newLock(f.Name())
	if err == nil || l != nil {
		t.Fatal("expected error creating lock on file")
	}

	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating tmpdir: %v", err)
	}
	defer os.Remove(d)

	l, err = newLock(d)
	if err != nil {
		t.Fatalf("error creating newLock: %v", err)
	}

	l.Close()
	if err != nil {
		t.Fatalf("error unlocking lock: %v", err)
	}

	if err = os.Remove(d); err != nil {
		t.Fatalf("error removing tmpdir: %v", err)
	}

	l, err = newLock(d)
	if err == nil {
		t.Fatalf("expected error creating lock on nonexistent path")
	}
}

func TestExclusiveLock(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating tmpdir: %v", err)
	}
	defer os.Remove(dir)

	// Set up the initial exclusive lock
	l, err := ExclusiveLock(dir)
	if err != nil {
		t.Fatalf("error creating lock: %v", err)
	}

	// Now try another exclusive lock, should fail
	_, err = TryExclusiveLock(dir)
	if err == nil {
		t.Fatalf("expected err trying exclusive lock")
	}

	// Unlock the original lock
	err = l.Close()
	if err != nil {
		t.Fatalf("error closing lock: %v", err)
	}

	// Now another exclusive lock should succeed
	_, err = TryExclusiveLock(dir)
	if err != nil {
		t.Fatalf("error creating lock: %v", err)
	}
}

func TestSharedLock(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating tmpdir: %v", err)
	}
	defer os.Remove(dir)

	// Set up the initial shared lock
	l1, err := SharedLock(dir)
	if err != nil {
		t.Fatalf("error creating new shared lock: %v", err)
	}

	// Subsequent shared locks should succeed
	l2, err := TrySharedLock(dir)
	if err != nil {
		t.Fatalf("error creating shared lock: %v", err)
	}
	l3, err := TrySharedLock(dir)
	if err != nil {
		t.Fatalf("error creating shared lock: %v", err)
	}

	// But an exclusive lock should fail
	_, err = TryExclusiveLock(dir)
	if err == nil {
		t.Fatal("expected exclusive lock to fail")
	}

	// Close the locks
	err = l1.Close()
	if err != nil {
		t.Fatalf("error closing lock: %v", err)
	}
	err = l2.Close()
	if err != nil {
		t.Fatalf("error closing lock: %v", err)
	}
	err = l3.Close()
	if err != nil {
		t.Fatalf("error closing lock: %v", err)
	}

	// Now try an exclusive lock, should succeed
	_, err = TryExclusiveLock(dir)
	if err != nil {
		t.Fatalf("error creating lock: %v", err)
	}
}
