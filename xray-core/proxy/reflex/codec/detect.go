package codec

import "bytes"

// LooksLikeHTTPPost returns true if peeked bytes resemble an HTTP POST request.
// We keep it intentionally conservative: if unsure, return false so fallback can handle.
func LooksLikeHTTPPost(peeked []byte) bool {
	// Typical start: "POST /..."
	if len(peeked) < 5 {
		return false
	}
	if !bytes.HasPrefix(peeked, []byte("POST ")) {
		return false
	}

	// Extra hint: must contain "HTTP/1." somewhere in the first ~64 bytes.
	// (We avoid heavy parsing here.)
	if bytes.Contains(peeked, []byte("HTTP/1.")) {
		return true
	}
	return false
}
