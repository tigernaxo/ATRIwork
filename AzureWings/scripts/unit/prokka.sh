#!/bin/bash
currentDir=$( dirname $( readlink -f ${0}  ))/
conda="${currentDir}/../../conda-envs/bin/conda"
condaActivate="${currentDir}/../../conda-envs/bin/activate"
condaDeactivate="${currentDir}/../../conda-envs/bin/deactivate"
env="${currentDir}/../../conda-envs/envs/prokka"
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

######## Set argument ########
##############################
while getopts ":n:d:i:h" opt; do
  case ${opt} in
    h )
      exit 0 ;;
    d )
      workdir=${OPTARG} ;;
    i )
      assembly=${OPTARG} ;;
    n )
      gffName=${OPTARG} ;;
    \? )
      echo "Invaiid Option: -${OPTARG}" 1>&2
      exit 1 ;;
  esac
done

#### Script command start ####
##############################

cmd="prokka --cpus=0 \
  --force --outdir=${workdir} \
  --prefix=prokka ${assembly};
  cd ${workdir};
  ls ${workdir} | grep prokka | grep -v \"gff\\|txt\" | xargs rm;
  mv ${workdir}/prokka.gff ${workdir}/${gffName}.gff"

##############################
#### Script command end ####

### Log: current cmd start
echo -e `${timeStamp}` "[prokka.sh]executing command:\n\t${cmd}"
`bash -c "${cmd}"`

#### deactivate conda environment
. ${condaDeactivate}
echo "`${timeStamp}` [${thisScript}]end."
