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

	// lavoro sui bytes per trovare il separatore ":" per migliorare le prestazioni (non alloca stringhe se non va bene l'header)
	lineBytes := data[:i]
	idx := bytes.IndexByte(lineBytes, ':')
	if idx == -1 {
		return 0, false, errors.New("wrong header format: missing colon")
	}

	// isolo la chiave prima del separatore
	keyBytes := lineBytes[:idx]

	if !validateKey(keyBytes) {
		return 0, false, errors.New("wrong KEY format")
	}

	kLower := strings.ToLower(string(keyBytes))

	// Isoliamo il valore (i byte DOPO il ':') e togliamo gli spazi
	val := strings.TrimSpace(string(lineBytes[idx+1:]))

	// controllo se esiste già la chiave negli headers
	if v, ok := h[kLower]; ok {
		val = strings.Join([]string{v, val}, ", ")
	}

	// e aggiungo il nuovo header
	h[kLower] = val

	return i + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

func validateKey(key []byte) bool {
	// lunghezza minima di 1 carattere (come richiesto dall'RFC)
	if len(key) == 0 {
		return false
	}

	for i := 0; i < len(key); i++ {
		b := key[i]

		switch {
		case b >= 'a' && b <= 'z':
		case b >= 'A' && b <= 'Z':
		case b >= '0' && b <= '9':
		case b == '!' || b == '#' || b == '$' || b == '%' || b == '&' || b == '\'' ||
			b == '*' || b == '+' || b == '-' || b == '.' || b == '^' || b == '_' ||
			b == '`' || b == '|' || b == '~':
			// Carattere valido, continua
		default:
			return false
		}
	}

	return true
}
