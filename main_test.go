package main

import (
	"testing"
)

const maxAttempts = 1000

func Test_riddler_shoot(t *testing.T) {
	// Quickly verify that all three scores are possible, and we don't have an off-by-one error or such
	r := newRiddler()

	var five, nine, ten bool = false, false, false

	for i := 0; i < maxAttempts; i++ {
		score := r.shoot()
		if score == 5 {
			five = true
		} else if score == 9 {
			nine = true
		} else if score == 10 {
			ten = true
		} else {
			t.Fatalf("Unexpected score: %d\n", score)
		}

		// All possible scores validated
		if five == true && nine == true && ten == true {
			return
		}
	}

	t.Fatalf("Not all desired score options rolled. Five (%t), Nine (%t), Ten (%t)\n", five, nine, ten)
}
