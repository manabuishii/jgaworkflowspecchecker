package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

type SimpleSchema struct {
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

type ReferenceSchema struct {
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

func outputReference(rss *ReferenceSchema) (string, error) {
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

/*
 * Check files inside config file are exists.
 * Return value:
 *   true: all files are exists
 *   false: some files are missing
 */
func CheckOutputReference(rss *ReferenceSchema) bool {
	result := true
	// reference file
	if !IsExistsFile(rss.Reference.Path) {
		fmt.Printf("Referenece file [%s] is missing\n", rss.Reference.Path)
		result = false
	}
	// dbsnp
	if !IsExistsFile(rss.Dbsnp.Path) {
		fmt.Printf("dbsnp file [%s] is missing\n", rss.Dbsnp.Path)
		result = false
	}
	// mills
	if !IsExistsFile(rss.Mills.Path) {
		fmt.Printf("mills file [%s] is missing\n", rss.Mills.Path)
		result = false
	}
	// known_indels
	if !IsExistsFile(rss.KnownIndels.Path) {
		fmt.Printf("known_indels file [%s] is missing\n", rss.KnownIndels.Path)
		result = false
	}
	// haplotypecaller_autosome_PAR_interval_bed
	if !IsExistsFile(rss.HaplotypecallerAutosomePARIntervalBed.Path) {
		fmt.Printf("haplotypecaller_autosome_PAR_interval_bed file [%s] is missing\n", rss.HaplotypecallerAutosomePARIntervalBed.Path)
		result = false
	}
	// haplotypecaller_autosome_PAR_interval_list
	if !IsExistsFile(rss.HaplotypecallerAutosomePARIntervalList.Path) {
		fmt.Printf("haplotypecaller_autosome_PAR_interval_list file [%s] is missing\n", rss.HaplotypecallerAutosomePARIntervalList.Path)
		result = false
	}
	// haplotypecaller_chrX_nonPAR_interval_bed
	if !IsExistsFile(rss.HaplotypecallerChrXNonPARIntervalBed.Path) {
		fmt.Printf("haplotypecaller_chrX_nonPAR_interval_bed file [%s] is missing\n", rss.HaplotypecallerChrXNonPARIntervalBed.Path)
		result = false
	}
	// haplotypecaller_chrX_nonPAR_interval_list
	if !IsExistsFile(rss.HaplotypecallerChrXNonPARIntervalList.Path) {
		fmt.Printf("haplotypecaller_chrX_nonPAR_interval_list file [%s] is missing\n", rss.HaplotypecallerChrXNonPARIntervalList.Path)
		result = false
	}
	// haplotypecaller_chrY_nonPAR_interval_bed
	if !IsExistsFile(rss.HaplotypecallerChrYNonPARIntervalBed.Path) {
		fmt.Printf("haplotypecaller_chrY_nonPAR_interval_bed file [%s] is missing\n", rss.HaplotypecallerChrYNonPARIntervalBed.Path)
		result = false
	}
	// haplotypecaller_chrY_nonPAR_interval_list
	if !IsExistsFile(rss.HaplotypecallerChrYNonPARIntervalList.Path) {
		fmt.Printf("haplotypecaller_chrY_nonPAR_interval_list file [%s] is missing\n", rss.HaplotypecallerChrYNonPARIntervalList.Path)
		result = false
	}

	return result
}

/*
 * Check and display files for workflow execution .
 * Return value:
 *   true: all files are exists
 *   false: some files are missing
 */
func CheckAndDisplayFilesForExecute(rss *ReferenceSchema) bool {
	result := true
	// Workflow file
	if !IsExistsWorkflowFile(rss.WorkflowFile.Path) {
		fmt.Printf("workflow file [%s] is missing\n", rss.WorkflowFile.Path)
		result = false
	}
	// Secondary files
	secondaryFilesFlag, _ := CheckSecondaryFilesExists(rss.Reference.Path)
	if !secondaryFilesFlag {
		fmt.Printf("Some secondary files are missing\n")
		result = false
	}
	//
	refFilesFlag := CheckOutputReference(rss)
	if !refFilesFlag {
		fmt.Printf("Some files are missing\n")
		result = false
	}

	return result
}

// call per sample
func outputJobFile(s *Sample, rss *ReferenceSchema) (string, error) {
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

func CreateJobFile(ss *SimpleSchema, rss *ReferenceSchema) error {
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

func BuildVersionString(version, revision, date string) string {
	result := fmt.Sprintf("Version: %s-%s (built at %s)\n", version, revision, date)
	return result
}

func IsExistsToilCWLRunner() bool {
	_, err := exec.LookPath("toil-cwl-runner")
	return err == nil
}

func IsExistsSbatch() bool {
	_, err := exec.LookPath("sbatch")
	return err == nil
}

func IsExistsSingularity() bool {
	_, err := exec.LookPath("singularity")
	return err == nil
}

func IsInVirtualenv() bool {
	result := false
	result = result || IsInPythonVirtualenv()
	result = result || IsInCondaEnv()
	return result
}

func IsInCondaEnv() bool {
	result := false
	condavenv := os.Getenv("CONDA_DEFAULT_ENV")
	if condavenv != "" {
		result = true
	}
	return result
}

func IsInPythonVirtualenv() bool {
	result := false
	venv := os.Getenv("VIRTUAL_ENV")
	if venv != "" {
		result = true
	}
	return result
}

/*
  Check whether workflow path is exists .
  MEMO: Path starts http:// or https:// , do not check.
*/
func IsExistsWorkflowFile(workflowFilePath string) bool {
	if !strings.HasPrefix(workflowFilePath, "http://") {
		if !strings.HasPrefix(workflowFilePath, "https://") {
			if _, err := os.Stat(workflowFilePath); os.IsNotExist(err) {
				return false
			}
		}
	}
	return true
}

func DisplayJobManagerRecoginition(rss *ReferenceSchema) {
	//fmt.Printf("Workflow file is exists [%t]\n", isExistsWorkflowFile(workflowFilePath))
	fmt.Printf("toil-cwl-runner is exists [%t]\n", IsExistsToilCWLRunner())
	fmt.Printf("Using Virtualenv if true set TOIL_CHECK_ENV=True [%t]\n", IsInVirtualenv())

	fmt.Printf("  Using Python virtualenv [%t]\n", IsInPythonVirtualenv())
	fmt.Printf("  Using Conda virtual env [%t]\n", IsInCondaEnv())
	fmt.Printf("sbatch(slurm) is exists [%t]\n", IsExistsSbatch())
	fmt.Printf("singularity is exists [%t]\n", IsExistsSingularity())
	result := CheckAndDisplayFilesForExecute(rss)
	if result {
		fmt.Println("All files for workflow Execution are found.")
	} else {
		fmt.Println("Some files for workflow Execution are missing.")
	}
}

//

func Md5File(filePath string) (string, error) {
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

/**
 * return value: true is fine, false is some thing wrong
 */
func CheckRunData(runData *RunData, fileExistsCheckFlag bool, fileHashCheckFlag bool) (bool, error) {
	result := false
	if runData.PEOrSE == "PE" {
		r1, _ := CheckRunDataFile(runData.FQ1, runData.FQ1_MD5, fileExistsCheckFlag, fileHashCheckFlag)
		r2, _ := CheckRunDataFile(runData.FQ2, runData.FQ2_MD5, fileExistsCheckFlag, fileHashCheckFlag)
		result = r1 && r2
	} else {
		result, _ = CheckRunDataFile(runData.FQ1, runData.FQ1_MD5, fileExistsCheckFlag, fileHashCheckFlag)
	}
	return result, nil
}

/*
 Check file is exists
 Return value:
  true  is exists
  false is not found
*/
func IsExistsFile(fn string) bool {
	// Check file is exist
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return false
	}
	return true
}
func CheckRunDataFile(fn string, fnmd5 string, fileExistsCheckFlag bool, fileHashCheckFlag bool) (bool, error) {
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
		md5, _ := Md5File(fn)
		if fnmd5 != md5 {
			result = false
			fmt.Printf("expected: [%s]\n", fnmd5)
			fmt.Printf("actual  : [%s]\n", md5)
			fmt.Println("md5 is not match")
		}
	}
	return result, nil
}

func IsExistsAllResultFilesPrefixRunId(outputDirectoryPath string, runId string) bool {
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
func IsExistsAllResultFilesPrefixSampleId(outputDirectoryPath string, sampleId string) bool {
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

func CheckSecondaryFilesExists(fn string) (bool, error) {
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

func ExecCWL(outputDirectoryPath string, workflowFilePath string, sampleId string) string {
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
