#!/bin/bash

# misc setting
projectID="MST107119"

# Task setting
queue=${1}
jobID=${2}
workDir=${3}

# Queue setting
# now can only send ngsTest to ngs96G
declare -A queueCore=( 
  ["ngsTest"]=1 
  ["ngs4G"]=1 
  ["ngs8G"]=2 
  ["ngs16G"]=4 
  ["ngs48G"]=10 
  ["ngs96G"]=20 
  )


constQueueStr=" 
  -q ${queue} -P ${projectID} \
  -W group_list=${projectID} \
  -l select=1:ncpus=${queueCore[$queue]} \
  -l place=pack \
  -o ${workDir}/pbs.out -e ${workDir}/pbs.err \
  -m e \
  -N ${jobID} "

echo ${constQueueStr}
