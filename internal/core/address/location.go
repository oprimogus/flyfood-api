package address

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type Location struct {
    Point geom.Point
}

func (l Location) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Lat float64 `json:"latitude"`
        Lng float64 `json:"longitude"`
    }{
        Lat: l.Point.Y(),
        Lng: l.Point.X(),
    })
}

func (l *Location) UnmarshalJSON(data []byte) error {
    var v struct {
        Lat float64 `json:"latitude"`
        Lng float64 `json:"longitude"`
    }
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }
    l.Point = *geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{v.Lng, v.Lat})
    return nil
}

func (l *Location) Scan(val any) error {
    b, ok := val.([]byte)
    if !ok {
        return fmt.Errorf("expected []byte, got %T", val)
    }
    got, err := ewkb.Unmarshal(b)
    if err != nil {
        return err
    }
    point, ok := got.(*geom.Point)
    if !ok {
        return fmt.Errorf("expected *geom.Point, got %T", got)
    }
    l.Point = *point
    return nil
}

func (l *Location) ScanText(v pgtype.Text) error {
    // pgx passa geography como text em formato hex WKB
    b, err := hex.DecodeString(v.String)
    if err != nil {
        return err
    }
    got, err := ewkb.Unmarshal(b)
    if err != nil {
        return err
    }
    point, ok := got.(*geom.Point)
    if !ok {
        return fmt.Errorf("expected *geom.Point, got %T", got)
    }
    l.Point = *point
    return nil
}

func (l Location) Lat() float64 { return l.Point.Y() }
func (l Location) Lng() float64 { return l.Point.X() }