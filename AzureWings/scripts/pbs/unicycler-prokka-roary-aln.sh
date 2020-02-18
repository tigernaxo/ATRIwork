#!/bin/bash
currentDir=$(dirname $(readlink -f ${0}))
env="${currentDir}/../../bin/env.sh"
. ${env} ${env}

unicyclerScript="${AWScriptsPbs}/unicycler-prokka.sh"
roaryScript="${AWScriptsPbs}/roary-aln.sh"
default_queue="ngs48G"

main(){
  getMyOpt "$@"
  mkdir -p ${outDir}
  unicyclerJobs=$(${unicyclerScript} -q ${queue} -d ${outDir}/assembly -l ${pairedlist})

  # get gffs and append to roaryJob arguments
  declare -a gffArr
  while IFS=$'\t' read r1 r2
  do
    id=$(basename ${r1}|sed 's/[\._-].*$//g' )
    gff=${outDir}/assembly/${id}/prokka/${id}.gff
    gffArr+=(${gff})
  done < ${pairedlist}

  roaryJob=$(${roaryScript} \
    -d ${outDir} \
    -q ${queue} -w ${unicyclerJobs} \
    ${gffArr[@]}) 

  echo "${roaryJob}"
}

getMyOpt(){
  if [ $# -eq 0 ] ;then
    usage
    exit 0
  fi

# get outdir form getopts
  while getopts ":q:d:l:h" opt; do
    case ${opt} in
      h )
        usage
        exit 0 ;;
      d )
        outDir=${OPTARG} ;;
      l )
        pairedlist=${OPTARG} ;;
      q )
        queue=${OPTARG} ;;
      \? )
        echo "Invalid Option: -${OPTARG}" 1>&2
        exit 1 ;;
    esac
  done

  logFile="${outDir}/assemble.log"
  [ ! -n "${queue}" ] && queue="${default_queue}" && echo "queue not specified, use ngs48G" >> ${logFile};
  [ ! -n "${outDir}" ] || [ ! -n "${pairedlist}"  ] && echo "out_dir and paired_list not set, exit." >> ${logFile} && exit 1 ;
}
usage(){
  echo 
  echo "Usage:"
  echo
  echo "  $(basename ${0}) [-h] [-q queue_name] [-d outdir] [-l paired_list.tsv]" 
  echo 
  echo "Option:"
  echo "  -h    Show this message."
  echo "  -l    paried_list.tsv"
  echo "  -q    Specify queue."
  echo "  -d    output dir"
  echo 
  echo "Example:"
  echo "  $(basename ${0}) -q ngsTest -d $(pwd) \\" 
  echo "    -l raw_reads.tsv "
}
main "$@"; exit
