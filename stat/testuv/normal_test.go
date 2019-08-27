package testuv

import (
	"math"
	"testing"
)

// It'd be better if we could use distuv.Normal but that would create a cyclic import
// Instead I created a random sample with the following python script
// ```
// import numpy as np
// np.random.seed(28041990)
// s = np.random.normal(0, 1, 100)
// for v in s:
//     print(str(v) + ",")
// ```
var norm = []float64{
	-0.42559355556264555,
	1.710539107218753,
	-0.33074948677055843,
	-0.5161417685463089,
	-0.6193456407278632,
	1.1351354008003518,
	1.3997307933412553,
	-0.029954251382664998,
	0.8434220441471606,
	2.1128080605675374,
	1.2484609152871478,
	1.8344987008482856,
	0.09931267796096076,
	1.428919892713464,
	1.0189271326046607,
	-0.2414299164875163,
	0.9066103078332822,
	0.3401187294324855,
	0.37733889116557606,
	-0.7323545583397159,
	-0.31958107231420857,
	-0.2512518199227104,
	-1.1383397596501,
	1.0035254463416154,
	1.1418421937745262,
	-1.2358788434742574,
	1.0100704641117955,
	-0.08036883653794552,
	-0.58513457961876,
	0.05084968136313145,
	0.5664086258955331,
	-0.43308548122088375,
	-3.1178336091684975,
	0.7603796596063905,
	0.2853512117407265,
	0.2176337503977105,
	2.152757471948443,
	-1.3087387278036695,
	0.9978177912374294,
	1.786766843343617,
	0.2306781426352178,
	1.4431459612084379,
	-0.6566580121689873,
	0.4684545103705484,
	-0.8951779888297282,
	0.822016499896265,
	-0.16799507924776486,
	-0.0020099175487792724,
	0.11204979393586693,
	1.5055567069562388,
	0.10391457337177848,
	-1.3699213591924695,
	-1.0072435369731148,
	0.9237191947896234,
	-0.019525907747287095,
	-1.0887672634881844,
	0.2751885920079707,
	-0.3434944093974466,
	0.4246733152123921,
	-0.6853964678846619,
	0.46427348833740134,
	-1.1769384002155117,
	0.1711116430457653,
	-1.7144590897345933,
	3.522255994195385,
	-0.37646202717640115,
	-2.433505931729811,
	1.256781201303011,
	-0.15655339743939028,
	0.4552009397642383,
	-0.18924946568753523,
	-1.3725653125945136,
	1.6436894604846044,
	0.597721827583972,
	-0.04796417108532695,
	0.5909740089894597,
	-0.32067675553881103,
	-0.9249846629433014,
	-0.5678094294504085,
	-0.6602157563037903,
	1.9230857720458412,
	-0.22916228851590298,
	-0.4570292925300327,
	-0.5904063162660567,
	0.9813766542000734,
	-0.15441746156260575,
	1.3686357504466162,
	-1.775095655630572,
	-1.0537085236430228,
	0.21705168855671522,
	0.10583883865411757,
	0.680143572120755,
	0.043548775401668044,
	-0.8204090347120757,
	-0.9641761107363781,
	-1.2735654936668594,
	0.6970301935177855,
	1.0141887564307395,
	0.49425150550631447,
	0.8237595013007176,
}

func TestNormalSkewTest(t *testing.T) {
	// Even though the skew test and kurtosis test are a copy from the scipy's the results
	// are not exactly the same because the skew/kurtosis corrections are not exactly the same
	// but the statistic value is close enough.
	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			expected: 1.018464355396213,
		},
		{
			values:   []float64{2, 8, 0, 4, 1, 9, 9, 0},
			expected: 0.5562585562766172,
		},
		{
			values:   []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			expected: 4.319816401673864,
		},
		{
			values:   []float64{100, 100, 100, 100, 100, 100, 100, 101},
			expected: 4.319820025201098,
		},
	} {
		z := NormalSkewTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalSkewTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}

func TestNormalKurtosisTest(t *testing.T) {
	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			expected: -1.6867033202073243,
		},
		{
			values:   norm,
			expected: 1.638462148403484,
		},
	} {
		z := NormalKurtosisTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalKurtosisTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}

func TestNormalTest(t *testing.T) {
	// normal test with scipy yield similar results:
	// ```
	// from scipy.stats import normaltest
	// import numpy as np
	// np.random.seed(28041990)
	// s = np.random.normal(0, 1, 100)
	// print(normaltest(list(range(9))).statistic)
	// print(normaltest([2, 8, 0, 4, 1, 9, 9, 0]).statistic)
	// print(normaltest([1, 2, 3, 4, 5, 6, 7, 8000]).statistic)
	// print(normaltest([100, 100, 100, 100, 100, 100, 100, 101]).statistic)
	// print(normaltest(list(range(20))).statistic)
	// print(normaltest(s).statistic)
	// ```

	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			expected: 1.7509245653153176,
		},
		{
			values:   []float64{2, 8, 0, 4, 1, 9, 9, 0},
			expected: 11.454757293481551,
		},
		{
			values:   []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			expected: 40.53534243515444,
		},
		{
			values:   []float64{100, 100, 100, 100, 100, 100, 100, 101},
			expected: 40.53539760601764,
		},
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			expected: 3.9272951079276743,
		},
		{
			values:   norm,
			expected: 2.695532408378151,
		},
	} {
		z := NormalTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}
