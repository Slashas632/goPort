package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// reset clears package-level state between tests, since results/JSONPath
// are package globals shared across the whole test binary.
func reset() {
	mu.Lock()
	results = nil
	mu.Unlock()
	JSONPath = ""
}

func TestInit_SetsJSONPath(t *testing.T) {
	reset()
	Init("scan.json")
	if JSONPath != "scan.json" {
		t.Errorf("JSONPath = %q, want scan.json", JSONPath)
	}
}

func TestAddResult_AppendsResult(t *testing.T) {
	reset()
	AddResult("10.0.0.1", 80, "HTTP/1.1 200 OK")
	AddResult("10.0.0.1", 22, "SSH-2.0-OpenSSH_9.2")

	mu.Lock()
	n := len(results)
	mu.Unlock()
	if n != 2 {
		t.Fatalf("len(results) = %d, want 2", n)
	}
}

func TestSaveJson_WritesValidJSON(t *testing.T) {
	reset()
	AddResult("10.0.0.1", 80, "HTTP/1.1 200 OK")
	AddResult("10.0.0.1", 22, "SSH-2.0-OpenSSH_9.2")

	dir := t.TempDir()
	path := filepath.Join(dir, "results.json")

	if err := SaveJson(path); err != nil {
		t.Fatalf("SaveJson returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}

	var got []Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}

	want := []Result{
		{IP: "10.0.0.1", Port: 80, Banner: "HTTP/1.1 200 OK"},
		{IP: "10.0.0.1", Port: 22, Banner: "SSH-2.0-OpenSSH_9.2"},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d results, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("result[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestSaveJson_EmptyResultsWritesEmptyArray(t *testing.T) {
	reset()

	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	if err := SaveJson(path); err != nil {
		t.Fatalf("SaveJson returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}

	var got []Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d results, want 0", len(got))
	}
}

func TestSaveJson_InvalidPathReturnsError(t *testing.T) {
	reset()
	AddResult("10.0.0.1", 80, "banner")

	err := SaveJson(filepath.Join(t.TempDir(), "no-such-dir", "results.json"))
	if err == nil {
		t.Fatal("expected an error when the parent directory does not exist, got none")
	}
}

func TestAddResult_ConcurrentSafe(t *testing.T) {
	reset()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			AddResult("10.0.0.1", n, "banner")
		}(i)
	}
	wg.Wait()

	mu.Lock()
	n := len(results)
	mu.Unlock()
	if n != 100 {
		t.Errorf("len(results) = %d, want 100 (possible race in AddResult)", n)
	}
}
