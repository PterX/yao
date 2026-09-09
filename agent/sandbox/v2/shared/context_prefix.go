package shared

import (
	"fmt"
	"strings"
	"time"
)

// BuildContextPrefix builds the session context metadata prefix
// for user message injection. Combines date/time with context variables.
// Includes a hint so the model uses the metadata directly without tool calls.
func BuildContextPrefix(vars map[string]string) string {
	now := time.Now()
	parts := []string{
		fmt.Sprintf("Time: %s (%s, %s)",
			now.Format("2006-01-02 15:04:05"),
			now.Weekday().String(),
			now.Location().String()),
	}

	if v := vars["LOCALE"]; v != "" {
		parts = append(parts, "Locale: "+v)
	}
	if v := vars["WORKSPACE_ID"]; v != "" {
		parts = append(parts, "WorkspaceID: "+v)
	}
	if v := vars["ASSISTANT_ID"]; v != "" {
		parts = append(parts, "AssistantID: "+v)
	}
	if v := vars["SKILLS_DIR"]; v != "" {
		parts = append(parts, "SkillsDir: "+v)
	}
	if v := vars["EXT_SKILLS_DIR"]; v != "" {
		parts = append(parts, "ExtSkillsDir: "+v)
	}

	return "[Session context — use directly: " + strings.Join(parts, ", ") + "]"
}
