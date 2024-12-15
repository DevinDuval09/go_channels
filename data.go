package main

import (
	"bytes"
	"fmt"
)

type Q struct {
	Message string
	Number  int
}

func (q Q) MarshallBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, q.Message, q.Number)
	return b.Bytes(), nil
}

func (q *Q) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer((data))

	_, err := fmt.Fscanln(b, &q.Message, &q.Number)

	return err
}

type R struct {
	Message  string
	Number   int
	Qmessage string
	Qnumber  int
}

func (r R) MarshallBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, r.Message, r.Number, r.Qmessage, r.Qnumber)
	return b.Bytes(), nil
}

func (r *R) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer((data))

	_, err := fmt.Fscanln(b, &r.Message, &r.Number, &r.Qmessage, &r.Qnumber)

	return err
}

func (q Q) Response(message string) R {
	return R{Message: message, Number: q.Number + 1, Qmessage: q.Message, Qnumber: q.Number}
}
