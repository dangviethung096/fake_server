package core

import "testing"

func TestOpenDBWithSuccessResponse(t *testing.T) {
	// Set up test cases
	testCase := struct {
		name    string
		dbInfo  DBInfo
		wantErr bool
	}{
		name: "Valid DBInfo",
		dbInfo: DBInfo{
			FilePath: "test.db",
		},
		wantErr: false,
	}

	// Call the function being tested
	db := openDBConnection(testCase.dbInfo)
	db.Close()
}
