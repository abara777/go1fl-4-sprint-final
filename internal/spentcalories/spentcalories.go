package spentcalories

import (
	"errors"
	"fmt"
	"log"
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
	d := strings.Split(data, ",")
	if len(d) != 3 {
		return 0, "", 0 * time.Hour, errors.New("неправильный формат строки")
	}

	steps, err := strconv.Atoi(d[0])

	if err != nil {
		return 0, "", 0 * time.Hour, err
	}

	if steps < 1 {
		return 0, "", 0 * time.Hour, errors.New("неверное количество шагов")
	}

	activity := d[1]
	t, err := time.ParseDuration(d[2])

	if err != nil {
		return 0, "", 0 * time.Hour, err
	}
	if t < 1 {
		return 0, "", 0 * time.Hour, errors.New("неверное время")
	}

	return steps, activity, t, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	distance := stepLen * float64(steps)
	return distance / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	return distance / float64(duration.Hours())
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	var dist, meanSp, cals float64

	steps, activity, t, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}

	switch activity {
	case "Бег":
		dist = distance(steps, height)
		meanSp = meanSpeed(steps, height, t)
		cals, err = RunningSpentCalories(steps, weight, height, t)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		dist = distance(steps, height)
		meanSp = meanSpeed(steps, height, t)
		cals, err = WalkingSpentCalories(steps, weight, height, t)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	res := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, float64(t.Hours()), dist, meanSp, cals)
	return res, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 1 {
		return 0, errors.New("Мало шагов")
	}

	if weight < 10 {
		return 0, errors.New("Низкий вес")
	}

	if height < 0.5 {
		return 0, errors.New("низкий рост")
	}

	if duration.Minutes() < 1 {
		return 0, errors.New("Мало время")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * meanSpeed * durationInMinutes) / minInH, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 1 || weight < 10 || height < 0.5 || duration < 0 {
		return 0, errors.New("")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	res := ((meanSpeed * weight * duration.Minutes()) / minInH) * walkingCaloriesCoefficient

	return res, nil
}
