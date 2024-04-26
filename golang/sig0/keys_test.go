package sig0

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadKey(t *testing.T) {
	keyName := createKey(t)

	signer, err := LoadKeyFile(keyName)
	if err != nil {
		t.Fatal(err)
	}
	if signer == nil {
		t.Fatal("signer is nil")
	}
}

func TestParseKeyFile(t *testing.T) {
	keyName := createKey(t)

	keyContent, err := os.ReadFile(keyName + ".key")
	if err != nil {
		t.Fatal(err)
	}

	privateContent, err := os.ReadFile(keyName + ".private")
	if err != nil {
		t.Fatal(err)
	}

	signer, err := ParseKeyData(string(keyContent), string(privateContent))
	if err != nil {
		t.Fatal(err)
	}
	if signer == nil {
		t.Fatal("signer is nil")
	}

	// signer.dnsKey.Hdr.Ttl = 0
	t.Logf("pk type: %T", signer.private)

	k := signer.dnsKey.DNSKEY

	out := k.String()
	out = strings.ReplaceAll(out, "\t", " ")

	assert.Equal(t, string(keyContent), out)

	pk := signer.dnsKey.PrivateKeyString(signer.private)
	assert.Equal(t, string(privateContent), pk)
}

func createKey(t *testing.T) string {
	var buf bytes.Buffer
	cmd := exec.Command("dnssec-keygen", "-K", "/tmp", "-a", "ED25519", "-n", "HOST", "-T", "KEY", "go.te.st")
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	keyName := filepath.Join("/tmp", strings.TrimSpace(buf.String()))

	t.Log("created key file:", keyName)

	t.Cleanup(func() {
		os.Remove(keyName + ".key")
		os.Remove(keyName + ".private")
	})

	return keyName
}

func createAndLoadKey(t *testing.T) *Signer {
	keyName := createKey(t)

	signer, err := LoadKeyFile(keyName)
	if err != nil {
		t.Fatal(err)
	}
	if signer == nil {
		t.Fatal("signer is nil")
	}

	return signer
}
