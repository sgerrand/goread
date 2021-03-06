/*
 * Copyright (c) 2013 Matt Jibson <matt.jibson@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package goapp

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"io"
	"regexp"
	"strings"
)

func Sanitize(s string) (string, string) {
	r := bytes.NewReader([]byte(s))
	z := html.NewTokenizer(r)
	buf := &bytes.Buffer{}
	snip := &bytes.Buffer{}
	scripts := 0
	for {
		if z.Next() == html.ErrorToken {
			if err := z.Err(); err == io.EOF {
				break
			} else {
				return s, snipper(s)
			}
		}
		t := z.Token()
		if t.DataAtom == atom.Script {
			if t.Type == html.StartTagToken {
				scripts++
			} else if t.Type == html.EndTagToken {
				scripts--
			}
		} else if scripts == 0 {
			buf.WriteString(t.String())
			if t.Type == html.TextToken {
				snip.WriteString(t.String())
			}
		}
	}

	return buf.String(), snipper(snip.String())
}

const snipLen = 100

var snipRe = regexp.MustCompile("[\\s]+")

func snipper(s string) string {
	s = snipRe.ReplaceAllString(strings.TrimSpace(s), " ")
	s = html.UnescapeString(s)
	if len(s) <= snipLen {
		return s
	}
	s = s[:snipLen]
	i := strings.LastIndexAny(s, " .-!?")
	if i != -1 {
		return s[:i]
	}
	return s
}
