package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func take(addr, name string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/take?name=%s", addr, url.QueryEscape(name)))
	if err != nil {
		return err
	}

	// Copy one byte at a time to show messages in realtime and avoid
	// missed messages.
	buff := make([]byte, 1)
	for {
		nr, er := resp.Body.Read(buff)
		if nr > 0 {
			nw, ew := os.Stdout.Write(buff[0:nr])
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return err
}

func give(addr, name string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/give?name=%s", addr, url.QueryEscape(name)))
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

func message(addr, name, msg string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/message?name=%s&message=%s",
		addr,
		url.QueryEscape(name),
		url.QueryEscape(msg)),
	)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

func dsControl(addr, name, state string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/ds?name=%s&state=%s",
		addr,
		url.QueryEscape(name),
		url.QueryEscape(state)),
	)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}
