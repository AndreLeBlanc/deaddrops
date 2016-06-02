import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import createReducer from './reducers/createReducer';

const loggerMiddleware = createLogger();

export default function configureCreateStore(initialState) {
  if (process.env.NODE_ENV === 'production')
  {
  	return createStore(
	    createReducer,
	    initialState,
	    applyMiddleware(
	      thunkMiddleware
	    )
	  )
  }
  else
  {
  	return createStore(
	    createReducer,
	    initialState,
	    applyMiddleware(
	      thunkMiddleware,
	      loggerMiddleware
	    )
	  )
  }  
}