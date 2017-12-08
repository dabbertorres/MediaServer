package websrv

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"io"
	"time"
)

const (
	newSessionCmd     = "insert into sessions (id, user, ipAddr, userAgent, expires, tunedTo) values ($1, $2, $3, $4, $5, $6)"
	isSessionValidCmd = "select expires from sessions where id = $1"
	getSessionCmd     = "select * from sessions where id = $1"
	endSessionCmd     = "delete from sessions where id = $1"
)

type SessionId string

type Session struct {
	Id        SessionId `json:"id"`
	User      string    `json:"username"`
	IPAddr    string    `json:"ipAddr"`
	UserAgent string    `json:"userAgent"`
	Expires   time.Time `json:"expires"`
	TunedTo   string    `json:"tunedTo"`
}

func NewSession(db *sql.DB, user, ipAddr, userAgent string, lifetime time.Duration) (SessionId, error) {
	hash := sha512.New()
	hash.Write([]byte(user))
	hash.Write([]byte(ipAddr))
	hash.Write([]byte(userAgent))

	entropy := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, entropy)
	if err != nil {
		return "", err
	}

	id := base64.StdEncoding.EncodeToString(hash.Sum(entropy))

	sess := Session{
		Id:        SessionId(id),
		User:      user,
		IPAddr:    ipAddr,
		UserAgent: userAgent,
		Expires:   time.Now().Add(lifetime),
	}

	_, err = db.Exec(newSessionCmd,
		&sess.Id,
		&sess.User,
		&sess.IPAddr,
		&sess.UserAgent,
		&sess.Expires,
		&sess.TunedTo)
	if err != nil {
		return "", err
	}

	return sess.Id, nil
}

func IsSessionValid(db *sql.DB, id SessionId) (bool, error) {
	res := db.QueryRow(isSessionValidCmd, &id)

	expires := time.Time{}
	if err := res.Scan(&expires); err != nil {
		return false, err
	}

	// TODO also check if the user for the queried session actually exists!

	return time.Now().Before(expires), nil
}

func GetSession(db *sql.DB, id SessionId) (*Session, error) {
	res := db.QueryRow(getSessionCmd, &id)

	sess := Session{}
	err := res.Scan(&sess.Id, &sess.User, &sess.IPAddr, &sess.UserAgent, &sess.Expires, &sess.TunedTo)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func EndSession(db *sql.DB, id SessionId) error {
	_, err := db.Exec(endSessionCmd, &id)
	return err
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.Expires)
}
