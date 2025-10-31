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
	lenStep                    = 0.65 // средняя длина шага (м)
	mInKm                      = 1000 // метров в километре
	minInH                     = 60   // минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент длины шага от роста
	walkingCaloriesCoefficient = 0.5  // поправка для ходьбы
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	stepsStr := parts[0]
	trType := parts[1]
	durationStr := parts[2]

	if stepsStr != strings.TrimSpace(stepsStr) || trType != strings.TrimSpace(trType) || durationStr != strings.TrimSpace(durationStr) {
		return 0, "", 0, fmt.Errorf("invalid spacing in input")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps value: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps must be positive")
	}

	if trType == "" {
		return 0, "", 0, errors.New("training type is required")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration value: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}

	return steps, trType, duration, nil
}

// meanSpeed — средняя скорость (км/ч)
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Hours() <= 0 {
		return 0
	}
	distanceKm := distance(steps, height)
	return distanceKm / duration.Hours()
}

// distance — пройденная дистанция (в км)
func distance(steps int, height float64) float64 {
	stepDistance := height * stepLengthCoefficient
	stepTotalDistance := float64(steps) * stepDistance
	return stepTotalDistance / mInKm
}

// TrainingInfo — основная функция
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error parsing data: %v", err)
	}

	if steps <= 0 {
		return "", fmt.Errorf("no steps recorded")
	}

	var calories float64
	switch strings.ToLower(trainingType) {
	case "running", "бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "walking", "ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainingType)
	}
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, duration.Hours(), dist, speed, calories,
	)
	return result, nil
}

// RunningSpentCalories — калории при беге
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

	calories := (weight * avgSpeed * duration.Minutes()) / minInH
	return calories, nil
}

// WalkingSpentCalories — калории при ходьбе
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
	return calories, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, _, duration, err := parseTraining(data)
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}
	distanceKm := float64(steps) * lenStep / mInKm

	calories := (0.035 * weight) + ((distanceKm / duration.Hours()) * (0.029 * weight * 1000 / height))
	return fmt.Sprintf("Количество шагов: %d\nДистанция: %.2f км\nВы сожгли: %.2f ккал\n", steps, distanceKm, calories)
}
