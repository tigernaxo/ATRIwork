#!/bin/bash

currentDir=$( dirname $( readlink -f ${0}  ))
. ${currentDir}/env.sh

# enabled job setting
declare -A enabledMap=(
  ["unicycler-prokka"]="unicycler-prokka.sh"
  ["roary-aln"]="roary-aln.sh"
  )

# switch $#
# if no argument received, hint user to choose job.
# if asigned program do not exist, error.
# if argument received, pass direct to pbs control script.
if [ $# -eq 0 ] ;then
  echo ""
  echo "select a job to get description, eg: ${0} unicycler-prokka"
  echo -e "\njob list:\n"
  for task in "${!enabledMap[@]}"; do
    echo -e "\t"${task}
  done
  echo ""
elif [ -z ${enabledMap[${1}]} ]; then
  echo ""
  echo "do not have an valid choice called ${1}"
  echo ""
else
  ${AWScriptsPbs}/${enabledMap[${1}]} ${@:2} 
fi

