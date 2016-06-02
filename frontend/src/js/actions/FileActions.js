import {CONFIG} from '../utils/Utils';
import Axios from 'axios';


export const FileActions = {
	FILES_ADDED: "FILES_ADDED",
	FILES_DROPPED: "FILES_DROPPED",
	UPLOAD_FILES: "UPLOAD_FILES",
	UPLOAD_FILE: "UPLOAD_FILE",
	UPDATE_FILE_STATE: "UPDATE_FILE_STATE",
	DELETE_FILE: "DELETE_FILE",
	TOGGLE_IS_OPEN: "TOGGLE_IS_OPEN"
}

export const FileState = {
	UPLOADING: "UPLOADING",
	UPLOAD_COMPLETE: "UPLOAD_COMPLETE",
	UPLOAD_FAILED: "UPLOAD_FAILED",
	UPLOAD_WAITING: "UPLOAD_WAITING",
	FILE_TOLARGE: "FILE_TOLARGE",
	FILE_READY: "FILE_READY",
	FILE_DOWNLOADING: "FILE_DOWNLOADING"
}


export function deleteFile(fileName) {
	return {
		type: FileActions.DELETE_FILE,
		fileName
	}
}

export function updateFileState(file, state) {
	return {
		type: FileActions.UPDATE_FILE_STATE,
		file, 
		state
	}
}

export function updateOpenFileState(fileID, state) {
	return {
		type: FileActions.UPDATE_FILE_STATE,
		fileID, 
		state
	}
}

export function toggleIsOpen() {
	return {
		type: FileActions.TOGGLE_IS_OPEN
	}
}

function uploadFile (file) {
	return {
		type: FileActions.UPLOAD_FILE,
		file
	}
}

function dropDatafetched (dropData) {
	return {
		type: "DROP_RECIEVED", 
  	 	dropData: dropData
	}
}

export function filesAdded(files) {
	return (dispatch, getState) => {
		const currFileCount = getState().files.files.size;
		const max = getState().files.maxFiles;
		const nrToAdd = max - currFileCount;

		let newFiles = [];
		let i = 0;

		for (let key in files) {
			if (i === max - currFileCount) {
				break;
			}
			if (files.hasOwnProperty(key)) {
				newFiles.push(files[key]);
			}
			i++;
		}
		dispatch({
			type: FileActions.FILES_ADDED,
			files: newFiles
		});
	}
}

export function uploadFiles() {
	return (dispatch, getState) => {
		//dispatch(toggleIsOpen());
		getState().files.files.forEach((value, key) => {
			if (value.state === FileState.UPLOAD_COMPLETE || value.data.size > getState().files.maxSize) {
				return;
			}
			
			let formData = new FormData();

			formData.append("uploadfile", value.data);
			formData.append("token", getState().essentials.key);
			formData.append("submit", "upload");

			dispatch(updateFileState(value, FileState.UPLOADING));

			Axios({
				  method: 'post',
				  url: CONFIG.baseURL + '/upload',
				  data: formData
			}).then((response) => {
				dispatch(updateFileState(value, FileState.UPLOAD_COMPLETE));
			}).catch((response) => {
				dispatch(updateFileState(value, FileState.UPLOAD_FAILED));
			});
							
		})
	}
}

export function fetchDropData(token){
	return (dispatch, getState) => {

		Axios({
			  method: 'get',
			  url: CONFIG.baseURL + '/download/' + token
		}).then((response) => {
			dispatch(dropDatafetched(response.data));
			dispatch({type: "UPDATE_OPEN_STATE", state: "loaded"}); 
		}).catch((response) => {
			dispatch({type: "UPDATE_OPEN_STATE", state: "failed"});
			dispatch(
				{
					type: "SET_ERROR_MESSAGE", 
					message: "Unable to connect to server, Try again"
				}
			);	
		});
	}
}

export function downloadFile(downloadData) {
	return (dispatch, getState) => {
	
		//getState().essentials.dropData.Files.forEach((value, key) => {			
		//	if(value.Id === downloadData.FileID)
		//	{
		//
		//		//	}					
		//})
		//console.log(downloadData.StashID, downloadData.FileID);
		//const file = getState().essentials.dropData.Files.get(downloadData.fileName);
		//dispatch(updateFileState(value, FileState.FILE_DOWNLOADING));
		downloadFileFunc(CONFIG.baseURL + "/" + downloadData.StashID + "/" + downloadData.FileID);

	}
}

//SRC: https://github.com/PixelsCommander/Download-File-JS/blob/master/src/download.js
window.downloadFileFunc = function (sUrl) {

    //iOS devices do not support downloading. We have to inform user about this.
    if (/(iP)/g.test(navigator.userAgent)) {
        alert('Your device does not support files downloading. Please try again in desktop browser.');
        return false;
    }

    //If in Chrome or Safari - download via virtual link click
    if (window.downloadFileFunc.isChrome || window.downloadFileFunc.isSafari) {
        //Creating new link node.
        var link = document.createElement('a');
        link.href = sUrl;

        if (link.download !== undefined) {
            //Set HTML5 download attribute. This will prevent file from opening if supported.
            var fileName = sUrl.substring(sUrl.lastIndexOf('/') + 1, sUrl.length);
            link.download = fileName;
        }

        //Dispatching click event.
        if (document.createEvent) {
            var e = document.createEvent('MouseEvents');
            e.initEvent('click', true, true);
            link.dispatchEvent(e);
            return true;
        }
    }

    // Force file download (whether supported by server).
    if (sUrl.indexOf('?') === -1) {
        sUrl += '?download';
    }

    window.open(sUrl, '_self');
    return true;
}

window.downloadFileFunc.isChrome = navigator.userAgent.toLowerCase().indexOf('chrome') > -1;
window.downloadFileFunc.isSafari = navigator.userAgent.toLowerCase().indexOf('safari') > -1;