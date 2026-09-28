package scan

import (
	"os"

	"github.com/sysrqio/telesieve/internal/pii"
)

// RunScan parses path and fills stats; returns file size for cost estimate.
func RunScan(path string, piiScanner *pii.Scanner, st *Stats) (fileSize int64, allPII []pii.Finding, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, nil, err
	}
	fileSize = info.Size()
	err = ParseFile(path, func(rawLine, line string) error {
		waste := ClassifyLine(rawLine)
		if len(waste) == 0 {
			waste = ClassifyLine(line)
		}
		st.RecordLine(line, waste)
		if piiScanner != nil {
			found := piiScanner.ScanLine(rawLine)
			if len(found) == 0 {
				found = piiScanner.ScanLine(line)
			}
			if len(found) > 0 {
				allPII = append(allPII, found...)
				st.PIIFindings += len(found)
				if pii.HasCritical(found) {
					st.PIICritical = true
				}
			}
		}
		return nil
	})
	return fileSize, allPII, err
}
