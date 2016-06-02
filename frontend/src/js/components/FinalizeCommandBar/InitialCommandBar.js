import React from 'react';
import { connect } from 'react-redux';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import { initializeDrop, DropState } from '../../actions/Finalize';
import RippleButton from '../RippleButton';

const InitialCommandBar = (props) => {

  const buttonStyle = {
    'height': '30px'
  };

  return (
    <div class="initialCommandBar">
      <div class="summary">
        <h1>Files: {props.nrOfFiles + " / " + props.maxFiles}</h1>
        <h1>Size: {props.size}</h1> 
      </div>
      <div class="finalizeButtonContainer">
        <ReactCSSTransitionGroup transitionName="file" transitionEnterTimeout={350} transitionLeaveTimeout={200}>
        {getPotentialDropButton(props.nrOfFiles, props.initializeDrop, props.dropState)}
        </ReactCSSTransitionGroup>      
        <label class="makeButton stylized green" style={buttonStyle} for="fileinput" title="Add files">Add</label>
      </div>
    </div>
  )
}

function getPotentialDropButton(nrOfFiles, action, dropState) {
  const buttonStyle = {
      'minHeight': '30px',
      'marginRight': '10px'
    };

  if (nrOfFiles > 0 && dropState === DropState.DROP_FINALIZED_FAILED) {
    return (<RippleButton cssClass="makeButton stylized amber" customStyle={buttonStyle} action={action} label="Retry"></RippleButton>);
  } else if (nrOfFiles > 0) {
    return (<RippleButton cssClass="makeButton stylized" customStyle={buttonStyle} action={action} label="Drop"></RippleButton>);
  }
}


const mapDispatchToProps = (dispatch) => {
  return {
    initializeDrop: (e) => {
      dispatch(initializeDrop());
    },
  }
}

const Initial = connect(
  null,
  mapDispatchToProps
)(InitialCommandBar);

export default Initial;