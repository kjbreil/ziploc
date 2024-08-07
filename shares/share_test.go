package shares

import "testing"

func Test_listShares(t *testing.T) {
	s, err := New("\\\\PHYSREG\\LANELINK", "POS_NCBP", "*Loc!Sms901")
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

func Test_readFile(t *testing.T) {
	s, err := New("\\\\PHYSREG\\LANELINK", "POS_NCBP", "*Loc!Sms901")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Connect()
	defer s.Close()

	if err != nil {
		t.Fatal(err)
	}
	data, err := s.ReadFile("startup.ini")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
