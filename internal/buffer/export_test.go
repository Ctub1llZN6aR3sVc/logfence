package buffer

// Cap exposes the internal capacity for white-box tests.
func (r *Ring) Cap() int {
	return r.cap
}
