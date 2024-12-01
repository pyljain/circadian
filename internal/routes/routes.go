package routes

import (
	"circadian/internal/config"
	"circadian/internal/db"
	"circadian/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, struct{ Result bool }{Result: true})
}

func Uptime(c echo.Context) error {
	cfg := c.Get("config").(*config.Config)
	db := c.Get("db").(*db.Db)

	hcr, err := db.ListHealthCheckResults()
	if err != nil {
		return err
	}

	healthCheckRecordsByTarget := organiseRecordsByTarget(hcr, cfg)

	var result []UptimeResponse

	for targetName, recordsByTarget := range healthCheckRecordsByTarget {
		var sumOfTimeTaken float64
		var downtime float64
		summaryByDay := []UptimeSummaryByDay{}

		previousTimestamp := time.Now().AddDate(0, 0, -60).Truncate(24 * time.Hour)
		currentDate := time.Now().AddDate(0, 0, -60).Truncate(24 * time.Hour)
		dayTotalTimeTaken := float64(0)
		dayTotalRecords := float64(0)
		downtimeForDay := float64(0)
		downtimeForDayPreviousTimestamp := time.Now().AddDate(0, 0, -60).Truncate(24 * time.Hour)

		for _, rec := range recordsByTarget {
			sumOfTimeTaken += rec.TimeTaken
			if rec.ResponseCode >= 400 {
				downtime += rec.CalloutTime.Sub(previousTimestamp).Minutes()
			}
			previousTimestamp = rec.CalloutTime

			dayTotalRecords += 1
			dayTotalTimeTaken += rec.TimeTaken
			if rec.ResponseCode >= 400 {
				downtimeForDay += rec.CalloutTime.Sub(downtimeForDayPreviousTimestamp).Minutes()
			}
			downtimeForDayPreviousTimestamp = rec.CalloutTime

			cdY, cdM, cdD := rec.CalloutTime.Date()
			curY, curM, curD := currentDate.Date()

			if !(cdY == curY && cdM == curM && cdD == curD) {
				uptimeSumByDay := UptimeSummaryByDay{}
				uptimeSumByDay.Date = currentDate
				uptimeSumByDay.AverageLatency = dayTotalTimeTaken / dayTotalRecords
				uptimeSumByDay.Downtime = downtimeForDay
				uptimeSumByDay.Status = endpointStatusForDay(downtimeForDay)
				summaryByDay = append(summaryByDay, uptimeSumByDay)

				dayTotalRecords = 0
				dayTotalTimeTaken = 0
				downtimeForDay = 0
				currentDate = rec.CalloutTime.Truncate(24 * time.Hour)
			}
		}

		summaryByDay = append(summaryByDay, UptimeSummaryByDay{
			Date:           currentDate,
			AverageLatency: dayTotalTimeTaken / dayTotalRecords,
			Downtime:       downtimeForDay,
			Status:         endpointStatusForDay(downtimeForDay),
		})

		currentStatus := "healthy"
		if recordsByTarget[len(recordsByTarget)-1].ResponseCode >= 400 {
			currentStatus = "down"
		}

		u := UptimeResponse{
			TargetName:     targetName,
			Endpoint:       recordsByTarget[0].TargetEndpoint,
			Interval:       getIntervalForTarget(targetName, cfg),
			CurrentStatus:  currentStatus,
			AverageLatency: sumOfTimeTaken / float64(len(recordsByTarget)),
			Uptime:         ((86400 - downtime) / 86400) * 100,
			Checks:         summaryByDay,
		}

		result = append(result, u)
	}

	return c.JSON(http.StatusOK, result)
}

func getIntervalForTarget(targetName string, cfg *config.Config) int {
	for _, t := range cfg.Targets {
		if t.Name == targetName {
			return t.Interval
		}
	}
	return 0
}

func organiseRecordsByTarget(hcr []models.HealthCheckResult, cfg *config.Config) map[string][]models.HealthCheckResult {

	urlToTargetNameLookup := map[string]string{}

	for _, t := range cfg.Targets {
		urlToTargetNameLookup[t.URL] = t.Name
	}

	healthCheckRecordsByTarget := map[string][]models.HealthCheckResult{}
	for _, h := range hcr {
		name, exists := urlToTargetNameLookup[h.TargetEndpoint]
		if !exists {
			continue
		}
		_, exists = healthCheckRecordsByTarget[name]
		if !exists {
			healthCheckRecordsByTarget[name] = []models.HealthCheckResult{}
		}

		healthCheckRecordsByTarget[name] = append(healthCheckRecordsByTarget[name], h)
	}

	return healthCheckRecordsByTarget
}

func endpointStatusForDay(downtime float64) string {
	if downtime > 400 {
		return "down"
	} else if downtime > 300 {
		return "degraded"
	} else {
		return "healthy"
	}
}

type UptimeResponse struct {
	TargetName     string               `json:"name"`
	Endpoint       string               `json:"url"`
	Interval       int                  `json:"interval"`
	CurrentStatus  string               `json:"currentStatus"`
	AverageLatency float64              `json:"avgLatency"`
	Uptime         float64              `json:"uptime"`
	Checks         []UptimeSummaryByDay `json:"checks"`
}

type UptimeSummaryByDay struct {
	Date           time.Time `json:"date"`
	AverageLatency float64   `json:"averageLatency"`
	Downtime       float64   `json:"downtime"`
	Status         string    `json:"status"`
}
