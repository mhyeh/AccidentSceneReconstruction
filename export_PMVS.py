import os
import sys
import cv2
import shutil

def exportPMVS(filepath):
    id = 0
    fs = cv2.FileStorage("{}/camera_parameter.yml".format(filepath), cv2.FILE_STORAGE_READ)
    fn = fs.getNode("camera_matrix")
    cameraMat = fn.mat()
    fn = fs.getNode("distortion_coefficients")
    distortMat = fn.mat()
    with open("{}/KeyFrameTrajectory.txt".format(filepath), "r") as f:
        keyFrame = f.readline().strip()
        while keyFrame:
            f.readline()
            f.readline()
            f.readline()
            f.readline()
            f.readline()
            f.readline()

            with open("{}/txt/{:08d}.txt".format(filepath, id), "w") as w:
                w.write("CONTOUR\n")
                w.write(f.readline())
                w.write(f.readline())
                w.write(f.readline())

            img = cv2.imread("{}/images/{}.jpg".format(filepath, keyFrame))
            img = cv2.undistort(img, cameraMat, distortMat)
            cv2.imwrite("{}/visualize/{:08d}.jpg".format(filepath, id), img)

            id += 1
            keyFrame = f.readline().strip()
    
    with open("{}/model".format(filepath), "w") as w:
        w.write("level 1\n")
        w.write("csize 2\n")
        w.write("threshold 0.7\n")
        w.write("wsize 7\n")
        w.write("minImageNum 3\n")
        w.write("CPU 8\n")
        w.write("useVisData 0\n")
        w.write("sequence {:d}\n".format(id / 2))
        w.write("quad 2.5\n")
        w.write("maxAngle 10\n")
        w.write("timages -1 0 {:d}\n".format(id))
        w.write("oimages 0\n")

if __name__ == "__main__":
    exportPMVS(sys.argv[1])    