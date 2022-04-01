package main

// Molarmass by Thijmen Voskuilen (https://github.com/thijmer)

// License: GPLv3.0
// Contact me at thijmenvoskuilen@gmail.com if you have any questions.
// Source code: https://github.com/Thijmer/molarmass-go

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

var running bool = true // Set to false when an EOF is encountered.
var silent bool = false
var interactive bool = false

var weights = map[string]float32{
	"H":   1.008,
	"He":  4.0026022,
	"Li":  6.94,
	"Be":  9.01218315,
	"B":   10.81,
	"C":   12.011,
	"N":   14.007,
	"O":   15.999,
	"F":   18.9984031636,
	"Ne":  20.17976,
	"Na":  22.989769282,
	"Mg":  24.305,
	"Al":  26.98153857,
	"Si":  28.085,
	"P":   30.9737619985,
	"S":   32.06,
	"Cl":  35.45,
	"Ar":  39.9481,
	"K":   39.09831,
	"Ca":  40.0784,
	"Sc":  44.9559085,
	"Ti":  47.8671,
	"V":   50.94151,
	"Cr":  51.99616,
	"Mn":  54.9380443,
	"Fe":  55.8452,
	"Co":  58.9331944,
	"Ni":  58.69344,
	"Cu":  63.5463,
	"Zn":  65.382,
	"Ga":  69.7231,
	"Ge":  72.6308,
	"As":  74.9215956,
	"Se":  78.9718,
	"Br":  79.904,
	"Kr":  83.7982,
	"Rb":  85.46783,
	"Sr":  87.621,
	"Y":   88.905842,
	"Zr":  91.2242,
	"Nb":  92.906372,
	"Mo":  95.951,
	"Tc":  98,
	"Ru":  101.072,
	"Rh":  102.905502,
	"Pd":  106.421,
	"Ag":  107.86822,
	"Cd":  112.4144,
	"In":  114.8181,
	"Sn":  118.7107,
	"Sb":  121.7601,
	"Te":  127.603,
	"I":   126.904473,
	"Xe":  131.2936,
	"Cs":  132.905451966,
	"Ba":  137.3277,
	"La":  138.905477,
	"Ce":  140.1161,
	"Pr":  140.907662,
	"Nd":  144.2423,
	"Pm":  145,
	"Sm":  150.362,
	"Eu":  151.9641,
	"Gd":  157.253,
	"Tb":  158.925352,
	"Dy":  162.5001,
	"Ho":  164.930332,
	"Er":  167.2593,
	"Tm":  168.934222,
	"Yb":  173.0451,
	"Lu":  174.96681,
	"Hf":  178.492,
	"Ta":  180.947882,
	"W":   183.841,
	"Re":  186.2071,
	"Os":  190.233,
	"Ir":  192.2173,
	"Pt":  195.0849,
	"Au":  196.9665695,
	"Hg":  200.5923,
	"Tl":  204.38,
	"Pb":  207.21,
	"Bi":  208.980401,
	"Po":  209,
	"At":  210,
	"Rn":  222,
	"Fr":  223,
	"Ra":  226,
	"Ac":  227,
	"Th":  232.03774,
	"Pa":  231.035882,
	"U":   238.028913,
	"Np":  237,
	"Pu":  244,
	"Am":  243,
	"Cm":  247,
	"Bk":  247,
	"Cf":  251,
	"Es":  252,
	"Fm":  257,
	"Md":  258,
	"No":  259,
	"Lr":  266,
	"Rf":  267,
	"Db":  268,
	"Sg":  269,
	"Bh":  270,
	"Hs":  269,
	"Mt":  278,
	"Ds":  281,
	"Rg":  282,
	"Cn":  285,
	"Nh":  286,
	"Fl":  289,
	"Mc":  289,
	"Lv":  293,
	"Ts":  294,
	"Og":  294,
	"Uue": 315,
}

func input(message string) string {
	fmt.Print(message)
	// var ans string
	// _, err := fmt.Scanln(&ans)
	in := bufio.NewReader(os.Stdin)
	ans, err := in.ReadString('\n')
	if errors.Is(err, io.EOF) {
		running = false
	}
	ans = ans[:len(ans)-1]
	return ans
}

func deletespaces(txt string) string {
	output := ""
	for _, letter := range txt {
		if letter != ' ' {
			output += string(letter)
		}
	}
	return output
}

func help() {
	fmt.Println(`Molarmass GO 1.0

Usage:
$> molarmass C6H12O6          #calculates the mass of C6H12O6
$> molarmass -i               #Enter interactive mode. (In this mode, you can get the weight of as many molecules as you want.)
$> molarmass -s               #Enter silent mode. (In this mode, molarmass won't say anything but the molecule weights.)

Limitations:
Molarmass is almost as stupid as an iPhone user. It doesn't know what to do with parentheses or molecule charges. You will need to make things simple to understand for molarmass.
Examples and fixes:
$> molarmass CH₃-COOH     ➜ $> molarmass CH3COOH
$> molarmass 3CO2         ➜ $> molarmass C3O6
$> molarmass (CH3)2C=O    ➜ $> molarmass C2H6CO
$> molarmass NH4+         ➜ $> molarmass NH4`)
}

type atomandcount struct {
	atom  string
	count uint16
}

func calc_molecule(molecule string) float32 {
	var separated_atoms []string
	{
		current_atom := ""
		for _, char := range molecule {
			if unicode.IsUpper(char) {
				if current_atom != "" {
					separated_atoms = append(separated_atoms, current_atom)
				}
				current_atom = string(char)
			} else {
				current_atom += string(char)
			}
		}
		if current_atom != "" {
			separated_atoms = append(separated_atoms, current_atom)
		}
	}
	var atoms []atomandcount
	for _, atomac := range separated_atoms {
		current_atom := atomandcount{}
		var count uint64
		var splitindex int = len(atomac)
		for index := range atomac {
			var err error
			count, err = strconv.ParseUint(atomac[index:], 10, 16)
			if err == nil {
				splitindex = index
				break
			}
		}
		if splitindex == len(atomac) {
			count = 1
		}
		current_atom.atom = atomac[:splitindex]
		current_atom.count = uint16(count)
		atoms = append(atoms, current_atom)
	}

	var weight float32 = 0
	for _, atom := range atoms {
		w, ok := weights[atom.atom]
		if !ok {
			if !silent {
				os.Stderr.WriteString("Atom \"" + atom.atom + "\" doesn't exist. (for as far as I know...)\n")
			}
			return -1
		}
		weight += w * float32(atom.count)
	}
	return weight
}

func main() {
	var args []string = os.Args[1:]
	if len(args) == 0 {
		help()
		os.Exit(0)
	}
	for _, arg := range args {
		if arg == "-s" || arg == "--silent" {
			silent = true
		} else if arg == "-i" || arg == "--interactive" {
			interactive = true
		} else if arg == "-si" || arg == "-is" {
			interactive = true
			silent = true
		} else if arg == "help" || arg == "-h" || arg == "--help" || arg == "-help" {
			help()
			os.Exit(0)
		} else {
			w := calc_molecule(arg)
			if w != -1 {
				fmt.Println(w)
			}
		}
	}
	if interactive {
		interactive_mode()
	}
}

func interactive_mode() {
	if !silent {
		fmt.Println(`Molarmass interactive mode
Type a molecule and Molarmass will tell you how much it weighs.
Type "help" for help, "license" for the license and "exit" to exit.`)
	}
	for {
		var command string
		if silent {
			command = input("")
		} else {
			command = input("Molecule formula: ")
		}
		if command == "exit" {
			if !silent {
				fmt.Println("Bye!")
			}
			break
		} else if command == "help" {
			fmt.Println(`===========================
 Molarmass interactive mode help
 Type a molecule and Molarmass will tell you how much it weighs.
 Type "help" to see this text and type "exit" to exit.
 
 Molarmass is a simple tool to calculate the molar mass of a molecule.
 You are now in Molarmass interactive mode.
 That means that you can type molecule formulas and molarmass will answer with the molar mass of that molecule.
 
 Keep in mind that Molarmass is almost as stupid as an Airpod owner. It doesn't know what to do with parentheses or molecule charges. You will need to make things simple to understand for molarmass.
 Examples and fixes:
 CH₃-COOH     ➜ CH3COOH
 3CO2         ➜ C3O6
 (CH3)2C=O    ➜ C2H6CO
 NH4+         ➜ NH4
===========================`)
		} else if command == "license" {
			fmt.Println(
				`===========================
 Molarmass is distributed under the GPLv3.0 license.
 More information about the license: https://raw.githubusercontent.com/Thijmer/molarmass-go/master/LICENSE
 Source code: https://github.com/Thijmer/molarmass-go
===========================`)
		} else if command == "" {
		} else {
			w := calc_molecule(deletespaces(command))
			if w != -1 {
				fmt.Println(w)
			}
		}
		if !running {
			break
		}
	}
}
