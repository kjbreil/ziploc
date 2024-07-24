package install

// func Test_makeSrc(t *testing.T) {
// 	m := makeSrc("@RUNOPTIONS\\Plugin_NCBP_LaneLink\\XchDev\\LaneLink\\*")
//
// 	lfm, err := script.ProcessReader(strings.NewReader(`
// @WIZRPL(SRC_PATH=@RUNOPTIONS\Plugin_NCBP_SCO_Universal\XchDev\SmartCheckLink\*);
// @WIZRPL(TAR_PATH=@RUNXchDev\SmartCheckLink\);
// @FMT(CMP,@DBHOT(FINDFIRST,@RUNXchDev\SmartCheckLink\SmartCheckLink.exe,#)=,'@EXEC(XCH=COPY)');
//
// @WIZRPL(SRC_PATH=@RUNOPTIONS\Plugin_NCBP_SCO_Universal\XchDev\SmartCheckLink\runtimes\win\lib\net6.0\*);
// @WIZRPL(TAR_PATH=@RUNXchDev\SmartCheckLink\runtimes\win\lib\net6.0\);
// @FMT(CMP,@DBHOT(FINDFIRST,@RUNXchDev\SmartCheckLink\runtimes\win\lib\net6.0\,#)=,'®EXEC(XCH=COPY)');
// `))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	fmt.Println(m.String())
// 	fmt.Println(lfm.String())
// }
//
// func Test_Install(t *testing.T) {
// 	i := New("Plugin_NCBP_SCO_Universal")
// 	i.CopyFolder("XchDev\\SmartCheckLink")
// 	i.CopyFolder("XchDev\\SmartCheckLink\\runtimes\\win\\lib\\net6.0")
//
// 	fmt.Println(i.processed.String())
// }
