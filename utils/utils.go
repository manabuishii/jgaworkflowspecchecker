package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

func DisplayJobManagerRecoginition(workflowFilePath string) {
	//fmt.Printf("Workflow file is exists [%t]\n", isExistsWorkflowFile(workflowFilePath))
	fmt.Printf("toil-cwl-runner is exists [%t]\n", IsExistsToilCWLRunner())
	fmt.Printf("Using Virtualenv if true set TOIL_CHECK_ENV=True [%t]\n", IsInVirtualenv())

	fmt.Printf("  Using Python virtualenv [%t]\n", IsInPythonVirtualenv())
	fmt.Printf("  Using Conda virtual env [%t]\n", IsInCondaEnv())
	fmt.Printf("sbatch(slurm) is exists [%t]\n", IsExistsSbatch())
	fmt.Printf("singularity is exists [%t]\n", IsExistsSingularity())
}
