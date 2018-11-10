package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	bashColorRegex = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]`)
	emailRegex     = regexp.MustCompile("([a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)")
	conn, _        = net.Dial("udp", "127.0.0.1:1234")
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	f, err := os.OpenFile("log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Println(txt)

		s := bashColorRegex.ReplaceAllString(maskEmail(txt), "")
		fmt.Fprintln(conn, s)
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func maskEmail(txt string) string {
	matches := emailRegex.FindAllStringSubmatch(txt, -1)
	for _, m := range matches {
		if len(m) > 0 {
			email := m[0]
			parts := strings.Split(email, "@")
			masked := strings.Repeat("*", len(parts[0])) + "@" + parts[1]
			txt = strings.Replace(txt, email, masked, -1)
		}
	}
	return txt
}
