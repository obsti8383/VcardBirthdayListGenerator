// VcardBirthdayListGenerator
// Copyright (C) 2020 Florian Probst
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mapaiva/vcard-go"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(versionCmd, csvCmd, textCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "VcardBirthdayListGenerator",
	Short: "VcardBirthdayListGenerator generates a birthday list as csv or text (to stdout) from vcf files",
	Long:  "VcardBirthdayListGenerator generates a birthday list as csv or text (to stdout) from vcf files",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var csvCmd = &cobra.Command{
	Use:   "csv <path>",
	Short: "generates a birthday list as csv (to stdout) from vcf files",
	Long:  "generates a birthday list as csv (to stdout) from vcf files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("name;month;day;year;error")

		// walk a files in directory
		err := filepath.Walk(args[0], evaluateVCardsCsv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var textCmd = &cobra.Command{
	Use:   "text <path>",
	Short: "generates a birthday list as text (to stdout) from vcf files",
	Long:  "generates a birthday list as text (to stdout) from vcf files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// walk a files in directory
		err := filepath.Walk(args[0], evaluateVCardsTxt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of VcardBirthdayListGenerator",
	Long:  "Print the version of VcardBirthdayListGenerator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VcardBirthdayListGenerator v0.2")
	},
}

// function for filepath.Walk() which does all the work of parsing VCF files
func evaluateVCardsCsv(path string, info os.FileInfo, err error) error {
	return evaluateVCards(path, info, err, true)
}

func evaluateVCardsTxt(path string, info os.FileInfo, err error) error {
	return evaluateVCards(path, info, err, false)
}

func evaluateVCards(path string, info os.FileInfo, err error, csvFlag bool) error {
	cards, err := collectVCards(path, info, err, false)
	if err != nil {
		return err
	}
	sort.Sort(VCards(cards))
	printVCards(cards, csvFlag)
	return nil
}

func printVCards(cards []vcard.VCard, csvFlag bool) {
	// iterate over all found cards
	for _, card := range cards {
		if card == (vcard.VCard{}) {
			if csvFlag {
				fmt.Println(";;;;VCard seems to be empty")
			} else {
				fmt.Println("VCard seems to be empty")
			}
		} else {
			bd := card.BirthDay
			nameSplit := strings.Split(card.StructuredName, ";")
			name := nameSplit[0] + " " + nameSplit[1]
			if bd == "" {
				if csvFlag {
					fmt.Println(name + ";;;;None")
				} else {
					fmt.Println(name + ": None")
				}
			} else {
				// check the different date formats which are used in VCARDs
				bdTime, err := time.Parse("20060102", bd)
				if err != nil {
					bdTime, err = time.Parse("2006-01-02", bd)
					if err != nil {
						if strings.HasPrefix(bd, "--") {
							// year of birth unknown
							bd = strings.TrimPrefix(bd, "--")
							bd = "0001" + bd
							bdTime, err = time.Parse("20060102", bd)
							if err != nil {
								if csvFlag {
									fmt.Println(name + ";;;;Could not parse birthday date with suffix -- correctly: " + bd)
								} else {
									fmt.Println(name + ": Could not parse birthday date with suffix -- correctly: " + bd)
								}
							}
						} else {
							if csvFlag {
								fmt.Println(name + ";;;;BirthDay has unknown format: " + bd)
							} else {
								fmt.Println(name + ": BirthDay has unknown format: " + bd)
							}
						}
					}
				}

				// print birthday
				if csvFlag {
					if !bdTime.IsZero() {
						if bdTime.Year() != 1 {
							fmt.Printf("%s;%02d;%02d;%d;\n", name, int(bdTime.Month()), bdTime.Day(), bdTime.Year())
						} else {
							fmt.Printf("%s;%02d;%02d;;\n", name, int(bdTime.Month()), bdTime.Day())
						}
					} else {
						fmt.Println(name + ";;;;Could not evaluate birthday")
					}
				} else {
					if !bdTime.IsZero() {

						if bdTime.Year() != 1 {
							fmt.Printf("%s: %02d.%02d.%d\n", name, bdTime.Day(), int(bdTime.Month()), bdTime.Year())
						} else {
							fmt.Printf("%s: %02d.%02d.\n", name, bdTime.Day(), int(bdTime.Month()))
						}
					} else {
						fmt.Println(name + ": Could not evaluate birthday")
					}
				}
			}
		}
	}
}

func collectVCards(path string, info os.FileInfo, err error, csvFlag bool) ([]vcard.VCard, error) {
	var allCards []vcard.VCard

	_, err = os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		// parse file for VCARDs
		cards, err := vcard.GetVCards(path)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// verify if any card was in the file
		if len(cards) == 0 {
			if csvFlag {
				fmt.Println(";;;;No vcard found in file ", path)
			} else {
				fmt.Println("No vcard found in file ", path)
			}
		} else {
			allCards = append(allCards, cards...)
		}

	}

	return allCards, nil
}
