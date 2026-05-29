package headers

import (
	"bytes"
	"errors"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	i := bytes.Index(data, []byte(crlf))
	// non c'è il carattere di fine riga
	if i == -1 {
		return 0, false, nil
	}
	// è ad inizio riga quindi siamo alla fine degli headers, unico caso in cui "done = true"
	if i == 0 {
		return 2, true, nil // 2 sono i byte di CRLF
	}

	line := string(data[:i])
	idx := strings.Index(line, ":")
	if idx == -1 {
		return 0, false, errors.New("wrong header format: missing colon")
	}

	kWithColon := line[:idx+1]
	val := strings.TrimSpace(line[idx+1:])

	if ok := validateKey(kWithColon); !ok {
		return 0, false, errors.New("wrong KEY format")
	}

	// pulisco la Key per salvare senza ":"
	kClean := line[:idx]

	// e aggiungo il nuovo header
	h[kClean] = val

	return i + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

func validateKey(s string) bool {
	// 1. La stringa deve finire con ":" ed essere lunga almeno 2 caratteri (es. "A:")
	if len(s) < 2 || !strings.HasSuffix(s, ":") {
		return false
	}

	// 2. Isoliamo tutto ciò che c'è PRIMA del ":"
	rest := s[:len(s)-1]

	// 3. Controlliamo il PRIMO carattere del resto (deve essere MAIUSCOLO 'A'-'Z')
	if rest[0] < 'A' || rest[0] > 'Z' {
		return false
	}

	// 4. Controlliamo tutti i caratteri SUCCESSIVI (non devono esserci spazi)
	for i := 1; i < len(rest); i++ {
		if rest[i] == ' ' {
			return false
		}
	}

	// Se ha superato tutti i controlli, la stringa è valida
	return true
}
