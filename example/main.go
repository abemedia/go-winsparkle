//go:build windows
// +build windows

package main

import (
	"log"
	"time"

	sparkle "github.com/abemedia/go-winsparkle"
)

func main() {
	log.Println("starting app")

	sparkle.SetAppcastURL("https://winsparkle.org/example/appcast.xml")
	sparkle.SetAppDetails("winsparkle.org", "WinSparkle Go Example", "1.0")

	if err := sparkle.SetDSAPubPEM(pem); err != nil {
		log.Fatal(err)
	}

	c := make(chan struct{})

	sparkle.SetErrorCallback(func() {
		log.Println("received an error")
	})

	sparkle.SetDidFindUpdateCallback(func() {
		log.Println("did find update")
	})

	sparkle.SetDidNotFindUpdateCallback(func() {
		log.Println("did not find update")
	})

	sparkle.SetShutdownRequestCallback(func() {
		log.Println("installing update")
		close(c)
	})

	sparkle.SetUpdateCancelledCallback(func() {
		log.Println("cancelled update")
		close(c)
	})

	sparkle.SetUpdateSkippedCallback(func() {
		log.Println("skipped update")
	})

	sparkle.SetUpdateDismissedCallback(func() {
		log.Println("dismissed update")
	})

	sparkle.SetUpdatePostponedCallback(func() {
		log.Println("postponed update")
	})

	sparkle.SetUserRunInstallerCallback(func(s string) (bool, error) {
		log.Println("installer callback: " + s)
		return true, nil
	})

	config := &store{data: make(map[string]string)}
	sparkle.SetConfigMethods(config)

	sparkle.Init()
	defer sparkle.Cleanup()

	time.Sleep(time.Second)

	sparkle.CheckUpdateWithUI()

	// waits until update is installed or cancelled (10min timeout)
	select {
	case <-c:
	case <-time.After(10 * time.Minute):
	}
	log.Println("shutting down")
}

type store struct {
	data map[string]string
}

func (c *store) Read(name string) (string, bool) {
	log.Println("reading config: " + name)
	v, ok := c.data[name]
	if !ok {
		return "", false
	}
	return v, true
}

func (c *store) Write(name, value string) bool {
	log.Printf("writing config: %s=%q\n", name, value)
	c.data[name] = value
	return true
}

func (c *store) Delete(name string) bool {
	log.Println("deleting config: " + name)
	delete(c.data, name)
	return true
}

const pem = `-----BEGIN PUBLIC KEY-----
MIIGRjCCBDkGByqGSM44BAEwggQsAoICAQCS7frykSttT+BkjHTeY/BZMWN6Fxg3
nYg1gDHb+QxQElEikI/70f7oJ9f/UIijyFMZgUcdP2D28X2Wg/gFkvJkWric2HJL
/QiFEo0SWfcDu8ViDYnHqvGuZkQ1qZxs+cx6PCMU1Wej68JVj3jYhnL0j+hlUmI5
HodCMSjaIXA3MnX6VJ3qpG0dvz2E9dNP7SX52X0f1qJZU7LDRt+b+mhJifQM18Wq
pzH4GiJ1YwPMJ7jsz7kFx34qwa/WMvcpcpPvRZQnm5tNS6fW/VzWU9zoB9ANEWcE
7m0dsy2RPvPOu6DbHEK8HQngApIxshrX3QHBh/P1n48IEfctNMH3p9ZReGg7P5MU
McRemULs+uYxr846dNLp/SRplkOcoftCAdZpTdO+7kcZNpMxX8c0JuuOkI6+pNMn
iCw2kRmP7ktnPs1g19rbRB0G5wqXU3hoIVV8JJzgWMkj1sQdNWqvhdFuq8zzx5xQ
JlOHxDJIDINtZQc4p08oPcK5EFh2G0Zd3ykpn4TbKqBzbkFnqsSbwFeGKg+JJCT5
T9lommwz3Rtq0JxgJXizGA9OBLY49FYbPQBYUfOgakT6gJhgUoVqyrCj7OFLzH/p
O9ro+CcQgcsaaApHBLPgCX/kcBFDNuVuU/4Rwru2yi468hbMf7ySQJ8sLOPbYutZ
dKizdhJZ+qteGwIhAPVTS1dFaqLnENUXXC3Rf/9B9mZ1qd69pRaH16x94beXAoIC
ADtQsQV1F8jwYIQpkDIkypQFoMdXrkctnB+54SIw1PCE08t16MOI16SQhQ45bIY0
8wP6VXo/QtIswTG2OKiEAOprvhTheJR8/glXWqxAsbDhgcZ7OScsHbCywysLoCNh
FCbLzVt1WYnKxBdeOjdjPFq5gP/5CVgd7wdia3eqAz/tmgnttq0shVEP9hHh13jo
BKZrQcXap67Iluf1fcsmxHa2f4k2h2zz1AO7Zz/Fk6tNr7jGDLRmig4uii3x2tYz
wcmBT38Y0/9UWdyRTzZZ8ePiTxSJJhKRU8X8ARhvBjZ3teezKqD8vs7LThBOLODu
967p9Kf0q2pHNZyUuJFFlAeDp/H4/VJlIQwar8h7X1OWxxrA+b3NMuGTi+BMA6rS
7sEN6ZbO6/HJLeW3xtxsQ96miyQlyd68be9zbziGTKveT0ctdrztENtRi6wFf7/F
p7GUOMXBwkDluJw0fw8cy5t0PcIzzgjf7bSrbPfgDUDeaIHNKPjilCeWEaD3ahCO
pl3T0/E+RY+TEFkL0mhPvVHwxqaNSfrs977pEbxZJ+K9v68Rp6oUi3WmUQii6oqG
c/77JCoH8qVe8kgzj5BB6Z01oq+jZjrlqIofEsnxewi+Tm9RZ4affp9k7TtitUXp
YvxNXp/gOrReQV+1QLxLTiY4DjQ4UMi0nd4IeAhScWyDA4ICBQACggIAfngtVUTJ
r2sTChxOWLqHnIEi/Hbdm0wiTz5LXnvu/kRjXvG2YaRa2CYS1082jK7ywq2vf8g/
R0nMxWxpmB0+vW0NZY6hP3mnCaGn7v3W+i6tb0wUtdHPNRRGSj0zQJFJ7yA/nPLg
Aeku5wMlqUwIJvIRoPQKJ5eE/5dJrVz+wFcDsv3L3IrGR1S0rRbgoEAYpni5u7Bj
JaPJJ8pGVn0XqHPSyeSQh7So5cxE0YFYBUj/Bb+3etEcwv99pmsy7pWUUehofZt2
ZEBAdca6/AHIFbs5TqVMsVYXQha8/U3rENhzv+X0canxPeziu2n1Op5zIeUmQHqI
uu8w4YeEYo+cWO/jkssSAFvDAbtU2o0naMM+8deIscFkbqoRDu9iGL/X3O1g7Kb7
Df1fCnxSIon5LredUYqZekFAGPfdFPOfx243GAYEysSdmHwCd8muGQ1RvED63Hc6
rMEDBfwUjZMC4LhlMKICfcL5C6x0QaEtUhxKCKoxnuUCDpTv0FLVuD3ROXOkpWb9
8asN6oduRYvwk8XgJywdWlEgqZ3q0UYdq9lkPX1k2qcySWtuOtCCRgQvi/5n7jEN
Pw64lVzuSHSeHyFt7DmuVPBGI8I8QNQ0Zdn7/ciChYWjgfkewVxIzraeo23EVqp4
gHD0Yknt+5zQ1vtA0zKwfjJRYNzH9oVAYuc=
-----END PUBLIC KEY-----
`
