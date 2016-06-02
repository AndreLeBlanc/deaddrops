/**
 * @file        Defines the App component.
 */

import React from 'react';
/**
 * Base component of the application.
 * @param  {Object} props Properties for the App component.
 * @return {React.Component} Returns a react component.
 */
const App = (props) => (
	<div class="appContainer">
		{props.children}
	</div>
)

export default App