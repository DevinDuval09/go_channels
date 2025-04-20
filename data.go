package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type Q struct {
	Number  uint32
	Qsize   uint32
	Message string
}

func NewQ(buff []byte) (Q, error) {
	var q Q
	err := q.UnmarshalBinary(buff)

	return q, err
}

func (q Q) MarshallBinary() ([]byte, error) {

	//figure out buffer size
	intSize := 4
	intBufferSize := (intSize * 2) + len(q.Message)
	buff := make([]byte, intBufferSize)
	binary.BigEndian.PutUint32(buff, q.Number)
	q.Qsize = uint32(len(q.Message))
	binary.BigEndian.PutUint32(buff[intSize:], q.Qsize)
	copy(buff[8:], []byte(q.Message))
	return buff, nil
}

func (q *Q) UnmarshalBinary(data []byte) error {
	//b := bytes.NewBuffer((data))
	log.Println("Incoming unmarshal buffer: ", data)
	numSize := 4
	var qnumber uint32
	var messagesize uint32
	var message string
	//err := binary.Read(b, binary.BigEndian, qnumber)
	//if err != nil {
	//	return fmt.Errorf("failed to ready query number: %s", err)
	//}
	log.Println("q data size: ", len(data))
	copyBuff := make([]byte, len(data))
	copy(copyBuff, data)
	qnumber = binary.BigEndian.Uint32(copyBuff[:numSize])
	log.Println("Parsed query number: ", qnumber)
	messagesize = binary.BigEndian.Uint32(copyBuff[numSize:(numSize * 2)])
	log.Println("Parsed message size: ", messagesize)
	//stringbuffer := make([]byte, messagesize)
	//b.Next(numSize * 2)
	//size, err := io.ReadFull(b, stringbuffer)
	//if size != int(messagesize) {
	//	return fmt.Errorf("size mismatch parsing query message: %d, %d", messagesize, size)
	//}

	message = string(copyBuff[(numSize * 2):]) //string(stringbuffer[:size])

	q.Number = qnumber
	q.Qsize = messagesize
	q.Message = message

	log.Println("Unmarshaled q: ", q)

	return nil
}

type R struct {
	Number  uint32
	Qnumber uint32
	//RMessageSize int
	Message string
	//QMessageSize int
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
	return R{Message: message, Number: q.Number, Qmessage: q.Message, Qnumber: q.Number}
}
