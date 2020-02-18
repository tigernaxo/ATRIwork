#!/bin/bash
bash="/bin/bash"

currentDir=$( dirname $( readlink -f ${0}  ))/
conda=$(readlink -f "${currentDir}/../../conda-envs/bin/conda")
condaActivate=$(readlink -f "${currentDir}/../../conda-envs/bin/activate")
condaDeactivate=$(readlink -f "${currentDir}/../../conda-envs/bin/deactivate")
env=$(readlink -f "${currentDir}/../../conda-envs/envs/unicycler")
timeStamp=$(readlink -f "${currentDir}/../util/timeStamp.sh")
queue_ncpus=`${currentDir}/../util/queue_ncpus.sh`
thisScript=$(basename ${0})
PATH=${PATH}:${PBS_O_PATH}

#### Test Start ####
#### Test end ####

######## Set Arguments ########
###############################
while getopts ":d:f:r:h" opt; do
  case ${opt} in
    h )
      exit 0 ;;
    f )
      R1=${OPTARG} ;;
    r )
      R2=${OPTARG} ;;
    d )
      outputFolder=${OPTARG} ;;
    \? )
      echo "Invaiid Option: -${OPTARG}" 1>&2
      exit 1 ;;
  esac
done

#[ ! "${R1}" ] || [ ! "${R2}" ] || [ ! "${outputFolder}" ] && echo "need -1 r1 -2 r2 -d outdir" && exit 1;

### Log: current script start
echo "`${timeStamp}` [${thisScript}]start:"
qstat -f ${PBS_JOBID}
echo "`${timeStamp}` [${thisScript}]Activating conda env: $(basename ${env})"
. ${condaActivate} ${env}

#### Script command start ####
##############################

echo "`${timeStamp}` [${thisScript}]raw-reads:"
echo -e "\t${R1}"
echo -e "\t${R2}"

cmd="unicycler -t ${queue_ncpus} \
  --min_kmer_frac 0.5 --max_kmer_frac 0.99 --kmer_count 10 \
  --keep 0 --mode bold \
  -1 ${R1} \
  -2 ${R2} \
  -o ${outputFolder};
  quast -t ${queue_ncpus} -m100 -o ${outputFolder}/quast ${outputFolder}/assembly.fasta;
  cp ${outputFolder}/quast/quast.log  ${outputFolder};
  cp ${outputFolder}/quast/report.tsv ${outputFolder};
  rm -rf ${outputFolder}/quast;"

##############################
#### Script command end ####

### Log: current cmd start
echo -e `${timeStamp}` "[${thisScript}]executing command:\n\t${cmd}"
`bash -c "$cmd"`

#### deactivate conda environment
. ${condaDeactivate}
echo "`${timeStamp}` [${thisScript}]end."
