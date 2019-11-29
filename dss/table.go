package dss

// dssTable tells LOC about the files being installed and hashes for the files
type dssTable struct {
	Priority       int    `sil:"F2727" default:"30"` // Priority level of install, used for layering, usually should be 30
	Author         string `sil:"F2728" default:"KYGL"`
	Option         string `sil:"F2729"` // name of the option
	Destination    string `sil:"F2730"` // folder the file goes into
	Script         string `sil:"F2731"` // filename of the file
	FileDate       string `sil:"F2732"`
	LastChangeDate string `sil:"F253" default:"@DJSF @FMT(T6F,@NOW)"`
	Signature      string `sil:"F2733"`
}
