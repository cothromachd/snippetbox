package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/cothromachd/snippetbox/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SnippetModel - Определяем тип который обертывает пул подключения pgxpool.Pool
type SnippetModel struct {
	DB *pgxpool.Pool
}

// Insert - Метод для создания новой заметки в базе данных.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := fmt.Sprintf(`INSERT INTO snippets (title, content, created, expires)
 	VALUES ($1, $2, NOW() AT TIME ZONE ('UTC'), NOW() AT TIME ZONE ('UTC') + INTERVAL '%s' DAY);`, expires)
	
	result, err := m.DB.Exec(context.Background(), stmt, title, content)
	if err != nil {
		return 0, err
	}
	
	rows_affected := result.RowsAffected()
	// TODO:
	// логика rows_affected работает как id заметки, что потом переходит в showSnippet
	// и показывает будто id заметки, а на самом деле число измененных строк
	// нужно исправить
	return int(rows_affected), nil

}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > NOW() AT TIME ZONE ('UTC') AND id = $1`
	
	// Используем метод QueryRow() для выполнения SQL запроса, 
	// передавая ненадежную переменную id в качестве значения для плейсхолдера
	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(context.Background(), stmt, id)

	// Инициализируем указатель на новую структуру Snippet.
	s := &models.Snippet{}

	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в 
	// соответствующее поле в структуре Snippet. Обратите внимание, что аргументы 
	// для row.Scan - это указатели на место, куда требуется скопировать данные
	// и количество аргументов должно быть точно таким же, как количество 
	// столбцов в таблице базы данных.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект Snippet.
	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() AT TIME ZONE ('UTC') ORDER BY created DESC LIMIT 10`
	
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}


		snippets = append(snippets, s)

		if err = rows.Err(); err != nil {
			return nil, err
		}

	}
	return snippets, nil
}