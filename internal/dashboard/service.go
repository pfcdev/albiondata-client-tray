// internal/dashboard/service.go
package dashboard

import "github.com/ao-data/albiondata-client/internal/autostart"

// DashboardService is bound to the Wails frontend (see application.NewService
// in albiondata-client.go). Its methods back-fill the dashboard window's
// state when it's shown; ongoing updates arrive via the "status:changed",
// "counters:snapshot", and "log:line" events instead.
type DashboardService struct{}

// GetStatus returns the current status snapshot.
func (s *DashboardService) GetStatus() Status {
	return GetStatus()
}

// GetUploadCounts returns the current upload counters by topic.
func (s *DashboardService) GetUploadCounts() map[string]int64 {
	return GetUploadCounts()
}

// GetRecentLogs returns the currently buffered recent log lines.
func (s *DashboardService) GetRecentLogs() []LogLine {
	return GetRecentLogs()
}

// StartupSupported reports whether the Windows logon task is available.
func (s *DashboardService) StartupSupported() bool {
	return autostart.Available()
}

// GetStartupEnabled reports whether the logon task is enabled.
func (s *DashboardService) GetStartupEnabled() (bool, error) {
	return autostart.Status()
}

// SetStartupEnabled changes the logon task, prompting for UAC consent if
// the app was launched without administrator privileges.
func (s *DashboardService) SetStartupEnabled(enabled bool) (bool, error) {
	return autostart.Set(enabled)
}
