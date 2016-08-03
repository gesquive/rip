package format

import "fmt"
import "github.com/dustin/go-humanize"

// RawBytes formats a byte count without the units
func RawBytes(bytes uint64) string {
	prettySize, _ := humanize.ComputeSI(float64(bytes))
	return fmt.Sprintf("%3.0f", prettySize)
}

// Bytes formats a byte count
func Bytes(bytes uint64) string {
	prettySize, prettyUnit := humanize.ComputeSI(float64(bytes))
	return fmt.Sprintf("%3.0f %sB", prettySize, prettyUnit)
}

// Percent formats a percent
func Percent(current uint64, total uint64) string {
	return fmt.Sprintf("%3.0f%%", float64(current)/float64(total)*100)
}

// Progress formats a progress line
func Progress(current uint64, total uint64) string {
	// TODO: What we want to do is have a right adjusted " 50kB/100GB" type of message

	totalPrettySize, totalPrettyUnit := humanize.ComputeSI(float64(total))
	currentPrettySize, currentPrettyUnit := humanize.ComputeSI(float64(current))
	if currentPrettyUnit == totalPrettyUnit {
		currentPrettyUnit = ""
	}

	return fmt.Sprintf("%11s", fmt.Sprintf("%.0f%s/%.0f%sB", currentPrettySize,
		currentPrettyUnit, totalPrettySize, totalPrettyUnit))
}
