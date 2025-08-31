package views

// Unified scrolling methods that handle both metrics and regular content

// scrollUp scrolls up by the specified number of lines
func (m *ContentViewModel) scrollUp(lines int) {
	if m.showMetrics && m.metricsDisplay != nil {
		m.metricsDisplay.Scroll(-lines)
	} else {
		m.scrollY = max(0, m.scrollY-lines)
	}
}

// scrollDown scrolls down by the specified number of lines
func (m *ContentViewModel) scrollDown(lines int) {
	if m.showMetrics && m.metricsDisplay != nil {
		m.metricsDisplay.Scroll(lines)
	} else {
		m.scrollY = min(m.maxScroll, m.scrollY+lines)
	}
}

// scrollToTop scrolls to the top of the content
func (m *ContentViewModel) scrollToTop() {
	if m.showMetrics && m.metricsDisplay != nil {
		m.metricsDisplay.Scroll(-1000) // Large negative number to reach top
	} else {
		m.scrollY = 0
	}
}

// scrollToBottom scrolls to the bottom of the content
func (m *ContentViewModel) scrollToBottom() {
	if m.showMetrics && m.metricsDisplay != nil {
		m.metricsDisplay.Scroll(1000) // Large positive number to reach bottom
	} else {
		m.scrollY = m.maxScroll
	}
}
