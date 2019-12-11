package fingerprint

import "testing"

func TestFingerprint(t *testing.T) {
	newFingerprint := Fingerprint{}
	if newFingerprint.IsValid() {
		t.Errorf("expect empty fingerprint invalid, but result valid")
	}

	// TODO: CRC8 checksum
	newFingerprint.Write([]byte{'O', 'S', 'R', 'N', 5, 21, 12, 0})
	if !newFingerprint.IsValid() {
		t.Errorf("expect fingerprint valid after filled in, but result invalid")
	}
	newString := newFingerprint.String()
	expectString := "OSRN v5.21.12"
	if newString != expectString {
		t.Errorf("expect %s but got %s", expectString, newString)
	}
}
