package usecase

import (
	"bonded/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (u *calendarUsecase) FindCalendar(ctx context.Context, calendarID string) (*models.Calendar, error) {
	return u.calendarRepo.FindByCalendarID(ctx, calendarID)
}

func (u *calendarUsecase) FindPublicCalendars(ctx context.Context) ([]*models.Calendar, error) {
	//全件取得してフィルタリング
	calendars, err := u.calendarRepo.FindAllCalendars(ctx)
	if err != nil {
		return nil, err
	}
	publicCalendars := []*models.Calendar{}
	for _, calendar := range calendars {
		if *calendar.IsPublic {
			publicCalendars = append(publicCalendars, calendar)
		}
	}
	return publicCalendars, nil
}

func (u *calendarUsecase) CreateCalendar(ctx context.Context, calendar *models.CreateCalendar) error {
	if calendar.OwnerName == "" {
		user, err := u.userRepo.FindByUserID(ctx, calendar.OwnerUserID)
		if err != nil {
			return err
		}
		calendar.OwnerName = user.DisplayName
	}
	user := models.User{
		UserID:      calendar.OwnerUserID,
		DisplayName: calendar.OwnerName,
		Email:       calendar.OwnerUserID + "@example.com",
		Password:    "password",
		AccessLevel: "OWNER",
	}
	calendar.Users = []models.User{user}

	calendar.CalendarID = uuid.New().String()

	// CreateCalendarのフィールドをCalendarに変換
	calendarReq := models.Calendar{
		CalendarID:  calendar.CalendarID,
		SortKey:     "CALENDAR",
		Name:        calendar.Name,
		IsPublic:    calendar.IsPublic,
		OwnerUserID: calendar.OwnerUserID,
		Users:       calendar.Users,
		Events:      calendar.Events,
	}

	return u.calendarRepo.Create(ctx, &calendarReq)
}

func (u *calendarUsecase) EditCalendar(ctx context.Context, calendar *models.Calendar, input *models.Calendar) error {
	return u.calendarRepo.Edit(ctx, calendar, input)
}

func (u *calendarUsecase) DeleteCalendar(ctx context.Context, calendarID string) error {
	return u.calendarRepo.Delete(ctx, calendarID)
}

func (u *calendarUsecase) FindCalendars(ctx context.Context, userID string) ([]*models.Calendar, error) {
	return u.calendarRepo.FindByUserID(ctx, userID)
}

func (u *calendarUsecase) FollowCalendar(ctx context.Context, calendar *models.Calendar, userID string) error {
	user, err := u.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return u.calendarRepo.FollowCalendar(ctx, calendar, user)
}

func (u *calendarUsecase) UnfollowCalendar(ctx context.Context, calendar *models.Calendar, userID string) error {
	user, err := u.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return u.calendarRepo.UnfollowCalendar(ctx, calendar, user)
}
