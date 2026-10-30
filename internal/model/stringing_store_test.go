package model

import (
	"gorm.io/gorm/schema"
	"sync"
	"testing"
)

func TestStoreColumnsRoundTripAndClear(t *testing.T) {
	original := &CourtLocation{Address: "穿线店地址", Latitude: 0, Longitude: 108.1234567, CoordinateSystem: "gcj02"}
	columns := StoreColumnsFromLocation(original)
	original.Latitude = 5
	restored := columns.Location()
	if restored == nil || restored.Latitude != 0 || restored.Longitude != 108.1234567 || restored.Address != "穿线店地址" || restored.CoordinateSystem != "gcj02" {
		t.Fatalf("location round trip failed: %+v", restored)
	}
	record := RacketStringingRecord{StoreName: "穿线店", StoreColumns: columns}
	response := NewStringingRecordResponse(record)
	if response.StoreLocation == nil || response.StoreLocation.Latitude != 0 {
		t.Fatal("response must map independent columns")
	}
	record.StoreColumns = StoreColumnsFromLocation(nil)
	record.StoreName = ""
	response = NewStringingRecordResponse(record)
	if response.StoreLocation != nil || record.StoreAddress != "" || record.StoreCoordinateSystem != "" || record.StoreLatitude != nil || record.StoreLongitude != nil {
		t.Fatal("clearing must remove every location field")
	}
}

func TestStringingStoreColumnSchema(t *testing.T) {
	parsed, err := schema.Parse(&RacketStringingRecord{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatal(err)
	}
	for name, typ := range map[string]string{"store_address": "string", "store_latitude": "decimal(10,7)", "store_longitude": "decimal(10,7)", "store_coordinate_system": "string"} {
		field := parsed.FieldsByDBName[name]
		if field == nil || string(field.DataType) != typ {
			t.Fatalf("invalid column %s: %+v", name, field)
		}
	}
	if parsed.FieldsByDBName["store_location"] != nil {
		t.Fatal("JSON storage must not remain")
	}
}
