/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/manabuishiii/jgaworkflowspecchecker/utils"
	"github.com/spf13/cobra"
)

//
var pullSingularityImages bool
var pullDockerImages bool

// pullContainerImagesCmd represents the pullContainerImages command
var pullContainerImagesCmd = &cobra.Command{
	Use:   "pull-container-images",
	Short: "Pull container images, require --singularity or --docker",
	Long: `Pull container images
Please specify --sigularity or --docker.

Cached container images are save in container_cache_directory.
singularity image has suffix .sif
docker image has suffix .tar
`,
	Run: func(cmd *cobra.Command, args []string) {
		if pullDockerImages == false && pullSingularityImages == false {
			fmt.Println("One of --docker or --singularity is required")
			os.Exit(1)
		}
		if pullDockerImages == true && pullSingularityImages == true {
			fmt.Println("One of --docker or --singularity is required")
			os.Exit(1)
		}
		if pullDockerImages {
			if !utils.IsExistsDocker() {
				fmt.Println("docker is not found.")
				os.Exit(1)
			}
		}
		if pullSingularityImages {
			if !utils.IsExistsSingularity() {
				fmt.Println("singularity is not found.")
			}
		}
		pullResult := execPullContainerImages(args)
		if !pullResult {
			fmt.Println("Some error happens at pull images.")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(pullContainerImagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullContainerImagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullContainerImagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pullContainerImagesCmd.Flags().BoolVarP(&pullSingularityImages, "singularity", "", false, "Save Singularity images")
	pullContainerImagesCmd.Flags().BoolVarP(&pullDockerImages, "docker", "", false, "Save as Docker images")
}

func cacheDockerImages(cwldir string, container_cache_directory string, scriptCode string) bool {
	c1 := exec.Command("/bin/bash")
	scriptEnv := append(os.Environ(), "CWLDIR="+cwldir)
	if pullDockerImages {
		scriptEnv = append(scriptEnv, "CWL_DOCKER_CACHE="+container_cache_directory)
	}
	if pullSingularityImages {
		scriptEnv = append(scriptEnv, "CWL_SINGULARITY_CACHE="+container_cache_directory)
	}

	c1.Env = scriptEnv
	stdin, _ := c1.StdinPipe()
	io.WriteString(stdin, scriptCode)
	stdin.Close()
	output, _ := c1.CombinedOutput()

	c1.Start()
	c1.Wait()
	fmt.Println(string(output))
	exitCode := c1.ProcessState.ExitCode()
	fmt.Printf("Pull ExitCode is [%d]\n", exitCode)
	return exitCode == 0
}

func execPullContainerImages(args []string) bool {
	loadSampleSheetAndConfigFile(args)
	// check in sample sheet data
	// if !checkSampleSheet(&ss) {
	// 	return
	// }
	// check in config data
	if !checkConfigFile(&rss) {
		return false
	}
	if !utils.IsExistsWorkflowFile(rss.WorkflowFile.Path) {
		fmt.Printf("Workflow file [%s] is missing.\n", rss.WorkflowFile.Path)
		fmt.Println("Stop pull container images")
		return false
	}
	// assueme workflow file is
	//  jga-analysis/per-sample/Workflows/per-sample.cwl
	// become
	//  jga-analysis/per-sample
	cwldir := filepath.Dir(filepath.Dir(rss.WorkflowFile.Path))
	// cache create this directory
	//  docker image cache has suffix ".tar"
	//  singularity image cache has suffix ".sif"
	container_cache_directory := rss.ContainerCacheDirectory.Path
	fmt.Println(container_cache_directory)
	//
	result := false
	if pullDockerImages {
		result = cacheDockerImages(cwldir, container_cache_directory, string(createDockerImageScript))
	}
	if pullSingularityImages {
		result = cacheDockerImages(cwldir, container_cache_directory, string(createSingularityImageScript))
	}
	return result
}
