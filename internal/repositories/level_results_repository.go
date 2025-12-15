package repositories

import (
	"github.com/rajkovrga/1000words-game/internal/db"
	dbmodel "github.com/rajkovrga/1000words-game/internal/models/db"
)

func getLevelResultByUserId(userId int) (*dbmodel.LevelResult, error) {
	var levelResult dbmodel.LevelResult

	err := db.DB.QueryRow(`
	SELECT 
	id,
    user_id,
    target_language_id,
    level,
    success,
    level_completed_at,
    total_completed_at,
    created_at 
	FROM level_results
	WHERE user_id = ?
	`, userId).Scan(
		&levelResult.Id,
		&levelResult.UserId,
		&levelResult.Level,
		&levelResult.TargetLanguageId,
		&levelResult.Success,
		&levelResult.LevelCompletedAt,
		&levelResult.TotalCompletedAt,
		&levelResult.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &levelResult, nil

}
