package quoteapi

import (
	"os"
	"testing"
)

type testRecord struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// TestLoadBoltstore Kicks the tires on the bolt store
func TestLoadBoltstore(t *testing.T) {
	boltPath := "./test.boltdb"

	if _, err := os.Stat(boltPath); os.IsExist(err) {
		t.Logf("Deleting %v before test", boltPath)
		os.Remove(boltPath)
	}

	boltConnection, err := GetBoltStore(boltPath)

	if err != nil {
		t.Fatalf("Error connectiong to bolt store: %v", err)
	}

	defer boltConnection.Close()
	defer os.Remove(boltPath)

	myRecord := testRecord{Name: "myName", Comment: "my comment"}
	if err := boltConnection.SaveRecord("testRecord", &myRecord, myRecord.Name); err != nil {
		t.Fatalf("Failed to save record: %v", err)
	}

	var loadedRecord testRecord
	if err := boltConnection.LoadRecord("testRecord", &loadedRecord, myRecord.Name); err != nil {
		t.Fatalf("Failed to load record: %v", err)
	}

	if loadedRecord.Name != myRecord.Name || loadedRecord.Comment != myRecord.Comment {
		t.Fatalf("Records don't match")
	}
}
