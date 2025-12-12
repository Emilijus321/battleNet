package templates

import (
	"battleNet/models"
	"fmt"
)

// Helper function to format integers
func formatInt(n int) string {
	return fmt.Sprintf("%d", n)
}

// Helper function to count movies by status
func countByStatus(movies []models.Movie, status string) int {
	count := 0
	for _, movie := range movies {
		if movie.Status != nil && *movie.Status == status {
			count++
		}
	}
	return count
}

// Helper function to get average rating
func getAverageRating(movies []models.Movie) string {
	total := 0.0
	count := 0
	for _, movie := range movies {
		if movie.VoteAverage != nil {
			total += *movie.VoteAverage
			count++
		}
	}
	if count == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.1f", total/float64(count))
}

// Helper function to get released count
func getReleasedCount(movies []models.Movie) int {
	return countByStatus(movies, "Released")
}

// Helper function to slice strings for preview
func slice(s string, start, end int) string {
	if len(s) <= end {
		return s
	}
	return s[start:end]
}
func Slice(s string, start, end int) string {
	if start > len(s) {
		return ""
	}
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}

func Printf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
