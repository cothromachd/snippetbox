package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	//"github.com/jackc/pgx/v5"
	"github.com/cothromachd/snippetbox/pkg/models/postgresql"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *postgresql.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	dsn := flag.String("dsn", "postgres://web:190204@localhost:5432/snippetbox", "Название PostgreSQL источника данных")
	flag.Parse()
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	
	defer db.Close()
 
	app := &application {
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &postgresql.SnippetModel{DB: db},
	}

	srv := &http.Server {
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
}

/*
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
*/