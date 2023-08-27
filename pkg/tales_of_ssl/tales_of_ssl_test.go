package tales_of_ssl

import (
	"testing"
)

func BenchmarkRun(b *testing.B) {
	tof := new(TalesOfSsl)
	input := `{"private_key": "MIICXgIBAAKBgQC3flnJg/nQIPO8LYiv6D/PDXFpU2aEWsmDtbjnClhcx9QTSMY2VvsX45icIHjFPtWYkhdUZkhLNFfKbbPFhZfo7doYBVkEbPnMjYP+bagkRG6u/kV8g3Oc8uqgor83HOiJarhUm7jRKyUpmh/3RoW/4GUl2gtA37fynmbK41Q5JQIDAQABAoGBAJbmieB3cJIBB3DR4X8IeLFGVOZReBWQs/hlEdKvZ0ul6nqRdrBph8V1nqOG8MsFiZMXUahPFVUDjs8iuJzP6MR5Shzmh+xF7urecVgICrc3wq7URMCzZcCAW8Jh+hdHaDm5GwFxLh/VXOipRj3g4MMDkHrgm+uokSiWCk3E6HRhAkEA4hYGQurX2rPTummZAo/Dq6GL+YcQTqWOuBwaG7zSMBDdLwDV/EqLCn64D8aElTMcRI0wJ9bgc9SbqfQa18SjHwJBAM/FosyNn8HBzqh8cOzs1Xhqw4ggs+jtk4og7rtQGLEcN8J6JuakpVY01eeavAj+Qzj0DGoWhu5kYRECdU1iPzsCQQDbYhfNU2kF1b28g79wFuT+ZJIZHhCa4FfvG4XSKJWlApg+hgce+46ULoaY+D+rL59cfcyevAmtHD+42SA7A5ptAkBvMSSz9eMWbWLsllRU51ZC8FgeROojcOlxlbhNdEsSlbhdksO40UpOfU4weYXCBljxZOIz8ERb1kqAMOGX/Uk7AkEAxttZSrHeqPBzV6l7JPPp32C//NQyll8SRW/oN1hmdIHZDGNkaGurDkfJmUBplXlJJKklBqJ5XL739V8YG4uBfQ==", "required_data": {"domain": "summer-fire-4799.gov", "serial_number": "0x7a771e7a1e", "country": "Cocos Islands"}}`
	for n := 0; n < b.N; n++ {
		tof.Solve(input)
	}
}
