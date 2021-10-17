package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

const (
	extTmpl = ".tmpl"
	extGo   = ".go"
)

const (
	timeout = 5 * time.Second
)

//go:embed "errors.json"
var errorsJSONData []byte

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
	GRPCStatus  int    `json:"grpc"`
	HTTPStatus  int    `json:"http"`
	Description string `json:"description"`
	Permanent   bool   `json:"permanent"`
}

func walkFn(errDefs map[string]errorDefinition) fs.WalkDirFunc {
	return func(path string, _ fs.DirEntry, lastErr error) (err error) {
		if lastErr != nil {
			return lastErr
		}

		if filepath.Ext(path) != extTmpl {
			return nil
		}

		name := filepath.Base(path)
		baseDir := filepath.Dir(path)

		tmpl, err := template.New(name).ParseFiles(path)
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
	var errDefs map[string]errorDefinition
	err = json.Unmarshal(errorsJSONData, &errDefs)
	if err != nil {
		return fmt.Errorf("decoding definitions: %w", err)
	}

	err = filepath.WalkDir(path, walkFn(errDefs))
	if err != nil {
		return fmt.Errorf("walking dir: %w", err)
	}

	return nil
}
