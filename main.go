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
	"strings"
	"time"

	"github.com/manabuishiii/jgaworkflowspecchecker/cmd"
	"github.com/manabuishiii/jgaworkflowspecchecker/utils"

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

// path and format
type PathObject struct {
	Path   string `json:"path"`
	Format string `json:"format"`
}

// path only
type PathOnlyObject struct {
	Path string `json:"path"`
}

type referenceSchema struct {
	WorkflowFile                           *PathOnlyObject `json:"workflow_file"`
	OutputDirectory                        *PathOnlyObject `json:"output_directory"`
	Reference                              *PathObject     `json:"reference"`
	SortsamMaxRecordsInRam                 int             `json:"sortsam_max_records_in_ram"`
	SortsamJavaOptions                     string          `json:"sortsam_java_options"`
	BwaNumThreads                          int             `json:"bwa_num_threads"`
	BwaBasesPerBatch                       int             `json:"bwa_bases_per_batch"`
	UseBqsr                                bool            `json:"use_bqsr"`
	Dbsnp                                  *PathObject     `json:"dbsnp"`
	Mills                                  *PathObject     `json:"mills"`
	KnownIndels                            *PathObject     `json:"known_indels"`
	SamtoolsNumThreads                     int             `json:"samtools_num_threads"`
	Gatk4HaplotypeCallerNumThreads         int             `json:"gatk4_HaplotypeCaller_num_threads"`
	BgzipNumThreads                        int             `json:"bgzip_num_threads"`
	HaplotypecallerAutosomePARIntervalBed  *PathObject     `json:"haplotypecaller_autosome_PAR_interval_bed"`
	HaplotypecallerAutosomePARIntervalList *PathOnlyObject `json:"haplotypecaller_autosome_PAR_interval_list"`
	HaplotypecallerChrXNonPARIntervalBed   *PathObject     `json:"haplotypecaller_chrX_nonPAR_interval_bed"`
	HaplotypecallerChrXNonPARIntervalList  *PathOnlyObject `json:"haplotypecaller_chrX_nonPAR_interval_list"`
	HaplotypecallerChrYNonPARIntervalBed   *PathObject     `json:"haplotypecaller_chrY_nonPAR_interval_bed"`
	HaplotypecallerChrYNonPARIntervalList  *PathOnlyObject `json:"haplotypecaller_chrY_nonPAR_interval_list"`
}

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

func isExistsAllResultFilesPrefixRunId(outputDirectoryPath string, runId string) bool {
	result := true
	fn := filepath.Join(outputDirectoryPath, runId)
	for _, extension := range []string{".bam", ".bam.log"} {
		if _, err := os.Stat(fn + extension); os.IsNotExist(err) {
			fmt.Printf("Missing file [%s]\n", fn+extension)
			result = false
		}

	}
	return result
}
func isExistsAllResultFilesPrefixSampleId(outputDirectoryPath string, sampleId string) bool {
	result := true
	fn := filepath.Join(outputDirectoryPath, sampleId)
	for _, extension := range []string{".autosome_PAR_ploidy_2.g.vcf.gz",
		".autosome_PAR_ploidy_2.g.vcf.gz.bcftools-stats",
		".autosome_PAR_ploidy_2.g.vcf.gz.bcftools-stats.log",
		".autosome_PAR_ploidy_2.g.vcf.gz.log",
		".autosome_PAR_ploidy_2.g.vcf.gz.tbi",
		".autosome_PAR_ploidy_2.g.vcf.gz.tbi.log",
		".autosome_PAR_ploidy_2.g.vcf.log",
		".bam.log",
		".chrX_nonPAR_ploidy_1.g.vcf.gz",
		".chrX_nonPAR_ploidy_1.g.vcf.gz.bcftools-stats",
		".chrX_nonPAR_ploidy_1.g.vcf.gz.bcftools-stats.log",
		".chrX_nonPAR_ploidy_1.g.vcf.gz.log",
		".chrX_nonPAR_ploidy_1.g.vcf.gz.tbi",
		".chrX_nonPAR_ploidy_1.g.vcf.gz.tbi.log",
		".chrX_nonPAR_ploidy_1.g.vcf.log",
		".chrX_nonPAR_ploidy_2.g.vcf.gz",
		".chrX_nonPAR_ploidy_2.g.vcf.gz.bcftools-stats",
		".chrX_nonPAR_ploidy_2.g.vcf.gz.bcftools-stats.log",
		".chrX_nonPAR_ploidy_2.g.vcf.gz.log",
		".chrX_nonPAR_ploidy_2.g.vcf.gz.tbi",
		".chrX_nonPAR_ploidy_2.g.vcf.gz.tbi.log",
		".chrX_nonPAR_ploidy_2.g.vcf.log",
		".chrY_nonPAR_ploidy_1.g.vcf.gz",
		".chrY_nonPAR_ploidy_1.g.vcf.gz.bcftools-stats",
		".chrY_nonPAR_ploidy_1.g.vcf.gz.bcftools-stats.log",
		".chrY_nonPAR_ploidy_1.g.vcf.gz.log",
		".chrY_nonPAR_ploidy_1.g.vcf.gz.tbi",
		".chrY_nonPAR_ploidy_1.g.vcf.gz.tbi.log",
		".chrY_nonPAR_ploidy_1.g.vcf.log",
		".cram",
		".cram.autosome_PAR_ploidy_2.wgs_metrics",
		".cram.autosome_PAR_ploidy_2.wgs_metrics.log",
		".cram.chrX_nonPAR_ploidy_1.wgs_metrics",
		".cram.chrX_nonPAR_ploidy_1.wgs_metrics.log",
		".cram.chrX_nonPAR_ploidy_2.wgs_metrics",
		".cram.chrX_nonPAR_ploidy_2.wgs_metrics.log",
		".cram.chrY_nonPAR_ploidy_1.wgs_metrics",
		".cram.chrY_nonPAR_ploidy_1.wgs_metrics.log",
		".cram.collect_base_dist_by_cycle",
		".cram.collect_base_dist_by_cycle.chart.pdf",
		".cram.collect_base_dist_by_cycle.chart.png",
		".cram.crai",
		".cram.crai.log",
		".cram.flagstat",
		".cram.idxstats",
		".cram.log",
		".log",
		".metrics.txt"} {
		// outputDirectoryPath/sampleId/sampleId.*
		// outputDirectoryPath/XX00000/XX00000.*
		targetFile := fn + "/" + sampleId + extension
		if _, err := os.Stat(targetFile); os.IsNotExist(err) {
			fmt.Printf("Missing file [%s]\n", targetFile)
			result = false
		} else {
			// ".log" file is not need to check filesize.
			// other files MUST be check filesize is not 0
			if !strings.HasSuffix(extension, ".log") {
				fileinfo, _ := os.Stat(targetFile)
				if fileinfo.Size() == 0 {
					fmt.Printf("File size is zero [%s]\n", targetFile)
					result = false
				}
			}
		}
	}
	return result
}

func getFileNameWithoutExtension(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func checkSecondaryFilesExists(fn string) (bool, error) {
	// true is exist all files
	// false is some secodary files missing
	result := true
	for _, extension := range []string{".amb", ".ann", ".bwt", ".pac", ".sa", ".alt", ".fai"} {
		// Check file is exist
		if _, err := os.Stat(fn + extension); os.IsNotExist(err) {
			fmt.Printf("Missing file [%s]\n", fn+extension)
			result = false
		}
	}
	// ^.dict
	dictfile := filepath.Join(filepath.Dir(fn), getFileNameWithoutExtension(fn)+".dict")
	if _, err := os.Stat(dictfile); os.IsNotExist(err) {
		fmt.Printf("Missing file [%s]\n", dictfile)
		result = false
	}
	return result, nil
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

/**
 * return value: true is fine, false is some thing wrong
 */
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
	byteBuf.WriteString("haplotypecaller_autosome_PAR_interval_bed:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerAutosomePARIntervalBed.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3584\n")
	byteBuf.WriteString("haplotypecaller_autosome_PAR_interval_list:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerAutosomePARIntervalList.Path))
	byteBuf.WriteString("haplotypecaller_chrX_nonPAR_interval_bed:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerChrXNonPARIntervalBed.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3584\n")
	byteBuf.WriteString("haplotypecaller_chrX_nonPAR_interval_list:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerChrXNonPARIntervalList.Path))
	byteBuf.WriteString("haplotypecaller_chrY_nonPAR_interval_bed:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerChrYNonPARIntervalBed.Path))
	byteBuf.WriteString("  format: http://edamontology.org/format_3584\n")
	byteBuf.WriteString("haplotypecaller_chrY_nonPAR_interval_list:\n")
	byteBuf.WriteString("  class: File\n")
	byteBuf.WriteString(fmt.Sprintf("  path: %s\n", rss.HaplotypecallerChrYNonPARIntervalList.Path))

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

func execCWL(outputDirectoryPath string, workflowFilePath string, sampleId string) string {
	// execute toil
	//p, _ := os.Getwd()
	// c1 := exec.Command("toil-cwl-runner", "--maxDisk", "248G", "--maxMemory", "64G", "--defaultMemory", "32000", "--defaultDisk", "32000", "--workDir", p, "--disableCaching", "--jobStore", "./"+sampleId+"-jobstore", "--outdir", "./"+sampleId, "--stats", "--cleanWorkDir", "never", "--batchSystem", "slurm", "--retryCount", "1", "--singularity", "--logFile", sampleId+".log", "per-sample/Workflows/per-sample.cwl", sampleId+"_jobfile.yaml")
	commandArgs := createToilCwlRunnerArguments(outputDirectoryPath, sampleId, workflowFilePath)
	c1 := exec.Command("toil-cwl-runner", commandArgs...)
	// set environment value if needed
	//c1.Env = append(os.Environ(), "TOIL_SLURM_ARGS=\"-w node[1-9]\"")
	//
	stdoutfile, _ := os.Create(outputDirectoryPath + "/toil-outputs/" + sampleId + "-stdout.txt")
	defer stdoutfile.Close()
	c1.Stdout = stdoutfile
	//
	stderrfile, _ := os.Create(outputDirectoryPath + "/toil-outputs/" + sampleId + "-stderr.txt")
	defer stderrfile.Close()
	c1.Stderr = stderrfile
	//
	c1.Start()
	c1.Wait()
	// output exitcode
	exitcodefile, _ := os.Create(outputDirectoryPath + "/toil-outputs/" + sampleId + "-exitcode.txt")
	defer exitcodefile.Close()
	exitCode := c1.ProcessState.ExitCode()
	exitcodefile.WriteString(fmt.Sprintf("%d\n", exitCode))
	//
	stdoutwriter := bufio.NewWriter(stdoutfile)
	defer stdoutwriter.Flush()
	//
	stderrwriter := bufio.NewWriter(stderrfile)
	defer stderrwriter.Flush()
	//
	defer func() {
		// display messages depending on exitCode
		if exitCode == 0 {
			fmt.Printf("SampleId: %s is successfully finished\n", sampleId)
		} else {
			stdoutfileabs, _ := filepath.Abs(stdoutfile.Name())
			stderrfileabs, _ := filepath.Abs(stderrfile.Name())

			fmt.Printf("SampleId: %s is fail. exitcode = %d\n", sampleId, exitCode)
			fmt.Println("  See stdout: ", stdoutfileabs)
			fmt.Println("  See stderr: ", stderrfileabs)
		}
	}()
	//
	return ""
}

func getCurrentTime() time.Time {
	return time.Now()
}

func createJobStoreDir(outputDirectoryPath string, sampleId string, currentTime time.Time) string {
	return outputDirectoryPath + "/jobstores/" + sampleId + "-jobstore-" + currentTime.Format("20060102150405")
}
func createLogFilePath(outputDirectoryPath string, sampleId string, currentTime time.Time) string {
	return outputDirectoryPath + "/logs/" + sampleId + "-" + currentTime.Format("20060102150405") + ".log"
}

func createToilCwlRunnerArguments(outputDirectoryPath string, sampleId string, workflowFilePath string) []string {
	currentTime := getCurrentTime()
	jobStoreDir := createJobStoreDir(outputDirectoryPath, sampleId, currentTime)
	logFilePath := createLogFilePath(outputDirectoryPath, sampleId, currentTime)
	commandArgs := []string{"--maxDisk", "248G", "--maxMemory", "64G", "--defaultMemory", "32000", "--defaultDisk", "32000", "--disableCaching", "--jobStore", jobStoreDir, "--outdir", outputDirectoryPath + "/" + sampleId, "--stats", "--batchSystem", "slurm", "--retryCount", "1", "--singularity", "--logFile", logFilePath, workflowFilePath, sampleId + "_jobfile.yaml"}
	return commandArgs
}

var dryrunFlag bool
var helpFlag bool
var versionFlag bool
var fileExistsCheckFlag bool
var fileHashCheckFlag bool

func main() {
	cmd.Execute()
}

func main2() {
	flag.BoolVarP(&dryrunFlag, "dry-run", "n", false, "Dry-run, do not execute acutal command")
	flag.BoolVarP(&helpFlag, "help", "h", false, "Show help message")
	flag.BoolVarP(&versionFlag, "version", "v", false, "Show version")
	flag.BoolVarP(&fileExistsCheckFlag, "file-exists-check", "", true, "Check file exists")
	flag.BoolVarP(&fileHashCheckFlag, "file-hash-check", "", true, "Check file hash value")
	flag.Parse()

	if helpFlag {
		flag.PrintDefaults()
		return
	}
	if versionFlag {
		//utils.displayVersionString()
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

	// validate
	checkResult := false
	for _, s := range ss.SampleList {
		//fmt.Printf("Check index: %d, SampleId: %s\n", i, s.SampleId)
		for j, t := range s.RunList {
			r1, _ := checkRunData(&t.RunData)
			checkResult = checkResult || r1
			if !r1 {
				fmt.Println("Some error found. Not exist or Hash value error")
				fmt.Printf("Check index: %d, RunId: %s\n", j, t.RunId)
				fmt.Printf("pe or se: [%s]\n", t.RunData.PEOrSE)
				fmt.Printf("fq1: [%s]\n", t.RunData.FQ1)
				fmt.Printf("fq2: [%s]\n", t.RunData.FQ2)
				fmt.Printf("result=%t\n", r1)
			}
		}
	}
	if !checkResult {
		fmt.Println("some thing wrong. do not execute")
		return
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

	// if dryrunFlag {
	// 	fmt.Println("Dry-run flag is set")
	// 	return
	// }
	//
	secondaryFilesCheck, err := checkSecondaryFilesExists(rss.Reference.Path)
	if !secondaryFilesCheck {
		fmt.Println("Some secondary file is missing")
		return
	}

	// Set output directory path
	workflowFilePath := rss.WorkflowFile.Path

	// currently check local filesystem only
	if !utils.IsExistsWorkflowFile(workflowFilePath) {
		fmt.Printf("Missing workflow file [%s]\n", workflowFilePath)
		os.Exit(1)
	}

	// Set output directory path
	outputDirectoryPath := rss.OutputDirectory.Path
	// Create output directory
	// if not create , show error message and exit
	if err := os.MkdirAll(outputDirectoryPath, 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create output directory")
		return
	}
	if err := os.MkdirAll(outputDirectoryPath+"/toil-outputs", 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create toil outputs directory")
		return
	}
	if err := os.MkdirAll(outputDirectoryPath+"/logs", 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create logs directory")
		return
	}
	if err := os.MkdirAll(outputDirectoryPath+"/jobstores", 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create jobstores directory")
		return
	}
	// copy sample_sheet file
	original_sample_sheet, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer original_sample_sheet.Close()
	copied_sample_sheet, err := os.Create(outputDirectoryPath + "/" + original_sample_sheet.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer copied_sample_sheet.Close()
	_, err = io.Copy(copied_sample_sheet, original_sample_sheet)
	if err != nil {
		log.Fatal(err)
	}
	// copy config file
	original_configfile, err := os.Open(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	defer original_configfile.Close()
	copied_configfile, err := os.Create(outputDirectoryPath + "/" + original_configfile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer copied_configfile.Close()
	_, err = io.Copy(copied_configfile, original_configfile)
	if err != nil {
		log.Fatal(err)
	}

	// create job file for CWL
	createJobFile(&ss, &rss)

	// check toil-cwl-runner is exists or not
	foundToilCWLRunner := utils.IsExistsToilCWLRunner()

	// dry-run

	// exec and wait
	var eg errgroup.Group
	executeCount := 0
	for i, s := range ss.SampleList {
		isExecute := false
		// Check SampleId result directory is exist
		if _, err := os.Stat(outputDirectoryPath + "/" + s.SampleId); os.IsNotExist(err) {
			// SampleId result directory is missing
			// so this id must be executed
			isExecute = true
		} else {
			// check all result file is found or not
			// SampleId prefix files check
			check1 := isExistsAllResultFilesPrefixSampleId(outputDirectoryPath, s.SampleId)
			if !check1 {
				isExecute = true
			}
			// RunID prefix files check
			for _, r := range s.RunList {
				check2 := isExistsAllResultFilesPrefixRunId(outputDirectoryPath+"/"+s.SampleId, r.RunId)
				if !check2 {
					isExecute = true
				}
			}
			//fmt.Printf("index: %d, SampleId: %s will be Execute new.\n", i, s.SampleId)
		}
		if isExecute {
			executeCount += 1

			fmt.Printf("index: %d, SampleId: %s will be Execute new.\n", i, s.SampleId)
			sampleId := s.SampleId
			if !dryrunFlag {
				if foundToilCWLRunner {
					// only exec when toil-cwl-runner is found
					eg.Go(func() error {
						execCWL(outputDirectoryPath, workflowFilePath, sampleId)
						return nil
					})
				}
			}
		}
	}
	if dryrunFlag {
		fmt.Printf("[%d/%d] task will be executed.\n", executeCount, len(ss.SampleList))
	}

	// for i, s := range ss.SampleList {
	// 	fmt.Printf("index: %d, SampleId: %s\n", i, s.SampleId)
	// 	sampleId := s.SampleId
	// 	eg.Go(func() error {
	// 		// TODO check exit status
	// 		execCWL(outputDirectoryPath, workflowFilePath, sampleId)
	// 		return nil
	// 	})
	// }

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	if foundToilCWLRunner {
		fmt.Println("Can not find toil-cwl-runner")
	}

	fmt.Println("fin")
}
