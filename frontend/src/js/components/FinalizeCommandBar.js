import React from 'react';
import { connect } from 'react-redux';
import { updateDropState, DropState, initializeDrop } from '../actions/Finalize';
import { uploadFiles } from '../actions/FileActions';
import { uploadCompleteAmount, calcSize } from '../utils/Utils';
import InitialCommandBar from './FinalizeCommandBar/InitialCommandBar';
import FinalizedCommandBar from './FinalizeCommandBar/FinalizedCommandBar';


function getCommandBar(dropState, nrOfFiles, maxFiles, size) {
  if (dropState === DropState.DROP_FINALIZED) {
    return <FinalizedCommandBar />
  } else if(dropState === DropState.DROP_FINALIZED_FAILED) {
    return <InitialCommandBar nrOfFiles={nrOfFiles} maxFiles={maxFiles} size={size} dropState={dropState} />
  } else {
    return <InitialCommandBar nrOfFiles={nrOfFiles} maxFiles={maxFiles} size={size} dropState={dropState} />
  }

}

const FinalizeCommandBar = (props) => {      
  return (
    <div>
      {getCommandBar(props.dropState, props.nrOfFiles, props.maxFiles, props.size)}
    </div>
	)
}

const mapStateToProps = (state) => {
  return {
  	dropState: state.essentials.dropState
  }
}

const Finalize = connect(
  mapStateToProps
)(FinalizeCommandBar);

export default Finalize;