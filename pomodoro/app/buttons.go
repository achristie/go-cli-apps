package app

import (
	"context"
	"fmt"

	"achristie.net/pomodoro/pomodoro"
	"github.com/mum4k/termdash/private/button"
)

type buttonSet struct {
	btStart *button.Button
	btPause *button.Button
}

func newButtonSet(ctx context.Context, config *pomodoro.IntervalConfig, w *widgets, redrawCh chan<- bool, errorCh chan<- error) (*buttonSet, error) {

	startInterval := func() {
		i, err := pomodoro.GetInterval(config)
		errCh <- err

		start := func(i pomodoro.Interval) {
			message := "Take a break"
			if i.Category == pomodoro.CategoryPomodoro {
				message = "Focus on your task"
			}
			w.update([]int{}, i.Category, message, "", redrawCh)
		}

		end := func(pomodoro.Interval) {
			w.update([]int{}, "", "Nothing running...", "", redrawCh)
		}

		periodic := func(i pomodoro.Interval) {
			w.update([]int{int(i.ActualDuration), int(i.PlannedDuration)}, "", "", fmt.Sprint(i.PlannedDuration-i.ActualDuration), redrawCh)
		}

		errorCh <- i.Start(ctx, config, start, periodic, end)
	}

	pauseInterval := func() {
		i, err := pomodoro.GetInterval(config)
		if err != nil {
			errorCh <- err
			return
		}
		if err := i.Pause(config); err != nil {
			if err == pomodoro.ErrIntervalNotRunning {
				return
			}

			errorCh <- err
			return
		}
		w.update([]int{}, "", "Paused..", "", redrawCh)
	}

	btStart, err := button.New("(s)tart", func() error {
		go startInterval()
		return nil
	}, button.GlobalKey('s'), button.WidthFor("(p)ause"), button.Height(2))

	if err != nil {
		return nil, err
	}
}
