package migration

import (
	"math"
	"regexp"
	"strings"

	"tuneshift/internal/csvimport"
	"tuneshift/internal/tidal"
)

var (
	// Remove parenthetical and bracketed content: (Remaster), [Deluxe], etc.
	parentheticalRe = regexp.MustCompile(`\s*[\(\[].*?[\)\]]`)
	// Remove " - Remaster", " - 2004 Remaster", " - Original Mix", etc.
	dashSuffixRe = regexp.MustCompile(`\s*-\s*([\d]{4}\s*)?(remaster(ed)?|original.*|deluxe|bonus.*|single.*|mono|stereo|live).*$`)
	specialCharsRe = regexp.MustCompile(`[^\w\s]`)
	whitespaceRe   = regexp.MustCompile(`\s+`)
)

type Matcher struct {
	tidalClient *tidal.Client
}

func NewMatcher(tidalClient *tidal.Client) *Matcher {
	return &Matcher{tidalClient: tidalClient}
}

func (m *Matcher) MatchCSV(track csvimport.Track) (*tidal.Track, error) {
	// Strategy 1: ISRC lookup
	if track.ISRC != "" {
		result, err := m.tidalClient.SearchTrackByISRC(track.ISRC)
		if err == nil && result != nil {
			return result, nil
		}
	}

	// Strategy 2: Name + Artist search with fuzzy matching
	query := track.TrackName + " " + track.FirstArtist()
	candidates, err := m.tidalClient.SearchTrack(query, 10)
	if err != nil {
		return nil, err
	}

	return bestMatch(track.TrackName, track.FirstArtist(), track.DurationMS, candidates), nil
}

func bestMatch(sourceName, sourceArtist string, sourceDurationMS int, candidates []tidal.Track) *tidal.Track {
	var bestTrack *tidal.Track
	bestScore := 0.0

	normSource := normalize(sourceName)

	for i, candidate := range candidates {
		if shouldExclude(sourceName, candidate.Title) {
			continue
		}

		normCandidate := normalize(candidate.Title)

		// Multiple scoring strategies, take the best
		titleScore := maxFloat(
			wordOverlap(normSource, normCandidate),
			containsScore(normSource, normCandidate),
			// Also try without spaces for cases like "Reggaemylitis" vs "reggae mylitis"
			noSpaceMatch(normSource, normCandidate),
		)

		durationScore := durationMatch(sourceDurationMS, candidate.Duration*1000)

		score := titleScore*0.7 + durationScore*0.3

		if score > bestScore {
			bestScore = score
			bestTrack = &candidates[i]
		}
	}

	if bestScore < 0.45 {
		return nil
	}

	return bestTrack
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = parentheticalRe.ReplaceAllString(s, "")
	s = dashSuffixRe.ReplaceAllString(s, "")
	s = specialCharsRe.ReplaceAllString(s, " ")
	s = whitespaceRe.ReplaceAllString(strings.TrimSpace(s), " ")
	return s
}

func shouldExclude(a, b string) bool {
	aLower := strings.ToLower(a)
	bLower := strings.ToLower(b)

	exclusions := []string{"instrumental", "remix", "acapella", "a cappella", "karaoke"}
	for _, term := range exclusions {
		aHas := strings.Contains(aLower, term)
		bHas := strings.Contains(bLower, term)
		if aHas != bHas {
			return true
		}
	}
	return false
}

// wordOverlap: fraction of source words found in candidate
func wordOverlap(a, b string) float64 {
	if a == b {
		return 1.0
	}
	if a == "" || b == "" {
		return 0.0
	}

	aWords := strings.Fields(a)
	bWords := strings.Fields(b)

	matches := 0
	for _, aw := range aWords {
		for _, bw := range bWords {
			if aw == bw {
				matches++
				break
			}
		}
	}

	if len(aWords) == 0 {
		return 0
	}

	// Divide by source length, not max — "beat it" in "beat it state of shock" = 2/2 = 1.0
	return float64(matches) / float64(len(aWords))
}

// containsScore: if one string contains the other entirely
func containsScore(a, b string) float64 {
	if strings.Contains(b, a) || strings.Contains(a, b) {
		return 0.9
	}
	return 0
}

// noSpaceMatch: compare strings with all spaces removed (handles "Reggaemylitis" vs "reggae mylitis")
func noSpaceMatch(a, b string) float64 {
	aCompact := strings.ReplaceAll(a, " ", "")
	bCompact := strings.ReplaceAll(b, " ", "")
	if aCompact == bCompact {
		return 1.0
	}
	if strings.Contains(bCompact, aCompact) || strings.Contains(aCompact, bCompact) {
		return 0.85
	}
	return 0
}

func durationMatch(msA, msB int) float64 {
	diff := math.Abs(float64(msA - msB))
	if diff <= 3000 {
		return 1.0
	}
	if diff <= 10000 {
		return 0.5
	}
	return 0.0
}

func maxFloat(vals ...float64) float64 {
	m := vals[0]
	for _, v := range vals[1:] {
		if v > m {
			m = v
		}
	}
	return m
}
