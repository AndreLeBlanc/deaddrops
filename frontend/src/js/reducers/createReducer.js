import { combineReducers } from 'redux';
import { DropState, DropActions } from '../actions/Finalize';
import files from './fileReducer';


const initialState = {
	key: null,
	label: "",
	state: "initial",
	errorStateMsg: "",
	dropState: DropState.NONE
} 

function createReducer(state=initialState, action) {
	switch (action.type) {
		case "SET_KEY":
			return Object.assign({}, state, {
				key: action.key
			})
		case "UPDATE_CREATE_STATE":
			return Object.assign({}, state, {
				state: action.state
			})
		case "SET_ERROR_MESSAGE":
		return Object.assign({}, state, {
			errorStateMsg: action.message
		})
		case DropActions.UPDATE_DROP_STATE:
			return Object.assign({}, state, {
				dropState: action.state
			})
		case DropActions.UPDATE_LABEL:
			return Object.assign({}, state, {
				label: action.label
			})
		default:
			return state;
	}
}

const create = combineReducers({
	essentials: createReducer,
	files
});

export default create;