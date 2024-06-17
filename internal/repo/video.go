package repo

import (
	"context"
	"database/sql"
	"stockfoilo_test/internal/db"
	"stockfoilo_test/internal/model"

	"github.com/rs/zerolog/log"
)

func InsertVideo(tx *sql.Tx, ctx context.Context, video model.Video) error {
	query := `
		INSERT INTO video 
		(
			uuid, 
			path,
			video_name, 
			extension, 
			upload_time,
			is_trimed,
			trim_time,
			is_concated,
			concat_time,
			is_encoded,
			encode_time
		)
		VALUES
		(
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
		`

	_, err := tx.ExecContext(ctx, query,
		video.Id,
		video.Path,
		video.VideoName,
		video.Extension,
		video.UploadTime,
		video.IsTrimed,
		video.TrimTime,
		video.IsConcated,
		video.ConcatTime,
		video.IsEncoded,
		video.EncodeTime,
	)

	if err != nil {
		log.Error().Msgf("failed to insert video: %s", err.Error())
		return err
	}

	return nil
}

func FetchVideoWithVideoId(dbCtx *db.DbCtx, videoId string) (model.Video, error) {
	query := `
		SELECT 
			uuid, 
			path,
			video_name, 
			extension, 
			upload_time,
			is_trimed,
			trim_time,
			is_concated,
			concat_time
		FROM video
		WHERE uuid = ?
	`

	stmt, err := dbCtx.CreatePrepareStmt(query)
	if err != nil {
		log.Error().Msgf("failed to create prepare statement: %s", err.Error())
		return model.Video{}, err
	}

	defer stmt.Close()

	var video model.Video
	err = stmt.QueryRow(videoId).Scan(
		&video.Id,
		&video.Path,
		&video.VideoName,
		&video.Extension,
		&video.UploadTime,
		&video.IsTrimed,
		&video.TrimTime,
		&video.IsConcated,
		&video.ConcatTime,
	)

	if err != nil {
		log.Error().Msgf("failed to query row: %s", err.Error())
		return model.Video{}, err
	}

	return video, nil
}

func InsertTrimHistory(tx *sql.Tx, ctx context.Context, originVideoId string, trimVideo model.TrimVideo) error {
	query := `
		INSERT INTO trim_history 
		(
			uuid, 
			origin_video_uuid,
			start_time, 
			end_time
		)
		VALUES
		(
			?,
			?,
			?,
			?
		)
		`

	_, err := tx.ExecContext(
		ctx,
		query,
		trimVideo.VideoId,
		originVideoId,
		trimVideo.StartTime,
		trimVideo.EndTime,
	)

	if err != nil {
		log.Error().Msgf("failed to insert trim history: %s", err.Error())
		return err
	}

	return nil
}

func InsertConcatHistory(tx *sql.Tx, ctx context.Context, concatVideoId string, concatFilePath string) error {
	query := `
		INSERT INTO concat_history 
		(
			uuid, 
			concat_video_uuid_list
		)
		VALUES
		(
			?,
			?
		)
		`

	_, err := tx.ExecContext(
		ctx,
		query,
		concatVideoId,
		concatFilePath,
	)

	if err != nil {
		log.Error().Msgf("failed to insert concat history: %s", err.Error())
		return err
	}

	return nil

}

func DeleteVideo(dbCtx *db.DbCtx, videoId string) error {
	query := `
		DELETE FROM video
		WHERE uuid = ?
	`

	stmt, err := dbCtx.CreatePrepareStmt(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(dbCtx.Ctx, videoId)
	if err != nil {
		log.Error().Msgf("failed to delete video: %s", err.Error())
		return err
	}

	return nil
}

func FetchVideoInfoList(dbCtx *db.DbCtx) ([]model.Video, error) {
	query := `
		SELECT 
			uuid, 
			path,
			video_name, 
			extension, 
			upload_time,
			is_trimed,
			trim_time,
			is_concated,
			concat_time
		FROM video
	`

	stmt, err := dbCtx.CreatePrepareStmt(query)
	if err != nil {
		log.Error().Msgf("failed to create prepare statement: %s", err.Error())
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Error().Msgf("failed to query: %s", err.Error())
		return nil, err
	}

	defer rows.Close()

	videoList := []model.Video{}
	for rows.Next() {
		var video model.Video
		err := rows.Scan(
			&video.Id,
			&video.Path,
			&video.VideoName,
			&video.Extension,
			&video.UploadTime,
			&video.IsTrimed,
			&video.TrimTime,
			&video.IsConcated,
			&video.ConcatTime,
		)
		if err != nil {
			log.Error().Msgf("failed to scan: %s", err.Error())
			return nil, err
		}
		videoList = append(videoList, video)
	}

	return videoList, nil
}
