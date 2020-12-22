package wrapper

func (w *Wrapper) onLog(line string) {
	for status, reg := range logToStatus {
		if reg.MatchString(line) {
			w.Status = status
		}
	}
}
