#!/bin/bash
currentDir=$(dirname $(readlink -f ${0}))
env="${currentDir}/../../bin/env.sh"
. ${env} ${env}

main(){

  getMyOpt "$@"
  echo "$(${timeStamp}) [$(basename ${0})]start"
  . ${AWCondaActivae} ${AWConda}/envs/beast2

  trees=$(echo $(basename ${xml}) | sed 's/\..*$/.trees/')
  tree=$(echo $(basename ${xml}) | sed 's/\..*$/.tree/')
  cd ${workDir}
  beast -beagle -threads $(${queue_ncpus}) ${xml}
  treeannotator -burnin 10 ${trees} ${tree}

  . ${AWCondaDeactivae}
}

getMyOpt(){
  while getopts ":d:i:h" opt; do
    case ${opt} in
      d )
        workDir=${OPTARG} ;;
      i )
        xml=${OPTARG} ;;
      \? )
        echo "Invalid Option: -${OPTARG}" 1>&2
        exit 1 ;;
    esac
  done
} 

main "$@"; exit 0;
