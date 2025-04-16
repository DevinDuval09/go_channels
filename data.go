package main

import (
	"bytes"
	"fmt"
	"log"
)

type Q struct {
	Number  int
	Message string
}

func NewQ(buff []byte) (Q, error) {
	var q Q
	err := q.UnmarshalBinary(buff)

	return q, err
}

func (q Q) MarshallBinary() ([]byte, error) {

	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, q.Number, q.Message)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (q *Q) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer((data))

	log.Println("Empty q: ", q)

	//log.Println("Received buffer: ", b)

	n, err := fmt.Fscanln(b, &q.Number, &q.Message)

	log.Println("Unmarshaled q: ", q, " size ", n)

	return err
}

type R struct {
	Number   int
	Qnumber  int
	Message  string
	Qmessage string
}

func (r *R) MarshallBinary() ([]byte, error) {
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, r.Message, r.Number, r.Qmessage, r.Qnumber)
	if err != nil {
		return nil, err
	}
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
