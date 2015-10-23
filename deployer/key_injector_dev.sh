#!/bin/sh

echo "***************************************"
echo "Beginning key injector...."
echo "***************************************"


while getopts :ot opts; do 
  case ${opts} in 
    o) OVERWRITE=true ;;
    t) TEST=true ;;
  esac
done

if [[ $TEST == true ]]; then
  ENV="test"
  CORE_ETCD_DIR="foxcommerce.com/test/core/config"
  CORE_PG_URL="postgres://localhost:5432/foxcomm_test?sslmode=disable"
else
  ENV="development"
  CORE_ETCD_DIR="foxcommerce.com/core/config"
  CORE_PG_URL="postgres://localhost:5432/foxcomm?sslmode=disable"
fi

# etcd_ls=$(etcdctl ls | grep "Error" | awk '{print $1}')
# if [[ $etcd_ls == *Cannot* ]]; then
#   echo "Etcd is not running.  Excitng"
#   exit
# else 
#   echo "Etcd is running; continuing..."
# fi

#TODO: Add the ability to target a 'target' or rebuild this in Ansible
if [[ $OVERWRITE == true  ]]; then
  echo "Removing core directory"
  etcdctl rm --recursive $CORE_ETCD_DIR
fi
etcdctl mkdir "$CORE_ETCD_DIR"

etcdctl set $CORE_ETCD_DIR/FC_CORE_DB_URL $CORE_PG_URL
echo "Core DB URL Set"

etcdctl mkdir "endpoints/origin_frontend"
if [[ $OVERWRITE == true  ]]; then
  echo "Removing origin_frontend"
  etcdctl rm "endpoints/origin_frontend/localhost:8080"
fi

etcdctl set "endpoints/origin_frontend/localhost:8080" ""
echo "Origin_Frontend endpoint set to localhost:8080"

etcdctl mkdir "endpoints/origin_backend"
if [[ $OVERWRITE == true  ]]; then
  echo "Removing origin_backend"
  etcdctl rm "endpoints/origin_backend/localhost:8080"
fi
etcdctl set "endpoints/origin_backend/localhost:8080" ""
echo "Origin_Backend endpoint set to localhost:8080"

if [[ $OVERWRITE == true  ]]; then
  echo "Removing environment"
  etcdctl rm "$CORE_ETCD_DIR/FC_ENV"
fi
etcdctl set $CORE_ETCD_DIR/FC_ENV $ENV
echo "Core Environment set"


echo "Done..."
