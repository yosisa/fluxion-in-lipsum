// This code borrowed from https://golang.org/doc/codewalk/markov/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

type Prefix []string

func (p Prefix) String() string {
	return strings.Join(p, " ")
}

func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type Chain struct {
	chain     map[string][]string
	prefixLen int
}

func newChain(prefixLen int) *Chain {
	return &Chain{
		chain:     make(map[string][]string),
		prefixLen: prefixLen,
	}
}

func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			return
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func NewChain(text []byte) *Chain {
	c := newChain(2)
	for _, p := range bytes.Split(bytes.Trim(text, "\n"), []byte{'\n'}) {
		for _, s := range bytes.Split(p, []byte{'.'}) {
			s = bytes.Trim(s, " ")
			if len(s) == 0 {
				continue
			}
			c.Build(bytes.NewReader(append(s, '.')))
		}
	}
	return c
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
