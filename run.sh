#!/bin/bash

cat <<EOF

   ..:: cloudvet testing services ::..

Thank you for choosing cloudvet for all your
testing needs.  We are honored to serve you.

First, we're going to run a pre-flight check,
to validate the sanity of this environment.

Here we go!

EOF

rc=0
echo -n "Checking CF_USERNAME... "
if [[ -z ${CF_USERNAME:-} ]]; then
  echo "MISSING!"
  rc=1
else
  echo "${CF_USERNAME}"
fi

echo -n "Checking CF_PASSWORD... "
if [[ -z ${CF_PASSWORD:-} ]]; then
  echo "MISSING!"
  rc=1
else
  echo "(present)"
fi

echo -n "Checking CF_API_URL... "
if [[ -z ${CF_API_URL:-} ]]; then
  echo "MISSING!"
  rc=1
else
  echo "${CF_API_URL}"
fi

echo -n "Checking CF_DOMAIN ... "
if [[ -z ${CF_DOMAIN:-} ]]; then
  echo "MISSING!"
  rc=1
else
  echo "${CF_DOMAIN}"
fi

echo -n "Checking CF_ORG... "
if [[ -z ${CF_ORG:-} ]]; then
  echo "MISSING!"
  rc=1
else
  echo "${CF_ORG}"
fi

echo -n "Checking CF_SPACE... "
if [[ -z ${CF_SPACE:-} ]]; then
  echo "MISSING!"
  rc=1
else
  echo "${CF_SPACE}"
fi

n=0
echo "Checking if we should run MySQL tests... "
if [[ -z ${MYSQL:-} ]]; then
  echo "no"
else
  echo "yes"
  n=1
fi
echo "Checking if we should run Redis tests... "
if [[ -z ${REDIS:-} ]]; then
  echo "no"
else
  echo "yes"
  n=1
fi
echo "Checking if we should run RabbitMQ tests... "
if [[ -z ${RABBITMQ:-} ]]; then
  echo "no"
else
  echo "yes"
  n=1
fi

if [[ ${n} -eq 0 ]]; then
  echo
  echo "No test suites have been specified!"
  rc=1
fi

if [[ ${rc} -ne 0 ]]; then
  echo
  echo "Pre-flight checks FAILED."
  echo "Please correct these errors and try again."
  exit 1
fi

cat >Proc <<EOF
web: ./cloudvet
EOF

set -e
cf api ${CF_API_URL}
cf login -u ${CF_USERNAME} -p ${CF_PASSWORD}
cf target -o ${CF_ORG} -s ${CF_SPACE}
cf delete cloudvet || true
cf push -b binary_buildpack -d ${CF_DOMAIN} --no-start cloudvet
if [[ -n ${REDIS:-} ]]; then
  cf delete-service cloudvet-redis || true
  cf create-service ${REDIS} cloudvet-redis
  cf bind-service cloudvet cloudvet-redis
fi
if [[ -n ${MYSQL:-} ]]; then
  cf delete-service cloudvet-mysql || true
  cf create-service ${MYSQL} cloudvet-mysql
  cf bind-service cloudvet cloudvet-mysql
fi
cf start cloudvet
curl --fail https://cloudvet.${CF_DOMAIN}/ping
if [[ -n ${REDIS:-} ]]; then
  curl --fail https://cloudvet.${CF_DOMAIN}/mysql
fi
if [[ -n ${MYSQL:-} ]]; then
  curl --fail https://cloudvet.${CF_DOMAIN}/redis
fi

cf delete cloudvet
cf delete-service cloudvet-redis || true
cf delete-service cloudvet-mysql || true
exit 0
