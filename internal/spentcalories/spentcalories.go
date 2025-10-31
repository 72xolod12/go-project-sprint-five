package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000.0
	minInH                     = 60.0
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	stepsStr := strings.TrimSpace(parts[0])
	trainingType := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps value: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps must be positive")
	}

	if strings.Contains(durationStr, " h") || strings.Contains(durationStr, " m") {
		return 0, "", 0, fmt.Errorf("invalid duration format: space between number and unit")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration value: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}

	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepDistance := height * stepLengthCoefficient // метров
	stepTotalDistance := float64(steps) * stepDistance
	return stepTotalDistance / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Hours() <= 0 {
		return 0
	}
	distanceKm := distance(steps, height)

	speed := distanceKm / duration.Hours()

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {

		log.Println(err)
		return "", fmt.Errorf("error parsing data: %v", err)
	}
	if steps == 0 {
		return "", fmt.Errorf("no steps recorded")
	}

	var calories float64
	switch trainingType {
	case "running", "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "walking", "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainingType)
	}
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, duration.Hours(), dist, speed, calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 {
		return 0, errors.New("duration cannot be zero or negative")
	}
	if steps <= 0 {
		return 0, errors.New("steps cannot be zero or negative")
	}
	if weight <= 0 {
		return 0, errors.New("weight cannot be zero or negative")
	}
	if height <= 0 {
		return 0, errors.New("height cannot be zero or negative")
	}

	avgSpeed := meanSpeed(steps, height, duration) // км/ч
	calories := (weight * avgSpeed * duration.Minutes()) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps cannot be zero or negative")
	}
	if weight <= 0 {
		return 0, errors.New("weight cannot be zero or negative")
	}
	if height <= 0 {
		return 0, errors.New("height cannot be zero or negative")
	}
	if duration.Hours() <= 0 {
		return 0, errors.New("duration cannot be zero or negative")
	}

	avgSpeed := meanSpeed(steps, height, duration)

	calories := (weight * avgSpeed * duration.Minutes()) / minInH

	calories *= walkingCaloriesCoefficient
	s := fmt.Sprintf("%.2f", calories)
	calories, _ = strconv.ParseFloat(s, 64)

	return calories, nil
}
