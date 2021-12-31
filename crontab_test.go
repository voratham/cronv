package cronv

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name        string
		args        args
		wantCrontab *Crontab
		wantExtra   *Extra
		wantErr     bool
	}{
		{
			name: "parse successfully",
			args: args{
				line: "01 04 1 2 3	/usr/bin/somedirectory/somecommand1",
			},
			wantCrontab: &Crontab{
				Line: "01 04 1 2 3	/usr/bin/somedirectory/somecommand1",
				Schedule: &Schedule{
					Minute:     "01",
					Hour:       "04",
					DayOfMonth: "1",
					Month:      "2",
					DayOfWeek:  "3",
				},
				Job: "/usr/bin/somedirectory/somecommand1",
			},
		},
		{
			name: "invalid task",
			args: args{
				line: "MAILTO=example.com",
			},
			wantErr: true,
		},
		{
			name: "alias",
			args: args{
				line: "@hourly /path/to/do/something arg1",
			},
			wantCrontab: &Crontab{
				Line: "@hourly /path/to/do/something arg1",
				Schedule: &Schedule{
					Alias: "@hourly",
				},
				Job: "/path/to/do/something arg1",
			},
		},
		{
			name: "extra",
			args: args{
				line: "@reboot root /path/to/do/something arg1 arg2 arg3",
			},
			wantExtra: &Extra{
				Line:  "@reboot root /path/to/do/something arg1 arg2 arg3",
				Label: "@reboot",
				Job:   "root /path/to/do/something arg1 arg2 arg3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crontab, extra, err := parse(context.Background(), tt.args.line)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantCrontab, crontab)
			assert.Equal(t, tt.wantExtra, extra)
		})
	}
}

func TestCrontab_isRunningEveryMinutes(t *testing.T) {
	tests := []struct {
		name    string
		crontab string
		want    bool
	}{
		{
			name:    "false (only min=3)",
			crontab: "3 * * * *",
			want:    false,
		},
		{
			name:    "false (every minutes but only hour=1)",
			crontab: "* 1 * * *",
			want:    false,
		},
		{
			name:    "every minutes by wildcard",
			crontab: "* * * * *",
			want:    true,
		},
		{
			name:    "every minutes",
			crontab: "*/1 * * * *",
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _, err := parse(context.Background(), tt.crontab)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, c.isRunningEveryMinutes())
		})
	}
}
