package helper

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

// sessionStorage to store all sessions
var sessionStorage []session

// session is to define structure
type session struct {
	timestamp int64
	remoteIP  string
	token     string
	name      string
	id        string
}

func ValidateSession(id string, r *http.Request) (string, bool) {
	// Check in sessionStorage
	for i, eachSession := range sessionStorage {
		if eachSession.id == id && strings.Split(r.RemoteAddr, ":")[0] == eachSession.remoteIP {
			// Check time expire
			nowUnix := time.Now().Unix()
			if (nowUnix - eachSession.timestamp) <= 24*60*60*1000 { // Expires in a day
				return eachSession.name, true
			}
			delSession(i) // Delete outdated session
			return "", false
		}
	}
	return "", false
}

func AddSession(name string, token string, r *http.Request) string {
	var s session
	s.timestamp = time.Now().Unix()
	s.remoteIP = strings.Split(r.RemoteAddr, ":")[0]
	s.token = token
	s.name = name
	// Calculate sessionID
	obj := md5.New()
	bin := make([]byte, 8)
	binary.LittleEndian.PutUint64(bin, uint64(s.timestamp))
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
