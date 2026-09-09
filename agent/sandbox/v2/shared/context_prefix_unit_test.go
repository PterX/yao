//go:build unit

package shared_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yaoapp/yao/agent/sandbox/v2/shared"
)

func TestBuildContextPrefix_AllVars(t *testing.T) {
	vars := map[string]string{
		"LOCALE":         "zh-cn",
		"WORKSPACE_ID":   "ws-123",
		"ASSISTANT_ID":   "yao.general",
		"SKILLS_DIR":     "/workspace/.yao/assistants/yao.general/skills",
		"EXT_SKILLS_DIR": "/workspace/.yao/skills",
	}
	got := shared.BuildContextPrefix(vars)
	if !strings.HasPrefix(got, "[Session context — use directly: Time: ") {
		t.Errorf("missing prefix: %q", got)
	}
	if !strings.Contains(got, "Locale: zh-cn") {
		t.Errorf("missing Locale: %q", got)
	}
	if !strings.Contains(got, "WorkspaceID: ws-123") {
		t.Errorf("missing WorkspaceID: %q", got)
	}
	if !strings.Contains(got, "AssistantID: yao.general") {
		t.Errorf("missing AssistantID: %q", got)
	}
	if !strings.Contains(got, "SkillsDir: /workspace/.yao/assistants/yao.general/skills") {
		t.Errorf("missing SkillsDir: %q", got)
	}
	if !strings.Contains(got, "ExtSkillsDir: /workspace/.yao/skills") {
		t.Errorf("missing ExtSkillsDir: %q", got)
	}
	if !strings.HasSuffix(got, "]") {
		t.Errorf("missing closing bracket: %q", got)
	}
}

func TestBuildContextPrefix_NilVars(t *testing.T) {
	got := shared.BuildContextPrefix(nil)
	if !strings.HasPrefix(got, "[Session context — use directly: Time: ") {
		t.Errorf("nil vars should still produce prefix: %q", got)
	}
	for _, key := range []string{"Locale", "WorkspaceID", "AssistantID", "SkillsDir", "ExtSkillsDir"} {
		if strings.Contains(got, key) {
			t.Errorf("nil vars should not produce %s: %q", key, got)
		}
	}
}

func TestBuildContextPrefix_EmptyVars(t *testing.T) {
	got := shared.BuildContextPrefix(map[string]string{})
	for _, key := range []string{"Locale", "WorkspaceID", "AssistantID", "SkillsDir", "ExtSkillsDir"} {
		if strings.Contains(got, key) {
			t.Errorf("empty vars should not produce %s: %q", key, got)
		}
	}
}

func TestBuildContextPrefix_TimeFormat(t *testing.T) {
	got := shared.BuildContextPrefix(nil)
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	if !strings.Contains(got, dateStr) {
		t.Errorf("should contain today's date %q: %q", dateStr, got)
	}
	weekday := now.Weekday().String()
	if !strings.Contains(got, weekday) {
		t.Errorf("should contain weekday %q: %q", weekday, got)
	}
}

func TestBuildContextPrefix_PartialVars(t *testing.T) {
	got := shared.BuildContextPrefix(map[string]string{"LOCALE": "en-us"})
	if !strings.Contains(got, "Locale: en-us") {
		t.Errorf("should contain Locale: %q", got)
	}
	for _, key := range []string{"WorkspaceID", "AssistantID", "SkillsDir", "ExtSkillsDir"} {
		if strings.Contains(got, key) {
			t.Errorf("should not contain %s with only LOCALE set: %q", key, got)
		}
	}
}

func TestBuildContextPrefix_AssistantAndSkills(t *testing.T) {
	got := shared.BuildContextPrefix(map[string]string{
		"ASSISTANT_ID":   "yao.smith",
		"SKILLS_DIR":     "/ws/.yao/assistants/yao.smith/skills",
		"EXT_SKILLS_DIR": "/ws/.yao/skills",
	})
	if !strings.Contains(got, "AssistantID: yao.smith") {
		t.Errorf("missing AssistantID: %q", got)
	}
	if !strings.Contains(got, "SkillsDir: /ws/.yao/assistants/yao.smith/skills") {
		t.Errorf("missing SkillsDir: %q", got)
	}
	if !strings.Contains(got, "ExtSkillsDir: /ws/.yao/skills") {
		t.Errorf("missing ExtSkillsDir: %q", got)
	}
	if strings.Contains(got, "Locale") || strings.Contains(got, "WorkspaceID") {
		t.Errorf("should not contain unset keys: %q", got)
	}
}
