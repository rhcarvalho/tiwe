package commutative

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCommutative(t *testing.T) {
	plaintext := "4 blue"

	key1 := GenerateKey()
	key2 := GenerateKey()
	if reflect.DeepEqual(key1, key2) {
		t.Fatal("Generated keys are the same!")
	}

	msg := NewMessage([]byte(plaintext))
	key1.Encrypt(msg)
	ciphertext1 := fmt.Sprintf("%x", msg.Bytes)
	// Ciphertext should be different each time because of random
	// initialization vector, but it should not be equal to plaintext.
	if ciphertext1 == plaintext {
		t.Fatalf("Encrypt: got plaintext %q", ciphertext1)
	}
	key2.Encrypt(msg)
	ciphertext12 := fmt.Sprintf("%x", msg.Bytes)
	if ciphertext12 == ciphertext1 {
		t.Fatalf("Encrypt: double encryption did not change message %q", ciphertext12)
	}

	msgCopy := NewMessage(append([]byte(nil), msg.Bytes...))
	for k, v := range msg.Nonces {
		msgCopy.Nonces[k] = v
	}

	dec12 := key2.Decrypt(key1.Decrypt(msg))
	dec21 := key1.Decrypt(key2.Decrypt(msgCopy))
	if !reflect.DeepEqual(dec12, dec21) {
		t.Fatalf("Decrypt: not commutative: %x != %x", dec12.Bytes, dec21.Bytes)
	}
	if got := string(dec12.Bytes); got != plaintext {
		t.Fatalf("Decrypt: got %q, want %q", got, plaintext)
	}
}
