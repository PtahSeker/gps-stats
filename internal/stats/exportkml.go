// +build testkml

package stats

import (
    "fmt"
    "os"
)

// ExportManeuversToKML exports a slice of Tracks (subtracks) to a KML file with 4 styles:
// portJibe, portTack, starboardJibe, starboardTack.
// Each Track should have tackType and turnType set.
func ExportManeuversToKML(tracks []Track, windDir float64, filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    // Write KML header and styles
    fmt.Fprintln(f, `<?xml version="1.0" encoding="UTF-8"?>`)
    fmt.Fprintln(f, `<kml xmlns="http://www.opengis.net/kml/2.2"><Document>`)
    fmt.Fprintln(f, `<Style id="portJibe"><LineStyle><color>ff00ffff</color><width>1</width></LineStyle></Style>`)
    fmt.Fprintln(f, `<Style id="portTack"><LineStyle><color>ff0000ff</color><width>1</width></LineStyle></Style>`)
    fmt.Fprintln(f, `<Style id="starboardJibe"><LineStyle><color>ff00ff00</color><width>1</width></LineStyle></Style>`)
    fmt.Fprintln(f, `<Style id="starboardTack"><LineStyle><color>ffff0000</color><width>1</width></LineStyle></Style>`)
    fmt.Fprintln(f, `<Style id="arrow"><IconStyle><scale>1</scale><Icon><href>http://maps.google.com/mapfiles/kml/shapes/arrow.png</href></Icon></IconStyle></Style>`)

    for i, tr := range tracks {
        style := "unknown"
        name := "Unknown"
        var timeStr string
        if len(tr.ps) > 0 {
            timeStr = tr.ps[0].ts.Format("2006-01-02 15:04:05")
        } else {
            timeStr = "no_time"
        }
        switch {
        case tr.tackType == PortTack && detectTurnType(tr.ps, windDir) == JibeTurn:
            style = "portJibe"
            name = fmt.Sprintf("#%d Port Jibe (%s)", i+1, timeStr)
        case tr.tackType == PortTack && detectTurnType(tr.ps, windDir) == TackTurn:
            style = "portTack"
            name = fmt.Sprintf("#%d Port Tack (%s)", i+1, timeStr)
        case tr.tackType == StarboardTack && detectTurnType(tr.ps, windDir) == JibeTurn:
            style = "starboardJibe"
            name = fmt.Sprintf("#%d Starboard Jibe (%s)", i+1, timeStr)
        case tr.tackType == StarboardTack && detectTurnType(tr.ps, windDir) == TackTurn:
            style = "starboardTack"
            name = fmt.Sprintf("#%d Starboard Tack (%s)", i+1, timeStr)
        }
        fmt.Fprintf(f, `<Placemark><name>%s</name><styleUrl>#%s</styleUrl><LineString><coordinates>`, name, style)
        for _, p := range tr.ps {
            fmt.Fprintf(f, "%f,%f ", p.lon, p.lat)
        }
        fmt.Fprintln(f, `</coordinates></LineString></Placemark>`)

        if len(tr.ps) >= 2 {
            first := tr.ps[0]
            second := tr.ps[1]
            h := heading(second, first)
            fmt.Fprintf(f, `<Placemark><name>%s</name><styleUrl>#arrow</styleUrl><Point><coordinates>%f,%f</coordinates></Point><Style><IconStyle><heading>%.1f</heading></IconStyle></Style></Placemark>`, name, first.lon, first.lat, h)
        }
    }

    fmt.Fprintln(f, `</Document></kml>`)
    return nil
}