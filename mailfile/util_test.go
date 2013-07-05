package mailfile

import (
	"testing"
)

func TestDecodeRFC2047String(t *testing.T) {
	original := "=?utf-8?Q?JWorld@TW=E8=A9=B1=E9=A1=8C=E6=9B=B4=E6=96=B0=E9?= =?utf-8?Q?=80=9A=E7=9F=A5:Functional?= Programming for Java Developers"
	expected := "JWorld@TW話題更新通知:Functional Programming for Java Developers"

	result, err := DecodeRFC2047String(original)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
