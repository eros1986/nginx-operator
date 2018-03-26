#!/bin/bash

if [ -z ${DEBUG_LEVEL}]; then
    DEBUG_LEVEL = 5
fi

resyncSeconds=${RESYNC_SECONDS#0}

cmd="/app/operator server -l ${DEBUG_LEVEL}
    --watchNamespace ${WATCH_NAMESPACE}
    --resyncSeconds ${resyncSeconds}
"

echpo "cmd: " ${cmd}
eval ${cmd}