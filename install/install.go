package install

// Install holds information to build out an INSTALL.SQL file
// type Install struct {
// 	optionName string
// 	processed  *script.Processed
// }
//
// func New(optionName string) *Install {
// 	return &Install{
// 		optionName: optionName,
// 		processed:  script.New(),
// 	}
// }
//
// func (i *Install) CopyFolder(folder string) {
// 	srcPath := strings.Join([]string{"@RUNOPTIONS", i.optionName, folder, "*"}, "\\")
// 	i.processed.Add(macros.CreateWizrpl(macro.StringRaw("SRC_PATH"), macro.StringRaw(srcPath)))
// 	dstPath := strings.Join([]string{"@RUN", folder, "*"}, "\\")
// 	i.processed.Add(macros.CreateWizrpl(macro.StringRaw("TAR_PATH"), macro.StringRaw(dstPath)))
// 	i.processed.Add(macro.StringCRLF())
// }
//
// // @WIZRPL(SRC_PATH=@RUNOPTIONS\Plugin_NCBP_LaneLink\XchDev\LaneLink\*);
// // @WIZRPL(TAR_PATH=@RUNXchDev\LaneLink\);
// // @FMT(CMP,@DBHOT(FINDFIRST,@RUNXchDev\LaneLink\LaneLink.exe,#)=,'®EXEC(XCH=COPY)');
//
// func makeSrc(srcPath string) *macro.Main {
// 	m := macros.CreateWizrpl(macro.StringRaw("SRC_PATH"), macro.StringRaw(srcPath))
//
// 	return m
// }
