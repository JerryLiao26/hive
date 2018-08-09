package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// sessionStorage to store all sessions
var sessionStorage []session

// session is to define structure
type session struct {
	timeStamp int64
	remoteIP  string
	token     string
	id        string
}

func validateSession(id string, r *http.Request) bool {
	// Check in sessionStorage
	for i, eachSession := range sessionStorage {
		if eachSession.id == id && strings.Split(r.RemoteAddr, ":")[0] == eachSession.remoteIP {
			// Check time expire
			nowUnix := time.Now().Unix()
			fmt.Println("minus:", nowUnix-eachSession.timeStamp)
			if (nowUnix - eachSession.timeStamp) <= 24*60*60*1000 { // Expires in a day
				return true
			}
			delSession(i) // Delete outdated session
			return false
		}
	}
	return false
}

func addSession(token string, r *http.Request) string {
	var s session
	s.timeStamp = time.Now().Unix()
	s.remoteIP = strings.Split(r.RemoteAddr, ":")[0]
	s.token = token
	// Calculate sessionID
	obj := md5.New()
	bin := make([]byte, 8)
	binary.LittleEndian.PutUint64(bin, uint64(s.timeStamp))
	obj.Write(bin)
	obj.Write([]byte(s.remoteIP))
	obj.Write([]byte(s.token))
	s.id = hex.EncodeToString(obj.Sum(nil))
	sessionStorage = append(sessionStorage, s)

	return s.id
}

func delSession(i int) {
	sessionStorage = append(sessionStorage[:i], sessionStorage[i+1:]...)
}
