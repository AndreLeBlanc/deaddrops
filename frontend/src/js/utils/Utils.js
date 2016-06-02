import React from 'react';
import File from '../components/File';
import { FileState } from '../actions/FileActions';
/*** Byte handling ***/


export function formatBytes(bytes,decimals) {
   if(bytes == 0) return '0 Byte';
   var k = 1024; // or 1024 for binary
   var dm = decimals + 1 || 3;
   var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
   var i = Math.floor(Math.log(bytes) / Math.log(k));
   return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

/**
 * File helper functions
 */

export function constructFileObjects(files) {
  let fileObjects = [];
  for(let v of files.keys()) {
    const file = files.get(v);
    fileObjects.push(
      <File key={file.data.name} name={file.data.name} state={file.state} size={file.data.size} />
    );
  }
  return fileObjects
}

export function constructFileObjectsFromServer(stashID, files) {
  if(!files || !stashID)
  {
    return;
  }

  let fileObjects = [];
  files.forEach((value, key) => {
    let file_state = FileState.FILE_READY;
    if(typeof value.state !== "undefined" && value.state != null)
      file_state = value.state;
    fileObjects.push(
      <File key={value.Id} stashID={stashID} fileID={value.Id} downloads={value.Download} name={value.Fname} state={file_state} size={value.Size} />
    );
  });
  return fileObjects
}

export function uploadCompleteAmount(files)
{
  let uploaded = 0;
  files.forEach(function (file) {
    if(file.state == FileState.UPLOAD_COMPLETE)
      uploaded++;  
  });

  return uploaded;
}

export function calcSize(files) {
  let size = 0;
  for(let v of files.keys()) {
      const file = files.get(v);
      if (file.state !== FileState.FILE_TOLARGE) {
        size += file.data.size;
      }
    }
    return formatBytes(size);
}

export function countValidFiles(files) {
  let count = 0;
  for(let v of files.keys()) {
      const file = files.get(v);
      if (file.state !== FileState.FILE_TOLARGE) {
        count += 1;
      }
    }
    return count;
}


function getConfig() {
 let CONFIG = {};
 if (process.env.NODE_ENV === 'production') {
   CONFIG.baseURL = 'http://deadrop.win:9090';
   CONFIG.baseLocalURL = 'http://deadrop.win';
 } else {
   CONFIG.baseURL = 'http://localhost:9090';
   CONFIG.baseLocalURL = 'http://localhost:8080';
 }
 return CONFIG;
}

export const CONFIG = getConfig();