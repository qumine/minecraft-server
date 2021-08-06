package wrapper

func (w *Wrapper) onLog(line string) {
	for reg, status := range logToStatus {
		if reg.MatchString(line) {
			w.Status = status
		}
	}
}
