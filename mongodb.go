package middleware


import (
"net/http"

"github.com/codegangsta/negroni"
	"gopkg.in/mgo.v2"
	"context"
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
	session, err := mgo.Dial(url)
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
