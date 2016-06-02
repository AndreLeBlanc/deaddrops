import { FileState } from './FileActions';
import {CONFIG, uploadCompleteAmount} from '../utils/Utils';
import Axios from 'axios';
import { hashHistory } from 'react-router';
import { updateFileState } from './FileActions';

export const DropActions = {
	UPDATE_DROP_STATE: "UPDATE_DROP_STATE",
	UPDATE_LABEL: "UPDATE_LABEL"
}

export const DropState = {
	NONE: 'NONE',
	INITIALIZING_DROP: 'INITIALIZING_DROP',
	FINALIZING_DROP: "FINALIZING_DROP",
	DROP_FINALIZED: "DROP_FINALIZED",
	DROP_FINALIZED_FAILED: "DROP_FINALIZED_FAILED"
}

export function updateDropState(state) {
	return {
		type: DropActions.UPDATE_DROP_STATE,
		state
	}
}

export function updateLabel(label) {
	return {
		type: DropActions.UPDATE_LABEL,
		label
	}
}

function createUploadRequest(value, token, dispatch) {
			
	let formData = new FormData();

	formData.append("uploadfile", value.data);
	formData.append("token", token);
	formData.append("submit", "upload");
	dispatch(updateFileState(value, FileState.UPLOADING));

	return Axios({
		method: 'post',
		url: CONFIG.baseURL + '/upload',
		data: formData
	}).then((response) => {
		dispatch(updateFileState(value, FileState.UPLOAD_COMPLETE));
	}).catch((response) => {
		dispatch(updateFileState(value, FileState.UPLOAD_FAILED));
	});					
}


export function initializeDrop(){
	return (dispatch, getState) => {
		let nrOfFiles = 0;
		let requests = [];
		dispatch(updateDropState(DropState.INITIALIZING_DROP));

		getState().files.files.forEach((value, key) => {
			if (value.state === FileState.UPLOAD_COMPLETE || value.data.size > getState().files.maxSize) {
				return;
			}
			requests.push(createUploadRequest(value, getState().essentials.key, dispatch));
			nrOfFiles++;					
		});
		Axios.all(requests).then(Axios.spread((acct, perms) => {
			if (nrOfFiles === uploadCompleteAmount(getState().files.files)) {
				dispatch(finalizeDrop());
			}
		}));
	}
}

export function finalizeDrop() {
	return (dispatch, getState) => {
		let finalizeJSON = {};
		finalizeJSON['Token'] = getState().essentials.key;
		finalizeJSON['StashName'] = getState().essentials.label;
		finalizeJSON['Lifetime'] = 60;
		finalizeJSON['Files'] = [];
		//dispatch(toggleIsOpen());
		getState().files.files.forEach((value, key) => {
			if (value.data.size > getState().files.maxSize) {
				return;
			}
			let fileJSON = {};
			fileJSON['Fname'] = value.data.name;
			fileJSON['Size'] = value.data.size;
			fileJSON['Type'] = value.data.type;
			fileJSON['Download'] = 3;

			finalizeJSON.Files.push(fileJSON);
		});

		dispatch(updateDropState(DropState.FINALIZING_DROP));
		Axios({
			  method: 'post',
			  url: CONFIG.baseURL + '/finalize',
			  data: JSON.stringify(finalizeJSON),
			  headers: {
			    'content-type': 'application/json'
			  }
		}).then((response) => {
			if(response.data.Token == getState().essentials.key)
			{
				dispatch(updateDropState(DropState.DROP_FINALIZED));
			}
			else
			{
				console.log(response);
				dispatch(updateDropState(DropState.DROP_FINALIZED_FAILED));
			}			
		}).catch((response) => {
			console.log(response);
			dispatch(updateDropState(DropState.DROP_FINALIZED_FAILED));
		});
	}
}