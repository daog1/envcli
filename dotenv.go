package envcli

import (
	"bytes"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const doubleQuoteSpecialChars = "\\\n\r\"!$`"

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader) (*orderedmap.OrderedMap[string, string], error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	return UnmarshalBytes(buf.Bytes())
}

// Unmarshal reads an env file from a string, returning a map of keys and values.
func Unmarshal(str string) (envMap *orderedmap.OrderedMap[string, string], err error) {
	return UnmarshalBytes([]byte(str))
}

// UnmarshalBytes parses env file from byte slice of chars, returning a map of keys and values.
func UnmarshalBytes(src []byte) (*orderedmap.OrderedMap[string, string], error) {
	//out := make(map[string]string)
	out := orderedmap.New[string, string]()
	err := parseBytes(src, out)

	return out, err
}

// Write serializes the given environment and writes it to a file.
func Write(envMap *orderedmap.OrderedMap[string, string], filename string) error {
	content, err := Marshal(envMap)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return file.Sync()
}

// Marshal outputs the given environment as a dotenv-formatted environment file.
// Each line is in the format: KEY="VALUE" where VALUE is backslash-escaped.
func Marshal(envMap *orderedmap.OrderedMap[string, string]) (string, error) {
	lines := make([]string, 0, envMap.Len())
	for pair := envMap.Oldest(); pair != nil; pair = pair.Next() {
		k := pair.Key
		v := pair.Value
		if d, err := strconv.Atoi(v); err == nil {
			lines = append(lines, fmt.Sprintf(`%s=%d`, k, d))
		} else {
			lines = append(lines, fmt.Sprintf(`%s="%s"`, k, doubleQuoteEscape(v)))
		}
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n"), nil
}

func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

func doubleQuoteEscape(line string) string {
	for _, c := range doubleQuoteSpecialChars {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.Replace(line, string(c), toReplace, -1)
	}
	return line
}
