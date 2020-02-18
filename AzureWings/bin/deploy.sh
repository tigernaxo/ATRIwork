#!/bin/bash
currentDir=$(dirname $(readlink -f ${0}))
. ${currentDir}/env.sh ${currentDir}/env.sh

main(){
  deployBin
  deployPkg
}

# give all AW bin excutable
deployBin(){
  cd ${AWBin}
  chmod a+x *
}

# TODO:
# Install envs by ymls
deployEnv(){

}

# Unpackage packages that come with AW (not in conda)
deployPkg(){
  cd ${AWPkg}
  pkgs=($(ls *.tgz))
  for pkg in ${pkgs[@]}
  do
    tar zxf ./${pkg}
  done
}

main "$@"; exit 0
