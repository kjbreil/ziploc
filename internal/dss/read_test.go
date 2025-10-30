package dss

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	dssContent := `@CREATE(DSS_TAB,DSS);
INSERT INTO HEADER_DCT VALUES
('HM','05411875','MANUAL','PAL',,,2019266,0000,2019266,0000,,'ADDRPL','ADDRPL FROM GO',,,,,,,,,);


CREATE VIEW dss_LOAD AS SELECT F2727,F2728,F2729,F2730,F2731,F2732,F253,F2733 FROM dss_DCT;

INSERT INTO dss_LOAD VALUES
(45,'ZIPLOC','A_SAMPLE','OFFICE','Setting.ini','2021001 16:24:30','@DJSF @FMT(T6F,@NOW)','568243DD')
(45,'ZIPLOC','A_SAMPLE','SSM','TRS_CLT_Pos.ssm','2021001 15:54:24','@DJSF @FMT(T6F,@NOW)','BA144D2D')
(45,'ZIPLOC','A_SAMPLE','LANE','INSTALL.SQL','2021001 15:54:24','@DJSF @FMT(T6F,@NOW)','62A3AA7C');


DELETE FROM DSS_TAB WHERE F2729 IN (SELECT F2729 FROM DSS_LOAD);
@UPDATE_BATCH(JOB=ADDRPL,TAR=DSS_TAB,KEY=F2729=:F2729 AND F2730=:F2730 AND F2731=:F2731,
SRC=SELECT * FROM DSS_LOAD);
DROP TABLE DSS_LOAD;
`

	reader := strings.NewReader(dssContent)
	entries, err := Read(reader)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Read() returned %d entries, expected 3", len(entries))
	}

	// Test first entry
	if entries[0].Priority != 45 {
		t.Errorf("Entry 0 Priority = %d, expected 45", entries[0].Priority)
	}
	if entries[0].Author != "ZIPLOC" {
		t.Errorf("Entry 0 Author = %s, expected ZIPLOC", entries[0].Author)
	}
	if entries[0].Option != "A_SAMPLE" {
		t.Errorf("Entry 0 Option = %s, expected A_SAMPLE", entries[0].Option)
	}
	if entries[0].Destination != "OFFICE" {
		t.Errorf("Entry 0 Destination = %s, expected OFFICE", entries[0].Destination)
	}
	if entries[0].Script != "Setting.ini" {
		t.Errorf("Entry 0 Script = %s, expected Setting.ini", entries[0].Script)
	}
	if entries[0].Signature != "568243DD" {
		t.Errorf("Entry 0 Signature = %s, expected 568243DD", entries[0].Signature)
	}

	// Test second entry
	if entries[1].Destination != "SSM" {
		t.Errorf("Entry 1 Destination = %s, expected SSM", entries[1].Destination)
	}
	if entries[1].Script != "TRS_CLT_Pos.ssm" {
		t.Errorf("Entry 1 Script = %s, expected TRS_CLT_Pos.ssm", entries[1].Script)
	}
	if entries[1].Signature != "BA144D2D" {
		t.Errorf("Entry 1 Signature = %s, expected BA144D2D", entries[1].Signature)
	}

	// Test third entry
	if entries[2].Destination != "LANE" {
		t.Errorf("Entry 2 Destination = %s, expected LANE", entries[2].Destination)
	}
	if entries[2].Script != "INSTALL.SQL" {
		t.Errorf("Entry 2 Script = %s, expected INSTALL.SQL", entries[2].Script)
	}
	if entries[2].Signature != "62A3AA7C" {
		t.Errorf("Entry 2 Signature = %s, expected 62A3AA7C", entries[2].Signature)
	}
}
