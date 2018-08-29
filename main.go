package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha256"

	"github.com/labstack/echo"
)

var (
	ip       = flag.String("ip", "localhost", "server ip address")
	port     = flag.String("port", "20000", "server port to listen on")
	maxDelay = flag.Int("max-delay", 1000, "largest delay in loop compare (ms)")
	maxUser  = flag.Int("max-user", 8, "largest size of a user name")
	tagLen   = flag.Int("tag-len", 4, "number of bytes in a tag")
	secret   = flag.String("secret", "secret", "secret used to produce tags")
	logFile  = flag.String("log-file", "", "path which is logged to (empty=off)")
)

const (
	path = "/auth/:delay/:user/:tag"
)

func main() {
	flag.Parse()

	// log to file and stdout
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(io.MultiWriter(os.Stdout, f))
	}

	// serve clients
	e := echo.New()
	e.GET(path, func(c echo.Context) error {
		d, u, t, err := parseArgs(c.Param("delay"), c.Param("user"), c.Param("tag"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err = authenticate(d, u, t); err != nil {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		log.Printf("[Success] %d,%s,%s,%s", d, u, c.RealIP(), agent(c))
		return c.String(http.StatusOK, "Access granted!\n")
	})
	e.Logger.Fatal(e.Start(*ip + ":" + *port))
}

func authenticate(delay int, user, tag []byte) error {
	mac := hmac.New(sha256.New, mkKey(*secret, delay))
	mac.Write(user)
	for i, b := range mac.Sum(nil)[:*tagLen] {
		if b != tag[i] {
			return fmt.Errorf("tag mis-match\n")
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	return nil
}

func parseArgs(delay, user, tag string) (d int, u, t []byte, err error) {
	d, err = strconv.Atoi(delay)
	if err != nil {
		err = fmt.Errorf("%v\n", err.Error())
		return
	}
	if d <= 0 || d > *maxDelay {
		err = fmt.Errorf("expected delay in (0, %v], got %v\n", *maxDelay, d)
		return
	}

	u = []byte(user)
	if n := len(u); n <= 0 || n > *maxUser {
		err = fmt.Errorf("expected user of length (0, %v], got %v\n", *maxUser, n)
		return
	}

	n := hex.DecodedLen(len(tag))
	if n != *tagLen {
		err = fmt.Errorf("expected tag of length %v, got %v\n", *tagLen, n)
		return
	}
	t = make([]byte, n)
	n, err = hex.Decode(t, []byte(tag))
	if err != nil {
		err = fmt.Errorf("expected hex tag, got bad byte on position %v\n", n)
	}
	return
}

func agent(c echo.Context) string {
	a, ok := c.Request().Header["User-Agent"]
	if !ok {
		return "unkown"
	}
	return strings.Join(a, " ")
}

func mkKey(secret string, delay int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(delay))
	return append([]byte(secret), b...)
}
