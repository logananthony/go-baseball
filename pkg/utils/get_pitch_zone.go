package utils

const (
    leftEdge   = -0.71
    rightEdge  = 0.71
    topEdge    = 3.5
    bottomEdge = 1.5

    bufferX = 0.25 
    bufferZ = 0.25
)


func GetPitchZone(plateX, plateZ float64) int {
    // Inside the 9-zone strike zone (1-9)
    if plateX >= leftEdge && plateX <= rightEdge && plateZ >= bottomEdge && plateZ <= topEdge {
        var col, row int
        if plateX < leftEdge/3 {
            col = 1
        } else if plateX > rightEdge/3 {
            col = 3
        } else {
            col = 2
        }
        if plateZ > (topEdge - (topEdge-bottomEdge)/3) {
            row = 1
        } else if plateZ < (bottomEdge + (topEdge-bottomEdge)/3) {
            row = 3
        } else {
            row = 2
        }
        return (row-1)*3 + col // zones 1-9
    }

    // Infinite outer zones
    if plateZ > topEdge && plateX < leftEdge {
        return 11 // upper left (zone 11)
    }
    if plateZ > topEdge && plateX > rightEdge {
        return 12 // upper right (zone 12)
    }
    if plateZ < bottomEdge && plateX < leftEdge {
        return 13 // lower left (zone 13)
    }
    if plateZ < bottomEdge && plateX > rightEdge {
        return 14 // lower right (zone 14)
    }

    return 14 // FIX THIS
}


