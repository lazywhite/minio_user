#!/bin/sh

set -x

ENDPOINT="${ENDPOINT:-http://minio.storage:9000}"
# POLICY="${POLICY:-getonly}"
POLICY="${POLICY:-1}"
USER="${MC_USER:-avlcloud}"
ADMIN="avl-admin"

MINIO_GROUP=${MINIO_GROUP:-default}

CONFIG_DIR="/home/jovyan/.mc"
MC="mc -C ${CONFIG_DIR}"

if [ "${USER}" == "avlcloud" ]; then
    if [ -n "${JUPYTERHUB_USER}" ]; then
        USER="${JUPYTERHUB_USER}"
    fi
fi


# When start single user (jupyterlab), we reconfig the mc config file
rm -rf /home/jovyan/.mc
${MC} config host add ${ADMIN} ${ENDPOINT} "${ACCESS_KEY}" "${SECRET_KEY}" --api s3v4
PASSWD=$(openssl rand -hex 18)
# ${MC} admin user add ${ADMIN} ${USER} ${PASSWD} ${POLICY}
#python3 /create.py ${USER} ${PASSWD} ${POLICY}
/usr/local/bin/minio_user -a add -u ${USER} -p ${PASSWD} -g ${MINIO_GROUP}
# echo "${MC} config host add ${USER} ${ENDPOINT} ${USER} ${PASSWD} --api s3v4"
${MC} config host add ${USER} ${ENDPOINT} ${USER} ${PASSWD} --api s3v4
for i in gcs local play s3 ${ADMIN}; do ${MC} config host rm "${i}"; done
rm -rf /home/jovyan/.mc/config.json.old
# For avlfs mount
echo "${USER}:${PASSWD}" > /home/jovyan/.mc/.passwd-s3fs

cat >> /home/jovyan/.mc/.credentials << EOF
[default]
aws_access_key_id = ${USER}
aws_secret_access_key = ${PASSWD}
EOF

chown 1000:100 -R ${CONFIG_DIR}


