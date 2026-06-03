package repository

import (
	"database/sql"
	"fmt"

	"github.com/gowild/server/internal/model"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *model.Team) error {
	return r.db.QueryRow(
		`INSERT INTO teams (name, creator_id, code) VALUES ($1,$2,$3) RETURNING id, created_at, updated_at`,
		team.Name, team.CreatorID, team.Code,
	).Scan(&team.ID, &team.CreatedAt, &team.UpdatedAt)
}

func (r *TeamRepository) FindByCode(code string) (*model.Team, error) {
	team := &model.Team{}
	err := r.db.QueryRow(
		`SELECT id, name, creator_id, code, created_at, updated_at FROM teams WHERE code=$1`, code,
	).Scan(&team.ID, &team.Name, &team.CreatorID, &team.Code, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询车队失败: %w", err)
	}
	return team, nil
}

func (r *TeamRepository) FindByID(id int) (*model.Team, error) {
	team := &model.Team{}
	err := r.db.QueryRow(
		`SELECT id, name, creator_id, code, created_at, updated_at FROM teams WHERE id=$1`, id,
	).Scan(&team.ID, &team.Name, &team.CreatorID, &team.Code, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询车队失败: %w", err)
	}
	return team, nil
}

func (r *TeamRepository) FindByUserID(userID int) (*model.Team, error) {
	team := &model.Team{}
	err := r.db.QueryRow(
		`SELECT t.id, t.name, t.creator_id, t.code, t.created_at, t.updated_at 
		 FROM teams t JOIN team_members tm ON t.id = tm.team_id WHERE tm.user_id = $1`, userID,
	).Scan(&team.ID, &team.Name, &team.CreatorID, &team.Code, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("用户尚未加入车队: %w", err)
	}
	return team, nil
}

func (r *TeamRepository) AddMember(teamID, userID int, color string) error {
	_, err := r.db.Exec(
		`INSERT INTO team_members (team_id, user_id, color) VALUES ($1,$2,$3) ON CONFLICT (team_id, user_id) DO NOTHING`,
		teamID, userID, color,
	)
	return err
}

func (r *TeamRepository) RemoveMember(teamID, userID int) error {
	_, err := r.db.Exec(
		`DELETE FROM team_members WHERE team_id=$1 AND user_id=$2`, teamID, userID,
	)
	return err
}

func (r *TeamRepository) GetMembers(teamID int) ([]model.TeamMember, error) {
	rows, err := r.db.Query(
		`SELECT tm.id, tm.team_id, tm.user_id, tm.joined_at, tm.color, u.username, u.avatar 
		 FROM team_members tm JOIN users u ON tm.user_id = u.id 
		 WHERE tm.team_id=$1 ORDER BY tm.joined_at`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.TeamMember
	for rows.Next() {
		var m model.TeamMember
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.JoinedAt,
			&m.Color, &m.Username, &m.Avatar); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *TeamRepository) AddMessage(msg *model.TeamMessage) error {
	return r.db.QueryRow(
		`INSERT INTO team_messages (team_id, user_id, content, type) VALUES ($1,$2,$3,$4) RETURNING id, created_at`,
		msg.TeamID, msg.UserID, msg.Content, msg.Type,
	).Scan(&msg.ID, &msg.CreatedAt)
}

func (r *TeamRepository) GetMessages(teamID, limit, offset int) ([]model.TeamMessage, error) {
	rows, err := r.db.Query(
		`SELECT tm.id, tm.team_id, tm.user_id, u.username, tm.content, tm.type, tm.created_at 
		 FROM team_messages tm JOIN users u ON tm.user_id = u.id 
		 WHERE tm.team_id=$1 ORDER BY tm.created_at DESC LIMIT $2 OFFSET $3`, teamID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []model.TeamMessage
	for rows.Next() {
		var m model.TeamMessage
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.Username,
			&m.Content, &m.Type, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *TeamRepository) UpdateLocation(loc *model.TeamLocation) error {
	_, err := r.db.Exec(
		`INSERT INTO team_locations (team_id, user_id, latitude, longitude, speed, distance, updated_at) 
		 VALUES ($1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP)
		 ON CONFLICT (team_id, user_id) 
		 DO UPDATE SET latitude=$3, longitude=$4, speed=$5, distance=$6, updated_at=CURRENT_TIMESTAMP`,
		loc.TeamID, loc.UserID, loc.Latitude, loc.Longitude, loc.Speed, loc.Distance,
	)
	return err
}

func (r *TeamRepository) GetLocations(teamID int) ([]model.TeamLocation, error) {
	rows, err := r.db.Query(
		`SELECT tl.id, tl.team_id, tl.user_id, u.username, tl.latitude, tl.longitude, tl.speed, tl.distance, tl.updated_at 
		 FROM team_locations tl JOIN users u ON tl.user_id = u.id 
		 WHERE tl.team_id=$1 AND tl.updated_at > NOW() - INTERVAL '10 minutes'`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locs []model.TeamLocation
	for rows.Next() {
		var l model.TeamLocation
		if err := rows.Scan(&l.ID, &l.TeamID, &l.UserID, &l.Username,
			&l.Latitude, &l.Longitude, &l.Speed, &l.Distance, &l.UpdatedAt); err != nil {
			return nil, err
		}
		locs = append(locs, l)
	}
	return locs, nil
}

func (r *TeamRepository) IsMember(teamID, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id=$1 AND user_id=$2)`, teamID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *TeamRepository) GetMemberCount(teamID int) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM team_members WHERE team_id=$1`, teamID).Scan(&count)
	return count, err
}
