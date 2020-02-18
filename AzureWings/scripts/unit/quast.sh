#!/bin/bash
workdir=${1}
assembly=${2}

currentDir=$( dirname $( readlink -f ${0}  ))/
conda="${currentDir}/../../conda-envs/bin/conda"
condaActivate="${currentDir}/../../conda-envs/bin/activate"
condaDeactivate="${currentDir}/../../conda-envs/bin/deactivate"
env="${currentDir}/../../conda-envs/envs/unicycler"
timeStamp="${currentDir}/../util/timeStamp.sh"
queue_ncpus=`${currentDir}/../util/queue_ncpus.sh`
thisScript=$(basename ${0})
PATH=${PATH}:${PBS_O_PATH}

#### Test Start ####
#### Test end ####

### Log: current script start
echo "`${timeStamp}` [${thisScript}]start:"
qstat -f ${PBS_JOBID}
echo `${timeStamp}` "[${thisScript}]Activating conda env: $(basename ${env})"
. ${condaActivate} ${env}

#### Script command start ####

cmd="quast -t ${queue_ncpus} ${2}"

#### Script command end ####

### Log: current cmd start
echo -e `${timeStamp}` "[${thisScript}]executing command:\n\t${cmd}"
`bash -c \"${cmd}\"`

#### deactivate conda environment
. ${condaDeactivate}
echo "`${timeStamp}` [${thisScript}]end."
