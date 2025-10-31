package daysteps

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	if strings.TrimSpace(data) == "" {
		return 0, 0, errors.New("empty input")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid data format")
	}

	if parts[0] != strings.TrimSpace(parts[0]) {
		return 0, 0, errors.New("invalid steps format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("invalid steps")
	}

	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		return 0, 0, errors.New("invalid duration")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil || steps <= 0 || duration <= 0 {
		log.Printf("Ошибка при разборе данных: %v", err)
		return ""
	}

	distance := float64(steps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка при обработке данных: %v", err)
		return ""
	}

	calories = math.Round(calories*100) / 100

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, calories,
	)
}
