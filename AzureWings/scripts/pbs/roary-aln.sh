#!/bin/bash

currentDir=$( dirname $( readlink -f ${0}  ))
env="${currentDir}/../../bin/env.sh"
. ${env} ${env}

roary="${AWScriptsUnit}/roary.sh"

main(){
  getMyOpt "$@"

  # mkdir -p ${workDir}
  jobID=$(qsub \
    $(${pbsSnippet} ${queue} $(basename ${workDir}) ${workDir}) \
    ${pbsAppend} \
    -- ${roary} ${workDir} ${gffs})

  echo ${jobID}
}

####################################
# set arguments
getMyOpt(){
  if [ $# -eq 0 ] ;then
    usage
    exit 0
  fi

  while getopts ":d:q:w:h" opt; do
    case ${opt} in
      h )
        usage
        exit 0 ;;
      d )
        workDir=$( readlink -f ${OPTARG} )/roary_$(date +"%Y%m%d%H%M%S") ;;
      q )
        queue=${OPTARG} ;;
      w )
        afterokString="-W depend=afterok:${OPTARG}" ;;
      \? )
        echo "Invalid Option: -${OPTARG}" 1>&2
        exit 1 ;;
    esac
  done

  ######################################
  # 給queue, workDir, logFile 設置預設值
  logFile="${workDir}/roary.log"
  defaultQueue="ngs16G"
  # [ "${workDir}" ] || workDir=${currentDir}
  # [ "${queue}"  ] || queue=${defaultQueue}

  # 取得gffs
  shift $(($OPTIND - 1))
  gffs=${@}

  # 如果都沒有gff就退出
  [ $(echo ${gffs[@]} | wc -c) -eq 1 ] && echo "no gff found in input, exit." && exit 1

  # 設置額外的pbs命令
  [ -n "${afterokString}"  ] && pbsAppend="${afterokString}"
}

########################################
# give argument hint if no argument gived.
usage() {
      echo ""
      echo "Usage:"
      echo "  $(basename ${0}) [-h] [-q queue_name] [-d outdir_prefix] [-w afterok_string]"
      echo ""
      echo "Option:"
      echo "  -h    Show this message."
      echo "  -q    Specify queue."
      echo "  -d    output dir(time stamp will be add to avoid name confict)"
      echo "  -w    Wait the job on the list all done before start. (ex job1:job2)"
      echo ""
      echo "Example:"
      echo "  $(basename ${0}) -q ngsTest \\" 
      echo "    -d $(pwd) \\"
      echo "    -w srv.12345:srv:12346:srv.12347 \\"
      echo "    a.gff b.gff c.gff"
      echo ""
}

main "$@"; exit 0
