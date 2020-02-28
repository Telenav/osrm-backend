package oasis

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
)

var fakeRoute1 route.Route = route.Route{
	Legs: []*route.Leg{
		&route.Leg{
			Distance: 37386.2,
			Steps: []*route.Step{
				&route.Step{
					Distance: 449.5,
					Geometry: "y~ecF|_qgVm@@kAJw@RqAf@a@PmBzAqA|Ae@v@a@|@aAzC",
				},
				&route.Step{
					Distance: 243,
					Geometry: "irfcFdqqgVEp@Ib@MZe@x@k@j@c@R_@Fm@DiBA",
				},
				&route.Step{
					Distance: 207.7,
					Geometry: "u|fcF|xqgV?zA?xA?h@?j@@tA?hA?h@",
				},
				&route.Step{
					Distance: 344.6,
					Geometry: "s|fcFrgrgV@h@tKHdA@t@@v@@",
				},
				&route.Step{
					Distance: 28988.5,
					Geometry: "gjfcFlirgVb@PZHZHj@PRJTJVRRTNXNd@DZ@\\A\\C`@Ih@Ih@Mn@Ot@OfA_@dDShAgA|FQ~@e@bCa@xBq@pDqB|K{@zEw@|Eq@vEO~@ANKr@EZGh@G`@[|Bm@hE_@fCIn@E^EZYpBw@bGM|@_AdHgBzN_@rCa@`CETCPIn@yAxKq@bFy@hGIh@e@fDKx@UbBCRId@Id@G`@}@dHy@tGKp@Kv@UdBIl@k@bEgAnIIt@Ih@Ip@If@YjBa@tCGf@Il@q@zEKv@g@rDuAfKg@xDsA~JyAxKOlAIp@QrAObAOjAKt@Kz@_@tCm@bEYpBUpA_@xBg@xCWtAa@|BYxAIb@AJMx@UtAa@~BUdBQjA_AdHqBhOYvBq@dFc@hDwAdKMz@Ir@Ib@Kt@QvAKp@]xBS~@g@nBIVIZYdAMb@q@bCK^k@jBcA`DcAbDaA`D{@jCQl@KXy@rC_A`Dc@|AOh@g@`Bw@fCK^cAxDSr@cAnDQj@Sn@w@dCu@jBGNiAzB[j@[j@aAdBINGJgAdBc@p@s@bAiAdBeAbBuB`DMPkDpFsCjEe@t@aCpD_@l@Yb@aBdCoAnBOR_AnASVOTGHMPu@jAiAdBm@`A_FzH_CnDqAlBaAvAy@bA_AjAGF_CtCcFlFeBjBkCrC[ZUVu@v@_@^mArAg@h@oCvC_@`@{C~CoArAqArASTqBtBgAjA}@`A{A~Ag@h@YZq@r@gCjCwB|BkDtDcBfBuHdIyA|AkAnAOPsAtAqAtA]^WXWX}@`Ae@f@QPcGnGoApAoCxCeHtHeJrJg@h@qAvA{FdGONMLoDvDs@v@oCzCg@l@w@~@cArAeAzA[f@]l@QXcB|Ci@hAs@zAm@lAs@xAs@xAaBjDeChFiD`HaDjGUb@qBhEsAtC_@t@iBpD[n@wArCwAtCwCjGS`@oCpFoCrFgAxBc@|@cAlBu@vA}@pBs@xAuEjJyAxCgG`MmBzDaAnBc@z@q@rAo@rAsAjCy@`BKXITIPkA|BkAbCeBjD_AjBe@~@y@dBk@fAmBzDqEdJcKpSWh@Sb@kA~BINGLs@vAq@tAwDvHeBlD{@dBmBzDe@`AKRo@tA_@v@Yh@cAnBcApBi@fA_CvEyA|Cm@pAuBlEm@tAu@lBo@|AkAlDABOh@CJCJCHGPADITSr@CJGTENe@dBi@zBw@|D[hBG\\[bB{@dFO~@G`@c@zDUhBQdBO`B]bEC^OxBcAzOiBbZIjAc@bH}@lN{@~MMtBw@bMUxCKpAY~Ck@|EGh@a@xCOfAIf@Id@Ib@Mx@Kn@Q`AG\\G\\Mh@I^Qv@Ml@i@tBW`AYfAQn@i@fBm@pBk@nB_@hAk@fBOd@eArCwBlGo@hBOb@Qf@}@jC_AlCs@tBs@~Bs@zBe@bBUx@U`AW`AU`AUdAU`A[vA[xA]hBId@Id@I`@UdAIl@]jBqC~Ne@bCKh@Kj@Q~@Mj@Kj@a@nBYtAMr@Mh@Oj@e@~Aw@`Cg@pAo@vAgAvBgAlBg@p@s@~@cBrBeE~EmAvAcAjAsA|AwChDiBvB{AhBuFvG{CpDaAhAcF`GeJrKyJhLu@|@ePnR]`@WZoAvAsEpFyFxGg@l@}@dAoAvA_BlByAdBqAzAcE|EoAzAcAjAqEjFyEtF{@bAe@h@",
				},
				&route.Step{
					Distance: 467.2,
					Geometry: "mx_dFbfgiVcAv@s@j@WV[\\eA`Am@h@]V[R[JYF_@D]@yAG[?W@UDUFSHSLQLc@`@",
				},
				&route.Step{
					Distance: 6035.5,
					Geometry: "so`dF`sgiVGf@NVR\\RZv@jAZf@LTV^RXl@~@b@r@`@n@TXRV~@dA^^d@b@VRFFTPf@\\BB^VXN`@TRJ\\R`Ah@`@TRJ`@NVLTLd@b@\\`@HJVXNJVLLPNPX^X^RVTVRVTVLPNPNPLPX\\RRLNNRZ^XZd@n@LPZ^DFBDRVr@z@RPZRJDPHXHj@JnATF@h@Jd@PPHd@ZPRPVR^LXNl@@HHr@Bb@Bb@@^B\\JbCBb@D`ADl@Dp@Dn@@XDf@DTHTLZV^JNPRVR\\VFDf@ZLHXNVXb@|@Rx@@DLp@Ln@ZbBXxADPRfAb@zBNv@Jj@FZR~@RpA`@bCLx@Nx@FVP~@^rBZ|AJh@DRf@fCJf@BNJd@RdA?P?DDRDVF`@Jf@FZJf@FZFZLn@FVFXX~AHb@Hb@Pz@Pz@Jd@PbARbALz@BNFd@B`@@f@@z@?~B?R?\\?x@?p@?lC?Z@l@Dz@Dx@?l@A|@Et@Kv@Kh@Md@Ob@MXEHMTQZ]n@KRKROZEHABQf@I\\GVETCNG^En@A`@Ad@@b@Dx@D|@@\\@b@@^@^?j@@dA?\\@f@Ax@A^Cj@Iz@APE`@E`@Gj@Gh@I`AEj@Cj@AnABdAFfAFz@Fl@Fp@BRH|@L~AFl@BVFl@BXDb@HbA^`FHjAz@hKDh@Fr@Ht@Hv@P~AFv@Hx@HfADj@B`@Dr@F|@Bl@Bj@@f@@`A?bAAx@Ez@IhAIx@UfBKx@Gr@UxCUxCCXEl@Eb@C^E`@C`@Gl@KdAGn@",
				},
				&route.Step{
					Distance: 650.3,
					Geometry: "aq}cFfesiVK^OXQLQDy@@O?OBYLeAr@[N_@Fa@DUAu@IcAOSCe@Em@Ei@?mA@u@Bo@Dg@Jk@Vw@PcAZcA`@KF",
				},
				&route.Step{
					Distance: 0,
					Geometry: "os~cFvmsiV",
				},
			},
		},
	},
}

func TestFindChargeLocation4Route1(t *testing.T) {
	cases := []struct {
		currEnergy  float64
		preferLevel float64
		maxRange    float64
		chargeCount int
	}{
		{
			15000.0,
			10000.0,
			50000.0,
			1,
		},
		{
			1000.0,
			1500.0,
			15000.0,
			3,
		},
	}

	for _, c := range cases {
		var expect []*stationfinder.StationCoordinate
		expect, _ = findChargeLocation4Route(&fakeRoute1, expect, c.currEnergy, c.preferLevel, c.maxRange)
		if len(expect) != c.chargeCount {
			t.Errorf("expect to charge %v times, but got %v times", c.chargeCount, len(expect))
		}
	}

	// var expect []*stationfinder.StationCoordinate
	// currEnergy := 15000.0
	// preferLevel := 10000.0
	// maxRange := 50000.0
	// expect, _ = findChargeLocation4Route(&fakeRoute1, expect, currEnergy, preferLevel, maxRange)
	// fmt.Printf("%d \n", len(expect))

	// s := "gjfcFlirgVb@PZHZHj@PRJTJVRRTNXNd@DZ@\\A\\C`@Ih@Ih@Mn@Ot@OfA_@dDShAgA|FQ~@e@bCa@xBq@pDqB|K{@zEw@|Eq@vEO~@ANKr@EZGh@G`@[|Bm@hE_@fCIn@E^EZYpBw@bGM|@_AdHgBzN_@rCa@`CETCPIn@yAxKq@bFy@hGIh@e@fDKx@UbBCRId@Id@G`@}@dHy@tGKp@Kv@UdBIl@k@bEgAnIIt@Ih@Ip@If@YjBa@tCGf@Il@q@zEKv@g@rDuAfKg@xDsA~JyAxKOlAIp@QrAObAOjAKt@Kz@_@tCm@bEYpBUpA_@xBg@xCWtAa@|BYxAIb@AJMx@UtAa@~BUdBQjA_AdHqBhOYvBq@dFc@hDwAdKMz@Ir@Ib@Kt@QvAKp@]xBS~@g@nBIVIZYdAMb@q@bCK^k@jBcA`DcAbDaA`D{@jCQl@KXy@rC_A`Dc@|AOh@g@`Bw@fCK^cAxDSr@cAnDQj@Sn@w@dCu@jBGNiAzB[j@[j@aAdBINGJgAdBc@p@s@bAiAdBeAbBuB`DMPkDpFsCjEe@t@aCpD_@l@Yb@aBdCoAnBOR_AnASVOTGHMPu@jAiAdBm@`A_FzH_CnDqAlBaAvAy@bA_AjAGF_CtCcFlFeBjBkCrC[ZUVu@v@_@^mArAg@h@oCvC_@`@{C~CoArAqArASTqBtBgAjA}@`A{A~Ag@h@YZq@r@gCjCwB|BkDtDcBfBuHdIyA|AkAnAOPsAtAqAtA]^WXWX}@`Ae@f@QPcGnGoApAoCxCeHtHeJrJg@h@qAvA{FdGONMLoDvDs@v@oCzCg@l@w@~@cArAeAzA[f@]l@QXcB|Ci@hAs@zAm@lAs@xAs@xAaBjDeChFiD`HaDjGUb@qBhEsAtC_@t@iBpD[n@wArCwAtCwCjGS`@oCpFoCrFgAxBc@|@cAlBu@vA}@pBs@xAuEjJyAxCgG`MmBzDaAnBc@z@q@rAo@rAsAjCy@`BKXITIPkA|BkAbCeBjD_AjBe@~@y@dBk@fAmBzDqEdJcKpSWh@Sb@kA~BINGLs@vAq@tAwDvHeBlD{@dBmBzDe@`AKRo@tA_@v@Yh@cAnBcApBi@fA_CvEyA|Cm@pAuBlEm@tAu@lBo@|AkAlDABOh@CJCJCHGPADITSr@CJGTENe@dBi@zBw@|D[hBG\\[bB{@dFO~@G`@c@zDUhBQdBO`B]bEC^OxBcAzOiBbZIjAc@bH}@lN{@~MMtBw@bMUxCKpAY~Ck@|EGh@a@xCOfAIf@Id@Ib@Mx@Kn@Q`AG\\G\\Mh@I^Qv@Ml@i@tBW`AYfAQn@i@fBm@pBk@nB_@hAk@fBOd@eArCwBlGo@hBOb@Qf@}@jC_AlCs@tBs@~Bs@zBe@bBUx@U`AW`AU`AUdAU`A[vA[xA]hBId@Id@I`@UdAIl@]jBqC~Ne@bCKh@Kj@Q~@Mj@Kj@a@nBYtAMr@Mh@Oj@e@~Aw@`Cg@pAo@vAgAvBgAlBg@p@s@~@cBrBeE~EmAvAcAjAsA|AwChDiBvB{AhBuFvG{CpDaAhAcF`GeJrKyJhLu@|@ePnR]`@WZoAvAsEpFyFxGg@l@}@dAoAvA_BlByAdBqAzAcE|EoAzAcAjAqEjFyEtF{@bAe@h@"
	// coords, _, _ := polyline.DecodeCoords([]byte((&fakeRoute1).Legs[0].Steps[4].Geometry))

	// fmt.Printf("++++ %v\n", coords[0][0])

	//fmt.Printf("### %v\n", []byte((&fakeRoute1).Legs[0].Steps[4].Geometry))
}
