// Derived from SciPy's special/cephes/lanczos.c
// https://github.com/scipy/scipy/blob/master/scipy/special/cephes/lanczos.c

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright ©2006 John Maddock
// Portions Copyright ©2003 Boost
// Portions Copyright ©2016 The gonum Authors. All rights reserved.

package cephes // import "gonum.org/v1/gonum/mathext/internal/cephes"

// Optimal values for G for each N are taken from
// http://web.mala.bc.ca/pughg/phdThesis/phdThesis.pdf,
// as are the theoretical error bounds.

// Constants calculated using the method described by Godfrey
// http://my.fit.edu/~gabdo/gamma.txt and elaborated by Toth at
// http://www.rskey.org/gamma.htm using NTL::RR at 1000 bit precision.

var lanczosNum = [...]float64{
	2.506628274631000270164908177133837338626,
	210.8242777515793458725097339207133627117,
	8071.672002365816210638002902272250613822,
	186056.2653952234950402949897160456992822,
	2876370.628935372441225409051620849613599,
	31426415.58540019438061423162831820536287,
	248874557.8620541565114603864132294232163,
	1439720407.311721673663223072794912393972,
	6039542586.35202800506429164430729792107,
	17921034426.03720969991975575445893111267,
	35711959237.35566804944018545154716670596,
	42919803642.64909876895789904700198885093,
	23531376880.41075968857200767445163675473,
}

var lanczosDenom = [...]float64{
	1,
	66,
	1925,
	32670,
	357423,
	2637558,
	13339535,
	45995730,
	105258076,
	150917976,
	120543840,
	39916800,
	0,
}

var lanczosSumExpgScaledNum = [...]float64{
	0.006061842346248906525783753964555936883222,
	0.5098416655656676188125178644804694509993,
	19.51992788247617482847860966235652136208,
	449.9445569063168119446858607650988409623,
	6955.999602515376140356310115515198987526,
	75999.29304014542649875303443598909137092,
	601859.6171681098786670226533699352302507,
	3481712.15498064590882071018964774556468,
	14605578.08768506808414169982791359218571,
	43338889.32467613834773723740590533316085,
	86363131.28813859145546927288977868422342,
	103794043.1163445451906271053616070238554,
	56906521.91347156388090791033559122686859,
}

var lanczosSumExpgScaledDenom = [...]float64{
	1,
	66,
	1925,
	32670,
	357423,
	2637558,
	13339535,
	45995730,
	105258076,
	150917976,
	120543840,
	39916800,
	0,
}

var lanczosSumNear1D = [...]float64{
	0.3394643171893132535170101292240837927725e-9,
	-0.2499505151487868335680273909354071938387e-8,
	0.8690926181038057039526127422002498960172e-8,
	-0.1933117898880828348692541394841204288047e-7,
	0.3075580174791348492737947340039992829546e-7,
	-0.2752907702903126466004207345038327818713e-7,
	-0.1515973019871092388943437623825208095123e-5,
	0.004785200610085071473880915854204301886437,
	-0.1993758927614728757314233026257810172008,
	1.483082862367253753040442933770164111678,
	-3.327150580651624233553677113928873034916,
	2.208709979316623790862569924861841433016,
}

var lanczosSumNear2D = [...]float64{
	0.1009141566987569892221439918230042368112e-8,
	-0.7430396708998719707642735577238449585822e-8,
	0.2583592566524439230844378948704262291927e-7,
	-0.5746670642147041587497159649318454348117e-7,
	0.9142922068165324132060550591210267992072e-7,
	-0.8183698410724358930823737982119474130069e-7,
	-0.4506604409707170077136555010018549819192e-5,
	0.01422519127192419234315002746252160965831,
	-0.5926941084905061794445733628891024027949,
	4.408830289125943377923077727900630927902,
	-9.8907772644920670589288081640128194231,
	6.565936202082889535528455955485877361223,
}

const lanczosG = 6.024680040776729583740234375

func lanczosSum(x float64) float64 {
	return ratevl(x,
		lanczosNum[:],
		len(lanczosNum)-1,
		lanczosDenom[:],
		len(lanczosDenom)-1)
}

func lanczosSumExpgScaled(x float64) float64 {
	return ratevl(x,
		lanczosSumExpgScaledNum[:],
		len(lanczosSumExpgScaledNum)-1,
		lanczosSumExpgScaledDenom[:],
		len(lanczosSumExpgScaledDenom)-1)
}

func lanczosSumNear1(dx float64) float64 {
	var result float64

	for i, val := range lanczosSumNear1D {
		k := float64(i + 1)
		result += (-val * dx) / (k*dx + k*k)
	}

	return result
}

func lanczosSumNear2(dx float64) float64 {
	var result float64
	x := dx + 2

	for i, val := range lanczosSumNear2D {
		k := float64(i + 1)
		result += (-val * dx) / (x + k*x + k*k - 1)
	}

	return result
}
