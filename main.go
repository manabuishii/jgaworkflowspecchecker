package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	// for CLI
	flag "github.com/spf13/pflag"
)

type RunData struct {
	PEOrSE  string `json:"se_or_pe"`
	FQ1     string `json:"fq1"`
	FQ1_MD5 string `json:"fq1_MD5"`
	FQ2     string `json:"fq2"`
	FQ2_MD5 string `json:"fq2_MD5"`
}

type Run struct {
	RunId string `json:"runid"`

	RunData `json:"data"`
}

type Sample struct {
	SampleId string `json:"sampleid"`
	Platform string `json:"platform"`
	RunList  []*Run `json:"runlist"`
}

type simpleSchema struct {
	Name string `json:"name"`

	SampleList []*Sample `json:"samplelist"`
}

//
type PathObject struct {
	Path   string `json:"path"`
	Format string `json:"format"`
}

type referenceSchema struct {
	Name                                         string      `json:"name"`
	Reference                                    *PathObject `json:"reference"`
	SortsamMaxRecordsInRam                       int         `json:"sortsam_max_records_in_ram"`
	SortsamJavaOptions                           string      `json:"sortsam_java_options"`
	BwaNumThreads                                int         `json:"bwa_num_threads"`
	BwaBasesPerBatch                             int         `json:"bwa_bases_per_batch"`
	UseBqsr                                      bool        `json:"use_bqsr"`
	Dbsnp                                        *PathObject `json:"dbsnp"`
	Mills                                        *PathObject `json:"mills"`
	KnownIndels                                  *PathObject `json:"known_indels"`
	SamtoolsNumThreads                           int         `json:"samtools_num_threads"`
	Gatk4HaplotypeCallerNumThreads               int         `json:"gatk4_HaplotypeCaller_num_threads"`
	BgzipNumThreads                              int         `json:"bgzip_num_threads"`
	HaplotypecallerAutosomePARPloidy2IntervalBed *PathObject `json:"haplotypecaller_autosome_PAR_ploidy_2_interval_bed"`
	HaplotypecallerChrXNonPARPloidy2IntervalBed  *PathObject `json:"haplotypecaller_chrX_nonPAR_ploidy_2_interval_bed"`
	HaplotypecallerChrXNonPARPloidy1IntervalBed  *PathObject `json:"haplotypecaller_chrX_nonPAR_ploidy_1_interval_bed"`
	HaplotypecallerChrYNonPARPloidy1IntervalBed  *PathObject `json:"haplotypecaller_chrY_nonPAR_ploidy_1_interval_bed"`
}

func md5File(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(hash.Sum(nil)[:16]), nil
}

func checkRunDataFile(fn string, fnmd5 string) (bool, error) {
	// Check file existance flag is set
	if fileExistsCheckFlag == false {
		return true, nil
	}
	// Check file is exist
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return false, err
	}
	// Check file hash
	if fileHashCheckFlag == false {
		return true, nil
	}
	// Check file hash value if specified
	result := true
	if fnmd5 != "" {
		md5, _ := md5File(fn)
		if fnmd5 != md5 {
			result = false
			fmt.Printf("expected: [%s]\n", fnmd5)
			fmt.Printf("actual  : [%s]\n", md5)
			fmt.Println("md5 is not match")
		}
	}
	return result, nil
}

func checkRunData(runData *RunData) (bool, error) {
	result := false
	if runData.PEOrSE == "PE" {
		r1, _ := checkRunDataFile(runData.FQ1, runData.FQ1_MD5)
		r2, _ := checkRunDataFile(runData.FQ2, runData.FQ2_MD5)
		result = r1 && r2
	} else {
		result, _ = checkRunDataFile(runData.FQ1, runData.FQ1_MD5)
	}
	return result, nil
}

func outputReference(rss *referenceSchema) (bool, error) {
	fmt.Println("")
	fmt.Printf("reference:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.Reference.Path)
	fmt.Printf("  format: http://edamontology.org/format_1929\n")
	fmt.Printf("sortsam_max_records_in_ram: %d\n", rss.SortsamMaxRecordsInRam)
	fmt.Printf("sortsam_java_options: %s\n", rss.SortsamJavaOptions)
	fmt.Printf("bwa_num_threads: %d\n", rss.BwaNumThreads)
	fmt.Printf("bwa_bases_per_batch: %d\n", rss.BwaBasesPerBatch)
	fmt.Printf("use_bqsr: %t\n", rss.UseBqsr)
	fmt.Printf("dbsnp:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.Dbsnp.Path)
	fmt.Printf("  format: http://edamontology.org/format_3016\n")
	fmt.Printf("mills:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.Mills.Path)
	fmt.Printf("  format: http://edamontology.org/format_3016\n")
	fmt.Printf("known_indels:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.KnownIndels.Path)
	fmt.Printf("  format: http://edamontology.org/format_3016\n")
	fmt.Printf("samtools_num_threads: %d\n", rss.SamtoolsNumThreads)
	fmt.Printf("gatk4_HaplotypeCaller_num_threads: %d\n", rss.Gatk4HaplotypeCallerNumThreads)
	fmt.Printf("bgzip_num_threads: %d\n", rss.BgzipNumThreads)
	fmt.Printf("haplotypecaller_autosome_PAR_ploidy_2_interval_bed:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.HaplotypecallerAutosomePARPloidy2IntervalBed.Path)
	fmt.Printf("  format: http://edamontology.org/format_3584\n")
	fmt.Printf("haplotypecaller_chrX_nonPAR_ploidy_2_interval_bed:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.HaplotypecallerChrXNonPARPloidy2IntervalBed.Path)
	fmt.Printf("  format: http://edamontology.org/format_3584\n")
	fmt.Printf("haplotypecaller_chrX_nonPAR_ploidy_1_interval_bed:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.HaplotypecallerChrXNonPARPloidy1IntervalBed.Path)
	fmt.Printf("  format: http://edamontology.org/format_3584\n")
	fmt.Printf("haplotypecaller_chrY_nonPAR_ploidy_1_interval_bed:\n")
	fmt.Printf("  class: File\n")
	fmt.Printf("  path: %s\n", rss.HaplotypecallerChrYNonPARPloidy1IntervalBed.Path)
	fmt.Printf("  format: http://edamontology.org/format_3584\n")

	return true, nil
}

var helpFlag bool
var fileExistsCheckFlag bool
var fileHashCheckFlag bool

func main() {
	flag.BoolVarP(&helpFlag, "help", "h", false, "show help message")
	flag.BoolVarP(&fileExistsCheckFlag, "file-exists-check", "", true, "Check file exists")
	flag.BoolVarP(&fileHashCheckFlag, "file-hash-check", "", true, "Check file hash value")
	flag.Parse()

	if helpFlag {
		flag.PrintDefaults()
		return
	}
	//return

	path, err := filepath.Abs("./")

	// MUST must be canonical
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[1])
	// MUST must be canonical
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[2])

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	fmt.Println("Load sample")
	raw, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var ss simpleSchema

	json.Unmarshal(raw, &ss)
	fmt.Println("Load end")
	fmt.Printf("Name is [%s]\n", ss.Name)

	//
	numOfPE := 0
	numOfSE := 0
	// validate
	for i, s := range ss.SampleList {
		fmt.Printf("index: %d, SampleId: %s\n", i, s.SampleId)
		for j, t := range s.RunList {
			if t.RunData.PEOrSE == "SE" {
				numOfSE = numOfSE + 1
			}
			if t.RunData.PEOrSE == "PE" {
				numOfPE = numOfPE + 1
			}
			fmt.Println(t)
			fmt.Printf("index: %d, RunId: %s\n", j, t.RunId)
			fmt.Printf("pe or se: [%s]\n", t.RunData.PEOrSE)
			fmt.Printf("fq1: [%s]\n", t.RunData.FQ1)
			fmt.Printf("fq2: [%s]\n", t.RunData.FQ2)
			r1, _ := checkRunData(&t.RunData)
			fmt.Printf("result=%t\n", r1)

		}
	}
	// reference config validate

	// MUST must be canonical
	rschemaLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[3])
	// MUST must be canonical
	rdocumentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[4])

	rresult, err := gojsonschema.Validate(rschemaLoader, rdocumentLoader)
	if err != nil {
		panic(err.Error())
	}

	if rresult.Valid() {
		fmt.Printf("The reference config document is valid\n")
	} else {
		fmt.Printf("The reference config document is not valid. see errors :\n")
		for _, desc := range rresult.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return
	}
	fmt.Println("Load sample")
	rraw, err := ioutil.ReadFile(os.Args[4])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var rss referenceSchema

	json.Unmarshal(rraw, &rss)
	fmt.Println("Load end")

	// output job file for CWL per each sampleID
	outputReference(&rss)
	for _, s := range ss.SampleList {
		fmt.Printf("sample_id: %s\n", s.SampleId)
		if numOfPE == 0 {
			fmt.Printf("runlist_pe: []\n")
		} else {
			fmt.Printf("runlist_pe:\n")
			for _, t := range s.RunList {
				if t.RunData.PEOrSE != "PE" {
					continue
				}
				fmt.Printf("  - run_id: %s\n", t.RunId)
				fmt.Printf("    platform_name: ILLUMINA\n")
				fmt.Printf("    fastq1:\n")
				fmt.Printf("      class: File\n")
				fmt.Printf("      path: %s\n", t.RunData.FQ1)
				fmt.Printf("      format: http://edamontology.org/format_1930\n")
				fmt.Printf("    fastq2:\n")
				fmt.Printf("      class: File\n")
				fmt.Printf("      path: %s\n", t.RunData.FQ2)
				fmt.Printf("      format: http://edamontology.org/format_1930\n")

				//r1, _ := checkRunData(&t.RunData)
				//fmt.Printf("result=%t\n", r1)
			}

		}
		if numOfSE == 0 {
			fmt.Printf("runlist_se: []\n")
		} else {
			fmt.Printf("runlist_se:\n")
			for _, t := range s.RunList {
				if t.RunData.PEOrSE != "SE" {
					continue
				}
				fmt.Printf("  - run_id: %s\n", t.RunId)
				fmt.Printf("    platform_name: ILLUMINA\n")
				fmt.Printf("    fastq1:\n")
				fmt.Printf("      class: File\n")
				fmt.Printf("      path: %s\n", t.RunData.FQ1)
				fmt.Printf("      format: http://edamontology.org/format_1930\n")

				//r1, _ := checkRunData(&t.RunData)
				//fmt.Printf("result=%t\n", r1)
			}

		}

	}

	// fmt.Printf("Md5 is [%s]\n", ss.Md5)

	// // validate
	// for i, s := range ss.SampleList {
	// 	fmt.Printf("index: %d, SampleId: %s\n", i, s.SampleId)
	// 	for j, t := range s.RunList {
	// 		fmt.Printf("index: %d, SampleId: %s\n", j, t.RunId)
	// 		r1, _ := checkRunData(&t.RunData)
	// 		fmt.Printf("result=%t\n", r1)

	// 	}
	// }
}
