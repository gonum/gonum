// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prng

import (
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

// Random values in tests are produced by 40 iterations of the C code.

var _ rand.Source = (*SplitMix64)(nil)

func TestSplitMix64(t *testing.T) {
	t.Parallel()
	want := []uint64{
		10451216379200822465, 13757245211066428519, 17911839290282890590, 8196980753821780235, 8195237237126968761,
		14072917602864530048, 16184226688143867045, 9648886400068060533, 5266705631892356520, 14646652180046636950,
		7455107161863376737, 11168034603498703870, 8392123148533390784, 9778231605760336522, 8042142155559163816,
		3081251696030599739, 11904322950028659555, 15040563541741120241, 12575237177726700014, 16312908901713405192,
		1216750802008901446, 1501835286251455644, 9147370558249537485, 2270958130545493676, 5292580334274787743,
		883620860755687159, 9509663594007654709, 13166747327335888811, 807013244984872231, 18405200023706498954,
		11028426030083068036, 10820770463232788922, 7326479631639850093, 8097875853865443356, 4672064935750269975,
		9772298966463872780, 10028955912863736053, 13802505617680978881, 15090054588401425688, 12333003614474408764,
	}

	sm := NewSplitMix64(1)
	for i := range want {
		got := sm.Uint64()
		if got != want[i] {
			t.Errorf("unexpected random value at iteration %d: got:%d want:%d", i, got, want[i])
		}
	}
}

func TestSplitMix64RoundTrip(t *testing.T) {
	t.Parallel()
	var src SplitMix64
	src.Seed(uint64(time.Now().Unix()))

	buf, err := src.MarshalBinary()
	if err != nil {
		t.Errorf("unexpected error marshaling state: %v", err)
	}

	var dst SplitMix64
	// Get dst into a non-zero state.
	dst.Seed(1)
	for i := 0; i < 10; i++ {
		dst.Uint64()
	}

	err = dst.UnmarshalBinary(buf)
	if err != nil {
		t.Errorf("unexpected error unmarshaling state: %v", err)
	}

	if dst != src {
		t.Errorf("mismatch between generator states: got:%+v want:%+v", dst, src)
	}
}

var _ rand.Source = (*Xoshiro256plus)(nil)

func TestXoshiro256plus(t *testing.T) {
	t.Parallel()
	want := []uint64{
		201453059313051084, 16342930563397888806, 2922809869868169223, 13315230553875954649, 6410977891529050008,
		2721661332018190285, 3769995280709464022, 17208995829377771030, 16938999919058283733, 8307416726322109393,
		13997290115667311691, 5498422487743993519, 13193129985428835789, 17178224140053183722, 3371202013665523682,
		6673444001875245482, 11649545741795472859, 4657392542380076879, 8631341306563158492, 16151880809814987639,
		15271080878658922261, 6998002807989632655, 11431762507643441726, 136605885039865329, 16072241235209520170,
		17064623797431990278, 6319393334343723778, 3599071131527455911, 14678971584471326753, 11566847267978507055,
		37242444495476935, 9767625399998905638, 14799351402198708144, 15147234459691564338, 10081976988475685812,
		12402022881820243150, 17939631254687971868, 15680836376982110901, 179319489669050051, 16194215847106809765,
	}

	xsr := NewXoshiro256plus(1)
	for i := range want {
		got := xsr.Uint64()
		if got != want[i] {
			t.Errorf("unexpected random value at iteration %d: got:%d want:%d", i, got, want[i])
		}
	}
}

func TestXoshiro256plusRoundTrip(t *testing.T) {
	t.Parallel()
	var src Xoshiro256plus
	src.Seed(uint64(time.Now().Unix()))

	src.Uint64() // Step PRNG once to makes sure states are mixed.

	buf, err := src.MarshalBinary()
	if err != nil {
		t.Errorf("unexpected error marshaling state: %v", err)
	}

	var dst Xoshiro256plus
	// Get dst into a non-zero state.
	dst.Seed(1)
	for i := 0; i < 10; i++ {
		dst.Uint64()
	}

	err = dst.UnmarshalBinary(buf)
	if err != nil {
		t.Errorf("unexpected error unmarshaling state: %v", err)
	}

	if dst != src {
		t.Errorf("mismatch between generator states: got:%+v want:%+v", dst, src)
	}
}

var _ rand.Source = (*Xoshiro256plusplus)(nil)

func TestXoshiro256plusplus(t *testing.T) {
	t.Parallel()
	want := []uint64{
		14971601782005023387, 13781649495232077965, 1847458086238483744, 13765271635752736470, 3406718355780431780,
		10892412867582108485, 18204613561675945223, 9655336933892813345, 1781989159761824720, 2477283028068920342,
		16978024111547606601, 6336475467619303347, 1336129645694042326, 7278725533440954441, 1650926874576718010,
		2884092293074692283, 10277292511068429730, 8723528388573605619, 17670016435951889822, 11847526622624223050,
		4869519043768407819, 14645621260580619786, 2927941368235978475, 7627105703721172900, 4384663367605854827,
		11119034730948704880, 3397900810577180010, 18115970067406137490, 11274606161466886392, 13467911786374401590,
		10949103424463861935, 11981483663808188895, 9358210361682609782, 11442939244776437245, 17602980262171424054,
		5959474180322755185, 1996769245947054333, 13544632058761996522, 16649296193330087156, 12760326241867116135,
	}

	xsr := NewXoshiro256plusplus(1)
	for i := range want {
		got := xsr.Uint64()
		if got != want[i] {
			t.Errorf("unexpected random value at iteration %d: got:%d want:%d", i, got, want[i])
		}
	}
}

func TestXoshiro256plusplusRoundTrip(t *testing.T) {
	t.Parallel()
	var src Xoshiro256plusplus
	src.Seed(uint64(time.Now().Unix()))

	src.Uint64() // Step PRNG once to makes sure states are mixed.

	buf, err := src.MarshalBinary()
	if err != nil {
		t.Errorf("unexpected error marshaling state: %v", err)
	}

	var dst Xoshiro256plusplus
	// Get dst into a non-zero state.
	dst.Seed(1)
	for i := 0; i < 10; i++ {
		dst.Uint64()
	}

	err = dst.UnmarshalBinary(buf)
	if err != nil {
		t.Errorf("unexpected error unmarshaling state: %v", err)
	}

	if dst != src {
		t.Errorf("mismatch between generator states: got:%+v want:%+v", dst, src)
	}
}

var _ rand.Source = (*Xoshiro256starstar)(nil)

func TestXoshiro256starstar(t *testing.T) {
	t.Parallel()
	want := []uint64{
		12966619160104079557, 9600361134598540522, 10590380919521690900, 7218738570589545383, 12860671823995680371,
		2648436617965840162, 1310552918490157286, 7031611932980406429, 15996139959407692321, 10177250653276320208,
		17202925169076741841, 17657558547222227110, 17206619296382044401, 12342657103067243573, 11066818095355039191,
		16427605434558419749, 1484150211974036615, 9063990983673329711, 845232928428614080, 1176429380546917807,
		8545088851120551825, 9158324580728115577, 11267126437916202177, 6452051665337041730, 7460617819096774474,
		3909615622106851260, 7148019177890935463, 15761474764570999248, 13856144421012645925, 18119237044791779759,
		202581184499657049, 16256128138147959276, 7894450248801719761, 7285265299121834259, 11974578372788407364,
		4350246478179107086, 4560570958642824732, 15448532239578831742, 7084622563335324071, 8654072644765974953,
	}

	xsr := NewXoshiro256starstar(1)
	for i := range want {
		got := xsr.Uint64()
		if got != want[i] {
			t.Errorf("unexpected random value at iteration %d: got:%d want:%d", i, got, want[i])
		}
	}
}

func TestXoshiro256starstarRoundTrip(t *testing.T) {
	t.Parallel()
	var src Xoshiro256starstar
	src.Seed(uint64(time.Now().Unix()))

	src.Uint64() // Step PRNG once to makes sure states are mixed.

	buf, err := src.MarshalBinary()
	if err != nil {
		t.Errorf("unexpected error marshaling state: %v", err)
	}

	var dst Xoshiro256starstar
	// Get dst into a non-zero state.
	dst.Seed(1)
	for i := 0; i < 10; i++ {
		dst.Uint64()
	}

	err = dst.UnmarshalBinary(buf)
	if err != nil {
		t.Errorf("unexpected error unmarshaling state: %v", err)
	}

	if dst != src {
		t.Errorf("mismatch between generator states: got:%+v want:%+v", dst, src)
	}
}
