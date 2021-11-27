package mcbe

import (
	"bytes"
	"encoding/binary"
	"net"
	"strings"
	"time"
)

func Query(address string) (map[string]string, error) {
	conn, err := net.Dial("udp", address)

	if err != nil {
		return nil, err
	}
	if err := conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
		return nil, err
	}
	defer conn.Close()
	var magic = []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78}
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, byte(0x01))
	binary.Write(&buf, binary.BigEndian, time.Now().Unix())
	binary.Write(&buf, binary.BigEndian, magic)
	binary.Write(&buf, binary.BigEndian, []byte{0, 0, 0, 0, 0, 0, 0, 0})

	if _, err = conn.Write(buf.Bytes()); err != nil {
		return nil, err
	}

	var b = make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		return nil, err
	}
	b = b[:n]
	splitted := strings.Split(string(b), ";")
	return map[string]string{
		"motd":             splitted[1],
		"protocol_version": splitted[2],
		"game_version":     splitted[3],
		"online_players":   splitted[4],
		"max_players":      splitted[5],
		"server_guid":      splitted[6],
		"default_world":    splitted[7],
		"game_mode":        splitted[8],
	}, nil
}
