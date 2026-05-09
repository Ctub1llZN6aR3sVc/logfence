package buffer

// Cap exposes the internal capacity for white-box tests.
func (r *Ring) Cap() int {
	return r.cap
}

// Len exposes the current number of entries stored for white-box tests.
func (r *Ring) Len() int {
	return r.len
}

// Head exposes the current head index for white-box tests.
func (r *Ring) Head() int {
	return r.head
}
