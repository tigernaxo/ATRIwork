#!/bin/bash
outputFolder=${1}

currentDir=$( dirname $( readlink -f ${0}  ))/
conda="${currentDir}/../../conda-envs/bin/conda"
condaActivate="${currentDir}/../../conda-envs/bin/activate"
condaDeactivate="${currentDir}/../../conda-envs/bin/deactivate"
env="${currentDir}/../../conda-envs/envs/roary"
timeStamp="${currentDir}/../util/timeStamp.sh"
queue_ncpus=`${currentDir}/../util/queue_ncpus.sh`
thisScript=$(basename ${0})
PATH=${PATH}:${PBS_O_PATH}

#### Test Start ####
#### Test end ####

### Log: current script start
echo "`${timeStamp}` [${thisScript}]start:"
qstat -f ${PBS_JOBID}
echo "`${timeStamp}` [${thisScript}]Activating conda env: $(basename ${env})"
. ${condaActivate} ${env}

#### Script command start ####
##############################


# if outputdir exist, roary will append a _number after out path
# so we need to make sure directionary is removed before roary start
cmd="${cmd}
  rm -rf ${outputFolder};
  roary -e --mafft -f ${outputFolder} -p ${queue_ncpus} ${@:2};
  snp-sites ${outputFolder}/core_gene_alignment.aln -o ${outputFolder}/cgSNP.fasta;
  mkdir -p ${outputFolder};"
    

##############################
#### Script command end ####

### Log: current cmd start
echo -e `${timeStamp}` "[${thisScript}]executing command:\n\t${cmd}"
echo ${cmd} >> /work1/u00ycc00/teseAW-unicycler-roary/log
`bash -c "$cmd"`

#### deactivate conda environment
. ${condaDeactivate}
echo "`${timeStamp}` [${thisScript}]end."
