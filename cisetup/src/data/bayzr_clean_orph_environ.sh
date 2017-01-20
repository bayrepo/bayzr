#!/bin/bash

export PATH=/sbin:/bin:/usr/sbin:/usr/bin

LOG_FILE="/var/log/bayzr-env-cleaner.log"

write_log()
{
 LOGTIME=`date "+%Y-%m-%d %H:%M:%S"`
 if [ "$LOG_FILE" == "" ]; then
  echo $LOGTIME": $1";
 else
  LOG=$LOG_FILE
  touch $LOG
  if [ ! -f $LOG ]; then echo "=====================================ERROR!! Cannot create log file $LOG. Exiting.====================================="; exit 1; fi
  echo $LOGTIME": $1" | tee -a $LOG;
 fi
}

write_log "========================================Begin of check orphaned ebvironments====================================================="

ls -1 /mnt/chroot/ | while read -r dir_nm
do
 FULL_NAME="/mnt/chroot/$dir_nm"
 write_log "Find file $FULL_NAME"
 if [ -d "$FULL_NAME" ]; then
  filename=$(basename "$FULL_NAME")
  extension="${filename##*.}"
  write_log "Check for process /usr/sbin/citool -task=$extension -task-run"
  is_found=$(ps -e -o pid= -o command= | grep "/usr/sbin/citool -task=$extension -task-run" | grep -v grep)
  if [ -z "$is_found" ]; then
   write_log "=====>Process /usr/sbin/citool -task=$extension -task-run not found. Clean environment $FULL_NAME"
   /usr/sbin/yumbootstrap --uninstall centos-7-mod /mnt/chroot/centos-7-mod."$extension"
   if [ -e "/mnt/chroot/centos-7-mod.$extension" ]; then
    rm -rf /mnt/chroot/centos-7-mod."$extension"
   fi
  else
   write_log "Process /usr/sbin/citool -task=$extension -task-run found"
  fi
 else
  write_log "===!>File $FULL_NAME not a directory. Remove it"
 fi
done

write_log "===========================================End of check orphaned ebvironments====================================================="
