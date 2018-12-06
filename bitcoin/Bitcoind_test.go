package bitcoin

import "testing"

func TestGetBlock(t *testing.T) {
	b, err := New("BSV")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	block, berr, err := b.GetBlock("000000000000000000e35af10dba3f792fee6792eb01d11a1c87fec81d4e54ce")
	if berr.Code != 0 {
		t.Error(err)
		t.Fail()
	}

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Logf("%+v", block)
	t.Fail()
}
