package cronv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCronvCommand(t *testing.T) {
	baseTime := time.Date(1985, 12, 8, 11, 30, 0, 0, time.UTC)

	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want *Option
	}{
		{
			name: "set default properties",
			args: args{
				t: baseTime,
			},
			want: &Option{
				OutputFilePath: optDefaultOutputPath,
				Duration:       optDefaultDuration,
				FromDate:       "1985/12/08",
				FromTime:       "11:30",
				Title:          optDefaultTitle,
				Width:          optDefaultWidth,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret := NewCronvOption(tt.args.t)
			assert.Equal(t, tt.want, ret)
		})
	}
}

func TestCommand_toDurationMinutes(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		want     float64
		wantErr  bool
	}{
		{
			name:     "valid (1m)",
			duration: "1m",
			want:     float64(1),
		},
		{
			name:     "valid (2h)",
			duration: "2h",
			want:     float64(1 * 60 * 2),
		},
		{
			name:     "valid (3d)",
			duration: "3d",
			want:     float64(1 * 60 * 24 * 3),
		},
		{
			name:     "invalid (not a duration string)",
			duration: "INVALID",
			wantErr:  true,
		},
		{
			name:     "invalid (unknown suffix)",
			duration: "1F",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCronvOption(time.Now())
			c.Duration = tt.duration
			ret, err := c.toDurationMinutes()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, ret)
		})
	}
}

func TestCommand_toFromTime(t *testing.T) {
	tests := []struct {
		name     string
		fromDate string
		fromTime string
		want     time.Time
		wantErr  bool
	}{
		{
			name:     "successfully",
			fromDate: "2016/11/08",
			fromTime: "01:30",
			want:     time.Date(2016, 11, 8, 1, 30, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCronvOption(time.Now())
			c.FromDate = tt.fromDate
			c.FromTime = tt.fromTime
			ret, err := c.toFromTime()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, ret)
		})
	}
}
