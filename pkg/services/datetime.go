package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func convertToDate(date string) (time.Time, error) {
	s := strings.Split(date, " ")
	hours := strings.Split(s[0], ":")[0]
	minutes := strings.Split(s[0], ":")[1]

	h, _ := strconv.Atoi(hours)
	m, _ := strconv.Atoi(minutes)

	if s[1] == "PM" && h != 12 {
		h = h + 12
	}

	if s[1] == "AM" && h == 12 {
		h = 0
	}

	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now().In(location)

	dateformat := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dateformat = dateformat.Add(time.Hour * time.Duration(h))
	dateformat = dateformat.Add(time.Minute * time.Duration(m))

	return dateformat, nil
}

func isValid(utime, opening_time, closing_time time.Time) error {
	if utime.Sub(opening_time).Minutes() < 0 {
		return errors.New("Restaurant will be closed at the selected time")
	}

	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now().In(location)
	dateformat := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	if utime.Sub(dateformat).Minutes() < 20 {
		return errors.New("Preparation time needed. Please Choose a time greater than current time + 20 minutes ")
	}

	return nil
}

func handleOpenCloseConflict(opentime, closetime time.Time) time.Time {
	if opentime.Sub(closetime) > 0 {
		closetime = closetime.AddDate(0, 0, 1)
	}

	return closetime
}

func ValidateRequestTime(opening_time, closing_time, user_time string) error {
	opentime_date, err := convertToDate(opening_time)
	closetime_date, err := convertToDate(closing_time)
	usertime_date, err := convertToDate(user_time)

	closetime_date = handleOpenCloseConflict(opentime_date, closetime_date)

	if err != nil {
		fmt.Println(err)
	}

	return isValid(usertime_date, opentime_date, closetime_date)
}
