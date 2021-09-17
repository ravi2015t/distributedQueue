package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/ravi2015t/distributedQueue/client"
)

const maxN = 10000000
const maxSize = 1024 * 1024

func main() {
	s := client.NewBasic([]string{"localhost"})
	want, err := Send(s)
	if err != nil {
		log.Fatalf("Send error: %v", err)
	}

	got, err := Receive(s)
	if err != nil {
		log.Fatalf("Receive error: %v", err)
	}

	if want != got {
		log.Fatalf("The send and received sum is not equal; want %v got %v", want, got)
	}

}

func Send(s *client.Basic) (sum int64, err error) {
	var b bytes.Buffer
	for i := 0; i <= maxN; i++ {
		sum += int64(i)
		fmt.Fprintf(&b, "%d\n", i)

		if b.Len() >= maxSize {
			if err := s.Send(b.Bytes()); err != nil {
				return 0, err
			}
			b.Reset()
		}
	}

	if b.Len() > 0 {
		if err := s.Send(b.Bytes()); err != nil {
			return 0, err
		}
	}
	return sum, nil
}

func Receive(s *client.Basic) (sum int64, err error) {

	buf := make([]byte, maxSize)
	for {

		res, err := s.Receive(buf)
		if err == io.EOF {
			return sum, nil
		} else if err != nil {
			return 0, err
		}

		ints := strings.Split(string(res), "\n")

		for _, str := range ints {
			if str == "" {
				continue
			}
			i, err := strconv.Atoi(str)
			if err != nil {
				return 0, err
			}
			sum += int64(i)
		}
	}
}
