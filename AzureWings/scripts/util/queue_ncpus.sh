#!/bin/bash
# Should be use in PBS job scripts.
# Will not run correctly while outside PBS job.
ncpus=$(/opt/pbs/bin/qstat -xf ${PBS_JOBID} | grep "Resource_List.select" | sed 's/^.*=\| //g')
echo ${ncpus}
