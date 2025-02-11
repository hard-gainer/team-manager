package template

import (
    "fmt"
    "github.com/jackc/pgx/v5/pgtype"
    "html/template"
)

// GetTemplateFuncs returns a map of template helper functions
func GetTemplateFuncs() template.FuncMap {
    return template.FuncMap{
        "formatDuration": formatDuration,
    }
}

// formatDuration converts duration in seconds to HH:MM:SS format
func formatDuration(duration pgtype.Int8) string {
    if !duration.Valid {
        return "00:00:00"
    }
    seconds := duration.Int64
    hours := seconds / 3600
    minutes := (seconds % 3600) / 60
    remainingSeconds := seconds % 60
    return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, remainingSeconds)
}