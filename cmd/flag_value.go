package cmd

import "strconv"

// =======================
type StringValue string

func (v *StringValue) String() string {
	return string(*v)
}

func (v *StringValue) Set(s string) error {
	*v = StringValue(s)
	return nil
}

func (v StringValue) Type() string {
	return "StringValue"
}

// =======================
type IntValue int64

func (v *IntValue) String() string {
	return strconv.Itoa(int(*v))
}

func (v *IntValue) Set(s string) error {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}

	*v = IntValue(i)
	return nil
}

func (v IntValue) Type() string {
	return "IntValue"
}

// =======================
type BoolValue bool

func (v *BoolValue) String() string {
	return strconv.FormatBool(bool(*v))
}

func (v *BoolValue) Set(s string) error {
	newv, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*v = BoolValue(newv)
	return nil
}

func (v BoolValue) Type() string {
	return "BoolValue"
}

// =======================
type FloatValue float64 // int64

func (v *FloatValue) String() string {
	return strconv.Itoa(int(*v))
}

func (v *FloatValue) Set(s string) error {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*v = FloatValue(i)
	return nil
}

func (v FloatValue) Type() string {
	return "FloatValue"
}

// =======================
