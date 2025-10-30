package macro

import (
	"bytes"
	"testing"
)

func TestConvertUTF8ToWindows1252_SimpleASCII(t *testing.T) {
	// Simple ASCII text with LF should convert to CRLF
	input := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x0a, 0x57, 0x6f, 0x72, 0x6c, 0x64}
	// "Hello\nWorld" -> "Hello\r\nWorld"
	expected := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x0d, 0x0a, 0x57, 0x6f, 0x72, 0x6c, 0x64}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestConvertUTF8ToWindows1252_RegisteredSymbol(t *testing.T) {
	// UTF-8 registered symbol ® (0xc2 0xae) should convert to Windows-1252 (0xae)
	// "Test ® Symbol\n"
	input := []byte{0x54, 0x65, 0x73, 0x74, 0x20, 0xc2, 0xae, 0x20, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x0a}
	// "Test ® Symbol\r\n" in Windows-1252
	expected := []byte{0x54, 0x65, 0x73, 0x74, 0x20, 0xae, 0x20, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x0d, 0x0a}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestConvertUTF8ToWindows1252_CopyrightSymbol(t *testing.T) {
	// UTF-8 copyright symbol © (0xc2 0xa9) should convert to Windows-1252 (0xa9)
	// "Copyright © 2024\n"
	input := []byte{0x43, 0x6f, 0x70, 0x79, 0x72, 0x69, 0x67, 0x68, 0x74, 0x20, 0xc2, 0xa9, 0x20, 0x32, 0x30, 0x32, 0x34, 0x0a}
	// "Copyright © 2024\r\n" in Windows-1252
	expected := []byte{0x43, 0x6f, 0x70, 0x79, 0x72, 0x69, 0x67, 0x68, 0x74, 0x20, 0xa9, 0x20, 0x32, 0x30, 0x32, 0x34, 0x0d, 0x0a}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestConvertUTF8ToWindows1252_MixedSpecialChars(t *testing.T) {
	// Mix of registered and copyright symbols
	// "® and ©\n"
	input := []byte{0xc2, 0xae, 0x20, 0x61, 0x6e, 0x64, 0x20, 0xc2, 0xa9, 0x0a}
	// "® and ©\r\n" in Windows-1252
	expected := []byte{0xae, 0x20, 0x61, 0x6e, 0x64, 0x20, 0xa9, 0x0d, 0x0a}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestConvertUTF8ToWindows1252_FirstCharNewline(t *testing.T) {
	// Test edge case where first character is a newline
	// "\nHello"
	input := []byte{0x0a, 0x48, 0x65, 0x6c, 0x6c, 0x6f}
	// "\r\nHello"
	expected := []byte{0x0d, 0x0a, 0x48, 0x65, 0x6c, 0x6c, 0x6f}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestConvertUTF8ToWindows1252_EmptyInput(t *testing.T) {
	// Empty input should return empty output
	input := []byte{}
	expected := []byte{}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestConvertUTF8ToWindows1252_MultipleLines(t *testing.T) {
	// Multiple lines with special characters
	// "Line 1 ®\nLine 2 ©\nLine 3\n"
	input := []byte{
		0x4c, 0x69, 0x6e, 0x65, 0x20, 0x31, 0x20, 0xc2, 0xae, 0x0a,
		0x4c, 0x69, 0x6e, 0x65, 0x20, 0x32, 0x20, 0xc2, 0xa9, 0x0a,
		0x4c, 0x69, 0x6e, 0x65, 0x20, 0x33, 0x0a,
	}
	// "Line 1 ®\r\nLine 2 ©\r\nLine 3\r\n" in Windows-1252
	expected := []byte{
		0x4c, 0x69, 0x6e, 0x65, 0x20, 0x31, 0x20, 0xae, 0x0d, 0x0a,
		0x4c, 0x69, 0x6e, 0x65, 0x20, 0x32, 0x20, 0xa9, 0x0d, 0x0a,
		0x4c, 0x69, 0x6e, 0x65, 0x20, 0x33, 0x0d, 0x0a,
	}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestConvertUTF8ToWindows1252_PreservesExistingCRLF(t *testing.T) {
	// Text that already has CRLF should not add extra CR
	// "Hello\r\nWorld"
	input := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x0d, 0x0a, 0x57, 0x6f, 0x72, 0x6c, 0x64}
	// Should remain "Hello\r\nWorld"
	expected := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x0d, 0x0a, 0x57, 0x6f, 0x72, 0x6c, 0x64}

	result, err := ConvertUTF8ToWindows1252(input)
	if err != nil {
		t.Fatalf("ConvertUTF8ToWindows1252() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ConvertUTF8ToWindows1252() = %v, want %v", result, expected)
	}
}

func TestCRLF_BackwardCompatibility(t *testing.T) {
	// Test that CRLF still works for backward compatibility
	input := []byte{0x54, 0x65, 0x73, 0x74, 0x0a}
	expected := []byte{0x54, 0x65, 0x73, 0x74, 0x0d, 0x0a}

	result, err := CRLF(input)
	if err != nil {
		t.Fatalf("CRLF() error = %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("CRLF() = %v, want %v", result, expected)
	}
}
