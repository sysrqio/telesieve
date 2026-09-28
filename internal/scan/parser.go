package scan

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// LineHandler is called for each log line: rawLine is the file line, line is the extracted message body.
type LineHandler func(rawLine, line string) error

// ParseFile streams a log dump from path (optional .gz) and invokes h per line.
func ParseFile(path string, h LineHandler) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return ParseReader(f, strings.HasSuffix(strings.ToLower(path), ".gz"), h)
}

// ParseReader streams from r, optionally gzip-decoded.
func ParseReader(r io.Reader, gz bool, h LineHandler) error {
	var src io.Reader = r
	if gz {
		gr, err := gzip.NewReader(r)
		if err != nil {
			return fmt.Errorf("gzip: %w", err)
		}
		defer gr.Close()
		src = gr
	}
	sc := bufio.NewScanner(src)
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 10*1024*1024)
	for sc.Scan() {
		raw := sc.Text()
		msg := extractMessage(raw)
		if err := h(raw, msg); err != nil {
			return err
		}
	}
	return sc.Err()
}

func extractMessage(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	if raw[0] == '{' {
		var m map[string]any
		if json.Unmarshal([]byte(raw), &m) == nil {
			for _, key := range []string{"message", "msg", "body", "log", "text"} {
				if v, ok := m[key]; ok {
					if s, ok := v.(string); ok && s != "" {
						return s
					}
				}
			}
			if v, ok := m["attributes"]; ok {
				if am, ok := v.(map[string]any); ok {
					if s, ok := am["message"].(string); ok {
						return s
					}
				}
			}
		}
	}
	if strings.Contains(raw, "=") && !strings.Contains(raw, "{") {
		return logfmtMessage(raw)
	}
	return raw
}

func logfmtMessage(line string) string {
	if i := strings.Index(line, `msg="`); i >= 0 {
		rest := line[i+5:]
		if j := strings.Index(rest, `"`); j >= 0 {
			return rest[:j]
		}
	}
	if i := strings.Index(line, "msg="); i >= 0 {
		rest := strings.TrimSpace(line[i+4:])
		rest = strings.Trim(rest, `"`)
		if space := strings.Index(rest, " "); space > 0 && !strings.Contains(line, `msg="`) {
			return rest[:space]
		}
		if rest != "" {
			return rest
		}
	}
	parts := strings.Fields(line)
	var msgParts []string
	for _, p := range parts {
		if idx := strings.Index(p, "="); idx > 0 {
			key := p[:idx]
			val := strings.Trim(p[idx+1:], `"`)
			if key == "msg" || key == "message" {
				return val
			}
			if key == "level" || key == "ts" || key == "time" {
				continue
			}
			msgParts = append(msgParts, val)
		} else {
			msgParts = append(msgParts, p)
		}
	}
	if len(msgParts) > 0 {
		return strings.Join(msgParts, " ")
	}
	return line
}
