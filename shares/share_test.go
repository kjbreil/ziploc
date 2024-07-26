package shares

import "testing"

func Test_listShares(t *testing.T) {
	s, err := New("\\\\PHYSREG\\REG", "POS_NCBP", "*Loc!Sms901")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Connect()
	defer s.Close()

	if err != nil {
		t.Fatal(err)
	}
	s.WalkDir("Storeman")

}
