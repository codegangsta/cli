package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

// BoolFlag is a flag with type bool
type BoolFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	EnvVars     []string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       bool
	DefaultText string
	Destination *bool
	HasBeenSet  bool
	Count       *int
}

// boolValue needs to implement the boolFlag internal interface in flag
// to be able to capture bool fields and values
// type boolFlag interface {
//	  Value
//	  IsBoolFlag() bool
// }
type boolValue struct {
	destination *bool
	count       *int
}

func newBoolValue(val bool, p *bool, count *int) *boolValue {
	*p = val
	return &boolValue{
		destination: p,
		count:       count,
	}
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		err = errors.New("parse error")
		return err
	}
	*b.destination = v
	if b.count != nil {
		*b.count = *b.count + 1
	}
	return err
}

func (b *boolValue) Get() interface{} { return *b.destination }

func (b *boolValue) String() string {
	if b.destination != nil {
		return strconv.FormatBool(*b.destination)
	}
	return strconv.FormatBool(false)
}

func (b *boolValue) IsBoolFlag() bool { return true }

// IsSet returns whether or not the flag has been set through env or file
func (f *BoolFlag) IsSet() bool {
	return f.HasBeenSet
}

// String returns a readable representation of this value
// (for usage defaults)
func (f *BoolFlag) String() string {
	return FlagStringer(f)
}

// Names returns the names of the flag
func (f *BoolFlag) Names() []string {
	return flagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *BoolFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise false
func (f *BoolFlag) TakesValue() bool {
	return false
}

// GetUsage returns the usage string for the flag
func (f *BoolFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *BoolFlag) GetValue() string {
	return ""
}

// IsVisible returns true if the flag is not hidden, otherwise false
func (f *BoolFlag) IsVisible() bool {
	return !f.Hidden
}

// Apply populates the flag given the flag set and environment
func (f *BoolFlag) Apply(set *flag.FlagSet) error {
	if val, ok := flagFromEnvOrFile(f.EnvVars, f.FilePath); ok {
		if val != "" {
			valBool, err := strconv.ParseBool(val)

			if err != nil {
				return fmt.Errorf("could not parse %q as bool value for flag %s: %s", val, f.Name, err)
			}

			f.Value = valBool
			f.HasBeenSet = true
		}
	}

	for _, name := range f.Names() {
		var value flag.Value
		if f.Destination != nil {
			value = newBoolValue(f.Value, f.Destination, f.Count)
		} else {
			t := new(bool)
			value = newBoolValue(f.Value, t, f.Count)
		}
		set.Var(value, name, f.Usage)
	}

	return nil
}

// Bool looks up the value of a local BoolFlag, returns
// false if not found
func (c *Context) Bool(name string) bool {
	if fs := c.lookupFlagSet(name); fs != nil {
		return lookupBool(name, fs)
	}
	return false
}

func lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return parsed
	}
	return false
}
