package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	// "io/ioutil"
	// "os"
	"strings"
)

type Data struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

func main() {
	file, err := os.Open("users.txt")
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(fileContents), "\n")
	regEmail := regexp.MustCompile("@")
	regAndroid := regexp.MustCompile("Android")
	regMSIE := regexp.MustCompile("MSIE")
	seenBrowsers := map[string]bool{}
	// uniqueBrowsers := 0
	foundUsers := ""

	for i, line := range lines {
		d := &Data{}
		data := []byte(line)
		err := json.Unmarshal(data, d)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// fmt.Printf("struct:\n\t%#v\n\n", d)

		isAndroid := false
		isMSIE := false

		for _, browser := range d.Browsers {
			if str := regAndroid.FindString(browser); str != "" {
				isAndroid = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = true
				}
			}
			if str := regMSIE.FindString(browser); str != "" {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = true
				}
			}
		}
		if !(isAndroid && isMSIE) {

			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := regEmail.ReplaceAllString(d.Email, " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, d.Name, email)
	}

	fmt.Fprintln(os.Stdout, "found users:\n"+foundUsers)
	fmt.Fprintln(os.Stdout, "Total unique browsers", len(seenBrowsers))
}
