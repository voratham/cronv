package cronv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCrontab(t *testing.T) {
	line := "01 04 1 2 3 /usr/bin/somedirectory/somecommand1"
	r, _ := parseCrontab(line)
	assert.NotNil(t, r)
	assert.Equal(t, r.Line, line)
	assert.Equal(t, r.Schedule.Minute, "01")
	assert.Equal(t, r.Schedule.Hour, "04")
	assert.Equal(t, r.Schedule.DayOfMonth, "1")
	assert.Equal(t, r.Schedule.Month, "2")
	assert.Equal(t, r.Schedule.DayOfWeek, "3")
	assert.Equal(t, r.Job, "/usr/bin/somedirectory/somecommand1")
}

func TestParseCrontabInvalidTask(t *testing.T) {
	line := "MAILTO=example.com"
	_, err := parseCrontab(line)
	assert.NotNil(t, err)
}

func TestIsRunningEveryMinutesFalseCase(t *testing.T) {
	c, _ := parseCrontab("3 * * * *")
	assert.False(t, c.isRunningEveryMinutes())

	c2, _ := parseCrontab("* * * * 1")
	assert.False(t, c2.isRunningEveryMinutes())
}

func TestIsRunningEveryMinutesTrueCase(t *testing.T) {
	c, _ := parseCrontab("* * * * *")
	assert.True(t, c.isRunningEveryMinutes())

	c2, _ := parseCrontab("*/1 * * * *")
	assert.True(t, c2.isRunningEveryMinutes())
}
