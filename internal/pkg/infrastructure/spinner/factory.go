package spinner

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

var defaultSpinnerColor string = "fgCyan"
var defaultSpinnerDuration int = 120
var defaultSpinnerType int = 11

type Spinner struct {
	spinner *spinner.Spinner
}

type SpinnerConfig struct {
	Color       *string
	Duration    *int
	SpinnerType *int
}

func (sc SpinnerConfig) getSpinnerType() int {
	if sc.SpinnerType == nil {
		return defaultSpinnerType
	}

	return *sc.SpinnerType
}

func (sc SpinnerConfig) getDuration() int {
	if sc.Duration == nil {
		return defaultSpinnerDuration
	}

	return *sc.Duration
}

func (sc SpinnerConfig) getColor() string {
	if sc.Color == nil {
		return defaultSpinnerColor
	}

	return *sc.Color
}

func getSpinnerConfig(config []SpinnerConfig) SpinnerConfig {
	for _, config := range config {
		return config
	}

	return SpinnerConfig{
		Color:       &defaultSpinnerColor,
		Duration:    &defaultSpinnerDuration,
		SpinnerType: &defaultSpinnerType,
	}
}

func NewSpinner(config ...SpinnerConfig) *Spinner {
	innerConfig := getSpinnerConfig(config)

	dotStyle := spinner.CharSets[innerConfig.getSpinnerType()]

	return &Spinner{
		spinner: spinner.New(
			dotStyle,
			time.Duration(innerConfig.getDuration())*time.Millisecond,
			spinner.WithWriter(os.Stderr),
			spinner.WithColor(innerConfig.getColor()),
		),
	}
}

func (s Spinner) Start() {
	s.spinner.Start()
}

func (s Spinner) Stop() {
	s.spinner.Stop()
}
