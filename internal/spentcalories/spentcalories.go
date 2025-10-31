package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps value: %w", err)
	}
	if steps < 0 {
		return 0, "", 0, errors.New("steps cannot be negative")
	}
	trainingType := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration value: %w", err)
	}
	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepDistance := height * stepLengthCoefficient / 100
	stepTotalDistance := float64(steps) * stepDistance
	return stepTotalDistance / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Hours() == 0 {
		return 0
	}
	distanceKm := distance(steps, height)
	return distanceKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("Error parsing data: %v", err)
	}
	if steps == 0 {
		return "", fmt.Errorf("No steps recorded")
	}
	var calories float64
	switch trainingType {

	case "running", "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)

	case "walking", "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)

	default:
		return "", fmt.Errorf("Unknown training type: %s", trainingType)
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Training type: %s, Calories spent: %.2f", trainingType, calories), nil
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
	avgSpeed := meanSpeed(steps, height, duration)
	calories := (weight * 18 * avgSpeed) / mInKm * duration.Minutes() * walkingCaloriesCoefficient
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
	calories := (0.035 * weight) + ((avgSpeed*avgSpeed)/height)*0.029*weight*duration.Minutes()
	caloriesTotal := calories * walkingCaloriesCoefficient
	return caloriesTotal, nil
}
