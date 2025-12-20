package repositories

import (
	"database/sql"

	"github.com/rajkovrga/1000words-game/internal/db"
	gamemodel "github.com/rajkovrga/1000words-game/internal/models/game"
)

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (wordRepository *WordRepository) GetGameWords(
	nativeLanguageId int,
	targetLanguageId int,
	level int,
	limit int,
) ([]gamemodel.WordPair, error) {
	var wordPairs []gamemodel.WordPair
	if level <= 0 {
		level = 1
	}
	if limit <= 0 {
		limit = 100
	}

	offset := (level - 1) * limit

	rows, err := db.DB.Query(`
	SELECT id, wt1.word_value as native_word, wt2.word_walue as target_word
	from words w
	inner join word_translations wt1 on wt1.word_id = w.id and language_id = ?
	inner join word_translations wt2 on wt2.word_id = w.id and language id = ?
	ORDER BY RANDOM()
	LIMIT ? OFFSET ?
	`, nativeLanguageId, targetLanguageId, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var wp gamemodel.WordPair

		err := rows.Scan(&wp.Id, &wp.NativeWord, &wp.TargetWord)

		if err != nil {
			return nil, err
		}

		wordPairs = append(wordPairs, wp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return wordPairs, nil
}
