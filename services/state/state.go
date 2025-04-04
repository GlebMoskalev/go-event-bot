package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
	"strings"
)

type state struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.State {
	return &state{db: db, log: log}
}

func (s *state) StartAddStaff(ctx context.Context, chatID int64) error {
	err := s.db.SetState(ctx, chatID, models.StateStaffRegisterFullName, []byte("{}"))
	return err
}

func (s *state) RegisterStaffFullName(ctx context.Context, chatID int64, fullName string) error {
	fullNameSplit := strings.Split(fullName, " ")
	fmt.Println(fullNameSplit)
	if len(fullNameSplit) < 3 {
		return errors.New("name is incomplete")
	}

	staff := models.Staff{
		FirstName:  fullNameSplit[0],
		LastName:   fullNameSplit[1],
		Patronymic: fullNameSplit[2],
	}
	data, err := json.Marshal(staff)
	if err != nil {
		s.log.Error("register full name", "error", err)
		return err
	}
	err = s.db.SetState(ctx, chatID, models.StateStaffRegisterPhoneNumber, data)
	s.log.Error("register full name", "error", err)
	return err
}

func (s *state) RegisterStaffNumberPhone(ctx context.Context, chatID int64, numberPhone string) error {
	_, data, err := s.db.GetStateAndData(ctx, chatID)
	if err != nil {
		return err
	}
	var staff models.Staff
	err = json.Unmarshal(data, &staff)
	if err != nil {
		return err
	}
	err = s.db.SetState(ctx, chatID, models.StateStaffRegisterConfirm, data)
	return err
}

func (s *state) ConfirmAddStaff(ctx context.Context, chatID int64) error {
	err := s.db.RemoveState(ctx, chatID)
	s.log.Error("ef", "error", err)
	return err
}

func (s *state) Get(ctx context.Context, chatID int64) (models.State, error) {
	return s.db.GetState(ctx, chatID)
}
