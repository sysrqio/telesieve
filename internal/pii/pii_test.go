package pii

import "testing"

const sampleJWT = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0In0.x"

func TestScanLine_JWT(t *testing.T) {
	s := NewScanner(true)
	f := s.ScanLine(sampleJWT)
	if len(f) == 0 {
		t.Fatal("jwt")
	}
	if !HasCritical(f) {
		t.Error("critical")
	}
	if f[0].Masked == sampleJWT {
		t.Error("should mask")
	}
}

func TestScanLine_Authorization(t *testing.T) {
	s := NewScanner(false)
	f := s.ScanLine("Authorization: Bearer secret-token-xyz")
	if len(f) == 0 {
		t.Fatal("auth header")
	}
}

func TestScanLine_CreditCard(t *testing.T) {
	s := NewScanner(true)
	f := s.ScanLine("charge 4111111111111111")
	found := false
	for _, x := range f {
		if x.Kind == "credit_card_like" {
			found = true
		}
	}
	if !found {
		t.Fatal("cc")
	}
}

func TestScanLine_IBAN(t *testing.T) {
	s := NewScanner(true)
	f := s.ScanLine("transfer DE89370400440532013000")
	found := false
	for _, x := range f {
		if x.Kind == "iban_like" {
			found = true
		}
	}
	if !found {
		t.Fatal("iban")
	}
}

func TestScanLine_Clean(t *testing.T) {
	s := NewScanner(true)
	if len(s.ScanLine("hello world")) != 0 {
		t.Error("clean")
	}
}

func TestHasCritical_InfoOnly(t *testing.T) {
	if HasCritical([]Finding{{Kind: "x", Severity: SeverityInfo}}) {
		t.Error("info not critical")
	}
}

func TestLuhn(t *testing.T) {
	if !luhnDigits("4111111111111111") {
		t.Error("luhn")
	}
	if luhnDigits("4111111111111112") {
		t.Error("bad luhn")
	}
}
