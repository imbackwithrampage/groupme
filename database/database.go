// mautrix-whatsapp - A Matrix-WhatsApp puppeting bridge.
// Copyright (C) 2019 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	log "maunium.net/go/maulogger/v2"
	"maunium.net/go/mautrix-whatsapp/database/upgrades"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
	log     log.Logger
	dialect string

	User    *UserQuery
	Portal  *PortalQuery
	Puppet  *PuppetQuery
	Message *MessageQuery
}

func New(dbType string, uri string, baseLog log.Logger) (*Database, error) {

	var conn gorm.Dialector

	if dbType == "sqlite3" {
		//_, _ = conn.Exec("PRAGMA foreign_keys = ON")
		conn = sqlite.Open(uri)
	} else {
		conn = postgres.Open(uri)
	}
	print("no")
	gdb, err := gorm.Open(conn, &gorm.Config{
		// Logger: baseLog,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db := &Database{
		DB:      gdb,
		log:     baseLog.Sub("Database"),
		dialect: dbType,
	}
	db.User = &UserQuery{
		db:  db,
		log: db.log.Sub("User"),
	}
	db.Portal = &PortalQuery{
		db:  db,
		log: db.log.Sub("Portal"),
	}
	db.Puppet = &PuppetQuery{
		db:  db,
		log: db.log.Sub("Puppet"),
	}
	db.Message = &MessageQuery{
		db:  db,
		log: db.log.Sub("Message"),
	}
	return db, nil
}

func (db *Database) Init() error {
	err := db.AutoMigrate(&Portal{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Puppet{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Message{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&mxRegistered{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&mxUserProfile{})
	if err != nil {
		return err
	}

	return upgrades.Run(db.log.Sub("Upgrade"), db.dialect, db.DB)
}

type Scannable interface {
	Scan(...interface{}) error
}
