package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid data format")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps value: %w", err)
	}
	if steps < 0 {
		return 0, 0, errors.New("steps cannot be negative")
	}

	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration value: %w", err)
	}

	if duration < 0 {
		return 0, 0, errors.New("duration cannot be negative")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	parsePackage(data)
	steps, duration, err := parsePackage(data)
	if err != nil {
		return fmt.Errorf("Error parsing data: %v", err).Error()
	}
	if steps == 0 {
		return fmt.Errorf("No steps recorded").Error()
	}
	distanceKm := float64(steps) * stepLength / mInKm
	calories := (0.035 * weight) + ((distanceKm / (duration.Hours())) * (0.029 * weight * 1000 / height))
	return fmt.Sprintf("Steps: %d.\nDistance: %.2f km.\nCalories burned: %.2f kcal.n", steps, distanceKm, calories)
}
