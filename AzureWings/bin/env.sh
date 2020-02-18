#!/bin/bash
export AWRoot=$(readlink -f $(dirname  ${1})/../)
export AWBin=${AWRoot}/bin
export AWPkg=${AWRoot}/pkg
export AWConda=${AWRoot}/conda-envs
export AWCondaActivae=${AWRoot}/conda-envs/bin/activate
export AWCondaDeactivae=${AWRoot}/conda-envs/bin/deactivate
export AWCondaYmls=${AWRoot}/conda-ymls
export AWScripts=${AWRoot}/scripts
export AWScriptsPbs=${AWScripts}/pbs
export AWScriptsPipe=${AWScripts}/pipe
export AWScriptsUnit=${AWScripts}/unit
export AWScriptsUtil=${AWScripts}/util

# Whole area pkg setting
export beastgen=${AWPkg}/BEASTGen_v1.0.2/bin/beastgen

# set util
export pbsSnippet=${AWScriptsUtil}/pbsSnippet.sh
export queue_ncpus=${AWScriptsUtil}/queue_ncpus.sh
export timeStamp=${AWScriptsUtil}/timeStamp.sh

