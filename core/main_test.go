package core

import (
	"testing"
)

func TestMain(m *testing.M) {
	Init("../VNA_new_redeem_voucher/core-config.yml")
	m.Run()
}
