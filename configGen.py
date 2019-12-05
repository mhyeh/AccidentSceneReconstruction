import os
import sys
import cv2
import yaml

def configGen(filepath):
    config = {}
    fs = cv2.FileStorage("{}/camera_parameter.yml".format(filepath), cv2.FILE_STORAGE_READ)
    with open("{}/ORB.yaml".format(filepath), "w") as w:
        w.write("%YAML:1.0\n")
        fn = fs.getNode("camera_matrix")
        m = fn.mat()
        config["Camera.fx"] = float(m[0][0])
        config["Camera.fy"] = float(m[1][1])
        config["Camera.cx"] = float(m[0][2])
        config["Camera.cy"] = float(m[1][2])

        fn = fs.getNode("distortion_coefficients")
        m = fn.mat()
        config["Camera.k1"] = float(m[0])
        config["Camera.k2"] = float(m[1])
        config["Camera.p1"] = float(m[2])
        config["Camera.p2"] = float(m[3])
        config["Camera.k3"] = float(m[4])

        config["Camera.fps"] = 30.0
        config["Camera.RGB"] = 1

        config["ORBextractor.nFeatures"] = 5000
        config["ORBextractor.scaleFactor"] = 1.2
        config["ORBextractor.nLevels"] = 15
        config["ORBextractor.iniThFAST"] = 20
        config["ORBextractor.minThFAST"] = 7

        config["Viewer.KeyFrameSize"] = 0.05
        config["Viewer.KeyFrameLineWidth"] = 1
        config["Viewer.GraphLineWidth"] = 0.9
        config["Viewer.PointSize"] = 2
        config["Viewer.CameraSize"] = 0.08
        config["Viewer.CameraLineWidth"] = 3
        config["Viewer.ViewpointX"] = 0
        config["Viewer.ViewpointY"] = -0.7
        config["Viewer.ViewpointZ"] = -1.8
        config["Viewer.ViewpointF"] = 500

        yaml.dump(config, w, default_flow_style=False)


if __name__ == "__main__":
    configGen(sys.argv[1])    