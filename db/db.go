package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	db *sql.DB
}

func CheckDB() *sql.DB {
	_, err := os.Stat("db/database-sqlite.db")
	if os.IsNotExist(err) {
		createFile()
	}
	var d database
	d.open("db/database-sqlite.db")
	d.createTable()
	return d.db
}

func createFile() {
	file, err := os.Create("db/database-sqlite.db")
	if err != nil {
		log.Fatalf("file doesn't create %v", err)
	}
	defer file.Close()
}

func (d *database) open(file string) {
	var err error
	d.db, err = sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("this error is in db/open() %v", err)
	}
}

func (d *database) createTable() {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS users
	(
		id    INTEGER NOT NULL UNIQUE,
		login    TEXT NOT NULL UNIQUE,
		password    	BLOB NOT NULL,
		email    TEXT NOT NULL UNIQUE,
		PRIMARY KEY(id AUTOINCREMENT)
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().users!!! %v", err)
	}

	// _, err = d.db.Exec(`INSERT INTO users (id, login, password, email)
	// VALUES (1, 'zangar', '12345678', 'zangarzaynesh'),
	// (2, 'batyr', '12345678', 'batyrbatyrovich'),
	// (3, 'magzhan', '12345678', 'magzhanmagzhanovich'),
	// (4, 'nurlan', '12345678', 'nurlannurlanovich');`)

	// if err != nil {
	// 	log.Fatalf("This error is in db.d.InsertIntoUsers!!! %v", err)
	// }

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS categories (
        "id"    INTEGER NOT NULL UNIQUE,
        "name"    TEXT NOT NULL UNIQUE,
        PRIMARY KEY("id")
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().categories!!! %v", err)
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS likes (
        "id"    INTEGER NOT NULL UNIQUE,
        "name"    TEXT NOT NULL UNIQUE,
        PRIMARY KEY("id")
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().likes!!! %v", err)
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS posts (
        "id"    INTEGER NOT NULL UNIQUE,
        "date"    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "user_id"    INTEGER NOT NULL,
        "category_id"    INTEGER,
        "post"    TEXT NOT NULL,
        FOREIGN KEY("user_id") REFERENCES "users"("id"),
        FOREIGN KEY("category_id") REFERENCES "categories"("id"),
        PRIMARY KEY("id" AUTOINCREMENT)
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().posts!!! %v", err)
	}

	// _, err = d.db.Exec(`INSERT INTO posts (id, user_id, category_id, post)
	// VALUES (11, 1, 2, 'fourth post zangar'),
	// (2, 1, 1, 'second post zangar');`)
	// // (3, 2, 3, 'first post batyr'),
	// // (4, 2, 1, 'second post batyr'),
	// // (5, 3, 3, 'first post magzhan'),
	// // (6, 3, 2, 'second post magzhan'),
	// // (7, 4, 3, 'first post nurlan'),
	// // (8, 4, 3, 'second post nurlan'),
	// // (9, 4, 6, 'third post nurlan'),
	// // (10, 1, 5, 'third post zangar');`)

	// if err != nil {
	// 	log.Fatalf("This error is in db.d.InsertIntoPosts!!! %v", err)
	// }

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS postslikes (
        "id"    INTEGER NOT NULL UNIQUE,
        "value"   INTEGER,
        "user_id"    INTEGER NOT NULL,
		"post_id"    INTEGER NOT NULL,
        FOREIGN KEY("user_id") REFERENCES "users"("id"),
        FOREIGN KEY("post_id") REFERENCES "posts"("id"),
		PRIMARY KEY("id" AUTOINCREMENT),
		CONSTRAINT postslikes_user_id_post_id_fk UNIQUE ("user_id", "post_id")
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().postslikes!!! %v", err)
	}

	_, err = d.db.Exec(`CREATE TABLE  IF NOT EXISTS comments (
        "id"    INTEGER NOT NULL UNIQUE,
        "user_id"    INTEGER NOT NULL,
        "post_id"    INTEGER NOT NULL,
        "date"    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "comment"    TEXT NOT NULL,
        FOREIGN KEY("user_id") REFERENCES "users"("id"),
        FOREIGN KEY("post_id") REFERENCES "posts"("id"),
        PRIMARY KEY("id" AUTOINCREMENT)
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().comments!!! %v", err)
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS commentslikes (
        "id"    INTEGER NOT NULL UNIQUE,
        "value"    INTEGER,
        "user_id"    INTEGER NOT NULL,
        "comment_id"    INTEGER NOT NULL,
        FOREIGN KEY("user_id") REFERENCES "users"("id"),
        FOREIGN KEY("comment_id") REFERENCES "comments"("id"),
		PRIMARY KEY("id" AUTOINCREMENT),
		CONSTRAINT commentlikes_user_id_comment_id_fk UNIQUE ("user_id", "comment_id")
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().commentslikes!!! %v", err)
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
        "id"    INTEGER NOT NULL UNIQUE,
		"user_id"     INTEGER NOT NULL UNIQUE,
		"key"		TEXT NOT NULL,
		"date"		TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY("user_id") REFERENCES "users"("id"),
        PRIMARY KEY("id" AUTOINCREMENT)
	);`)

	if err != nil {
		log.Fatalf("This error is in db.d.createTable().sessions!!! %v", err)
	}
}
