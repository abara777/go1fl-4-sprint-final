package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	if len(data) < 5 {
		return 0, 0 * time.Hour, errors.New("некорректная строка")
	}

	str := strings.Split(data, ",")

	if len(str) != 2 {
		return 0, 0 * time.Hour, errors.New("неправильный формат строки")
	}

	steps, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, 0 * time.Hour, err
	}

	t, err := time.ParseDuration(str[1])

	if err != nil {
		return 0, 0 * time.Hour, err
	}

	if steps < 1 {
		return 0, 0 * time.Hour, errors.New("отрийательные шаги")
	}

	if t <= 0 {
		return 0, 0 * time.Hour, errors.New("отрицательное время")
	}

	return steps, t, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, t, err := parsePackage(data)

	if err != nil {
		log.Println(err)
		return ""
	}

	distance := (float64(steps) * stepLength) / mInKm

	burnt, err := spentcalories.WalkingSpentCalories(steps, weight, height, t)

	if err != nil {
		return fmt.Sprint(err)
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, burnt)

}
