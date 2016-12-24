#!/usr/bin/python

import os
import time
import logging
import pwd
import grp
import yumbootstrap.yum
import yumbootstrap.log
from shutil import copyfile

def createUser(username,uid,gid):
    os.system("/usr/sbin/groupadd " + username + " -g " + str(gid))
    return  os.system("/usr/sbin/useradd -s " + "/bin/bash "+ "-d "+ "/home/" + username + " -m " + username + " -u " + str(uid) + " -g " + str(gid))

#-----------------------------------------------------------------------------

logger = logging.getLogger()
logger.addHandler(yumbootstrap.log.ProgressHandler())
if os.environ['VERBOSE'] == 'true':
  logger.setLevel(logging.INFO)

#-----------------------------------------------------------------------------
out_user = os.environ['OUT_USER']
uid = pwd.getpwnam(out_user).pw_uid
gid = grp.getgrnam(out_user).gr_gid

if uid != os.getuid():

    print "Prepare chroot for non root user"
    os.chown(os.environ['TARGET'], uid, gid)


    real_root = os.open("/", os.O_RDONLY)

    os.chroot(os.environ['TARGET'])

    createUser(out_user, uid, gid)

    with open("/etc/sudoers", "a") as myfile:
        myfile.write("%s ALL = NOPASSWD : /usr/bin/yum, /usr/bin/rpm\n" % out_user)

    os.fchdir(real_root)
    os.chroot(".")

    os.mkdir( os.environ['TARGET'] + "/home/" + out_user + "/.ssh/", 0700 );
    copyfile("/root/config", os.environ['TARGET'] + "/home/" + out_user + "/.ssh/config")
    os.chown(os.environ['TARGET'] + "/home/" + out_user + "/.ssh/", uid, gid)
    os.chown(os.environ['TARGET'] + "/home/" + out_user + "/.ssh/config", uid, gid)

    # Back to old root
    os.close(real_root)



