/**
 * @header       Defines the header component
 *             the ui. 
 */
import React from 'react';
import { connect } from 'react-redux';
import InfoCard from './InfoCard';


/**
 * HeaderObj component
 * @param  {Object} props Properties for the HeaderObj component.
 * @return {React.Component} Returns a react component.
 */
const HeaderObj = (props) => {
	return (
		<div class="header">
			<InfoCard iconName="cloud_upload" header="Upload" text="Select 5 files up to 25MB to upload and share." />
			<InfoCard iconName="work" header="Drop" text="Pleased with your selection of files, Drop it. It's now ready for sharing." />
			<InfoCard iconName="share" header="Share" text="Easily share your Drop with whomever you want. The files are waiting!" />
		</div>
	)
}

export default HeaderObj;