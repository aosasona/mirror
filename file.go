package mirror

import (
	"bytes"
	"os"
)

func (m *Mirror) areSameBytesContent(out string) bool {
	stat, err := os.Stat(m.config.OutputFileOrDefault())
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}

	file, err := os.ReadFile(m.config.OutputFileOrDefault())
	if err != nil {
		return false
	}

	return bytes.Equal(file, []byte(out))
}
