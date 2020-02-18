#!/bin/bash

. ./env.sh ./env.sh

# >>> conda initialize >>>
# !! Contents within this block are managed by 'conda init' !!
__conda_setup="$(${AWConda}/bin/conda shell.bash hook 2> /dev/null)"
if [ $? -eq 0 ]; then
    eval "$__conda_setup"
else
    if [ -f "${AWConda}/etc/profile.d/conda.sh" ]; then
        . "${AWConda}/etc/profile.d/conda.sh"
    else
        export PATH="${AWConda}/bin:$PATH"
    fi
fi
unset __conda_setup
# <<< conda initialize <<<

