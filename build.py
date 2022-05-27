#!/usr/bin/python

import os
import sys
# import os.path as path

osarg=sys.argv[1]
tags={
    "default": "",
    "default-static": "no_sqlite3 moderncsqlite",
    "most": "most",
    "most-static": "most no_sqlite3"
}
osname={
    # "ubuntu-18.04": "GOARCH=amd64 GOOS=linux",
    # "macos-latest": "GOARCH=amd64 GOOS=darwin",
    # "windows-latest": "GOARCH=amd64 GOOS=windows",
    "ubuntu-18.04": "linux",
    "macos-latest": "darwin",
    "windows-latest": "windows",
}
cgoenabled={
    "default-static": "CGO_ENABLED=0",
    "most-static": "CGO_ENABLED=0",
    "default": "",
    "most": ""
}


print(osarg)

def runcmd(cmd):
    print("\n>>>", cmd, flush=True)
    exitcode = os.system(cmd)
    if exitcode != 0:
        sys.exit(exitcode)

def runbuilds():
    for tag in tags:
        print("\n>>>Building ", tag, flush=True)
        cmd = cgoenabled[tag]
        cmd += " go build -v "
        if osname[osarg] != "windows":
            cmd += " -ldflags '-s -w' "
        if osname[osarg] == "windows" and (tag == "default-static" or tag == "most-static"):
            # github fails on windows with CGO_ENABLED
            continue

        cmd += " -tags '%s' " % tags[tag]
        exefile='out/%s/sqlui' % tag
        if osname[osarg] == "windows":
            exefile += ".exe"
        cmd += " -o " + exefile
        cmd += " ."
        runcmd(cmd)
        zipfile='out/sqlui-%s-%s.zip' %(osname[osarg], tag)
        zipcmd = "7z a -tzip -mx=9 %s %s" %(zipfile, exefile)
        runcmd(zipcmd)
        # rename, removes path from sqlui
        exename = "sqlui"
        if osname[osarg] == "windows":
            exename = "sqlui.exe"
        zipcmd = "7z rn %s %s %s" %(zipfile, exefile, exename)
        runcmd(zipcmd)

runcmd("go env")
runcmd("go mod download")
runbuilds()
runcmd('git log --pretty=format:"%s" > changelog.txt')
runcmd('cat changelog.txt')

