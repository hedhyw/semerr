package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v2"
)

const (
	extTmpl = ".tmpl"
	extGo   = ".go"
)

const (
	timeout = 5 * time.Second
)

//go:embed "errors.yaml"
var errorsYAMLData []byte

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("invalid count of arguments: %d", len(os.Args))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)

		err := generate(os.Args[1])
		if err != nil {
			log.Fatalf("generating: %v", err)
		}
	}()

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}

type errorDefinition struct {
	Name string `yaml:"-"`

	GRPCStatus  int    `yaml:"grpc"`
	HTTPStatus  int    `yaml:"http"`
	Description string `yaml:"description"`
	Temporary   bool   `yaml:"temporary"`
	Reverse     bool   `yaml:"reverse"`
}

func walkFn(errDefs []errorDefinition) fs.WalkDirFunc {
	return func(path string, _ fs.DirEntry, lastErr error) (err error) {
		if lastErr != nil {
			return lastErr
		}

		if filepath.Ext(path) != extTmpl {
			return nil
		}

		name := filepath.Base(path)
		baseDir := filepath.Dir(path)

		tmpl, err := template.New(name).Funcs(template.FuncMap{
			"multlineComment": multlineComment,
			"httpStatusText":  http.StatusText,
			"grpcStatusText":  grpcStatusText,
		}).ParseFiles(path)
		if err != nil {
			return fmt.Errorf("parsing template: %s: %w", name, err)
		}

		outFilePath := strings.TrimSuffix(name, extTmpl)
		outFile := filepath.Join(
			baseDir,
			outFilePath,
		)

		f, err := os.OpenFile(
			outFile,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			os.ModePerm,
		)
		if err != nil {
			return fmt.Errorf("opening out file: %w", err)
		}

		defer func() { err = semerr.NewMultiError(err, f.Close()) }()

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, errDefs)
		if err != nil {
			return fmt.Errorf("executing tmpl: %s: %w", name, err)
		}

		source := buf.Bytes()

		if filepath.Ext(outFilePath) == extGo {
			source, err = format.Source(source)
			if err != nil {
				return fmt.Errorf("formating code: %w", err)
			}
		}

		_, err = f.Write(source)
		if err != nil {
			return fmt.Errorf("writing buf to file: %s: %w", name, err)
		}

		log.Print(outFile, " generated")

		return nil
	}
}

func generate(path string) (err error) {
	var errDefs struct {
		Errors map[string]errorDefinition `yaml:"errors"`
	}

	err = yaml.Unmarshal(errorsYAMLData, &errDefs)
	if err != nil {
		return fmt.Errorf("decoding definitions: %w", err)
	}

	defs := make([]errorDefinition, 0, len(errDefs.Errors))
	for name, def := range errDefs.Errors {
		def.Name = name
		defs = append(defs, def)
	}

	sort.Slice(defs, func(i, j int) bool {
		left, right := defs[i], defs[j]

		if left.GRPCStatus == right.GRPCStatus {
			return left.HTTPStatus < right.HTTPStatus
		}

		return left.GRPCStatus < right.GRPCStatus
	})

	err = filepath.WalkDir(path, walkFn(defs))
	if err != nil {
		return fmt.Errorf("walking dir: %w", err)
	}

	return nil
}

func multlineComment(text string) string {
	sb := strings.Builder{}
	sb.Grow(len(text))

	ss := bufio.NewScanner(strings.NewReader(text))
	for ss.Scan() {
		if sb.Len() != 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString("// ")
		sb.WriteString(ss.Text())
	}

	return sb.String()
}

func grpcStatusText(code int) string {
	return codes.Code(code).String()
}
