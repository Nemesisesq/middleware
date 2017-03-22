package middleware

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"gopkg.in/mgo.v2"
	"context"
	"crypto/tls"
	"net"
)

type Database struct {
	dba DatabaseAccessor
}

func NewDatabase(databaseAccessor DatabaseAccessor) *Database {
	return &Database{databaseAccessor}
}

func (d *Database) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		reqSession := d.dba.Clone()
		defer reqSession.Close()
		d.dba.Set(request, reqSession)
		next(writer, request)
	}
}

type DatabaseAccessor struct {
	*mgo.Session
	url  string
	name string
	coll string
}

func NewDatabaseAccessor(url, name, coll string) (*DatabaseAccessor, error) {

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	dialInfo, err := mgo.ParseURL(url)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err == nil {
		return &DatabaseAccessor{session, url, name, coll}, nil
	} else {
		return &DatabaseAccessor{}, err
	}
}

func (da *DatabaseAccessor) Set(request *http.Request, session *mgo.Session) {
	db := session.DB(da.name)
	context.WithValue(request.Context(), "db", db)
	context.WithValue(request.Context(), "mgoSession", session)
}
