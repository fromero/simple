package simple_sessions

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

const sessionIdName = "session-id"
const oneHour = 3600
const sessionsDirectory = "./sessions/"

type Session struct {
	Values     map[string]interface{}
	LastAccess time.Time
}

var sessions = map[string]*Session{}

func getSession(sessionId string) *Session {
	sess, ok := sessions[sessionId]
	if !ok || expired(*sess) {
		sess = &Session{Values: map[string]interface{}{} /*, LastAccess: time.Now()*/}
		sessions[sessionId] = sess
	}
	sess.LastAccess = time.Now()
	return sess
}

func expired(sess Session) bool {
	seconds := time.Now().Sub(sess.LastAccess).Seconds()
	return seconds > oneHour
}

func getSessionId(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(sessionIdName)
	if err == http.ErrNoCookie {
		sessionIdValue := uuid.New()
		cookie = &http.Cookie{Name: sessionIdName, Value: sessionIdValue.String()}
		http.SetCookie(w, cookie)
	}
	return cookie.Value
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionId(w, r)
	sess := getSession(sessionId)
	counter, ok := sess.Values["counter"]
	if !ok {
		counter = 0
	}
	counter = counter.(int) + 1
	sess.Values["counter"] = counter
	go persistSession(sessionId, sess)
	fmt.Fprintf(w, strconv.Itoa(counter.(int)))
}

func persistSession(sessionId string, sess *Session) {
	err := writeGob(sessionsDirectory+sessionId, sess)
	if err != nil {
		log.Println("Can't save session: " + err.Error())
	}
}

func main() {
	http.HandleFunc("/", basicHandler)

	//fs := http.FileServer(http.Dir("static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}
