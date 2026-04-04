package youtube

import (
	"regexp"
	"strings"
)

var (
	// Matches (Official Video), [Official Audio], (Lyric Video), (Audio), (Visualizer), etc.
	videoSuffixRe = regexp.MustCompile(`(?i)\s*[\(\[]\s*(official\s+(video|audio|music\s+video|lyric\s+video|visualizer|hd)|lyric\s+video|lyrics|audio|visualizer|music\s+video|mv|hd|hq|4k|live)\s*[\)\]]`)

	// Matches "- Topic" suffix on YouTube Music auto-generated channels
	topicSuffixRe = regexp.MustCompile(`(?i)\s*-\s*topic\s*$`)

	// Matches "VEVO" suffix on artist channels
	vevoSuffixRe = regexp.MustCompile(`(?i)vevo\s*$`)
)

// ParseVideoTitle extracts track name and artist from a YouTube video title and channel name.
// YouTube Music videos commonly use "Artist - Track Name" format.
func ParseVideoTitle(title, channelTitle string) (trackName, artistName string) {
	// Strip video-type suffixes: (Official Video), [Official Audio], etc.
	cleaned := videoSuffixRe.ReplaceAllString(title, "")
	cleaned = strings.TrimSpace(cleaned)

	// Try splitting on " - " (most common YouTube Music format)
	if parts := strings.SplitN(cleaned, " - ", 2); len(parts) == 2 {
		artist := strings.TrimSpace(parts[0])
		track := strings.TrimSpace(parts[1])
		if artist != "" && track != "" {
			return track, artist
		}
	}

	// Fallback: title is track name, channel is artist
	artistName = cleanChannelName(channelTitle)
	return cleaned, artistName
}

func cleanChannelName(name string) string {
	name = topicSuffixRe.ReplaceAllString(name, "")
	name = vevoSuffixRe.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}
