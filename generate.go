package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

var sources = []string{
	"https://github.com/hagezi/dns-blocklists/raw/refs/heads/main/hosts/native.xiaomi.txt",
	"https://github.com/Kittyskj/FreeFromMi/raw/refs/heads/main/hosts",
}

func main() {

	lineChan := make(chan string, 1024)

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(100)

	for _, url := range sources {
		g.Go(func() error {
			return httpGetLines(ctx, url, func(line string) {
				lineChan <- line
			})
		})
	}

	go func() {
		checkErr(g.Wait())
		close(lineChan)
	}()

	hosts := make([]string, 0)

	for line := range lineChan {
		parseHostsLine(line, func(host string) {
			hosts = append(hosts, host)
		})
	}

	slices.Sort(hosts)
	hosts = unique(hosts)
	domains, ip4s, ip6s := separateHosts(hosts)

	buf := new(bytes.Buffer)
	writeTextRules(buf, domains, ip4s, ip6s)
	checkErr(os.WriteFile("rules.txt", buf.Bytes(), 0644))
	buf.Reset()

	buf.WriteString("payload:\n")
	writeYamlRules(buf, domains, ip4s, ip6s)
	checkErr(os.WriteFile("rules.yaml", buf.Bytes(), 0644))

	fmt.Println("done.")
}

func httpGetLines(
	ctx context.Context,
	url string,
	lineFunc func(line string),
) error {

	fmt.Println("fetching", url)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not make request: %w", err)
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		lineFunc(s.Text())
	}

	err = s.Err()
	if err != nil {
		return fmt.Errorf("could not scan body: %w", err)
	}

	return nil
}

func parseHostsLine(line string, hostFunc func(host string)) {
	firstField := true
	for field := range strings.FieldsSeq(line) {
		if firstField {
			if field[0] == '#' {
				break
			}
			firstField = false
			continue
		}
		hostFunc(field)
	}
}

// unique returns a deduplicated version of the sorted slice s.
// s is overwritten and the result uses the underlying array of s.
func unique[T comparable](s []T) []T {
	if len(s) == 0 {
		return s
	}
	r := s[:1]
	for _, v := range s[1:] {
		if v != r[len(r)-1] {
			r = append(r, v)
		}
	}
	return r
}

func separateHosts(hosts []string) (domains, ip4s, ip6s []string) {
	for _, v := range hosts {
		ip := net.ParseIP(v)
		if ip == nil {
			v = strings.TrimPrefix(v, "www.")
			if !strings.ContainsRune(v, '.') {
				continue
			}
			domains = append(domains, v)
			continue
		}
		if !ip.IsGlobalUnicast() || ip.IsPrivate() {
			continue
		}
		if ip.To4() != nil {
			ip4s = append(ip4s, v)
		} else {
			ip6s = append(ip6s, v)
		}
	}
	return
}

func writeTextRules(w io.StringWriter, domains, ip4s, ip6s []string) error {
	for _, v := range domains {
		if err := writeString(w, "DOMAIN-SUFFIX,", v, "\n"); err != nil {
			return err
		}
	}
	for _, v := range ip4s {
		if err := writeString(w, "IP-CIDR,", v, "/32\n"); err != nil {
			return err
		}
	}
	for _, v := range ip6s {
		if err := writeString(w, "IP-CIDR6,", v, "/128\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeYamlRules(w io.StringWriter, domains, ip4s, ip6s []string) error {
	for _, v := range domains {
		if err := writeString(w, "- 'DOMAIN-SUFFIX,", v, "'\n"); err != nil {
			return err
		}
	}
	for _, v := range ip4s {
		if err := writeString(w, "- 'IP-CIDR,", v, "/32'\n"); err != nil {
			return err
		}
	}
	for _, v := range ip6s {
		if err := writeString(w, "- 'IP-CIDR6,", v, "/128'\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeString(w io.StringWriter, s ...string) error {
	for _, v := range s {
		_, err := w.WriteString(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
