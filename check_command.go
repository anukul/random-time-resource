package resource

import (
	"math/rand"
	"time"

	"github.com/concourse/time-resource/lord"
	"github.com/concourse/time-resource/models"
)

type CheckCommand struct {
}

func (*CheckCommand) Run(request models.CheckRequest) ([]models.Version, error) {
	err := request.Source.Validate()
	if err != nil {
		return nil, err
	}

	previousTime := request.Version.Time
	currentTime := time.Now().UTC()

	specifiedLocation := request.Source.Location
	if specifiedLocation != nil {
		currentTime = currentTime.In((*time.Location)(specifiedLocation))
	}

	interval := request.Source.Interval

	if interval == nil && request.Source.MinInterval != nil && request.Source.MaxInterval != nil {
		randomInterval := models.Interval(rand.Intn(int(*request.Source.MaxInterval) - int(*request.Source.MinInterval) + int(*request.Source.MinInterval)))
		interval = &randomInterval
	}

	tl := lord.TimeLord{
		PreviousTime: previousTime,
		Location:     specifiedLocation,
		Start:        request.Source.Start,
		Stop:         request.Source.Stop,
		Interval:     interval,
		Days:         request.Source.Days,
	}

	versions := []models.Version{}

	if !previousTime.IsZero() {
		versions = append(versions, models.Version{Time: previousTime})
	} else if request.Source.InitialVersion {
		versions = append(versions, models.Version{Time: currentTime})
		return versions, nil
	}

	if tl.Check(currentTime) {
		versions = append(versions, models.Version{Time: currentTime})
	}

	return versions, nil
}
