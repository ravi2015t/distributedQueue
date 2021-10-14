package server

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestInitLastChunkIdx(t *testing.T) {
	dir := getTempDir(t)

	testCreateFile(t, filepath.Join(dir, "moscow-chunk1"))
	testCreateFile(t, filepath.Join(dir, "moscow-chunk10"))

	t.Cleanup(func() { os.RemoveAll(dir) })

	srv := testNewOnDisk(t, dir)

	want := uint64(11)
	got := srv.lastChunkIdx

	if got != want {
		t.Errorf("Last chunk idx = %d, want %d", got, want)
	}
}

func TestGetFileDescriptor(t *testing.T) {
	dir := getTempDir(t)

	testCreateFile(t, filepath.Join(dir, "moscow-chunk1"))

	t.Cleanup(func() { os.RemoveAll(dir) })

	srv := testNewOnDisk(t, dir)

	testcases := []struct {
		desc     string
		filename string
		write    bool
		wantErr  bool
	}{
		{
			desc:     "Read from a existing file shouldn't file",
			filename: "moscow-chunk1",
			write:    false,
			wantErr:  false,
		},
		{
			desc:     "Shouldn't override existing files",
			filename: "moscow-chunk1",
			write:    true,
			wantErr:  true,
		},
		{
			desc:     "Should not be able to read from files that don't exist",
			filename: "moscow-chunk2",
			write:    false,
			wantErr:  true,
		},
		{
			desc:     "Should be able to create files that don't exist",
			filename: "moscow-chunk2",
			write:    true,
			wantErr:  false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := srv.getFileDescriptor(tc.filename, tc.write)

			defer srv.forgetFileDescriptor(tc.filename)

			if tc.wantErr && err == nil {
				t.Errorf("wanted errors, got not errors")
			} else if !tc.wantErr && err != nil {
				t.Errorf("wanted no errors, got error: %v", err)
			}
		})
	}
}

func TestReadWrite(t *testing.T) {

	srv := testNewOnDisk(t, getTempDir(t))
	want := "first\nsecond\nthird\n"
	if err := srv.Write(context.Background(), []byte(want)); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	chunks, err := srv.ListChunks()
	if err != nil {
		t.Fatalf("List chunks failed: %v", err)
	}

	if want, got := 1, len(chunks); want != got {
		t.Fatalf("len(ListChunks()) = %d, want %d", got, want)
	}

	chunk := chunks[0].Name
	var b bytes.Buffer
	if err := srv.Read(chunk, 0, uint64(len(want)), &b); err != nil {
		t.Fatalf("Read(%q) = %v, want no errors", chunk, err)
	}

	got := b.String()

	if got != want {
		t.Fatalf("Read(%q) = %q want %q", chunk, got, want)
	}

	//Check last message is not truncated
	want = "first\nsecond\n"
	b.Reset()
	if err := srv.Read(chunk, 0, uint64(len("first\nsecond\n")+1), &b); err != nil {
		t.Fatalf("Read(%q) = %v, want no errors", chunk, err)
	}

	got = b.String()

	if got != want {
		t.Fatalf("Read(%q) = %q want %q", chunk, got, want)
	}
}

type nilHooks struct{}

func (n *nilHooks) BeforeCreatingChunk(ctx context.Context, category string, fileName string) error {
	return nil
}

func testNewOnDisk(t *testing.T, dir string) *OnDisk {
	t.Helper()

	srv, err := NewOnDisk(dir, "test", "moscow", &nilHooks{})
	if err != nil {
		t.Fatalf("New on Disk failed: %v", err)
	}
	return srv
}

func testCreateFile(t *testing.T, filename string) {
	t.Helper()

	if _, err := os.Create(filename); err != nil {
		t.Fatalf("couldn't create a file %q: %v", filename, err)
	}

}

func getTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp(os.TempDir(), "lastchunkidx")
	if err != nil {
		t.Fatalf("mkdir temp failed : %v", err)
	}

	t.Cleanup(func() { os.RemoveAll(dir) })

	return dir
}

func TestAckOfLastChunk(t *testing.T) {
	srv := testNewOnDisk(t, getTempDir(t))
	want := "first\nsecond\nthird\n"
	if err := srv.Write(context.Background(), []byte(want)); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	chunks, err := srv.ListChunks()
	if err != nil {
		t.Fatalf("List chunks failed: %v", err)
	}

	if want, got := 1, len(chunks); want != got {
		t.Fatalf("len(ListChunks()) = %d, want %d", got, want)
	}

	if err := srv.Ack(chunks[0].Name, chunks[0].Size); err == nil {
		t.Errorf("Ack() : got no errors, expected an error")
	}
}

func TestAckOfCompletedChunk(t *testing.T) {
	dir := getTempDir(t)
	srv := testNewOnDisk(t, dir)
	testCreateFile(t, filepath.Join(dir, "moscow-chunk1"))

	if err := srv.Ack("moscow-chunk1", 0); err != nil {
		t.Errorf("Ack(chunk1) = %v, expected no errors", err)
	}
}
