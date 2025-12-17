package identity_helper

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	_ "sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Agg struct {
	ShortName  string
	CPUCoreMax int
	RAMGBMax   int
	DiskGBMax  int
}

func main() {
	srcXLSXPath := flag.String("source-xlsx", "source.xlsx", "detailed source XLSX (private_ip,...,short_name,CPU,RAM,Hard,...)")
	srcSheet := flag.String("sheet", "", "worksheet name (defaults to first sheet)")
	tgtPath := flag.String("target", "target.csv", "existing target CSV (name,short_name,cpu,ram,disk,...)")
	outPath := flag.String("out", "updated.csv", "path for updated target")
	diffPath := flag.String("diff", "diff.csv", "path for differences")
	flag.Parse()

	// 1) aggregate from XLSX by short_name
	srcAgg, err := aggregateFromXLSX(*srcXLSXPath, *srcSheet)
	if err != nil {
		log.Fatal(err)
	}

	// 2) load target CSV
	tgtRows, tgtHeader, err := readCSV(*tgtPath)
	if err != nil {
		log.Fatal(err)
	}
	col := headerIndex(tgtHeader)
	required := []string{"name", "short_name", "cpu", "ram", "disk"}
	for _, f := range required {
		if _, ok := col[f]; !ok {
			log.Fatalf("target header missing column %q", f)
		}
	}

	// outputs
	updated := make([][]string, 0, len(tgtRows)+1)
	diff := [][]string{{"short_name", "old_cpu", "old_ram", "old_disk", "new_cpu", "new_ram", "new_disk"}}
	updated = append(updated, tgtHeader)

	for _, rec := range tgtRows {
		shortName := strings.TrimSpace(rec[col["short_name"]])
		oldCPU, oldRAM, oldDisk := rec[col["cpu"]], rec[col["ram"]], rec[col["disk"]]

		newCPU, newRAM, newDisk := oldCPU, oldRAM, oldDisk
		if a, ok := srcAgg[strings.ToUpper(shortName)]; ok {
			newCPU = fmt.Sprintf("%d Core", a.CPUCoreMax)
			newRAM = fmt.Sprintf("%d GB", a.RAMGBMax)
			newDisk = humanDisk(a.DiskGBMax)
		}

		changed := normalizeCPU(oldCPU) != normalizeCPU(newCPU) ||
			normalizeGB(oldRAM) != normalizeGB(newRAM) ||
			normalizeDisk(oldDisk) != normalizeDisk(newDisk)

		rec[col["cpu"]] = newCPU
		rec[col["ram"]] = newRAM
		rec[col["disk"]] = newDisk
		updated = append(updated, rec)

		if changed {
			diff = append(diff, []string{shortName, oldCPU, oldRAM, oldDisk, newCPU, newRAM, newDisk})
		}
	}

	if err := writeCSV(*outPath, updated); err != nil {
		log.Fatal(err)
	}
	if err := writeCSV(*diffPath, diff); err != nil {
		log.Fatal(err)
	}

	log.Printf("Updated: %s", *outPath)
	log.Printf("Diff:    %s (changed rows: %d)", *diffPath, len(diff)-1)
}

// ------ XLSX source aggregation ------

func aggregateFromXLSX(path, sheetName string) (map[string]Agg, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	if sheetName == "" {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("xlsx has no sheets")
		}
		sheetName = sheets[0]
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("sheet %q is empty", sheetName)
	}

	header := lowerHeader(rows[0])
	col := headerIndex(header)
	need := []string{"short_name", "cpu", "ram", "hard"}
	for _, k := range need {
		if _, ok := col[k]; !ok {
			return nil, fmt.Errorf("xlsx header missing column %q", k)
		}
	}

	agg := map[string]Agg{}
	for i := 1; i < len(rows); i++ {
		rec := padRow(rows[i], len(header))
		sn := strings.TrimSpace(rec[col["short_name"]])
		if sn == "" {
			continue
		}
		key := strings.ToUpper(sn)
		cpu := parseCPU(rec[col["cpu"]])
		ram := parseGB(rec[col["ram"]])
		disk := parseDiskToGB(rec[col["hard"]])

		a := agg[key]
		a.ShortName = sn
		if cpu > a.CPUCoreMax {
			a.CPUCoreMax = cpu
		}
		if ram > a.RAMGBMax {
			a.RAMGBMax = ram
		}
		if disk > a.DiskGBMax {
			a.DiskGBMax = disk
		}
		agg[key] = a
	}
	return agg, nil
}

// ------ CSV helpers ------

func readCSV(path string) (rows [][]string, header []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	header, err = r.Read()
	if err != nil {
		return nil, nil, err
	}
	for {
		rec, e := r.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, nil, e
		}
		for i := range rec {
			rec[i] = strings.TrimSpace(rec[i])
		}
		rows = append(rows, rec)
	}
	return rows, lowerHeader(header), nil
}

func writeCSV(path string, rows [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	for _, rec := range rows {
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	return nil
}

func padRow(r []string, n int) []string {
	if len(r) >= n {
		return r
	}
	out := make([]string, n)
	copy(out, r)
	return out
}

func lowerHeader(h []string) []string {
	out := make([]string, len(h))
	for i, s := range h {
		out[i] = strings.ToLower(strings.TrimSpace(s))
	}
	return out
}

func headerIndex(h []string) map[string]int {
	m := map[string]int{}
	for i, s := range h {
		m[strings.TrimSpace(strings.ToLower(s))] = i
	}
	return m
}

// ------ parsing / normalization ------

var numRe = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(TB|GB|Core)?`)

func parseCPU(s string) int {
	m := numRe.FindStringSubmatch(s)
	if len(m) == 0 {
		return 0
	}
	v, _ := strconv.ParseFloat(m[1], 64)
	return int(v + 0.5)
}

func parseGB(s string) int {
	m := numRe.FindStringSubmatch(s)
	if len(m) == 0 {
		return 0
	}
	v, _ := strconv.ParseFloat(m[1], 64)
	return int(v + 0.5)
}

func parseDiskToGB(s string) int {
	total := 0
	ms := numRe.FindAllStringSubmatch(s, -1)
	for _, m := range ms {
		if len(m) < 3 {
			continue
		}
		val, _ := strconv.ParseFloat(m[1], 64)
		unit := strings.ToUpper(strings.TrimSpace(m[2]))
		switch unit {
		case "TB":
			total += int(val*1024 + 0.5)
		default:
			total += int(val + 0.5)
		}
	}
	return total
}

func humanDisk(gb int) string {
	if gb >= 1024 {
		tb := float64(gb) / 1024
		if abs(tb-float64(int(tb))) < 0.05 {
			return fmt.Sprintf("%d TB", int(tb+0.5))
		}
		return fmt.Sprintf("%.1f TB", tb)
	}
	return fmt.Sprintf("%d GB", gb)
}

func normalizeCPU(s string) string { return fmt.Sprintf("%d Core", parseCPU(s)) }
func normalizeGB(s string) string  { return fmt.Sprintf("%d GB", parseGB(s)) }
func normalizeDisk(s string) string {
	return humanDisk(parseDiskToGB(s))
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
