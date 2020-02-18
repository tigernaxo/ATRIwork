#!/bin/bash
currentDir=$(dirname $(readlink -f ${0}))
. ${currentDir}/env.sh ${currentDir}/env.sh

pkgs=($(ls -d ${AWPkg}/*/ | xargs -n1 -I{} basename {}))

main(){
  cd ${AWPkg}
  for pkg in "${pkgs[@]}"; 
  do
    rmThenPack ${pkg}
  done
}

rmThenPack(){
  echo "[backupPkg] packing ${pkg}"
  [ -f ${1}.tgz ] && rm ${1}.tgz 
  tar zcf ${1}.tgz ${1}/
}

main "$@"; exit 0
