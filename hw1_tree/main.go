package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out *os.File, path string, printFiles bool) error {
	if printFiles {
		return WalkAll(out, path, "")
	}
	return WalkDir(out, path, "")
}

func WalkDir(out *os.File, path string, sufix string) error {
	dirname := path + string(filepath.Separator)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dirs := []os.FileInfo{}
	for _, dir := range fi {
		if dir.IsDir() {
			dirs = append(dirs, dir)
		}
	}
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})
	prefix := "├───"
	for i, dir := range dirs {
		if i == len(dirs)-1 {
			prefix = "└───"
		}
		fmt.Fprint(out, sufix+prefix+dir.Name()+"\n")
		if i != len(dirs)-1 {
			WalkDir(out, dirname+dir.Name(), sufix+"│	")
		} else {
			WalkDir(out, dirname+dir.Name(), sufix+"	")
		}

	}
	return nil
}

func WalkAll(out *os.File, path string, sufix string) error {
	dirname := path + string(filepath.Separator)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sort.Slice(fi, func(i, j int) bool {
		return fi[i].Name() < fi[j].Name()
	})
	prefix := "├───"
	for i, fix := range fi {
		if i == len(fi)-1 {
			prefix = "└───"
		}

		fmt.Fprint(out, sufix+prefix+fix.Name()+" ")
		if !fix.IsDir() {
			fmt.Fprint(out, bytesToString(fix.Size()))
		}
		fmt.Fprintln(out)
		if fix.IsDir() {
			if i != len(fi)-1 {
				WalkAll(out, dirname+fix.Name(), sufix+"│	")
			} else {
				WalkAll(out, dirname+fix.Name(), sufix+"	")
			}
		}
	}
	return nil
}

func bytesToString(bytes int64) string {
	if bytes == 0 {
		return "(empty)"
	}
	result := "("
	result += strconv.FormatInt(bytes, 10)
	result += "b)"
	return result
}
