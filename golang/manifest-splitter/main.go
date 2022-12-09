package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	// read multi document Kubernetes manifests
	yfile, err := os.ReadFile("install.yaml")
	if err != nil {
		log.Fatal(err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(yfile))
	for {
		var manifest map[string]interface{}
		if dec.Decode(&manifest) != nil {
			break
		}
		b, err := yaml.Marshal(manifest)
		if err != nil {
			panic(err)
		}
		if len(manifest) == 0 {
			fmt.Println("Warning: detected empty YAML doc - skipping!")
			continue
		}
		filename := fmt.Sprintf("%s.yaml", strings.ToLower(manifest["kind"].(string)))
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			err = os.WriteFile(filename, b, 0644)
			if err != nil {
				panic(err)
			}
		} else {
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = f.WriteString("---\n")
			if err != nil {
				panic(err)
			}
			_, err = f.Write(b)
			if err != nil {
				panic(err)
			}
		}
	}
}
