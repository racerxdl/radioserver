package server

func (rs *RadioServer) onSamples(samples []complex64) {
	rs.sessionLock.Lock()
	defer rs.sessionLock.Unlock()

	for _, v := range rs.sessions {
		v.CG.PushSamples(samples)
	}
}
