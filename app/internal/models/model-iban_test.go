package models

import "testing"

func TestIBAN(t *testing.T) {
	for _, s := range []string{"EE382200221020145685", "DE89370400440532013000", "GB82WEST12345698765432"} {
		if !ValidIBAN(s) {
			t.Errorf("%s is valid", s)
		}
	}
	for _, s := range []string{
		"", "EE38", "EE382200221020145686", "DE89370400440532013001", "EE38220022102014568",
		"XX82WEST12345698765432", "EE38 2200 2210 2014 5685", "ee382200221020145685", "EE38220022102014568$",
	} {
		if ValidIBAN(s) {
			t.Errorf("%q is not valid", s)
		}
	}
	if got := NormalizeIBAN(" ee38 2200\t2210 2014 5685 "); got != "EE382200221020145685" {
		t.Errorf("NormalizeIBAN = %q", got)
	}
	if got := FormatIBAN("ee382200221020145685"); got != "EE38 2200 2210 2014 5685" {
		t.Errorf("FormatIBAN = %q", got)
	}
}
