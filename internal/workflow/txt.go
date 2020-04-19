package workflow

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ToJsonString(v interface{}) string {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func ToYamlString(v interface{}) string {
	buf, err := yaml.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func Indent(prefix, str string) string {
	r := bufio.NewReader(bytes.NewBuffer([]byte(str)))
	w := &bytes.Buffer{}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				w.WriteString(prefix)
				w.WriteString(line)
			}
			break
		}
		w.WriteString(prefix)
		w.WriteString(line)
	}

	return strings.TrimRight(w.String(), "\n ")
}
