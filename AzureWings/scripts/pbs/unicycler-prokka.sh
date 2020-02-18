#!/bin/bash

##############################################
# set pbs environment
currentDir=$( dirname $( readlink -f ${0}  ))
env=$(readlink -f ${currentDir}/../../bin/env.sh)
. ${env} ${env}

unicycler="${AWScriptsUnit}/unicycler.sh"
prokka="${AWScriptsUnit}/prokka.sh"
quast="${AWScriptsUnit}/quast.sh"

main(){
  getMyOpts "$@"

  logFile="${workDir}/unicycler-prokka.log"
  declare -a ProkkaIdArr
  while IFS=$'\t' read r1 r2
  do
    id=$(basename ${r1}|sed 's/[\._-].*$//g' )
    mkdir -p ${workDir}/${id}/unicycler
    mkdir -p ${workDir}/${id}/prokka
  
    r1=${dataPool}${r1}
    r2=${dataPool}${r2}
  
    if [ -f ${workDir}/${id}/unicycler/assembly.fasta ] ; then
  
      echo $(${timeStamp})" skip ${id} assemble, because assembly already exists:" >> ${logFile}
      echo "  ${workDir}/${id}/unicycler/assembly.fasta" >> ${logFile}
  
      if [ ! -f ${workDir}/${id}/prokka/prokka.gff  ] ; then
        prokkaId=$(submitProkka ${id})
        echo $( ${timeStamp} )" start prokka annotation: ${prokkaId}"  >> ${logFile}
        ProkkaIdArr+=(${prokkaId})
        continue
      fi
  
      echo "skip ${id}.fasta annotation, because gff already exists:" >> ${logFile}
      echo "  ${workDir}/${id}/prokka/prokka.gff" >> ${logFile}
      continue;
  
    fi
  
    echo $(${timeStamp})" Start assemble: ${id}." >> ${logFile}
    #    1. submit assemble job, and get id
    asmId=$(submitUnicycler ${id} ${r1} ${r2})
    echo $( ${timeStamp} )" submit unicycler asasembly: ${asmId}" >> ${logFile}
  
    #    2. defer annotation job til assemble is  done
    #       -W depend=afterok:${asmId}
    prokkaId=$(submitProkka ${id} -W depend=afterok:${asmId})
    echo $( ${timeStamp} )" hold prokka annotation: ${prokkaId}."  >> ${logFile}

    # echo ${prokkaId}
    ProkkaIdArr+=(${prokkaId})
    #echo ${ProkkaIdArr[@]}
  
  done  < ${reads}
  echo ${ProkkaIdArr[@]} | sed 's/ /:/g'
}


##############################################
# set arguments
getMyOpts(){
  if [ $# -eq 0 ] ;then
    usage
    exit 0
  fi

  while getopts ":l:d:q:f:w:h" opt; do
    case ${opt} in
      h )
        usage
        exit 0 ;;
      l )
        reads=${OPTARG} ;;
      d )
        workDir=${OPTARG} ;;
      f )
        dataPool=${OPTARG} ;;
      q )
        queue=${OPTARG} ;;
      w )
        afterokString="-W depend=afterok:${OPTARG}" ;;
      \? )
        echo "Invalid Option: -${OPTARG}" 1>&2
        exit 1 ;;
    esac
  done
  
  [ ! "${reads}" ] || [ ! "${workDir}"  ] && echo "You must set -l reads_list or -d work_dir" && exit 1
  [ ! "${queue}"  ] && queue="ngs16G"
}

##############################################
# set derived variable

# $1=id $2=other pbs option
submitProkka(){
  qsub ${@:2} \
    $(${pbsSnippet} ${queue} ${1}-prokka ${workDir}/${1}/prokka) \
    -- ${prokka} \
      -n ${1} \
      -d ${workDir}/${1}/prokka \
      -i ${workDir}/${1}/unicycler/assembly.fasta 
}

submitUnicycler(){
  qsub \
    $(${pbsSnippet} ${queue} ${1}-unicycler ${workDir}/${1}/unicycler) \
    -- ${unicycler} \
     -d ${workDir}/${1}/unicycler \
     -f ${2} -r ${3}
}

###########################################
# give argument hint if no argument gived.
usage() {
  echo 
  echo "Usage:"
  echo "  $(basename ${0}) -l reads_list -d outdir [-h] [-q queue_name] [-w afterok_string]"
  echo 
  echo "Option:"
  echo "  -l    Tab delimited reads pair."
  echo "  -d    output dir."
  echo "  -q    Specify quene(default: ngs16G)."
  echo "  -w    Wait the job on the list all done before start. (ex job1:job2)"
  echo "  -h    Show this message."
  echo
  echo "Example: "
  echo "  $(basename ${0}) -l /work1/u00ycc00/test/reads\\" 
  echo "    -d /home/naxo/work \\"
  echo "    -q ngs48G \\"
  echo
  echo "This Program recept meta file contain paired reads delimited by tab each line. "
  echo "Output is prokka job id delimited by \":\"."
  echo ""
  echo "#/work1/u00ycc00/test/reads looks like:"
  echo "/pathTo/A39W18507_R1.tgz /pathTo/A39W18507_R2.tgz"
  echo "/pathTo/A39W18508_R1.tgz /pathTo/A39W18508_R2.tgz"
  echo "/pathTo/A39W18509_R1.tgz /pathTo/A39W18509_R2.tgz"
  echo "/pathTo/A39W18510_R1.tgz /pathTo/A39W18510_R2.tgz"
  echo
}

main "$@" ; exit 0
