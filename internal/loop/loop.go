package loop

import (
	"bytes"
	"circadian/internal/config"
	"circadian/internal/db"
	"circadian/internal/models"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg *config.Config, db *db.Db) error {
	client := http.Client{}
	eg, egCtx := errgroup.WithContext(ctx)
	for _, target := range cfg.Targets {
		eg.Go(func() error {
			for {
				log.Printf("Calling out to %s", target.URL)
				resp, err := makeCallout(egCtx, target, &client)
				if err != nil {
					if egCtx.Err() != nil {
						return fmt.Errorf("context was cancelled: %s", err)
					}

					return err
				}

				// Write to database
				err = db.InsertHealthCheckResult(&models.HealthCheckResult{
					TargetEndpoint: target.URL,
					HTTPMethod:     target.Method,
					CalloutTime:    resp.timestamp,
					ResponseCode:   resp.responseCode,
					Response:       resp.responseDetail,
					TimeTaken:      float64(resp.timeTaken),
				})
				if err != nil {
					return err
				}

				time.Sleep((time.Duration(target.Interval) * time.Second))
			}

		})
	}
	return nil
}

func makeCallout(ctx context.Context, target config.Target, client *http.Client) (*Response, error) {
	var err error
	var req *http.Request

	if target.Body != "" {
		buf := bytes.NewBufferString(target.Body)
		req, err = http.NewRequest(target.Method, target.URL, buf)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(target.Method, target.URL, nil)
		if err != nil {
			return nil, err
		}
	}

	req = req.WithContext(ctx)

	for name, value := range target.Headers {
		req.Header.Add(name, replaceEnvVars(value))
	}

	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return &Response{
			responseCode:   0,
			timestamp:      time.Now(),
			responseDetail: err.Error(),
			timeTaken:      time.Since(startTime).Milliseconds(),
		}, nil
	}

	timeTaken := time.Since(startTime).Milliseconds()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		responseCode:   resp.StatusCode,
		timestamp:      time.Now(),
		responseDetail: string(respBytes),
		timeTaken:      timeTaken,
	}, nil

}

func replaceEnvVars(input string) string {
	// Regular expression to match $VAR
	re := regexp.MustCompile(`\$\w+`)

	// Replace all matches with the corresponding environment variable value
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// Remove the leading $ to get the environment variable name
		envVarName := strings.TrimPrefix(match, "$")
		// Get the environment variable value
		value, exists := os.LookupEnv(envVarName)
		if !exists {
			// If the environment variable does not exist, leave it as is
			return match
		}
		return value
	})

	return result
}

type Response struct {
	timeTaken      int64
	timestamp      time.Time
	responseCode   int
	responseDetail string
}
