package handler

import (
	"database/sql"

	"paylist.server/infra/store/postgres/store"
)

type BaseHandler struct {
	Db    *sql.DB
	Store store.Storage
}
