package utils

const (
    leftEdge   = -0.71
    rightEdge  = 0.71
    topEdge    = 3.5
    bottomEdge = 1.5
)

// GetPitchZone maps Statcast-style plate coordinates (X, Z) to MLB 14 zones
func GetPitchZone(plateX, plateZ float64) int {
    zoneWidth := (rightEdge - leftEdge) / 3
    zoneHeight := (topEdge - bottomEdge) / 3

    // --- Inside strike zone: Zones 1–9 (3x3 grid) ---
    if plateX >= leftEdge && plateX <= rightEdge && plateZ >= bottomEdge && plateZ <= topEdge {
        col := int((plateX - leftEdge) / zoneWidth) + 1
        row := 3 - int((plateZ - bottomEdge) / zoneHeight)
        return (row-1)*3 + col
    }

    // --- Outside zones by quadrant (11–14) ---

    // Zone 11: Top-left (Z above, or X left, or both)
    if plateZ > topEdge && plateX < 0 || plateX < leftEdge && plateZ > bottomEdge {
        return 11
    }

    // Zone 12: Top-right
    if plateZ > topEdge && plateX >= 0 || plateX > rightEdge && plateZ > bottomEdge {
        return 12
    }

    // Zone 13: Bottom-left
    if plateZ < bottomEdge && plateX < 0 || plateX < leftEdge && plateZ < topEdge {
        return 13
    }

    // Zone 14: Bottom-right
    if plateZ < bottomEdge && plateX >= 0 || plateX > rightEdge && plateZ < topEdge {
        return 14
    }

    // --- Fallback (shouldn't happen) ---
    return 0
}

