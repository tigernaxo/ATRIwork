#!/bin/bash
currentDir=$(dirname $(readlink -f ${0}))
env="${currentDir}/../../bin/env.sh"
. ${env} ${env}

seedx="12345"
seedp="12345"
suffix="out"

main(){

  getMyOpt "$@"
  echo "$(${timeStamp}) [$(basename ${0})]start"
  . ${AWCondaActivae} ${AWConda}/envs/raxml

  name=$(echo $(basename ${alignment}) | sed 's/\..*$//')
  cd ${workDir}
  raxmlHPC-PTHREADS-SSE3 -f a -T $(${queue_ncpus}) \
    -x ${seedx} -p ${seedp} \
    -# 100 -m GTRGAMMA \
    -s ${alignment} -n ${suffix}

  . ${AWCondaDeactivae}
}

getMyOpt(){
  while getopts ":d:i:h" opt; do
    case ${opt} in
      d )
        workDir=${OPTARG} ;;
      i )
        alignment=${OPTARG} ;;
      \? )
        echo "Invalid Option: -${OPTARG}" 1>&2
        exit 1 ;;
    esac
  done
  [ ! -n ${alignment} ] && echo "must specify alignment use: -i" && exit 1;
  [ ! -n ${workDir} ] && echo "must specify work_dir use: -d" && exit 1;
} 

main "$@"; exit 0;
