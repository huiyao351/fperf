package main

import (
	"bytes"
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type option struct {
	output string
}

func gobuild(output string, imports []string) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("package main\n")
	for _, imp := range imports {
		buf.WriteString(`import _ "` + imp + `"` + "\n")
	}

	if err := ioutil.WriteFile("autoimports.go", buf.Bytes(), 0655); err != nil {
		log.Fatalln(err)
	}
	defer os.Remove("autoimports.go")

	cmd := exec.Command("go", "build", "-o", output)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	o := option{}
	flag.StringVar(&o.output, "o", "fperf", "build output")
	flag.Parse()

	paths := flag.Args()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	imports := make([]string, len(paths))
	for i := range paths {
		path := paths[i]
		if filepath.IsAbs(path) {
			var err error
			path, err = filepath.Rel(cwd, path)
			if err != nil {
				log.Fatalln(err)
			}
		}
		p, err := build.Import(path, cwd, 0)
		if err != nil {
			log.Fatalln(err)
		}
		imports[i] = p.ImportPath
	}
	if err := gobuild(o.output, imports); err != nil {
		log.Fatalln(err)
	}
}
