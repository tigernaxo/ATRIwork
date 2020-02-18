#!/bin/bash
currentDir=$(dirname $(readlink -f ${0}))
. ${currentDir}/env.sh ${currentDir}/env.sh

main(){
  echo "exporting conda envs"

  # export base
  exportEnv ${AWConda} ${AWCondaYmls}/base.yml

  # export envs
  envs=($(ls ${AWConda}/envs/))
  for env in "${envs[@]}"
  do
    exportEnv ${AWConda}/envs/${env} ${AWCondaYmls}/${env}.yml
  done
}

exportEnv(){
  cmd="conda env export -p ${1} -f ${2}"
  echo ${cmd}
  $(${cmd})
}

main "$@" ; exit 0
