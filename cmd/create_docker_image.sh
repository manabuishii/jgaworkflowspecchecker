#!/bin/bash
set -eu
# This script MUST be set 2 environment value
#  CWL_DOCKER_CACHE: docker cache save directory
#  CWLDIR: expects jga-analysis/per-sample/ directory
#   - JobManager currently resolve this value as following
#     - remove CWL file and 1 directory
#     - if "workflow_file" in config file has "per-sample/Workflows/per-sample.cwl" , 
#     - CWLDIR is set "per-sample".
mkdir -p ${CWL_DOCKER_CACHE}
RET=0
for DOCKERIMAGE in `grep -r dockerPull ${CWLDIR}  | awk '{print $NF}' | tr -d '"' | tr -d "'" | sort | uniq `
do
 DOCKER_IMAGE_FILE=`echo $DOCKERIMAGE| sed -e "s/\///g"`.tar
 docker pull ${DOCKERIMAGE}
 PULLRET=$?
 docker save -o ${CWL_DOCKER_CACHE}/${DOCKER_IMAGE_FILE} ${DOCKERIMAGE}
 SAVERET=$?
 #
 if [ ${PULLRET} -ne 0 ]; then
    echo "ERROR at pull [${DOCKERIMAGE}]"
 fi
 if [ ${SAVERET} -ne 0 ]; then
    echo "ERROR at save [${DOCKERIMAGE}]"
 fi
 RET=$((${RET}+${PULLRET}+${SAVERET}))
done
echo "EXIT STATUS[${RET}]"
exit ${RET}