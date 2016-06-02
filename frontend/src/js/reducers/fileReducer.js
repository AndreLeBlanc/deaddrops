import { FileActions, FileState } from '../actions/FileActions';

const initialFileState = {
	isOpen: true,
	maxFiles: 5,
	maxSize: (25*1024*1024),
	files: new Map()
}

function addFiles(fileContainer, files) {
	//console.log(files);
	let newFiles = fileContainer;
	for (let file of files) {
		if (typeof file === 'object') {
			let state = (file.size > initialFileState.maxSize) ? FileState.FILE_TOLARGE : FileState.UPLOAD_WAITING;
			newFiles.set(file.name, {
				state: state,
				data: file
			});	
		}		
	}
	return newFiles;
}

export default function fileReducer(state=initialFileState, action) {
	switch (action.type) {

		case FileActions.FILES_DROPPED:
		case FileActions.FILES_ADDED:
				let newFiles_add = addFiles(new Map(state.files), action.files);
				return Object.assign({}, state, {
					files: newFiles_add
				});
		case FileActions.DELETE_FILE:
				let newFiles_del = new Map(state.files);
				newFiles_del.delete(action.fileName);
				return Object.assign({}, state, {
					files: newFiles_del
				});
		case FileActions.UPDATE_FILE_STATE:
			let newFiles_Upl = new Map(state.files);
			newFiles_Upl.set(action.file.data.name, {state: action.state, data: action.file.data});
			return Object.assign({}, state, {
				files: newFiles_Upl
			});
		case FileActions.TOGGLE_IS_OPEN:
			return Object.assign({}, state, {
				isOpen: !state.isOpen
			});

		default:
			return state;
	}
	return state;
}