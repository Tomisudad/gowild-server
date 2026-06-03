package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/repository"
)

// 成员颜色数组
var memberColors = []string{
	"#E87D22", "#4A90D9", "#50C878", "#E74C3C", "#9B59B6",
	"#F39C12", "#1ABC9C", "#E91E63", "#00BCD4", "#FF5722",
}

type TeamService struct {
	repo *repository.TeamRepository
}

func NewTeamService(repo *repository.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

// generateCode 生成6位车队邀请码
func generateCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}

func (s *TeamService) Create(creatorID int, req model.CreateTeamRequest) (*model.Team, error) {
	// 检查是否已加入车队
	_, err := s.repo.FindByUserID(creatorID)
	if err == nil {
		return nil, errors.New("你已加入其他车队，请先退出")
	}

	code := generateCode()
	// 确保邀请码唯一
	for {
		_, err := s.repo.FindByCode(code)
		if err != nil {
			break
		}
		code = generateCode()
	}

	team := &model.Team{
		Name:      req.Name,
		CreatorID: creatorID,
		Code:      code,
	}

	if err := s.repo.Create(team); err != nil {
		return nil, err
	}

	// 自动加入车队
	color := memberColors[0]
	if err := s.repo.AddMember(team.ID, creatorID, color); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) Join(userID int, req model.JoinTeamRequest) (*model.Team, error) {
	// 检查是否已加入车队
	_, err := s.repo.FindByUserID(userID)
	if err == nil {
		return nil, errors.New("你已加入其他车队，请先退出")
	}

	team, err := s.repo.FindByCode(req.Code)
	if err != nil {
		return nil, errors.New("车队不存在或邀请码无效")
	}

	// 分配颜色
	count, _ := s.repo.GetMemberCount(team.ID)
	color := memberColors[count%len(memberColors)]

	if err := s.repo.AddMember(team.ID, userID, color); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) Leave(teamID, userID int) error {
	team, err := s.repo.FindByID(teamID)
	if err != nil {
		return err
	}
	if team.CreatorID == userID {
		return errors.New("车队创建者不能直接退出，请先转让或解散车队")
	}
	return s.repo.RemoveMember(teamID, userID)
}

func (s *TeamService) GetTeam(teamID int) (*model.Team, error) {
	return s.repo.FindByID(teamID)
}

func (s *TeamService) GetMyTeam(userID int) (*model.Team, error) {
	return s.repo.FindByUserID(userID)
}

func (s *TeamService) GetMembers(teamID int) ([]model.TeamMember, error) {
	return s.repo.GetMembers(teamID)
}

func (s *TeamService) SendMessage(teamID, userID int, req model.TeamMessageRequest) (*model.TeamMessage, error) {
	isMember, err := s.repo.IsMember(teamID, userID)
	if err != nil || !isMember {
		return nil, errors.New("你不是该车队成员")
	}

	msg := &model.TeamMessage{
		TeamID:  teamID,
		UserID:  userID,
		Content: req.Content,
		Type:    req.Type,
	}

	if err := s.repo.AddMessage(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *TeamService) GetMessages(teamID, limit, offset int) ([]model.TeamMessage, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	return s.repo.GetMessages(teamID, limit, offset)
}

func (s *TeamService) UpdateLocation(teamID, userID int, req model.TeamLocationRequest) error {
	isMember, err := s.repo.IsMember(teamID, userID)
	if err != nil || !isMember {
		return errors.New("你不是该车队成员")
	}

	loc := &model.TeamLocation{
		TeamID:    teamID,
		UserID:    userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Speed:     req.Speed,
		Distance:  req.Distance,
	}

	return s.repo.UpdateLocation(loc)
}

func (s *TeamService) GetLocations(teamID int) ([]model.TeamLocation, error) {
	return s.repo.GetLocations(teamID)
}
