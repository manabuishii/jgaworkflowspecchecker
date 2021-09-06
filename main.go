package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/sync/errgroup"

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

//
var version string
var revision string

//

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

func outputReference(rss *referenceSchema) (string, error) {
	var byteBuf bytes.Buffer
	byteBuf.WriteString("")
	byteBuf.WriteString("reference:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.Reference.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_1929\n")
	byteBuf.WriteString(fmt.Sprintf("sortsam_max_records_in_ram: %d\n", rss.SortsamMaxRecordsInRam))
	byteBuf.WriteString(fmt.Sprintf("sortsam_java_options: %s\n", rss.SortsamJavaOptions))
	byteBuf.WriteString(fmt.Sprintf("bwa_num_threads: %d\n", rss.BwaNumThreads))
	byteBuf.WriteString(fmt.Sprintf("bwa_bases_per_batch: %d\n", rss.BwaBasesPerBatch))
	byteBuf.WriteString(fmt.Sprintf("use_bqsr: %t\n", rss.UseBqsr))
	byteBuf.WriteString("dbsnp:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.Dbsnp.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3016\n")
	byteBuf.WriteString("mills:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.Mills.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3016\n")
	byteBuf.WriteString("known_indels:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.KnownIndels.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3016\n")
	byteBuf.WriteString(fmt.Sprintf("samtools_num_threads: %d\n", rss.SamtoolsNumThreads))
	byteBuf.WriteString(fmt.Sprintf("gatk4_HaplotypeCaller_num_threads: %d\n", rss.Gatk4HaplotypeCallerNumThreads))
	byteBuf.WriteString(fmt.Sprintf("bgzip_num_threads: %d\n", rss.BgzipNumThreads))
	byteBuf.WriteString("haplotypecaller_autosome_PAR_ploidy_2_interval_bed:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerAutosomePARPloidy2IntervalBed.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3584\n")
	byteBuf.WriteString("haplotypecaller_chrX_nonPAR_ploidy_2_interval_bed:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerChrXNonPARPloidy2IntervalBed.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3584\n")
	byteBuf.WriteString("haplotypecaller_chrX_nonPAR_ploidy_1_interval_bed:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerChrXNonPARPloidy1IntervalBed.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3584\n")
	byteBuf.WriteString("haplotypecaller_chrY_nonPAR_ploidy_1_interval_bed:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerChrYNonPARPloidy1IntervalBed.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3584\n")

	return byteBuf.String(), nil
}

// call per sample
func outputJobFile(s *Sample, rss *referenceSchema) (string, error) {
	//
	var byteBuf bytes.Buffer

	// count SE and PE entry
	numOfSE := 0
	numOfPE := 0
	for _, t := range s.RunList {
		if t.RunData.PEOrSE == "SE" {
			numOfSE = numOfSE + 1
		}
		if t.RunData.PEOrSE == "PE" {
			numOfPE = numOfPE + 1
		}
	}
	//
	byteBuf.WriteString(fmt.Sprintf("sample_id: %s\n", s.SampleId))
	if numOfPE == 0 {
		byteBuf.WriteString("runlist_pe: []\n")
	} else {
		byteBuf.WriteString("runlist_pe:\n")
		for _, t := range s.RunList {
			if t.RunData.PEOrSE != "PE" {
				continue
			}
			byteBuf.WriteString(fmt.Sprintf("  - run_id: %s\n", t.RunId))
			byteBuf.WriteString("    platform_name: ILLUMINA\n")
			byteBuf.WriteString("    fastq1:\n")
			byteBuf.WriteString("      class: File\n")
			byteBuf.WriteString(fmt.Sprintf("      path: %s\n", t.RunData.FQ1))
			byteBuf.WriteString("      format: http://edamontology.org/format_1930\n")
			byteBuf.WriteString("    fastq2:\n")
			byteBuf.WriteString("      class: File\n")
			byteBuf.WriteString(fmt.Sprintf("      path: %s\n", t.RunData.FQ2))
			byteBuf.WriteString("      format: http://edamontology.org/format_1930\n")
		}
	}
	if numOfSE == 0 {
		byteBuf.WriteString("runlist_se: []\n")
	} else {
		byteBuf.WriteString("runlist_se:\n")
		for _, t := range s.RunList {
			if t.RunData.PEOrSE != "SE" {
				continue
			}
			byteBuf.WriteString(fmt.Sprintf("  - run_id: %s\n", t.RunId))
			byteBuf.WriteString("    platform_name: ILLUMINA\n")
			byteBuf.WriteString("    fastq1:\n")
			byteBuf.WriteString("      class: File\n")
			byteBuf.WriteString(fmt.Sprintf("      path: %s\n", t.RunData.FQ1))
			byteBuf.WriteString("      format: http://edamontology.org/format_1930\n")
		}
	}

	return byteBuf.String(), nil
}

func createJobFile(ss *simpleSchema, rss *referenceSchema) error {
	for _, s := range ss.SampleList {
		// create file
		//
		filename := fmt.Sprintf("%s_jobfile.yaml", s.SampleId)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		// output reference data to job file per each sampleID
		referenceData, _ := outputReference(rss)
		if _, err := writer.WriteString(referenceData); err != nil {
			return err
		}
		sampleData, _ := outputJobFile(s, rss)
		if _, err := writer.WriteString(sampleData); err != nil {
			return err
		}

		// Flush
		writer.Flush()

	}
	return nil
}

func execCWL(sampleId string) string {
	// execute toil
	//p, _ := os.Getwd()
	// c1 := exec.Command("toil-cwl-runner", "--maxDisk", "248G", "--maxMemory", "64G", "--defaultMemory", "32000", "--defaultDisk", "32000", "--workDir", p, "--disableCaching", "--jobStore", "./"+sampleId+"-jobstore", "--outdir", "./"+sampleId, "--stats", "--cleanWorkDir", "never", "--batchSystem", "slurm", "--retryCount", "1", "--singularity", "--logFile", sampleId+".log", "per-sample/Workflows/per-sample.cwl", sampleId+"_jobfile.yaml")
	c1 := exec.Command("toil-cwl-runner", "--maxDisk", "248G", "--maxMemory", "64G", "--defaultMemory", "32000", "--defaultDisk", "32000", "--disableCaching", "--jobStore", "./"+sampleId+"-jobstore", "--outdir", "./"+sampleId, "--stats", "--batchSystem", "slurm", "--retryCount", "1", "--singularity", "--logFile", sampleId+".log", "per-sample/Workflows/per-sample.cwl", sampleId+"_jobfile.yaml")
	// set environment value if needed
	//c1.Env = append(os.Environ(), "TOIL_SLURM_ARGS=\"-w node[1-9]\"")
	c1.Start()
	c1.Wait()
	return ""
}

var dryrunFlag bool
var helpFlag bool
var versionFlag bool
var fileExistsCheckFlag bool
var fileHashCheckFlag bool

func main() {
	flag.BoolVarP(&dryrunFlag, "dry-run", "n", false, "Dry-run, do not execute acutal command")
	flag.BoolVarP(&helpFlag, "help", "h", false, "Show help message")
	flag.BoolVarP(&versionFlag, "version", "v", false, "Show version")
	flag.BoolVarP(&fileExistsCheckFlag, "file-exists-check", "", true, "Check file exists")
	flag.BoolVarP(&fileHashCheckFlag, "file-hash-check", "", true, "Check file hash value")
	flag.Parse()

	if helpFlag {
		fmt.Printf("Version: %s-%s\n", version, revision)
		flag.PrintDefaults()
		return
	}
	if versionFlag {
		fmt.Printf("Version: %s-%s\n", version, revision)
		return
	}

	if dryrunFlag {
		fmt.Println("Dry-run flag is set")
		return
	}
	fmt.Println("Dry-run flag is not set")
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

	// validate
	checkResult := false
	for _, s := range ss.SampleList {
		//fmt.Printf("Check index: %d, SampleId: %s\n", i, s.SampleId)
		for j, t := range s.RunList {
			// fmt.Println(t)
			// fmt.Printf("index: %d, RunId: %s\n", j, t.RunId)
			// fmt.Printf("pe or se: [%s]\n", t.RunData.PEOrSE)
			// fmt.Printf("fq1: [%s]\n", t.RunData.FQ1)
			// fmt.Printf("fq2: [%s]\n", t.RunData.FQ2)
			r1, _ := checkRunData(&t.RunData)
			checkResult = checkResult || r1
			if r1 {
				fmt.Println("Some error found. Not exist or Hash value error")
				fmt.Printf("Check index: %d, RunId: %s\n", j, t.RunId)
				fmt.Printf("pe or se: [%s]\n", t.RunData.PEOrSE)
				fmt.Printf("fq1: [%s]\n", t.RunData.FQ1)
				fmt.Printf("fq2: [%s]\n", t.RunData.FQ2)
				fmt.Printf("result=%t\n", r1)
			}
		}
	}
	if checkResult {
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

	// create job file for CWL
	createJobFile(&ss, &rss)

	// exec and wait
	var eg errgroup.Group
	for i, s := range ss.SampleList {
		fmt.Printf("index: %d, SampleId: %s\n", i, s.SampleId)
		sampleId := s.SampleId
		eg.Go(func() error {
			// time.Sleep(2 * time.Second) // 長い処理
			// if i > 90 {
			// 	fmt.Println("Error:", i)
			// 	return fmt.Errorf("Error occurred: %d", i)
			// }
			// fmt.Println("End:", i)
			execCWL(sampleId)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
