package repo

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

var testdata = []string{
	"VIEWS/index/Komplexität/190119e/00_190119e - Komplexität - 180522a, 190119d.txt",
	"VIEWS/index/Komplexität/190119e/01_190119d/00_190119d - Testing - clausen2021 87.txt",
	"VIEWS/index/Komplexität/190119e/02_180522a - Komplexität, Thermodynamische Tiefe - 210520var.png",
	"VIEWS/index/Komplexität/190119e/03_210520var - Varietät, Komplexität.txt",
	"VIEWS/index/Programmieren/220115p/00_220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5 - 220116s.pdf",
	"VIEWS/index/Programmieren/220115p/01_220116s - Spezifikation - Marco Fitz.pdf",
	"VIEWS/explore/220115p/210328obj - Objektorientiert, Programmierung - kernighan2016 155 - 170224a.pdf",
	"VIEWS/explore/220115p/190119d - Testing - clausen2021 87.txt",
	"VIEWS/explore/220115p/220116s - Spezifikation - Marco Fitz.pdf",
	"VIEWS/citations/clausen2021/190119d - Testing - clausen2021 87.txt",
	"VIEWS/citations/clausen2021/220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5 - 220116s.pdf",
	"VIEWS/context/Marco Fitz/220115p - Refactoring, Programmieren - Marco Fitz, clausen2021 5 - 220116s.pdf",
	"VIEWS/context/Marco Fitz/220116s - Spezifikation - Marco Fitz.pdf",
	"VIEWS/keywords/Komplexität/180522a - Komplexität, Thermodynamische Tiefe - 210520var.png",
	"VIEWS/keywords/Komplexität/190119e - Komplexität - 180522a, 190119d.txt",
	"VIEWS/keywords/Komplexität/210520var - Varietät, Komplexität.txt",
}

func BenchmarkCreateSyml(b *testing.B) {
	// Arrange
	wd, err := os.Getwd()
	if err != nil {
		b.Errorf("could not get the current working dir")
	}
	// Remove this directory, which might got created in a previous run
	var pathTestRepo = wd + "/testdata"
	fmt.Println("Clearing: " + pathTestRepo)
	clearPath(pathTestRepo) // This removing doesn't work yet, if it exists remove manually

	// Create about 2000 links for our benchmarking
	var testlinks []string
	for i := 1; i < 100; i++ {
		for _, testdatum := range testdata {
			testlinks = append(testlinks, "repo "+strconv.Itoa(i)+" - "+testdatum)
		}
	}

	// Run benchmark
	// This will test how fast a list of 2.000 links can be persisted
	for i := 0; i < b.N; i++ {
		for _, testlink := range testlinks {
			tl := "testdata/BenchmRun " + strconv.Itoa(i+1) + "/" + testlink
			persist(tl, tl) // we will just link to itself
		}
	}
}

func clearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Error occured: %v", err)
	}
}
