// Borrowed from https://github.com/anoopengineer/winscppasswd. Thank you!

package main

import "strconv"

const (
	PW_MAGIC = 0xA3
	PW_FLAG  = 0xFF
)

func Decrypt(host, username, password string) string {
	key := username + host
	passbytes := []byte{}
	for i := 0; i < len(password); i++ {
		val, _ := strconv.ParseInt(string(password[i]), 16, 8)
		passbytes = append(passbytes, byte(val))
	}
	var flag byte
	flag, passbytes = decNextChar(passbytes)
	var length byte = 0
	if flag == PW_FLAG {
		_, passbytes = decNextChar(passbytes)

		length, passbytes = decNextChar(passbytes)
	} else {
		length = flag
	}
	toBeDeleted, passbytes := decNextChar(passbytes)
	passbytes = passbytes[toBeDeleted*2:]

	clearpass := ""
	var (
		i   byte
		val byte
	)
	for i = 0; i < length; i++ {
		val, passbytes = decNextChar(passbytes)
		clearpass += string(val)
	}

	if flag == PW_FLAG {
		clearpass = clearpass[len(key):]
	}
	return clearpass
}

func decNextChar(passbytes []byte) (byte, []byte) {
	if len(passbytes) <= 0 {
		return 0, passbytes
	}
	a := passbytes[0]
	b := passbytes[1]
	passbytes = passbytes[2:]
	return ^(((a << 4) + b) ^ PW_MAGIC) & 0xff, passbytes
}
