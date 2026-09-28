package pii

import (
	"regexp"
	"strings"
)

// Severity of a PII match.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityCritical
)

// Finding describes one PII detection.
type Finding struct {
	Kind     string
	Severity Severity
	Masked   string
	Count    int
}

var (
	jwtBearer = regexp.MustCompile(`(?i)(Bearer\s+)?eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)
	authHdr   = regexp.MustCompile(`(?i)Authorization:\s*(Bearer\s+)?[^\s]+`)
	ibanLike  = regexp.MustCompile(`\b[A-Z]{2}[0-9]{2}[A-Z0-9]{11,30}\b`)
	ccLike    = regexp.MustCompile(`\b(?:\d[ -]*?){13,19}\b`)
)

// Scanner detects sensitive patterns in log text.
type Scanner struct {
	maskEnabled bool
}

// NewScanner creates a PII scanner.
func NewScanner(maskEnabled bool) *Scanner {
	return &Scanner{maskEnabled: maskEnabled}
}

// ScanLine returns findings for one line (may be empty).
func (s *Scanner) ScanLine(line string) []Finding {
	var out []Finding
	if m := jwtBearer.FindString(line); m != "" {
		out = append(out, Finding{Kind: "jwt_bearer", Severity: SeverityCritical, Masked: mask(m, s.maskEnabled), Count: 1})
	}
	if m := authHdr.FindString(line); m != "" {
		out = append(out, Finding{Kind: "authorization_header", Severity: SeverityCritical, Masked: mask(m, s.maskEnabled), Count: 1})
	}
	for _, m := range ibanLike.FindAllString(line, -1) {
		if len(digitsOnly(m)) >= 13 {
			out = append(out, Finding{Kind: "iban_like", Severity: SeverityCritical, Masked: mask(m, s.maskEnabled), Count: 1})
		}
	}
	for _, m := range ccLike.FindAllString(line, -1) {
		digits := digitsOnly(m)
		if len(digits) >= 13 && len(digits) <= 19 && luhnDigits(digits) {
			out = append(out, Finding{Kind: "credit_card_like", Severity: SeverityCritical, Masked: mask(m, s.maskEnabled), Count: 1})
		}
	}
	return out
}

func mask(s string, enabled bool) string {
	if !enabled {
		return s
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func luhnDigits(digits string) bool {
	sum := 0
	alt := false
	for i := len(digits) - 1; i >= 0; i-- {
		n := int(digits[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

// HasCritical reports whether any finding is critical severity.
func HasCritical(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == SeverityCritical {
			return true
		}
	}
	return false
}
