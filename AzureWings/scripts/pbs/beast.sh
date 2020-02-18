#!/bin/bash
currentDir=$(dirname $(readlink -f ${0}))
env="${currentDir}/../../bin/env.sh"
. ${env} ${env}

beastUnit="${AWScriptsUnit}/beast2.sh"
xmlGenerator="${AWPkg}/BEASTGen_v1.0.2_script/mkXML2.5.2-nt-HKY-Strict-Yule.sh"

main(){

  getMyOpts "$@"
  logFile="${workDir}/beast.log"
  name=$(basename $(echo ${alignment} | sed 's/\..*$//g'))

  mkdir -p ${workDir}
  cd ${workDir}

  # need to generate xml
  cmd="${xmlGenerator} -i ${alignment} -t ${timeTsv} -o ${name}.xml"
  $(${cmd})

  # up is ok
  # send xml to beast and annotator
  qsub \
    $(${pbsSnippet} ${queue} beast-${name} ${workDir}) \
    -- ${beastUnit} -d ${workDir} -i ${workDir}/${name}.xml
}

getMyOpts(){
  if [ $# -eq 0 ] ;then
    usage
    exit 0
  fi

  while getopts ":i:t:d:q:w:h" opt; do
    case ${opt} in
      h )
        usage
        exit 0 ;;
      t )
        timeTsv=${OPTARG} ;;
      d )
        workDir=${OPTARG} ;;
      i )
        alignment=${OPTARG} ;;
      q )
        queue=${OPTARG} ;;
      w )
        afterokString="-W depend=afterok:${OPTARG}" ;;
      \? )
        echo "Invalid Option: -${OPTARG}" 1>&2
        exit 1 ;;
    esac
  done
  [ ! -n "${workDir}" ] && echo "need option -d work_dir" && exit 1;
  [ ! -n "${alignment}" ] && echo "need option -i alignment" && exit 1;
  [ ! -n "${timeTsv}" ] && echo "need option -t time.tsv" && exit 1;
  [ ! -n "${queue}" ] && echo "using default queue: ngs48G" && queue="ngs48G";
}

usage(){
  echo 
  echo "Usage:"
  echo "  $(basename ${0}) [-h] [-q queue_name] [-d outdir_prefix] [-i nucleotide_alignment] [-t id_year.tsv] [-w afterok_string]"
  echo 
  echo "Option:"
  echo "  -h    Show this message."
  echo "  -i    alignment feed to beast2."
  echo "  -t    id_year.tsv"
  echo "  -q    Specify queue."
  echo "  -d    output dir(time stamp will be add to avoid name confict)"
  echo "  -w    Wait the job on the list all done before start. (ex job1:job2)"
  echo 
  echo "Example:"
  echo "  $(basename ${0}) -q ngsTest -d $(pwd) \\" 
  echo "    -i test.fasta -t test.tsv \\"
  echo "    -w srv.12345:srv:12346:srv.12347 "
}

main "$@" ; exit 0
