#!/usr/bin/python

import os
import sys

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
    "default-static": "CGO_ENABLED=1",
    "most-static": "CGO_ENABLED=1",
    "default": "",
    "most": ""
}


print(osarg)

def runcmd(cmd):
    print("\n>>>", cmd)
    exitcode = os.system(cmd)
    if exitcode != 0:
        sys.exit(exitcode)

def runbuilds():
    for tag in tags:
        print("\n>>>Building ", tag)
        # cmd = osenvs[osarg]
        cmd = cgoenabled[tag]
        cmd += " go build "
        cmd += " -ldflags '-s -w'"
        cmd += " -tags '%s'" % tags[tag]
        exefile='out/%s/sqlui' % tag
        if osname[osarg] == "windows":
            exefile += ".exe"
        cmd += " -o " + exefile
        cmd += " ."
        runcmd(cmd)
        zipfile='out/sqlui-%s-%s.zip' %(osname[osarg], tag)
        zipcmd = "zip -9 --junk-paths %s %s" %(zipfile, exefile)
        runcmd(zipcmd)

runcmd("go env")
runcmd("go mod download")
runbuilds()
runcmd('git log --pretty=format:"%s" > changelog.txt')
runcmd('cat changelog.txt')

