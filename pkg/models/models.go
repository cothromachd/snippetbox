package models

import (
	"errors"
	"time"
)

/*
CREATE TABLE snippets (
id BIGSERIAL NOT NULL PRIMARY KEY,
title VARCHAR(100) NOT NULL,
content TEXT NOT NULL,
created TIMESTAMP WITH TIME ZONE NOT NULL,
expires TIMESTAMP WITH TIME ZONE NOT NULL
);
*/

/*
INSERT INTO snippets (title, content, created, expires)
VALUES ('Не имей сто рублей',
'Не имей сто рублей,\nа имей сто друзей.',
NOW() AT TIME ZONE ('UTC'),
NOW() AT TIME ZONE ('UTC') + INTERVAL '365' DAY
);
*/

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Snippet struct {
	ID		int
	Title	string
	Content string
	Created time.Time
	Expires time.Time
}