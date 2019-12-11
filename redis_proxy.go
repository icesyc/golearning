package main

import (
	"fmt"
	"log"
	"net"
	"bufio"
	"io"
	"strconv"
	"bytes"
	"os"
	"strings"
)

const (
	MessageStatus = '+'
	MessageError = '-'
	MessageInt = ':'
	MessageBulk = '$'
	MessageMulti = '*'
	Crlf = "\r\n"
)

type Message struct{
	Type byte
	StringValue string
	IntValue int
	MultiValue []*Message
}
type RedisMessage struct{
	Msg *Message 
	Client net.Conn
}

func (self *Message) String() string {
	multiValue := make([]string, len(self.MultiValue))
	for i, msg := range self.MultiValue {
		multiValue[i] = msg.String()
	}
	multiValueStr := strings.Join(multiValue, ",")
	return fmt.Sprintf("Type=%c, StringValue=%s, IntValue=%d, MultiValue=[%s]", self.Type, self.StringValue, self.IntValue, multiValueStr)
}

var redisCh = make(chan *RedisMessage)

func main(){
	sockFile := "/tmp/redis_proxy.sock"
	os.Remove(sockFile)
	sock, err := net.Listen("unix", sockFile)
	if err != nil {
		log.Fatal(err)
	}


	go handleRedisRequest()

	for {
		conn, err := sock.Accept();
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleRedisRequest(){
	redis, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	for redisMsg := range redisCh {
		reqMsgPack, err := pack(redisMsg.Msg)
		if err != nil {
			fmt.Println(err.Error())
			continue;
		}
		_, err = redis.Write(reqMsgPack)
		respMsg, err := unpack(redis)
		fmt.Println("<<"  + respMsg.String())
		if err != nil {
			fmt.Println(err.Error())
			respMsg = &Message{
				Type: MessageError,
				StringValue: err.Error(),
			}
			continue;
		}
		respMsgPack, _ := pack(respMsg)
		redisMsg.Client.Write(respMsgPack)
	}
}

func pack(msg *Message) ([]byte, error) {
	var res bytes.Buffer
	res.WriteByte(msg.Type)
	switch(msg.Type){
	case MessageError:
		res.WriteString(msg.StringValue)
		res.WriteString(Crlf)
	case MessageStatus: 
		res.WriteString(msg.StringValue)
		res.WriteString(Crlf)
	case MessageInt:
		res.WriteString(strconv.Itoa(msg.IntValue))
		res.WriteString(Crlf)
	case MessageBulk: 
		len := len(msg.StringValue)
		lenStr := strconv.Itoa(len)
		res.WriteString(lenStr + Crlf)
		res.WriteString(msg.StringValue + Crlf)
	case MessageMulti:
		len := len(msg.MultiValue)
		lenStr := strconv.Itoa(len)
		res.WriteString(lenStr + Crlf)
		for _, multiMsg := range msg.MultiValue {
			packMsg, err := pack(multiMsg)
			if err != nil {
				return nil, err
			}
			res.Write(packMsg)
		}
	default: 
		return nil, fmt.Errorf("invalid message, msg.Type=%s", msg.Type)
	}
	return res.Bytes(), nil
}

func unpack(c net.Conn) (*Message, error){
	reader := bufio.NewReader(c)
	return unpackFromReader(reader)
}
func unpackFromReader(reader *bufio.Reader) (*Message, error){
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = line[0:len(line)-2]
	switch line[0] {
	case MessageError:
		return &Message{
			Type: MessageError,
			StringValue: line[1:],
		}, nil
	case MessageStatus:
		return &Message{
			Type: MessageStatus,
			StringValue: line[1:],
		}, nil
	case MessageInt:
		v, err := strconv.Atoi(line[1:]);
		if err != nil {
			return nil, err
		}
		return &Message{
			Type: MessageInt,
			IntValue: v,
		}, nil
	case MessageBulk:
		len, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		if len < 0 {
			return &Message{
				Type: MessageBulk,
			}, nil
		}
		buf := make([]byte, len + 2)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, err
		}
		return &Message{
			Type: MessageBulk,
			StringValue: string(buf[:len]),
		}, nil
	case MessageMulti:
		msgNum, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		if msgNum < 0 {
			return &Message{
				Type: MessageMulti,
			}, nil
		}
		msgList := make([]*Message, msgNum);
		for i := 0; i < msgNum; i++ {
			msg, err := unpackFromReader(reader)
			if err != nil {
				return nil, err
			}
			msgList[i] = msg
		}
		return &Message{
			Type: MessageMulti,
			MultiValue: msgList,
		}, nil
	default: 
		return nil, fmt.Errorf("invalid message, msg.Type=%s", line[0])
	}
}
func handleConnection(c net.Conn){
	defer c.Close()
	for {
		reqMsg, err := unpack(c)
		if err != nil {
			fmt.Println(err.Error())
			break;
		}
		fmt.Println(">>" + reqMsg.String())
		redisCh <- &RedisMessage{Msg: reqMsg, Client: c}
	}
}