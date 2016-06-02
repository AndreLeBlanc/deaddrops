import React from 'react';
import { connect } from 'react-redux';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import { downloadFiles, FileState } from '../actions/FileActions';
import { constructFileObjectsFromServer, calcSize, countValidFiles } from '../utils/Utils';

const OpenDrop = (props) => {

		const d = props.preventDefault;
		const f = props.files;
		return (
            <div class="openzone">
              <div class="fileWrapper">
                <ReactCSSTransitionGroup transitionName="file" transitionEnterTimeout={350} transitionLeaveTimeout={200}>
                  {constructFileObjectsFromServer(props.dropData.Token, props.dropData.Files)}
                </ReactCSSTransitionGroup>
              </div>
            </div>
		);
}

const mapStateToProps = (state) => {
  return {
    dropData: state.essentials.dropData
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    preventDefault: (e) => {
    	e.preventDefault();
    	e.stopPropagation();
    }
  }
}

const OpenDropZone = connect(
  mapStateToProps,
  mapDispatchToProps
)(OpenDrop);


export default OpenDropZone;
      