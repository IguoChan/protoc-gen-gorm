package model

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestUserModelTableNameAndIndexOptions(t *testing.T) {
	if got := (&UserModel{}).TableName(); got != "user" {
		t.Fatalf("UserModel.TableName() = %q, want %q", got, "user")
	}
	if got := (&StudentModel{}).TableName(); got != "student" {
		t.Fatalf("StudentModel.TableName() = %q, want %q", got, "student")
	}

	db := c.DB.Session(&gorm.Session{DryRun: true}).Model(&UserModel{})
	stmt := UserModelId(123).Apply(db).Find(&UserModel{}).Statement
	if sql := stmt.SQL.String(); !strings.Contains(sql, "id= ?") {
		t.Fatalf("UserModelId SQL = %q, want id predicate", sql)
	}
	if len(stmt.Vars) != 1 || stmt.Vars[0] != int32(123) {
		t.Fatalf("UserModelId vars = %#v, want [123]", stmt.Vars)
	}
}

func TestUserModelJSONValueAndScan(t *testing.T) {
	var zeroIDs UserModelIdsList
	if got, err := Value(zeroIDs); err != nil || got != nil {
		t.Fatalf("zero Value() = (%#v, %v), want nil value and nil error", got, err)
	}

	ids := UserModelIdsList{1, 2, 3}
	raw, err := ids.Value()
	if err != nil {
		t.Fatal(err)
	}
	var scannedIDs UserModelIdsList
	if err = scannedIDs.Scan(raw); err != nil {
		t.Fatal(err)
	}
	if len(scannedIDs) != len(ids) || scannedIDs[2] != ids[2] {
		t.Fatalf("scanned ids = %#v, want %#v", scannedIDs, ids)
	}

	m := UserModelM1Map{"k": "v"}
	raw, err = m.Value()
	if err != nil {
		t.Fatal(err)
	}
	var scannedMap UserModelM1Map
	if err = scannedMap.Scan(string(raw.([]byte))); err != nil {
		t.Fatal(err)
	}
	if scannedMap["k"] != "v" {
		t.Fatalf("scanned map = %#v, want key k", scannedMap)
	}
	if err = scannedMap.Scan(1); err == nil {
		t.Fatal("Scan accepted unsupported value type")
	}
}

func TestUserProtoModelRoundTrip(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	model := testUser("round_trip")
	model.Id = 101
	model.CreateTime = now
	model.UpdateTime = now
	model.DeletedTime = gorm.DeletedAt{Time: now}
	model.Date = datatypes.Date(now)
	model.T1 = now
	model.T2 = now
	model.T3 = now
	model.D1 = datatypes.Date(now)

	pb, err := model.User()
	if err != nil {
		t.Fatal(err)
	}
	if pb.Status != Status_PENDING || pb.Status1 != Status_UNKNOWN || !bytes.Equal(pb.Extra, model.Extra) {
		t.Fatalf("unexpected proto conversion: %+v", pb)
	}
	if pb.Author.GetName() != model.Author.Name {
		t.Fatalf("author name = %q, want %q", pb.Author.GetName(), model.Author.Name)
	}

	roundTrip, err := pb.UserModel()
	if err != nil {
		t.Fatal(err)
	}
	if roundTrip.Id != model.Id || roundTrip.Name != model.Name || roundTrip.Status != model.Status {
		t.Fatalf("round trip model = %+v, want id/name/status from %+v", roundTrip, model)
	}
	if len(roundTrip.Ids) != len(model.Ids) || roundTrip.M1["k"] != "v" {
		t.Fatalf("round trip repeated/map fields = %#v %#v", roundTrip.Ids, roundTrip.M1)
	}
}
