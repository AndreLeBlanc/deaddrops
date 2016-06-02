import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import openReducer from './reducers/openReducer';

const loggerMiddleware = createLogger();

export default function configureOpenStore(initialState) {
  if (process.env.NODE_ENV === 'production')
  {
  	return createStore(
	    openReducer,
	    initialState,
	    applyMiddleware(
	      thunkMiddleware
	    )
	  )
  }
  else
  {
  	return createStore(
	    openReducer,
	    initialState,
	    applyMiddleware(
	      thunkMiddleware,
	      loggerMiddleware
	    )
	  )
  }
}