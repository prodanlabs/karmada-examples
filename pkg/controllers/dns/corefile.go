package dns

import (
	"bufio"
	"fmt"
	"strings"

	"k8s.io/klog/v2"

	"github.com/prodanlabs/karmada-examples/pkg/util"
)

type Corefile struct {
	lines []string
}

func NewCorefile(config string) *Corefile {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(config))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return &Corefile{
		lines: lines,
	}
}

func (c *Corefile) trimSliceSpaces(oldSlice []string) []string {
	var newSlice []string
	for i := range oldSlice {
		if oldSlice[i] != "" {
			newSlice = append(newSlice, oldSlice[i])
		}
	}

	return newSlice
}

func (c *Corefile) isExists(hostname string) bool {
	for i := range c.lines {
		if strings.Contains(c.lines[i], hostname) {
			return true
		}
	}

	return false
}

func (c *Corefile) add(ip, hostname string) []byte {
	var insertIndex int
	for i, line := range c.lines {
		if strings.Contains(line, "hosts") {
			insertIndex = i + 1 // insert on the line below the keyword
			break
		}
	}

	// insert new row
	newLine := fmt.Sprintf("%s %s", ip, hostname)
	c.lines = append(c.lines[:insertIndex], append([]string{newLine}, c.lines[insertIndex:]...)...)

	return util.Format([]byte(strings.Join(c.lines, "\n")))
}

func (c *Corefile) update(ip, hostname string) []byte {
	for i, line := range c.lines {
		if strings.Contains(line, hostname) {
			l := strings.Split(c.lines[i], " ")
			host := c.trimSliceSpaces(l)
			if host[0] == ip {
				klog.V(6).Infof("The current %q is new for the corresponding %q", hostname, ip)
				return nil
			}

			c.lines[i] = strings.ReplaceAll(line, host[0], ip)
		}
	}

	return util.Format([]byte(strings.Join(c.lines, "\n")))
}

func (c *Corefile) AddOrUpdate(ip, hostname string) []byte {
	if !c.isExists(hostname) {
		return c.add(ip, hostname)
	}
	return c.update(ip, hostname)
}

func (c *Corefile) Delete(hostname string) []byte {
	// Iterate over the lines and remove the ones containing the hostname
	var newLines []string
	for _, line := range c.lines {
		if !strings.Contains(line, hostname) {
			newLines = append(newLines, line)
		}
	}

	// Join the remaining lines back into a single string
	newString := strings.Join(newLines, "\n")

	return util.Format([]byte(newString))
}
