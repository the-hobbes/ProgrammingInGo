package main 

import (
	"fmt"
	"strings"
)

// BitFlag declare a custome type called BitFlag, whose values are integers
type BitFlag int

// create 3 bit flags of custom type BitFlag
// Groups of Go constants work in that the first one is set to its 0 value, unless 
// explicity set (either to value or iota). Subsequent values are set to their 
// predecessors value (including iota). If iota is used, each subsequent iota value
// is one more than the previous one.
const (
	// Active is set to 1 and also and iota, meaning the other constants will have
	// successive values.
	Active BitFlag = 1 << iota						// 1 << 0 == 1
	Send 		// Implicity BitFlag = 1 << iota	// 1 << 1 == 2
	Receive 	//Implicitly BitFlag = 1 << iota 	// 1 << 2 == 4
)

func (flag BitFlag) String() string {
	var flags []string
	if flag&Active == Active {
		flags = append(flags, "Active")
	}
	if flag&Send == Send {
		flags = append(flags, "Send")
	}
	if flag&Receive == Receive {
		flags = append(flags, "Receive")
	}
	if len(flags) > 0 {  
		// int(flag) is used to prevent infitite recursion, as the print function would 
		// continue to call String() on the flag type infinitely as it tried to print.
		// Instead, we convert the BitFlag to its underlying type (int) so we can print.
		return fmt.Sprintf("%d(%s)", int(flag), strings.Join(flags, "|"))
		// we can change the %d to a %b to print out the binary value of the flag.
	}
	return "0()"
}

func main() {
	// flag is set to the bitwise OR of Active | Send, which is 3
	flag := Active | Send
	// we use our overridden String() function to print this out
	fmt.Println(flag)
	// print each of the flags out
	fmt.Println(Active)
	fmt.Println(Send)
	fmt.Println(Receive)
	// what happens when we print out an unknown flag?
	var badFlag BitFlag
	fmt.Println(badFlag)
	// a slice of flags?
	var flags [3]BitFlag
	flags[0] = Active
	flags[1] = Send
	flags[2] = Receive
	fmt.Println(flags)
	// more bitwise operations
	fmt.Println(Active | Receive)
	fmt.Println(Active | Receive | Send)
}