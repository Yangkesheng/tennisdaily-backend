package model

type StringingStore struct {
	StoreName     string         `json:"storeName"`
	StoreLocation *CourtLocation `json:"storeLocation"`
}

// StoreColumns mirrors court storage; API DTOs expose a nested location object.
type StoreColumns struct {
	StoreAddress          string   `json:"-" gorm:"size:512;not null;default:''"`
	StoreLatitude         *float64 `json:"-" gorm:"type:decimal(10,7)"`
	StoreLongitude        *float64 `json:"-" gorm:"type:decimal(10,7)"`
	StoreCoordinateSystem string   `json:"-" gorm:"size:16;not null;default:''"`
}

func (c StoreColumns) Location() *CourtLocation {
	if c.StoreLatitude == nil || c.StoreLongitude == nil {
		return nil
	}
	return &CourtLocation{Address: c.StoreAddress, Latitude: *c.StoreLatitude, Longitude: *c.StoreLongitude, CoordinateSystem: c.StoreCoordinateSystem}
}

func StoreColumnsFromLocation(l *CourtLocation) StoreColumns {
	if l == nil {
		return StoreColumns{}
	}
	lat, lng := l.Latitude, l.Longitude
	return StoreColumns{StoreAddress: l.Address, StoreLatitude: &lat, StoreLongitude: &lng, StoreCoordinateSystem: l.CoordinateSystem}
}
