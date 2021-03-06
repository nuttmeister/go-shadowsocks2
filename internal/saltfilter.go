package internal

// A shared instance used for checking salt repeat
var saltfilter *BloomRing

const (
	defBloomCapacity = 1000000  // Default Capacity
	defBloomFPR      = 0.000001 // Default False Positive Rate
	defBloomSlot     = 10       // Default Slot
)

func init() {
	setup(defBloomSlot, defBloomCapacity, defBloomFPR)
}

// setup will setup the bloom filter.
func setup(slot int, capacity int, fpr float64) {
	saltfilter = NewBloomRing(slot, capacity, fpr)
}

// TestSalt returns true if salt is repeated
func TestSalt(b []byte) bool {
	// If nil means feature disabled, return false to bypass salt repeat detection
	if saltfilter == nil {
		return false
	}
	return saltfilter.Test(b)
}

// AddSalt salt to filter
func AddSalt(b []byte) {
	// If nil means feature disabled
	if saltfilter == nil {
		return
	}
	saltfilter.Add(b)
}
