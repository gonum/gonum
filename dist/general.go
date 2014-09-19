package dist

// Parameter represents a parameter of a probability distribution
type Parameter struct {
	Name  string
	Value float64
}

// A ParameterMarshaler is a type that can marshal itself into a slice of Paramaters
// and unmarshal itself from the slice of parameters.
// ParameterMarshaler exists to support algorithms that modify parameters of arbitrary distributions.
// Typically, users should modify distributions using the fields of the specifc
// distribution.
// Marshal and Unmarshal are to be used as a pair, users should not attempt to construct
// the Parameter slice themselves.
//
// Both MarshalParameters and  UnmarshalParameters will
// panic if the length of the slice is not equal to the number of parameters.
// UnmarshalParameters will panic if the names of the parameters do not match.
// UnmarshalParameters tests names in the same order as they were created in
// MarshalParameters.
type ParameterMarshaler interface {
	MarshalParameters([]Parameter)
	UnmarshalParameters([]Parameter)
}
