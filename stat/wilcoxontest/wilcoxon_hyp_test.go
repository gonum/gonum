package wilcoxon

import (
	"reflect"
	"testing"
)

func TestRank(t *testing.T) {
	type args struct {
		in  []float64
		out []float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{name: "name1", args: args{
			in:  []float64{0, 5, 5, 7, 9, 10, 12, 15, 17, 20},
			out: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}, want: []float64{0, 1.5, 1.5, 3, 4, 5, 6, 7, 8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Rank(tt.args.in, tt.args.out); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWilcoxonSignedRankTest(t *testing.T) {
	type args struct {
		x           []float64
		y           []float64
		exactPValue bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test1",
			args: args{
				x:           []float64{1.83, 0.50, 1.62, 2.48, 1.68, 1.88, 1.55, 3.06, 1.30},
				y:           []float64{0.878, 0.647, 0.598, 2.050, 1.060, 1.290, 1.060, 3.140, 1.290},
				exactPValue: true,
			},
			want: 0.0390625,
		},
		{
			name: "test2",
			args: args{
				x:           []float64{2.0, 3.6, 2.6, 2.6, 7.3, 3.4, 14.9, 6.6, 2.3, 2.0, 6.8, 8.5},
				y:           []float64{3.5, 5.7, 2.9, 2.4, 9.9, 3.3, 16.7, 6.0, 3.8, 4.0, 9.1, 20.9},
				exactPValue: false,
			},
			want: 0.010757133119789652,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WilcoxonSignedRankTest(tt.args.x, tt.args.y, tt.args.exactPValue); got != tt.want {
				t.Errorf("WilcoxonSignedRankTest() = %v, want %v", got, tt.want)
			}
		})
	}
}
