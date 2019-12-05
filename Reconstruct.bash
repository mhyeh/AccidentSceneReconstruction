#!/bin/bash

ProjectName=$1

FilePath="project/${ProjectName}"
ORB_PATH="ORB_SLAM2"

python configGen.py $FilePath

export ROCON_RTSP_CAMERA_RELAY_URL=rtmp://localhost/stream
export ROCON_RTSP_CAMERA_RELAY_FILE_PATH=${ProjectName}

rosrun ORB_SLAM2 Mono ${ORB_PATH}/Vocabulary/ORBvoc.txt ${FilePath}/ORB.yaml ${ORB_PATH}/Examples/Monocular/GrubGyro.txt ${FilePath} &
roslaunch rocon_rtsp_camera_relay rtsp_camera_relay.launch &

mkdir ${FilePath}/txt
mkdir ${FilePath}/visualize
mkdir ${FilePath}/models

wait

python export_PMVS.py $FilePath
./CMVS-PMVS/program/build/main/pmvs2 ${FilePath}/ model