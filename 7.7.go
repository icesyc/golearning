package main

import (
	"flag"
	"fmt"
)

type Celsius float64
type Fahrenheit float64
type Kelvin float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC Celsius = 0
	BoilingC Celsius = 100
	AbsoluteZeroK Kelvin = 0
)

func (c Celsius) String() string {
	return fmt.Sprintf("%g°C", c)
}
func (f Fahrenheit) String() string {
	return fmt.Sprintf("%g°F", f)
}
func (k Kelvin) String() string {
	return fmt.Sprintf("%g°K", k)
}

func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c * 9 / 5 + 32)
}

func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

func CToK(c Celsius) Kelvin {
	return Kelvin(c + AbsoluteZeroC)
}

func KToC(k Kelvin) Celsius {
	return Celsius(k) - AbsoluteZeroC
}

type celsiusFlag struct {
	Celsius
}

func (f *celsiusFlag) Set(s string ) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	case "K", "°K":
		f.Celsius = KToC(Kelvin(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

func CelsiusFlag(name string, defaultValue Celsius, usage string) *Celsius {
	f := celsiusFlag{defaultValue}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}

//参数20.0没有包含°C的情况下，默认使用Celsius类型 返回值为*Celsius, 所以temp的类型为*Celsius,
//在fmt打印时会调用 temp的String()方法，由于String方法的receiver为Celsius,所以temp会隐式转换为Celsius类型
//然后调用String()方法，会带着单位一起输出
var temp = CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()
	fmt.Printf("temp=%s\n", temp)
}