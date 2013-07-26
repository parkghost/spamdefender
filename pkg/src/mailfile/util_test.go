package mailfile

import (
	"testing"
)

func TestDecodeRFC2047StringWithUTF8(t *testing.T) {
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

func TestDecodeRFC2047StringWithBig5(t *testing.T) {
	original := "=?BIG5?Q?=B2=C41717=A6^JWorld=AAZ=B9D=A4j=B7|-=A4C=C0s=AF]=A6A=B2{!=3F?="
	expected := "第1717回JWorld武道大會-七龍珠再現!?"

	result, err := DecodeRFC2047String(original)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
