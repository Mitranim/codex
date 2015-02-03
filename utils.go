package codex

// Public and private utility functions.

import (
	"math/rand"
	"time"
)

/***************************** Public Functions ******************************/

// Static generator functions exposed by the package.

// Takes a sample group of words, analyses their traits, and builds a set of all
// synthetic words that may derived from those traits. This should only be used
// for very small samples. More than just a handful of sample words causes a
// combinatorial explosion, takes a lot of time to calculate, and produces too
// many results to be useful. The number of results can easily reach hundreds of
// thousands for just a dozen of sample words.
func Words(words []string) (Set, error) {
	traits, err := NewTraits(words)
	if err != nil {
		return nil, err
	}
	return traits.Words(), nil
}

// Takes a sample group of words and a count limiter. Analyses the words and
// builds a random sample of synthetic words that may be derived from those
// traits, limited to the given count.
func WordsN(words []string, num int) (Set, error) {
	state, err := NewState(words)
	if err != nil {
		return nil, err
	}
	return state.WordsN(num), nil
}

/********************************* Utilities *********************************/

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Takes a word and splits it into a series of known glyphs representing sounds.
func getSounds(word string, known Set) ([]string, error) {
	sounds := make([]string, 0, len(word))
	// Loop over the word, matching known glyphs. Break if no match is found.
	for i := 0; i < len(word); i++ {
		// Check for a known digraph.
		if i+2 <= len(word) && known.Has(word[i:i+2]) {
			sounds = append(sounds, word[i:i+2])
			i++
			// Check for a known monograph.
		} else if known.Has(word[i : i+1]) {
			sounds = append(sounds, word[i:i+1])
			// Otherwise return an error.
		} else {
			return nil, errType("encountered unknown symbol")
		}
	}
	// Return the found glyphs.
	return sounds, nil
}

// Takes a sequence of sounds and returns the set of consequtive pairs that
// occur in this sequence.
func getPairs(sounds []string) (pairs PairSet) {
	for i := 0; i < len(sounds)-1; i++ {
		pairs.Add([2]string{sounds[i], sounds[i+1]})
	}
	return
}

// Takes a set of pairs of sounds and adds their reverses.
func addReversePairs(pairs PairSet) {
	for key := range pairs {
		pairs.Add([2]string{key[1], key[0]})
	}
}

// Checks if the given word is too short or too long.
func validLength(word string) bool {
	return len(word) > 1 && len(word) < 65
}

// Republished rand.Perm.
func permutate(length int) []int {
	return rand.Perm(length)
}

// Shuffles a slice of strings in-place, using the Fisher–Yates method.
func shuffle(values []string) {
	for i := range values {
		j := rand.Intn(i + 1)
		values[i], values[j] = values[j], values[i]
	}
}

// Returns the set of first values from the given pairs as a slice.
func firstValues(pairs PairSet) (results []string) {
	values := Set{}
	for pair := range pairs {
		values.Add(pair[0])
	}
	results = make([]string, 0, len(values))
	for value := range values {
		results = append(results, value)
	}
	return
}

// Returns the set of second values from the given pairs that begin with the
// given first value as a slice.
func secondMatching(pairs PairSet, first string) (results []string) {
	results = []string{}
	for pair := range pairs {
		if pair[0] != first {
			continue
		}
		results = append(results, pair[1])
	}
	return
}

// Version of firstValues() that shuffles the results.
func randFirsts(pairs PairSet) (results []string) {
	results = firstValues(pairs)
	shuffle(results)
	return
}

// Version of secondMatching() that shuffles the results.
func randSeconds(pairs PairSet, first string) (results []string) {
	results = secondMatching(pairs, first)
	shuffle(results)
	return
}

// Panic message used when breaking out from recursive iterations early.
const panicMsg = "early exit through panic"

// Wrapper for panic used when breaking out from recursive iterations early.
func interrupt() {
	panic(panicMsg)
}

// Wrapper for recovery from early iteration breakout through panic.
func aid() {
	msg := recover()
	if msg != nil && msg != panicMsg {
		panic(msg)
	}
}

/********************************** PairSet **********************************/

// PairSet behaves like a set of pairs of strings. Performance note: tried a
// slice version, and it significantly decreased the package's benchmark
// performance. Sticking with a map for now.
type PairSet map[[2]string]struct{}

// Creates a new set from the given keys. Usage:
//   PairSet.New(nil, [2]string{"one", "other"})
func (PairSet) New(keys ...[2]string) PairSet {
	set := make(PairSet, len(keys))
	for _, key := range keys {
		set.Add(key)
	}
	return set
}

// Adds the given element.
func (this *PairSet) Add(key [2]string) {
	if *this == nil {
		*this = PairSet{}
	}
	(*this)[key] = struct{}{}
}

// Deletes the given element.
func (this *PairSet) Del(key [2]string) {
	delete((*this), key)
}

// Checks for the presence of the given element.
func (this *PairSet) Has(key [2]string) bool {
	_, ok := (*this)[key]
	return ok
}

/*
// Commented out to avoid depending on fmt. If we include fmt at some point,
// this should be uncommented.

// Prints itself nicely in fmt(%#v).
func (this PairSet) GoString() string {
	keys := make([]string, 0, len(this))
	for key := range this {
		keys = append(keys, fmt.Sprintf("{%#v, %#v}", key[0], key[1]))
	}
	return "{" + strings.Join(keys, ", ") + "}"
}

// Prints itself nicely in println().
func (this PairSet) String() string {
	return this.GoString()
}
*/

/********************************** errType **********************************/

type errType string

func (this errType) Error() string {
	return string(this)
}
