import { combineReducers } from 'redux';
import { fileReducer } from './fileReducer'; 
import { FileState } from '../actions/FileActions'; 

const initialState = {
	dropData: {},
	openState: "initial",
	errorStateMsg: ""
} 

function openReducer(state=initialState, action) {
	switch (action.type) {
		case "DROP_RECIEVED":
			return Object.assign({}, state, {
				dropData: action.dropData
			});
		case "UPDATE_OPEN_STATE":
			return Object.assign({}, state, {
				openState: action.state
			});
		case "SET_ERROR_MESSAGE":
			return Object.assign({}, state, {
				errorStateMsg: action.message
			});
		case "UPDATE_FILE_STATE":
			let newFiles = state.dropData.Files;
			for(let i=0; i < newFiles.length; i++)
			{
				if(newFiles[i].Id === action.fileID)
				{
					newFiles[i].state = FileState.FILE_DOWNLOADING; 
				}
			}
			return Object.assign({}, state, {
				dropData: Object.assign({}, state.dropData, {Files: newFiles})
			});
		default:
			return state;
	}
}

const open = combineReducers({
	essentials: openReducer
});

export default open;