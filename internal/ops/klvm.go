package ops

import (
	"bytes"
	"strings"
)

// klvm format taken from Thibault Polge's "Write yourself a Git!" article
// real lifesaver

func parseKlvm(data []byte) (map[string][]string, string) {
	kvlm := make(map[string][]string)
	lines := bytes.Split(data, []byte{'\n'})

	var messageStartIndex int
	var currentKey string

	for i, line := range lines {
		if len(line) == 0 {
			messageStartIndex = i + 1
			break
		}

		if line[0] == ' ' {
			// Continuation of previous key
			kvlm[currentKey] = append(kvlm[currentKey], string(line[1:]))
		} else {
			parts := bytes.SplitN(line, []byte{' '}, 2)
			if len(parts) == 2 {
				key := string(parts[0])
				value := string(parts[1])
				currentKey = key
				if existingValue, ok := kvlm[key]; ok {
					kvlm[key] = append(existingValue, value)
				} else {
					kvlm[key] = []string{value}
				}
			}
		}
	}

	message := string(bytes.Join(lines[messageStartIndex:], []byte{'\n'}))
	return kvlm, message
}

func serializeKlvm(klvm map[string][]string, message string) []byte {
	var buffer bytes.Buffer

	for key, values := range klvm {
		for _, value := range values {
			buffer.WriteString(key)
			buffer.WriteByte(' ')
			buffer.WriteString(strings.Replace(value, "\n", "\n ", -1))
			buffer.WriteByte('\n')
		}
	}

	buffer.WriteByte('\n')
	buffer.WriteString(message)

	return buffer.Bytes()
}
