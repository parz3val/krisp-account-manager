package twilio_utils

import (
	"fmt"
	"strings"
	"time"
)

func GenerateFilename(prefix string) string {
	// Get current date
	currentTime := time.Now()

	// Format the date as "Month_Day_Year"
	dateString := currentTime.Format("January_01_2006")

	// Replace spaces with underscores and uppercase the prefix
	processedPrefix := strings.ToUpper(strings.ReplaceAll(prefix, " ", "_"))

	// Generate the final filename
	filename := fmt.Sprintf("%s_%s", dateString, processedPrefix)

	return filename
}
